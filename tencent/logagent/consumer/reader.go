package consumer

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-nsq"
	"github.com/dropbox/godropbox/database/sqlbuilder"
	"log"
	"strconv"
	"sync"
	"tencent/logagent/prod"
	"tencent/util"
	"time"
)

const (
	LOG_BUF_LEN  = 10000 // 每次缓冲的最大日志量
	MAX_DB_NUM   = 100   // 100个分库
	LOG_BASE_LEN = 512   // 日志的长度不应该超过256个字节
)

var (
	VERTICAL_VAR = []byte("|")
	DOT          = []byte(".")
)

// 插入集合， 在一个特定的表上面的插入语句
type TableSqlSet map[string]sqlbuilder.InsertStatement

// 代表了1张表的日志的输入
type logReader struct {
	table string // 表名
	q     *nsq.Consumer
	mu    sync.Mutex

	buf [][]byte // 缓存
	idx int      // 当前的下标

	timeout int       // 多少秒后超时，强行清空缓冲
	t       time.Time // 上一次清空日志时间

	sqlset []TableSqlSet

	// Recycler
	get  chan []byte
	give chan []byte

	jsonconf *prod.JsonConf

	sqlchan chan<- SqlItem
}

// 创建nsq输入日志流
func NewlogReader(table string, timeout int) *logReader {
	r := new(logReader)
	r.table = table

	r.sqlset = make([]TableSqlSet, MAX_DB_NUM)

	r.buf = make([][]byte, LOG_BUF_LEN)
	r.idx = 0
	r.timeout = timeout
	r.t = time.Now()

	r.get, r.give = makeRecycler(LOG_BASE_LEN)

	return r
}

func (r *logReader) Start(jsonconf *prod.JsonConf, inflight int, addrs []string,
	sqlChan chan SqlItem) {
	if r.q != nil {
		log.Printf("WAR: reader[%s], already run", r.table)
		return
	}
	r.jsonconf = jsonconf
	r.sqlchan = sqlChan

	var err error
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = inflight
	r.q, err = nsq.NewConsumer(r.table, "logcons", cfg)
	if err != nil {
		log.Fatal("ERR: reader[%s], nsq.NewConsumer, %s", r.table, err)
	}

	r.q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		r.mu.Lock()
		defer r.mu.Unlock()

		r.addLog(message.Body)

		// 缓冲已满
		if r.idx == LOG_BUF_LEN {
			log.Printf("INF, reader[%s], buffull", r.table)
			r.flush()
			return nil
		}
		// 超时
		if time.Since(r.t) > time.Duration(time.Second*time.Duration(r.timeout)) {
			log.Printf("INF, reader[%s], timeout, %d second", r.table, r.timeout)
			r.flush()
		}

		return nil
	}))

	err = r.q.ConnectToNSQLookupds(addrs)
	if err != nil {
		log.Fatal("ERR: reader[%s], nsq.ConnectToNSQLookupd, %s, %s", r.table, addrs, err)
	}
}

func (r *logReader) IsRun() bool {
	return r.q != nil
}

func (r *logReader) Stop() {
	r.q.Stop()
}

func (r *logReader) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	//r.q.Stop()
	r.flush()

	r.q = nil
	r.jsonconf = nil
	log.Printf("INF: reader[%s] stopped", r.table)
}

// 清空缓存
func (r *logReader) flush() {
	if r.idx == 0 {
		return
	}

	for i := 0; i < MAX_DB_NUM; i++ {
		r.sqlset[i] = nil // 释放之前的
		r.sqlset[i] = make(TableSqlSet, LOG_BUF_LEN/MAX_DB_NUM)
	}

	for i := 0; i < r.idx; i++ {
		logid, tableName, dbIdx, vals, err := r.splitLog(r.buf[i])
		if err != nil {
			log.Printf("WAR: reader[%s], splitLog, %s", r.table, err)
			continue
		}

		if _, ok := r.sqlset[dbIdx][tableName]; !ok {
			r.sqlset[dbIdx][tableName] = r.newInsertStat(logid, tableName)
		}
		stmt := r.sqlset[dbIdx][tableName]
		stmt.Add(vals...)
	}

	for i := 0; i < MAX_DB_NUM; i++ {
		db := fmt.Sprintf("petLog_%d", i)
		tblset := r.sqlset[i]
		for table, stmt := range tblset {
			sqlstr, err := stmt.String(db)
			if err != nil {
				log.Printf("WAR, stmt.String, db=%s, %s", db, err)
				continue
			}

			r.sqlchan <- SqlItem{sqlstr, i, table}
		}
	}

	log.Printf("INF, reader[%s], flush, %d", r.table, r.idx)

	// 回收内存
	for i := 0; i < r.idx; i++ {
		line := r.buf[i]
		b := line[:cap(line)]
		r.give <- b
	}
	r.idx = 0
	r.t = time.Now()

}

