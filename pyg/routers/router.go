package routers

import (
	"pyg/pyg/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	//路由过滤起 校验是否登录
	beego.InsertFilter("/user/*", beego.BeforeExec, guolvfuc)
	beego.Router("/", &controllers.MainController{})

	beego.Router("/register", &controllers.Usercontrollers{}, "get:ShowRegister;post:HandleRegister")

	beego.Router("/sendMsg", &controllers.Usercontrollers{}, "post:HandleSendMsg")
	//邮箱验证
	beego.Router("/register-email", &controllers.Usercontrollers{}, "get:ShowEmail;post:HandleEmail")
	//    激活用
	beego.Router("/active", &controllers.Usercontrollers{}, "get:Active")
	//登录
	beego.Router("/login", &controllers.Usercontrollers{}, "get:ShowLogin;post:HandleLogin")
	//首页
	beego.Router("/index", &controllers.Goodscontrollers{}, "get:ShowIndex")
	//	退出
	beego.Router("/user/loginout", &controllers.Usercontrollers{}, "get:LoginOut")
	//	用户中心
	beego.Router("/user/user_center_info", &controllers.Usercontrollers{}, "get:ShowUserCenterInfo")
	//	用户地址
	beego.Router("/user/user_center_site", &controllers.Usercontrollers{}, "get:ShowSite;post:HandleSite")
	//	生鲜模块
	beego.Router("/indexsx", &controllers.Goodscontrollers{}, "get:ShowIndexsx")
	//展示商品详情
	beego.Router("/goodsDetail", &controllers.Goodscontrollers{}, "get:ShowDetail")
	////展示商品清单
	beego.Router("/goodsType", &controllers.Goodscontrollers{}, "get:ShowList")
	//	加入购物车
	beego.Router("/addCart", &controllers.Cartcontrollers{}, "post:HandleAddCart")
	//展示购物车
	beego.Router("/user/showCart",&controllers.Goodscontrollers{},"get:ShowCart")
}
func guolvfuc(ctx *context.Context) {
	name := ctx.Input.Session("name")
	if name == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
