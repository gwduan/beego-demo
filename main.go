package main

import (
	_ "beego-demo/routers"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
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

func main() {
	setUserId()

	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	//beego.SetLevel(beego.LevelInformational)

	beego.Run()
}
