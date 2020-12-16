package routers

import (
	"cmdb/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//用户接口
	user := beego.NewNamespace("/user",
		beego.NSRouter("/login", &controllers.UserController{}, "*:Login"),
		beego.NSRouter("/info", &controllers.UserController{}, "*:Userinfo"),
		// beego.NSRouter("/GetUserInfo", &controllers.UserController{}, "*:GetUserInfo"),
		// beego.NSRouter("/count ", &controllers.UserController{}, "*:UserCount"),
		beego.NSRouter("/logout", &controllers.UserController{}, "*:Logout"),
		beego.NSInclude(
			&controllers.UserController{},
		),
	)

	rancher := beego.NewNamespace("/rancher",
		beego.NSRouter("/cluster", &controllers.RancherController{}, "*:Gettoken"),
		beego.NSRouter("/project", &controllers.RancherController{}, "*:Getproject"),
		beego.NSRouter("/worker", &controllers.RancherController{}, "*:Getworker"),
		beego.NSRouter("/changeworker", &controllers.RancherController{}, "*:Changeworker"),
		beego.NSInclude(
			&controllers.RancherController{},
		),
	)
	beego.AddNamespace(user)
	beego.AddNamespace(rancher)
}
