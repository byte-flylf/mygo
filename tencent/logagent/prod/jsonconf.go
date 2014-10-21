package prod

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

// github.com/dropbox/godropbox/database/sqlbuilder/column.go 存在的限制
var validIdentifierRegexp *regexp.Regexp = regexp.MustCompile("^[a-zA-Z_]\\w*$")

// This is a strict subset of the actual restrictions, because
// teisenberger is lazy
func validIdentifierName(name string) bool {
	return validIdentifierRegexp.MatchString(name)
}

type JsonConf struct {
	Default_db_port string
	Default_db_user string
	Default_db_pw   string
	Route_list      []RouteConf
	Log_list        []LogConf
	LogidSet        map[int]string // 缓存所有的日志id, logid --> table
}

type RouteConf struct {
	Uin_start   string
	Uin_end     string
	Db_host     string
	UinStartInt int // Uin_start
	UinEndInt   int // UinEnd
}

type LogConf struct {
	Logid      string
	Table      string
	Field_list []FieldConf
	LogidInt   int // Logid
}

const (
	FIELD_TYPE_VCHAR    = 0
	FIELD_TYPE_INT      = 1
	FIELD_TYPE_UINT     = 2
	FIELD_TYPE_DATETIME = 3
)

var FIELD_TYPE_NAME = map[string]int{
	"varchar":      FIELD_TYPE_VCHAR,
	"int":          FIELD_TYPE_INT,
	"int unsigned": FIELD_TYPE_UINT,
	"datetime":     FIELD_TYPE_DATETIME,
}

type FieldConf struct {
	Type      string
	Column    string
	Field_num string

	TypeInt     int // Type
	FieldNumInt int // Field_num
}

func NewJsonConf(filename string) *JsonConf {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("fail to open file %s", filename)
	}
	var cfg JsonConf
	json.Unmarshal(b, &cfg)
	cfg.LogidSet = make(map[int]string)

	// 这个配置有奇葩的设定：整数也用字符串表示， 例如"default_db_port":"3306"
	// 需要进行转换
	for i, v := range cfg.Route_list {
		cfg.Route_list[i].UinStartInt, _ = strconv.Atoi(v.Uin_start)
		cfg.Route_list[i].UinEndInt, _ = strconv.Atoi(v.Uin_end)
	}

	for i, v := range cfg.Log_list {
		p := &cfg.Log_list[i]
		p.LogidInt, _ = strconv.Atoi(v.Logid)
		for j, w := range p.Field_list {
			p.Field_list[j].FieldNumInt, _ = strconv.Atoi(w.Field_num)
			if _, ok := FIELD_TYPE_NAME[w.Type]; !ok {
				log.Panicf("ERR, table=%s, invalid field type: %d", p.Table, w.Type)
			}
			p.Field_list[j].TypeInt = FIELD_TYPE_NAME[w.Type]

			// 检查字段名的合法性
			if !validIdentifierName(p.Field_list[j].Column) {
				log.Panicf("ERR, table=%s, invalid column=%s", p.Table, p.Field_list[j].Column)
			}
		}
		cfg.LogidSet[p.LogidInt] = p.Table
	}

	return &cfg
}
