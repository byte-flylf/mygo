package prod

type ProdOption struct {
	JsFile   string `flag:"jsfile"`
	Duration int    `flag:"duration"` // 间隔多久读1次日志文件
	Record   string `flag:"record"`   // 文件记录了上次执行的时间戳
	LogDir   string `flag:"logdir"`   // 日志文件根目录
	NsqdAddr string `flag:"nsqdaddr"` // nsqd节点的地址
}

func NewProdOption() *ProdOption {
	o := &ProdOption{
		JsFile:   "/data/qqpet/conf/json_conf/logagent.json",
		Duration: 30,
		Record:   "/data/qqpet/logpub/conf/record.txt",
		LogDir:   "/data/qqpet/log/",
	}

	return o
}
