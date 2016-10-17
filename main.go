package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"beego-demo/controllers"
	_ "beego-demo/routers"

	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
)

func handleSignals(c chan os.Signal) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM:
		fmt.Println("Shutdown quickly, bye...")
	case syscall.SIGQUIT:
		fmt.Println("Shutdown gracefully, bye...")
		// do graceful shutdown
	}

	os.Exit(0)
}

func main() {
	graceful, _ := beego.AppConfig.Bool("graceful")
	if !graceful {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go handleSignals(sigs)
	}

	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	mode := beego.AppConfig.String("runmode")
	if mode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}

	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}
