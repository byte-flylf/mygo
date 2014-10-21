package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/mreiferson/go-options"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tencent/logagent/consumer"
	"tencent/util"
)

var (
	flagSet = flag.NewFlagSet("logcons", flag.ExitOnError)

	config = flagSet.String("config", "", "path to config file")
	pprof  = flagSet.String("pprof", "0.0.0.0:6060", "addr to net/http/pprof")

	jsfile = flagSet.String("jsfile", "", "path to json file")

	statfile = flagSet.String("statfile", "", "the file which have stat data")
	sqlDir   = flagSet.String("sqldir", "", "the directory which save sql file when mysql couldn't connect")

	duration    = flagSet.String("duration", "", "在多少时间间隔内输出数据到mysql")
	maxinflight = flagSet.Int("max-in-flight", 1, "nsqd的输入速率")

	addr = flagSet.String("httpaddr", "", "http pprof")

	nsqlookupd = util.StringArray{}
)

func init() {
	flagSet.Var(&nsqlookupd, "lookupd-http-address", "lookupd TCP address (may be given multiple times)")
}

func main() {
	flagSet.Parse(os.Args[1:])

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan int)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var cfg map[string]interface{}
	if *config != "" {
		_, err := toml.DecodeFile(*config, &cfg)
		if err != nil {
			log.Fatalf("ERR: failed to load config file %s - %s", *config, err.Error())
		}
	}
	opts := consumer.NewConsOption()
	options.Resolve(opts, flagSet, cfg)
	log.Println("DEB: opts", opts)

	daemon := consumer.NewConsServer(opts)
	daemon.Main()
	<-exitChan
	daemon.Exit()
}