// 分割日志成sqlbuilder的行数据
func (r *logReader) splitLog(line []byte) (logid uint, tableName string,
	dbIdx int, vals []sqlbuilder.Expression, err error) {

	parts := bytes.Split(line, VERTICAL_VAR)
	if len(parts) < 10 {
		err = fmt.Errorf("split, parts=%d, %d|%s", len(parts), len(line), line)
		return
	}
	pt2 := bytes.Split(parts[0], DOT)
	if len(pt2) != 2 {
		err = fmt.Errorf("split, invalid date|%s", line)
		return
	}

	// 解析本地时间
	var t time.Time
	t, err = time.Parse("2006-01-02 15:04:05", string(pt2[0]))
	if err != nil {
		err = fmt.Errorf("split, prase, %s|%s", err, line)
		return
	}
	locate, _ := time.LoadLocation("Asia/Shanghai")
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
		t.Second(), t.Nanosecond(), locate)

	logid, err = util.ByteToBase10(parts[2])
	if err != nil {
		err = fmt.Errorf("split, logid, %s|%s", err, line)
		return
	}
	tableName = r.log2Table(logid, &t)

	uin, err := util.ByteToBase10(parts[3])
	if err != nil {
		err = fmt.Errorf("split, uin, %s|%s", err, line)
		return
	}
	dbIdx = int(uin % 100)

	cfg := r.logConfByID(logid)
	vals = make([]sqlbuilder.Expression, len(cfg.Field_list))
	for i, field := range cfg.Field_list {
		if field.FieldNumInt > len(parts) {
			err = fmt.Errorf("field invalid, %d|%s", field.FieldNumInt, line)
			return
		}
		idx := field.FieldNumInt - 1

		if field.TypeInt == prod.FIELD_TYPE_INT {
			it, e2 := strconv.Atoi(string(parts[idx]))
			if e2 != nil {
				err = fmt.Errorf("int invalid, filed=%d, %s|%s", field.FieldNumInt, e2, line)
				return
			}
			vals[i] = sqlbuilder.Literal(it)
		} else if field.TypeInt == prod.FIELD_TYPE_UINT {
			uit, e2 := util.ByteToBase10(parts[idx])
			if e2 != nil {
				err = fmt.Errorf("unsigned invalid, filed=%d|%s", field.FieldNumInt, line)
				return
			}
			vals[i] = sqlbuilder.Literal(uit)
		} else if field.TypeInt == prod.FIELD_TYPE_DATETIME {
			// datetime类型, 位置只能在第1个字段
			if idx != 0 {
				err = fmt.Errorf("datetime invalid, %d|%s", field.FieldNumInt, line)
				return
			}
			vals[i] = sqlbuilder.Literal(pt2[0])
		} else {
			vals[i] = sqlbuilder.Literal(parts[idx])
		}
	}
	return
}

// 根据logid得到表名
func (r *logReader) log2Table(logid uint, t *time.Time) string {
	idx := int(logid)
	tablePre := r.jsonconf.LogidSet[idx]
	str := t.Format("20060102")
	return fmt.Sprintf("%s_%s", tablePre, str)
}

// 寻找指定logid的表配置
func (r *logReader) logConfByID(logid uint) *prod.LogConf {
	key := int(logid)
	for i := range r.jsonconf.Log_list {
		if r.jsonconf.Log_list[i].LogidInt == key {
			return &r.jsonconf.Log_list[i]
		}
	}
	log.Panicf("ERROR: never run here, logid: %d", logid)
	return nil
}

func (r *logReader) newInsertStat(logid uint, tblName string) sqlbuilder.InsertStatement {
	cfg := r.logConfByID(logid)

	columns := make([]sqlbuilder.NonAliasColumn, len(cfg.Field_list))
	for i := 0; i < len(cfg.Field_list); i++ {
		field := &cfg.Field_list[i]
		if field.TypeInt == prod.FIELD_TYPE_INT {
			columns[i] = sqlbuilder.IntColumn(field.Column, true)
		} else if field.TypeInt == prod.FIELD_TYPE_UINT {
			columns[i] = sqlbuilder.IntColumn(field.Column, true)
		} else {
			columns[i] = sqlbuilder.BytesColumn(field.Column, true)
		}
	}
	table := sqlbuilder.NewTable(tblName, columns...)
	return table.Insert(columns...)
}

func (r *logReader) addLog(msg []byte) {
	b := <-r.get
	n := copy(b, msg)
	if n != len(msg) {
		log.Printf("WAR: buf not enough, need=%d, has=%d, n=%d", len(msg), cap(b), n)
	}
	r.buf[r.idx] = b[:n]
	r.idx++
}
