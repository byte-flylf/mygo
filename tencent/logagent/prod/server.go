package prod

import (
	"bufio"
	"fmt"
	"github.com/bitly/go-nsq"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tencent/util"
	"time"
)

const (
	MAX_FILE_OFFSET_NUM = 2000
)

type ProdServer struct {
	waitGroup util.WaitGroupWrapper
	opts      *ProdOption
	quit      chan bool

	lastRun   time.Time
	jsonconf  *JsonConf
	jsonMtime time.Time
	w         *nsq.Producer
}

func NewProdServer(opt *ProdOption) *ProdServer {
	p := new(ProdServer)
	p.opts = opt
	p.quit = make(chan bool)
	p.lastRun = p.lastRunTime()

	config := nsq.NewConfig()
	var err error
	p.w, err = nsq.NewProducer(opt.NsqdAddr, config)
	if err != nil {
		log.Panicf("ERR: create nsq.NewProducer, %s", err)
	}

	p.reloadJson()

	return p
}

func (s *ProdServer) Main() {
	cb := func() {
		ticker := time.NewTicker(time.Duration(s.opts.Duration) * time.Second)
		for {
			select {
			case <-ticker.C:
				//log.Println("DEBUG: Tick at", t)
				s.loop()

			case <-s.quit:
				goto END
			}
		}
	END:
		log.Println("INF: server quit loop")
	}

	log.Println("INF: server run")
	s.waitGroup.Wrap(cb)
}

func (s *ProdServer) Exit() {
	s.quit <- true
	s.waitGroup.Wait()
	s.w.Stop()
	log.Println("INF: server exit")
}

func (s *ProdServer) loop() {
	s.reloadJson()
	now, files := s.logFileList()
	if len(files) == 0 {
		return
	}
	s.filterLogFiles(files, &now)
	s.saveRunTime(now)
}

// 检查配置是否更新
func (s *ProdServer) reloadJson() {
	fi, err := os.Stat(s.opts.JsFile)
	if err != nil {
		log.Printf("WAR: os.Stat %s, %s", s.opts.JsFile, err)
		return
	}
	if fi.ModTime() != s.jsonMtime {
		s.jsonconf = NewJsonConf(s.opts.JsFile)
		s.jsonMtime = fi.ModTime()
	}
}

// 最后1次执行的时间戳
func (s *ProdServer) lastRunTime() time.Time {
	b, err := ioutil.ReadFile(s.opts.Record)
	if err != nil {
		log.Panicf("ERR: fail to open file %s", s.opts.Record)
	}
	line := strings.TrimSuffix(string(b), "\n")
	i, err := strconv.Atoi(line)
	if err != nil {
		log.Panicf("ERR: fail to parse int, %s, %v", b, err)
	}
	if i == 0 {
		return time.Now()
	}
	return time.Unix(int64(i), 0)
}

// 更新执行时间，保存到文件
func (s *ProdServer) saveRunTime(t time.Time) {
	log.Printf("INF: lastRun=%d, now=%d", s.lastRun.Unix(), t.Unix())
	s.lastRun = t
	str := fmt.Sprintf("%d\n", s.lastRun.Unix())
	if err := ioutil.WriteFile(s.opts.Record, []byte(str), os.ModePerm); err != nil {
		log.Panicf("ERR: fail to write file, %s, %d", s.opts.Record, s.lastRun.Unix())
	}
}

// 上次执行时间, 到当前时间的产生的日志文件清单
func (s *ProdServer) logFileList() (time.Time, []string) {
	now := time.Now()
	hours := make([]string, 0)
	t := s.lastRun
	layout := "2006010215"
	for t.Before(now) {
		str := t.Format(layout)
		hours = append(hours, str)
		t = t.Add(time.Hour)
	}
	currHour := now.Format(layout)
	if hours[len(hours)-1] != currHour {
		hours = append(hours, currHour)
	}
	log.Printf("INF: hours, %d, %s, %s", len(hours), hours[0], hours[len(hours)-1])

	filenames := make([]string, 0, 100)
	var pattern string
	for _, v := range hours {
		pattern = s.opts.LogDir + "/*/info/*" + v + ".log"
		m, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("ERR: fail to glob, pattern=%s", pattern)
			continue
		}
		filenames = append(filenames, m...)
	}
	log.Printf("INF: match %d log file", len(filenames))

	return now, filenames
}

