package main

import (
	_ "beego-demo/routers"
	"fmt"
	"github.com/gwduan/beego"
	_ "github.com/gwduan/beego/session/redis"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
)

func setUserId() {
	userName := beego.AppConfig.String("user")
	u, err := user.Lookup(userName)
	if err != nil {
		fmt.Println("user config:", err)
		return
	}

	gid, _ := strconv.ParseInt(u.Gid, 0, 0)
	uid, _ := strconv.ParseInt(u.Uid, 0, 0)
	if err := syscall.Setregid(int(gid), int(gid)); err != nil {
		fmt.Println("setregid:", err)
	}
	if err := syscall.Setreuid(int(uid), int(uid)); err != nil {
		fmt.Println("setreuid:", err)
	}
}

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
	//setUserId()

	graceful, _ := beego.AppConfig.Bool("graceful")
	if !graceful {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM,
			syscall.SIGQUIT)
		go handleSignals(sigs)
	}

	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	mode := beego.AppConfig.String("runmode")
	if mode == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}

	beego.Run()
}
