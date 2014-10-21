package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tencent/logagent/prod"
	"tencent/util"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := cli.NewApp()
	app.Name = "logAdmin"

	app.Commands = []cli.Command{
		{
			Name:      "creat",
			ShortName: "c",
			Usage:     "为了支持企鹅的日志统计，分库分号段的建表",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "milestone", Value: "", Usage: "创建从那一天开始，默认是明天"},
				cli.IntFlag{Name: "days", Value: 3, Usage: "创建从milestone开始，几天内的表"},
				cli.StringFlag{Name: "temp", Value: "/data/qqpet/logcons/conf/create_table_temp.sql",
					Usage: "旧的表，建表语句会从sql模板中生成"},
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
			},
			Action: creatTbl,
		}, {
			Name:      "drop",
			ShortName: "d",
			Usage:     "为了支持企鹅的日志统计：删除不再使用的表",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "milestone", Value: "20", Usage: "删除表，从第几天开始"},
				cli.IntFlag{Name: "days", Value: 3, Usage: "删除表，从milestone开始，倒数前N天"},
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
			},
			Action: dropTbl,
		},
		{
			Name:  "delete",
			Usage: "删除指定时间范围内，指定的表",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "start", Value: "", Usage: "删除表，从那天开始, 格式为20140829"},
				cli.StringFlag{Name: "end", Value: "", Usage: "删除表，那天结束，包括这一天，格式为20140829"},
				cli.StringFlag{Name: "table", Value: "", Usage: "需要删除的表名"},
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
			},
			Action: deleteTbl,
		},

		{
			Name:      "select",
			ShortName: "s",
			Usage:     "查找指定用户的日志",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
				cli.StringFlag{Name: "start", Value: "", Usage: "日志的开始时间"},
				cli.StringFlag{Name: "end", Value: "", Usage: "日志的结束时间"},
				cli.StringFlag{Name: "table", Value: "", Usage: "需要查询的表名"},
				cli.IntFlag{Name: "uin", Value: 732834509, Usage: "查询的用户的QQ号码"},
			},
			Action: selectTbl,
		},
		{
			Name:      "count",
			ShortName: "t",
			Usage:     "统计指定天的一张表的日志总量",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
				cli.StringFlag{Name: "table", Value: "", Usage: "需要查询的表名"},
				cli.StringFlag{Name: "date", Value: "", Usage: "查询的日期，格式为20140807"},
			},
			Action: countTbl,
		},
		{
			Name:      "check",
			ShortName: "k",
			Usage:     "检查明天的表是否存在",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
				cli.StringFlag{Name: "table", Value: "tb_petquan_change", Usage: "需要查询的表名"},
			},
			Action: checkTbl,
		},
		{
			Name:      "load",
			ShortName: "l",
			Usage:     "读取文件中所有的insert语句，依次执行, 配合logcons异常时补录数据",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "json", Value: "/data/qqpet/conf/json_conf/logagent.json", Usage: "配置文件logagent.json"},
				cli.StringFlag{Name: "filename", Value: "", Usage: "sql文件路径"},
			},
			Action: loadSqlFile,
		},
	}

	app.Run(os.Args)
}

// 删除时间范围的表的数据
func deleteTbl(ctx *cli.Context) {
	layout := "20060102"
	var err error

	delTable := ctx.String("table")
	if delTable == "" {
		log.Fatal("ERR: arg 'table' empty")
	}

	var start time.Time
	if start, err = time.Parse(layout, ctx.String("start")); err != nil {
		log.Fatal("ERR: arg 'start' 格式非法，正确的格式是, %s", layout)
	}

	var end time.Time
	if end, err = time.Parse(layout, ctx.String("end")); err != nil {
		log.Fatal("ERR: arg 'end' 格式非法，正确的格式是, %s", layout)
	}
	if start.After(end) {
		log.Fatal("ERR:  'start' after 'end'")
	}

	dates := dateSlice(start, end)
	log.Printf("dates: %v", dates)
	jsonconf := prod.NewJsonConf(ctx.String("json"))

	sqlTemp := make(map[string]string)
	sqlTemp[delTable] = "DELETE FROM petLog_%d." + delTable + "_%s"

	DbExecute(jsonconf, dates, sqlTemp)
}

