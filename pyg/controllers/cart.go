package controllers

import (
	"github.com/astaxie/beego"
	"gomodule/redigo/redis"
	"pyg/pyg/models"
	"github.com/astaxie/beego/orm"
)

type Cartcontrollers struct {
	beego.Controller
}

//处理添加购物车
func (this *Cartcontrollers)HandleAddCart()  {
	//获取数据
	id,err:=this.GetInt("goodsId")
	num,err2:=this.GetInt("num")

	//返回ajax步骤
	//1.定义一个map容器
	resp:=make(map[string]interface{})

	//封装 集成 多态(3.把容器传递给前段,4.指定传递方式)
	defer  RespFunc(&this.Controller,resp)
	//校验数据
	if err!=nil||err2!=nil{
		//2.给容器赋值
		resp["errno"]=1
		resp["errmsg"]="输入数据不完整"
		return
	}

	name:=this.GetSession("name")

	if name==nil{
		resp["errno"]=2
		resp["errmsg"]="当前用户未登录,不能添加购物车"
		return
	}


	//处理数据
	//把数据存储在redis的hash中
	conn,err := redis.Dial("tcp","192.168.91.88:6379")
	if err != nil {
		resp["errno"] = 3
		resp["errmsg"] = "redis连接失败服务器异常"
		beego.Error("err:",err)
		return
	}
	defer conn.Close()

	oldNum,_ := redis.Int(conn.Do("hget","cart_"+name.(string),id))

	_,err = conn.Do("hset","cart_"+name.(string),id,oldNum + num)
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "添加商品到购物车失败"
		return
	}

	//返回数据
	resp["errno"] = 5
	resp["errmsg"] = "OK"
}

//展示购物车
func (this *Cartcontrollers)ShowCart()  {
	//获取数据


	conn,err:=redis.Dial("tcp","192.168.91.88:6379")
	if err!=nil{
		beego.Error(err)
		this.Redirect("/indexsx",302)
		return
	}
	defer conn.Close()

	//查询所有购物车数据
	name:=this.GetSession("name" )

	result,err:=redis.Ints(conn.Do("hgetall","cart_"+name.(string)))

	if err!=nil{
		this.Redirect("/indexsx",302)
		return
	}

	//定义大容器
	var goods []map[string]interface{}

	o:=orm.NewOrm()


	totalPrice:=0
	totalCount:=0


	for i:=0;i<len(result);i+=2{

		temp:=make(map[string]interface{})
		//获取商品信息 商品数量
		var goodsSku models.GoodsSKU

		goodsSku.Id=result[i]
		err:=o.Read(&goodsSku)
		if err !=nil{
			beego.Error("查询id不存在")
			this.Redirect("/indexsx",302)
			return
		}
		//给行容器赋值
		temp["goodsSku"]=goodsSku
		temp["count"]=result[i+1]

		littlePrice:=result[i+1]*goodsSku.Price
		temp["littlePrice"]=littlePrice

		totalPrice+=littlePrice
		totalCount++
		//把航容器添加到大容器里面
		goods=append(goods,temp)

	}
	this.Data["totalPrice"]=totalPrice
	this.Data["totalCount"]=totalCount
	this.Data["goods"]=goods

	this.TplName="cart.html"
}

//处理购物车数量
func (this *Cartcontrollers)HandleUpCart()  {
	id,err:=this.GetInt("goodsId")
	count,err2:=this.GetInt("count")
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)

	if err != nil || err2 != nil {
		resp["errno"] = 1
		resp["errmsg"] = "传输数据不完整"
		return
	}
	name := this.GetSession("name")
	if name == nil {
		resp["errno"] = 3
		resp["errmsg"] = "当前用户未登录"
		return
	}

	conn,err:=redis.Dial("tcp","192.168.91.88:6379")
	if err!=nil{
		resp["errno"]=2
		resp["errmsg"]="redis连接错误"
		return
	}
	defer conn.Close()

	_,err=conn.Do("hset","cart_"+name.(string),id,count)
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "redis写入失败"
		return
	}
	resp["errno"] = 5
	resp["errmsg"] ="OK"
}


//处理删除购物车数量
func (this * Cartcontrollers)HandleDeleteCart()  {
	id,err:=this.GetInt("goodsId")
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller,resp)
	//校验数据  宏定义  枚举
	if err != nil {
		resp["errno"] = 1
		resp["errmsg"] = "删除链接错误"
		return
	}

	//查看是否是登录状态
	name := this.GetSession("name")
	if name == nil {
		resp["errno"] = 2
		resp["errmsg"] = "当前用户不在登录状态"
		return
	}
	//向redis中写入数据
	conn,err :=redis.Dial("tcp","192.168.91.88:6379")
	if err != nil {
		resp["errno"] = 3
		resp["errmsg"] = "服务器异常"
		return
	}

	defer conn.Close()

	_,err = conn.Do("hdel","cart_"+name.(string),id)
	if err != nil {
		resp["errno"] = 4
		resp["errmsg"] = "数据库异常"
		return
	}

	resp["errno"] = 5
	resp["errmsg"] = "OK"
}
