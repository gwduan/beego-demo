package main

import (
	_ "beego-demo/routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
)

func main() {
	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.Run()
}

