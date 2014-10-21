package consumer

type ConsOption struct {
	JsFile             string   `flag:"jsfile"`
	NsqlookupdHttpAddr []string `flag:"lookupd-http-address" cfg:"lookupd_http_address"` // nsqlookupd的http连接地址

	StatFile string `flag:"statfile"` // 统计文件的目录, obsoleted

	SqlDir string `flag:"sqldir"` // 系统失败，保存日志到sql语句

	Duration    int `flag:"duration"`      // 系统输出频率，多少秒dump日志到数据库
	MaxInFlight int `flag:"max-in-flight"` // 系统输入频率， nsqd输入的速率

	HttpAddr string `flag:"httpaddr"`
}

func NewConsOption() *ConsOption {
	o := &ConsOption{
		JsFile:      "/data/qqpet/conf/json_conf/logagent.json",
		Duration:    15,
		StatFile:    "/data/qqpet/logcons/data",
		SqlDir:      "/data/qqpet/logcons/sql/",
		MaxInFlight: 1,
		HttpAddr:    "0.0.0.0:5001",
	}

	return o
}
