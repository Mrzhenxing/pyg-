package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"time"
	"strings"
	"fmt"
)

type OrderController struct {
	beego.Controller
}

//展示提交页面
func (this *OrderController) ShowOrder() {
	//获取数据
	goodsIds := this.GetStrings("checkGoods")

	//校验数据
	if len(goodsIds) == 0 {
		this.Redirect("/user/showCart", 302)
		return
	}
	//处理数据
	//获取当前用户的所有收货地址
	name := this.GetSession("name")

	o := orm.NewOrm()
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string)).All(&addrs)
	this.Data["addrs"] = addrs

	conn, _ := redis.Dial("tcp", "192.168.91.88:6379")

	//获取商品,获取总价和总件数
	var goods []map[string]interface{}
	var totalPrice, totalCount int

	for _, v := range goodsIds {
		temp := make(map[string]interface{})
		id, _ := strconv.Atoi(v)
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		o.Read(&goodsSku)

		//获取商品数量
		count, _ := redis.Int(conn.Do("hget", "cart_"+name.(string), id))

		//计算小计
		littlePrice := count * goodsSku.Price

		//把商品信息放到行容器
		temp["goodsSku"] = goodsSku
		temp["count"] = count
		temp["littlePrice"] = littlePrice

		totalPrice += littlePrice
		totalCount += 1

		goods = append(goods, temp)

	}

	//返回数据
	this.Data["goodsIds"] = goodsIds
	this.Data["totalPrice"] = totalPrice
	this.Data["totalCount"] = totalCount
	this.Data["truePrice"] = totalPrice + 10
	this.Data["goods"] = goods
	this.TplName = "place_order.html"
}

//处理提交页面
func (this *OrderController) HandlePushOrder() {
	//获取数据
	//param={"addrId":addrId,"payId":payId,"goodsIds":goodsIds,"totalCount":totalCount,"totalPrice":totalprice}
	addrId, err1 := this.GetInt("addrId")
	payId, err2 := this.GetInt("payId")
	goodsIds := this.GetString("goodsIds")
	totalCount, err3 := this.GetInt("totalCount")
	totalPrice, err4 := this.GetInt("totalPrice")
	//返回ajax四步骤
	resp := make(map[string]interface{})
	defer RespFunc(&this.Controller, resp)
	//获取用户id
	name := this.GetSession("name")
	if name == nil {
		resp["errno"] = 2
		resp["errmsg"] = "用户为登录,请登录"
		return
	}
	//fmt.Println(1111111111111111)
	//fmt.Println(name)

	//校验数据
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || goodsIds == "" {
		resp["errno"] = 1
		resp["errmsg"] = "获取数据不完整"
		this.Redirect("/user/addOrder", 302)
		return
	}
	//获取用户和地址
	o := orm.NewOrm()
	var user models.User
	user.Name = name.(string)
	o.Read(&user,"Name")

	var addr models.Address
	addr.Id = addrId
	o.Read(&addr)
	//获取订单表
	var orderinfo models.OrderInfo

	orderinfo.User = &user
	orderinfo.Address = &addr
	orderinfo.PayMethod = payId
	orderinfo.TotalCount = totalCount
	orderinfo.TotalPrice = totalPrice
	orderinfo.TransitPrice = 10
	orderinfo.TradeNo = time.Now().Format("20060102150405" + strconv.Itoa(user.Id))

	//在数据插入前开启事物,便于库存不足时回到最初的状态
	o.Begin()
	o.Insert(&orderinfo)
	conn, err := redis.Dial("tcp", "192.168.91.88:6379")
	if err != nil {
		beego.Error("redis连接服务器失败", err)
	}
	defer conn.Close()

	//插入订单商品表
	//获取商品id goodsid是[1 2 3]类型字符串
	goodsSlice := strings.Split(goodsIds[1:len(goodsIds)-1], " ")
	for _, v := range goodsSlice {



		var goodsSku models.GoodsSKU
		//把字符串id转换为intid
		id, _ := strconv.Atoi(v)
		//fmt.Println(id)
		goodsSku.Id = id
		o.Read(&goodsSku)

		//获取原始库存
		oldStock :=goodsSku.Stock


		//从redis获取商品数量,先连接,再获取,再用回复助手函数转成int
		count, err := redis.Int(conn.Do("hget", "cart_"+name.(string), id))
		if err != nil {
			beego.Error("redis获取商品数量失败", err)
		}
		//插入
		var ordergoods models.OrderGoods
		ordergoods.GoodsSKU = &goodsSku

		ordergoods.OrderInfo = &orderinfo

		ordergoods.Count = count

		ordergoods.Price = goodsSku.Price * count

		//fmt.Println(goodsSku.Stock)
		//插入前要更新商品库存和销量
		if goodsSku.Stock < count {
			fmt.Println(4444444444444)
			resp["errno"] = 4
			resp["errmsg"] = "库存不足"
			o.Rollback()
			return
		}

		//time.Sleep(time.Second*5)
		o.Read(&goodsSku)
		//库存
		goodsSku.Stock -= count
		//销量
		goodsSku.Sales += count

		qs:=o.QueryTable("GoodsSKU").Filter("Id",id).Filter("Stock",oldStock)
		qs.Update(orm.Params{"Stock":goodsSku.Stock-count,"Sales":goodsSku.Sales+count})
		o.Update(&goodsSku)
		fmt.Println(goodsSku)
		_, err = o.Insert(&ordergoods)
		fmt.Println(ordergoods)
		if err != nil {
			beego.Error(err)
			resp["errno"] = 3
			resp["errmsg"] = "服务器异常"
			o.Rollback()
			return
		}
		_, err = conn.Do("hdel", "cart_"+name.(string),id)
			if err!=nil{
				beego.Error(err)
				o.Rollback()
				return
				}

	}
	o.Commit()
	resp["errno"] = 5
	resp["errmsg"] = "ok"
	this.TplName = "index.html"
}