func loadSqlFile(ctx *cli.Context) {
	filename := ctx.String("filename")
	key := "petLog_"
	if strings.Index(filename, key) == -1 {
		log.Printf("WAR: substr '%s' not in filename", key)
		return
	}

	re, _ := regexp.Compile("petLog_([0-9]+)")
	slice := re.FindStringSubmatch(filename)
	if len(slice) < 2 {
		log.Printf("WAR, not match, %s", filename)
		return
	}
	log.Printf("INF: filename=%s, dbIdx=%s", filename, slice[1])
	dbIdx, err := strconv.Atoi(slice[1])
	if err != nil {
		log.Printf("WAR: Atoi fail, %s, %s", slice[1], err)
		return
	}
	dbname := fmt.Sprintf("petLog_%d", dbIdx)
	lines, err := util.ReadLines(filename)
	if err != nil {
		log.Printf("WAR: fail to read file, %s,, %s", filename, err)
		return
	}

	jsonconf := prod.NewJsonConf(ctx.String("json"))
	user := jsonconf.Default_db_user
	passwd := jsonconf.Default_db_pw
	port := jsonconf.Default_db_port

	var host string
	for _, conf := range jsonconf.Route_list {
		if conf.UinStartInt <= dbIdx && dbIdx <= conf.UinEndInt {
			host = conf.Db_host
			break
		}
	}
	if host == "" {
		fmt.Printf("WAR: uin=%d, host not found", dbIdx)
		return
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
	conn, err := sql.Open("mysql", dns)
	if err != nil {
		log.Panicf("ERR: sql.Open, %s, %s", dns, err)
	}

	for _, line := range lines {
		sqlstr := line
		stmt, err := conn.Prepare(sqlstr)
		if err != nil {
			log.Panicf("ERR, conn.Prepare, %s, sql=%s", err, sqlstr)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Panicf("ERR, stmt.Exec, %s, sql=%s", err, sqlstr)
		}
		stmt.Close()
	}
	log.Printf("INF: execute, host=%s, filename=%s, lines=%d", host, filename, len(lines))
	conn.Close()
}

// 统计一天的一个指定日志的总量
func countTbl(ctx *cli.Context) {
	table := ctx.String("table")

	date := ctx.String("date")
	if _, err := time.Parse("20060102", date); err != nil {
		log.Panic("ERR: arg 'date' invalid")
	}

	jsonconf := prod.NewJsonConf(ctx.String("json"))
	user := jsonconf.Default_db_user
	passwd := jsonconf.Default_db_pw
	port := jsonconf.Default_db_port

	sum := 0
	for _, conf := range jsonconf.Route_list {
		host := conf.Db_host
		for u := conf.UinStartInt; u <= conf.UinEndInt; u++ {
			dbname := fmt.Sprintf("petLog_%d", u)
			dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
			conn, err := sql.Open("mysql", dns)
			if err != nil {
				log.Panicf("ERR: sql.Open, %s, %s", dns, err)
			}
			sqlstr := fmt.Sprintf("SELECT count(*) FROM %s_%s", table, date)
			rows, err := conn.Query(sqlstr)
			if err != nil {
				log.Panic("ERR: db.Query, %s", err)
			}
			for rows.Next() {
				var cnt int
				_ = rows.Scan(&cnt)
				sum += cnt
			}

			log.Printf("INF: host=%s, db=%s, sql=%s", host, dbname, sqlstr)
			conn.Close()
		}
	}
	log.Printf("INF: table=%s, date=%s, sum=%d", table, date, sum)
}

// 检查表是否成功创建，如果失败，产生Agent告警
func checkTbl(ctx *cli.Context) {
	tomorrow := time.Now().AddDate(0, 0, 1)
	datestr := tomorrow.Format("20060102")

	jsonconf := prod.NewJsonConf(ctx.String("json"))

	table := ctx.String("table")
	// 最后1个号段的db
	dbname := "petLog_99"
	user := jsonconf.Default_db_user
	passwd := jsonconf.Default_db_pw
	port := jsonconf.Default_db_port
	// 配置文件中最后一台机器
	host := jsonconf.Route_list[len(jsonconf.Route_list)-1].Db_host

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
	conn, err := sql.Open("mysql", dns)
	if err != nil {
		log.Panicf("ERR: sql.Open, %s, %s", dns, err)
	}

	expectTable := fmt.Sprintf("%s_%s", table, datestr)
	sqlstr := fmt.Sprintf("select TABLE_NAME from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='%s' and TABLE_NAME='%s'",
		dbname, expectTable)
	rows, err := conn.Query(sqlstr)
	if err != nil {
		util.AgentWarn(fmt.Sprintf("logAdmin, check, conn.Query, %s, %s, %s", dbname, table, datestr))
		log.Panic("ERR: db.Query, %s", err)
	}

	log.Printf("INF: host=%s, sql=%s", host, sqlstr)
	result := false
	for rows.Next() {
		var tablename string
		err = rows.Scan(&tablename)
		if err != nil {
			log.Panic("ERR: rows.Scan, %s", err)
			break
		}
		if tablename == expectTable {
			result = true
			break
		}
	}
	if !result {
		util.AgentWarn("logAdmin, check, fail to create tomorrow table")
	}
	conn.Close()
}

// 统一的搜索语句
func selectTbl(ctx *cli.Context) {
	layout := "20060102"
	var start time.Time
	var err error

	if start, err = time.Parse(layout, ctx.String("start")); err != nil {
		log.Panicf("ERR: arg 'start' 格式非法，正确的格式是, %s", layout)
	}

	var end time.Time
	if end, err = time.Parse(layout, ctx.String("end")); err != nil {
		log.Panicf("ERR: arg 'end' 格式非法，正确的格式是, %s", layout)
	}
	if start.After(end) {
		log.Panic("ERR: 开始时间start，晚于结束时间end")
	}
	dates := dateSlice(start, end)

	var uin int
	if uin, err = strconv.Atoi(ctx.String("uin")); err != nil {
		log.Panic("ERR: arg 'uin' 非法")
	}
	mod := uin % 100

	jsonconf := prod.NewJsonConf(ctx.String("json"))
	var conf *prod.RouteConf
	for i, c := range jsonconf.Route_list {
		if c.UinStartInt <= mod && mod <= c.UinEndInt {
			conf = &jsonconf.Route_list[i]
			break
		}
	}
	if conf == nil {
		log.Panicf("ERR: 没有找到数据库配置，uin=%d", uin)
	}

	table := ctx.String("table")
	if table == "" {
		log.Panic("ERR: 需要参数'table'")
	}

	dbname := fmt.Sprintf("petLog_%d", uin%100)
	user := jsonconf.Default_db_user
	passwd := jsonconf.Default_db_pw
	port := jsonconf.Default_db_port
	host := conf.Db_host

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
	conn, err := sql.Open("mysql", dns)
	if err != nil {
		log.Panicf("ERR: sql.Open, %s, %s", dns, err)
	}
	for _, d := range dates {
		sqlstr := fmt.Sprintf("SELECT * FROM %s_%s WHERE uin=%d AND date >= '%s' AND  date < '%s'",
			table, d, uin, start, end)
		rows, err := conn.Query(sqlstr)
		if err != nil {
			log.Panic("ERR: db.Query, %s", err)
		}
		log.Printf("INF: execute, sql=%s", sqlstr)
		// 表的列数不固定
		cols, _ := rows.Columns()
		buff := make([]interface{}, len(cols)) // 临时slice，用来通过类型检查
		data := make([]string, len(cols))
		for i, _ := range buff {
			buff[i] = &data[i] // 把两个slice关联起来
		}
		uinCol := 0
		for i, col := range cols {
			if col == "uin" {
				uinCol = i
			}
		}
		for rows.Next() {
			rows.Scan(buff...)
			if cols[uinCol] == "" {
				// 跳过空结果集
				continue
			}
			fmt.Println()
			for k, col := range data {
				if cols[k] != "" {
					fmt.Printf("%10s:\t%10s\n", cols[k], col)
				}
			}
			fmt.Println()
		}

	}
	conn.Close()

}

// 根据模板文件，生成sql, 删除无用的表
func dropTbl(ctx *cli.Context) {
	var start time.Time
	var end time.Time

	layout := "20060102"
	//  第1个格式：是指定日期
	//  第2个格式是整数，起始点就是今天之前的第几天
	var err error
	if end, err = time.Parse(layout, ctx.String("milestone")); err != nil {
		if i, err := strconv.Atoi(ctx.String("milestone")); err != nil {
			log.Panicf("ERR: arg 'milestone' invalid, %s", ctx.String("milestone"))
		} else if i < 0 {
			log.Panicf("ERR: arg 'milestone' must >= 0")
		} else {
			end = time.Now().AddDate(0, 0, (-1)*i)
		}
	}

	days := ctx.Int("days")
	if days < 1 {
		log.Panicf("ERR: arg 'days' invalid, %d", days)
	}
	start = end.AddDate(0, 0, (-1)*days)
	log.Println(start.Format(layout), end.Format(layout))
	dates := dateSlice(start, end)
	log.Printf("INF: dates, %s", dates)

	jsonconf := prod.NewJsonConf(ctx.String("json"))
	sqlTemp := make(map[string]string)
	for _, conf := range jsonconf.Log_list {
		table := conf.Table
		// 新加的表的建表语句，直接由配置文件生成
		sqlTemp[table] = dropSql(table)
		log.Printf("INF: table=%s, sql=%s", table, sqlTemp[table])
	}

	DbExecute(jsonconf, dates, sqlTemp)
}

// 根据模板文件，json配置生成建表语言，并执行
func creatTbl(ctx *cli.Context) {
	sqlTemp := readSqlTemplate(ctx.String("temp"))
	jsonconf := prod.NewJsonConf(ctx.String("json"))
	for _, conf := range jsonconf.Log_list {
		table := conf.Table
		// 新加的表的建表语句，直接由配置文件生成
		if _, ok := sqlTemp[table]; !ok {
			sqlTemp[table] = createSql(&conf)
			log.Printf("INF: table=%s", table)
		}
	}

	var start time.Time
	var err error
	if ctx.String("milestone") == "" {
		start = time.Now().AddDate(0, 0, 1)
	} else {
		start, err = time.Parse("20060102", ctx.String("milestone"))
		if err != nil {
			log.Panicf("ERR: %s", err)
		}
	}
	days := ctx.Int("days")
	if days < 1 {
		log.Panicf("ERR: arg 'days' invalid, %d", days)
	}
	end := start.AddDate(0, 0, days)

	dates := dateSlice(start, end)
	log.Printf("INF: dates, %s", dates)
	DbExecute(jsonconf, dates, sqlTemp)
}

// 针对所有的DB, 所有的号段，执行建表和删表操作
func DbExecute(jsonconf *prod.JsonConf, dates []string, sqlTemp map[string]string) {
	for _, conf := range jsonconf.Route_list {
		//for u := conf.UinStartInt; u <= conf.UinEndInt; u++ {
		// 经分同事希望每台机器上都有100个DB
		log.Printf("host=%s, user=%s, passwd=%s, port=%s", conf.Db_host,
			jsonconf.Default_db_user, jsonconf.Default_db_pw, jsonconf.Default_db_port)
		for u := 0; u <= 99; u++ {
			execute(jsonconf.Default_db_user, jsonconf.Default_db_pw,
				jsonconf.Default_db_port, u, conf.Db_host, dates, sqlTemp)
		}
	}
}

// 针对1个号段，对指定日期内的所有表，执行create或者drop操作
func execute(user string, passwd string, port string, uin int, host string,
	dates []string, sqlTemp map[string]string) {

	dbname := fmt.Sprintf("petLog_%d", uin)
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=gbk", user, passwd, host, port, dbname)
	conn, err := sql.Open("mysql", dns)
	if err != nil {
		util.AgentWarn("logAdmin, sql.Open")
		log.Panicf("ERR: sql.Open, %s, %s", dns, err)
	}
	for _, date := range dates {

		for _, temp := range sqlTemp {
			sqlstr := fmt.Sprintf(temp, uin, date)
			stmt, err := conn.Prepare(sqlstr)
			if err != nil {
				util.AgentWarn("logAdmin, conn.Prepare")
				log.Panicf("ERR, conn.Prepare, %s, sql=%s", err, sqlstr)
			}
			_, err = stmt.Exec()
			if err != nil {
				util.AgentWarn("logAdmin, stmt.Exec")
				log.Panicf("ERR, stmt.Exec, %s, sql=%s", err, sqlstr)
			}
			log.Printf("INF, host=%s, sql=%s", host, sqlstr)
			stmt.Close()
		}
	}
	log.Printf("INF: execute, host=%s, uin=%d", host, uin)
	conn.Close()
}

// 从明天起的N天的日期
func dateSlice(start time.Time, end time.Time) []string {
	dates := make([]string, 0)

	layout := "20060102"
	for start.Before(end) {
		d := start.Format(layout)
		dates = append(dates, d)
		start = start.Add(time.Hour * 24)
	}
	return dates
}

func readSqlTemplate(file string) map[string]string {
	lines, err := util.ReadLines(file)
	if err != nil {
		log.Panicf("ERR: fail to read sql template file, %s, %s", file, err)
	}
	sqlTemp := make(map[string]string)
	const k1 = "%d."
	const k2 = "_%s"
	for _, line := range lines {
		start := strings.Index(line, k1)
		end := strings.Index(line, k2)
		table := line[start+len(k1) : end]
		sqlTemp[table] = line
		log.Printf("INF: table=%s", table)
	}
	return sqlTemp
}

func createSql(conf *prod.LogConf) string {
	sqlstr := "create table if not exists petLog_%d." + conf.Table + "_%s ("
	for i, field := range conf.Field_list {
		if i != 0 {
			sqlstr += ", "
		}
		if field.TypeInt == prod.FIELD_TYPE_VCHAR {
			sqlstr += field.Column + "  varchar(255)"
		} else {
			sqlstr += field.Column + " " + field.Type
		}
	}
	sqlstr += ") ENGINE=InnoDB DEFAULT CHARSET=gbk;"
	return sqlstr
}

func dropSql(table string) string {
	sqlstr := "drop table if exists petLog_%d." + table + "_%s"
	return sqlstr
}
