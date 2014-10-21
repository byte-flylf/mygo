package consumer

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"tencent/logagent/prod"
	"time"

	"tencent/util"
)

type SqlItem struct {
	sql   string
	dbIdx int
	table string
}

type DBPool struct {
	conns     []*sql.DB
	InputChan chan SqlItem

	FailChan chan bool // 通知主线程，pool无法工作
	dump     bool

	sqldir string
}

func NewDBPool(dir string) *DBPool {
	p := new(DBPool)
	p.conns = make([]*sql.DB, MAX_DB_NUM)
	p.FailChan = make(chan bool)
	p.InputChan = make(chan SqlItem)
	p.dump = false
	p.sqldir = dir

	return p
}

func (p *DBPool) Start(jsonconf *prod.JsonConf) {
	user := jsonconf.Default_db_user
	passwd := jsonconf.Default_db_pw
	port := jsonconf.Default_db_port

	for _, v := range jsonconf.Route_list {
		host := v.Db_host
		for i := v.UinStartInt; i <= v.UinEndInt; i++ {
			dbname := fmt.Sprintf("petLog_%d", i)
			dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
			c, err := sql.Open("mysql", dns)
			if err != nil {
				log.Fatal("ERR, sql.Open, %s, %s", dns, err)
			}

			if p.ping(dbname, c) {
				log.Printf("INF: db '%s' ping ok", dbname)
				p.conns[i] = c
			} else {
				log.Fatal("ERR, PingDB, dns=%s", dns)
			}
		}
	}

	go func() {
		for {
			it, ok := <-p.InputChan
			if !ok {
				log.Println("INF, dbpool loop break")
				// 正常停止
				break
			}

			if !p.dump {
				res, err := p.conns[it.dbIdx].Exec(it.sql)
				if err != nil {
					errstr := fmt.Sprintf("ERR: logcons, conn.Exec, %s, db=%d, sql=%s", err, it.dbIdx, it.sql)
					log.Printf(errstr)
					util.AgentWarn(errstr)
					p.dump = true
					p.FailChan <- true
					p.dump2file(it.dbIdx, it.sql)
					continue
				}

				rows, _ := res.RowsAffected()
				log.Printf("INF: db=%d, table=%s, RowsAffected=%d", it.dbIdx, it.table, rows)
			} else {
				p.dump2file(it.dbIdx, it.sql)
			}
		}
	}()
}

func (p *DBPool) Stop() {
	log.Println("INF: pool stop")
	if p.InputChan != nil {
		close(p.InputChan)
		p.InputChan = nil
	}
	for i := 0; i < MAX_DB_NUM; i++ {
		p.conns[i].Close()
	}
}

// 传输数据前就建立db连接
func (p *DBPool) ping(dbname string, conn *sql.DB) bool {
	today := time.Now().Format("20060102")
	//today := "20140819"
	lookupTable := fmt.Sprintf("t_lazywriter_online_%s", today)
	sqlstr := fmt.Sprintf("select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='%s' and TABLE_NAME='%s'",
		dbname, lookupTable)

	rows, err := conn.Query(sqlstr)
	if err != nil {
		log.Printf("WAR: PingDB, conn.Query, %s, %s", sqlstr, err)
		return false
	}
	defer rows.Close()

	res := false
	for rows.Next() {
		var str string
		err := rows.Scan(&str)
		if err == nil {
			if str == lookupTable {
				res = true
				break
			}
		}
	}

	return res
}

func (p *DBPool) dump2file(dbIdx int, sql string) {
	datestr := time.Now().Format("20060102")
	path := fmt.Sprintf("%s/petLog_%d_%s.sql", p.sqldir, dbIdx, datestr)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		log.Printf("WAR: os.OpenFile, %s, %s", path, err)
		return
	}
	defer file.Close()
	b := []byte(sql + "\n")
	file.Write(b)
}
