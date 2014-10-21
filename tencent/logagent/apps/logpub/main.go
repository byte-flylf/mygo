// 读取本地的日志文件，过滤传入nsq
package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/mreiferson/go-options"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tencent/logagent/prod"
)

var (
	flagSet = flag.NewFlagSet("logagent", flag.ExitOnError)

	config = flagSet.String("config", "", "path to config file")

	jsfile   = flagSet.String("jsfile", "", "path to json file")
	duration = flagSet.Int("duration", 60, "how many seconds to read log file")
	record   = flagSet.String("record", "", "path to svr path")
	logdir   = flagSet.String("logdir", "", "path to log root dir")
	nsqdaddr = flagSet.String("nsqdaddr", "", "addr to connect nsqd")
)

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
			log.Fatalf("ERROR: failed to load config file %s - %s", *config, err.Error())
		}
	}
	opts := prod.NewProdOption()
	options.Resolve(opts, flagSet, cfg)
	log.Println("DEBUG: opts", opts)

	daemon := prod.NewProdServer(opts)
	daemon.Main()
	<-exitChan
	daemon.Exit()
}
