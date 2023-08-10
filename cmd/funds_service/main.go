package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/DwGoing/OnlyPay/pkg/funds_service"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
)

func main() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path = filepath.Dir(path)
	if err := ioc.Load(config.WithSearchPath(path)); err != nil {
		panic(err)
	}
	app, err := funds_service.GetFundsServiceSingleton()
	if err != nil {
		panic(err)
	}
	err = app.Initialize()
	if err != nil {
		panic(err)
	}

	log.Println("进程已启动")
	//监听指定信号 ctrl+c kill
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Println("进程已结束")
}
