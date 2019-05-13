package controllers

import "github.com/astaxie/beego"

type Goodscontrollers struct {
	beego.Controller
}

func (this *Goodscontrollers)ShowIndex()  {
	name:=this.GetSession("name")
	if name!=nil{
		//name是interface类型  所以要类型断言
		this.Data["name"]=name.(string)
	}else {
		this.Data["name"]=""
	}
	this.TplName="index.html"

}