package main

import (
	_ "cmdb/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
)

var globalSessions *session.Manager

var certoken = func(ctx *context.Context) {
	token := ctx.Input.Cookie("Token")
	uri := ctx.Input.URL()
	if token == "" && uri != "/user/login/" {
		ctx.Redirect(302, "/login")
	}

}

func main() {

	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "192.168.3.105:6379"
	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.InsertFilter("/*", beego.BeforeRouter, certoken)
	beego.Run()
}