// 从文件中过滤出统计日志
func (s *ProdServer) filterLogFiles(files []string, now *time.Time) {
	allLogs := 0
	for _, file := range files {
		s.readOne(file, now, &allLogs)
	}
	log.Printf("INF: log lines %d", allLogs)
	return
}

// 读取1个日志文件
func (s *ProdServer) readOne(file string, now *time.Time, count *int) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("WAR: fail to open file: %s", file)
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	rows := 0
	validLog := 0
	// 各日志的数量
	stat := make(map[int]int)

	locate, _ := time.LoadLocation("Asia/Shanghai")

	for {
	LoopNextLine:

		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("WAR: fail to read file, %s, %s", file, err.Error())
			}
			break
		}

		line = strings.TrimSuffix(line, "\n")
		rows += 1
		parts := strings.Split(line, "|")
		// 至少有10个部分
		if len(parts) < 10 {
			continue
		}
		// 第1个是时间戳
		// 判断是否在时间范围内
		pt2 := strings.Split(parts[0], ".")
		if len(pt2) != 2 {
			log.Printf("WAR: invalid date, %s:%d, %s", file, rows, parts[0])
			continue
		}

		// 本地时间parse
		const layout = "2006-01-02 15:04:05"
		t, err := time.Parse(layout, pt2[0])
		if err != nil {
			log.Printf("WAR: parse date, %s:%d, %s", file, rows, parts[0])
			continue
		}
		// parse转换默认是UTC, 这个接口有点awful
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
			t.Second(), t.Nanosecond(), locate)

		if t.Before(s.lastRun) || t.After(*now) {
			//log.Printf("WAR: time check, %s:%d, %s, %s, %s", file, rows, t, s.lastRun, *now)
			continue
		}

		// 判断logid是否合法
		logid, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Printf("WAR: fail to get logid, %s:%d", file, rows)
			continue
		}
		if !s.isValidLogid(logid) {
			continue
		}

		// 辅助定位问题
		modName := parts[6]

		// 检查每行日志，是否满足表的字段的要求
		conf := s.LogConfByID(logid)
		for _, field := range conf.Field_list {
			if field.FieldNumInt > len(parts) {
				log.Printf("WAR: field invalid, %s:%d, field=%d, logid=%d, modName=%s", file, rows, field.FieldNumInt, logid, modName)
				goto LoopNextLine
			}
			idx := field.FieldNumInt - 1

			if field.TypeInt == FIELD_TYPE_INT {
				_, err := strconv.Atoi(parts[idx])
				if err != nil {
					log.Printf("WAR: int invalid, %s:%d, field=%d, logid=%d, modName=%s", file, rows, field.FieldNumInt, logid, modName)
					goto LoopNextLine
				}
			} else if field.TypeInt == FIELD_TYPE_UINT {
				_, err := strconv.ParseUint(parts[idx], 10, 32)
				if err != nil {
					log.Printf("WAR: uint invalid, %s:%d, field=%d, logid=%d, modName=%s", file, rows, field.FieldNumInt, logid, modName)
					goto LoopNextLine
				}
			} else if field.TypeInt == FIELD_TYPE_DATETIME {
				// 约定：datetime类型, 位置只能在第1个字段
				if idx != 0 {
					log.Printf("WAR: datetime invalid, %s:%d, field=%d, logid=%d, modName=%s", file, rows, field.FieldNumInt, logid, modName)
					goto LoopNextLine
				}
			}
		}

		// 计数+1
		validLog++
		tablePre := s.jsonconf.LogidSet[logid]
		err = s.w.Publish(tablePre, []byte(line))
		if err != nil {
			util.AgentWarn(fmt.Sprintf("WAR: nsqd publish, %s", err))
			log.Fatal("WAR: nsq publish, %s", err)
		}
		stat[logid]++
	}
	*count = *count + validLog
	log.Printf("INF: file=%s, logs=%d, %v", file, validLog, stat)

	return
}

// 检查日志id的合法性
func (s *ProdServer) isValidLogid(logid int) bool {
	if _, ok := s.jsonconf.LogidSet[logid]; ok {
		return true
	}
	return false
}

// 寻找指定logid的表配置
func (s *ProdServer) LogConfByID(id int) *LogConf {
	for i := range s.jsonconf.Log_list {
		if s.jsonconf.Log_list[i].LogidInt == id {
			return &s.jsonconf.Log_list[i]
		}
	}
	log.Panicf("ERR: never run here, logid: %d", id)
	return nil
}
