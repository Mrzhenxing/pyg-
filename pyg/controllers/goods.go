package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"

	"math"
	"fmt"
)

type Goodscontrollers struct {
	beego.Controller
}
//展示首页
func (this *Goodscontrollers)ShowIndex()  {
	name:=this.GetSession("name")
	if name!=nil{
		//name是interface类型  所以要类型断言
		this.Data["name"]=name.(string)
	}else {
		this.Data["name"]=""
	}

	//获取类型信息并传递给前段
	//获取一级菜单
	o := orm.NewOrm()
	//接受对象
	var oneClass []models.TpshopCategory
	//查询
	o.QueryTable("TpshopCategory").Filter("Pid",0).All(&oneClass)


	//获取第二级
	var types []map[string]interface{}//定义总容器
	for _,v := range oneClass{
		//行容器
		t := make(map[string]interface{})

		var secondClass []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid",v.Id).All(&secondClass)
		t["t1"] = v  //一级菜单对象
		t["t2"] = secondClass  //二级菜单集合
		//把行容器加载到总容器中
		types = append(types,t)
	}

	//获取第三季菜单
	for _,v1 := range types{
		//循环获取二级菜单
		var erji []map[string]interface{} //定义二级容器
		for _,v2 := range v1["t2"].([]models.TpshopCategory){
			t := make(map[string]interface{})
			var thirdClass []models.TpshopCategory
			//获取三级菜单
			o.QueryTable("TpshopCategory").Filter("Pid",v2.Id).All(&thirdClass)
			t["t22"] = v2  //二级菜单
			t["t23"] = thirdClass   //三级菜单
			erji = append(erji,t)
		}
		//把二级容器放到总容器中
		v1["t3"] = erji
	}


	this.Data["types"] = types


	this.TplName="index.html"

}

//展示生鲜首页
func (this *Goodscontrollers)ShowIndexsx()  {

	//获取生鲜首页内容
	//获取所有类型
	o:=orm.NewOrm()
	var GoodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&GoodsTypes)
	this.Data["GoodsType"]=GoodsTypes

	//获取轮播图片
	var IndexGoodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&IndexGoodsBanners)
	this.Data["IndexGoodsBanner"]=IndexGoodsBanners
	//fmt.Println(IndexGoodsBanners)

	//获取促销商品展示
	var IndexPromotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&IndexPromotionBanners)
	this.Data["IndexPromotionBanner"]=IndexPromotionBanners

	//获取首页分类商品展示
	var goods []map[string]interface{}

	for _,v:=range GoodsTypes{
		//var IndexTypeGoodsBanner []models.IndexTypeGoodsBanner
		//过滤类型等于v的 GoodsType__Id父类v.Id子类
		qs:=o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").Filter("GoodsType__Id",v.Id).OrderBy("Index")
		var textGoods []models.IndexTypeGoodsBanner
		var imageGoods []models.IndexTypeGoodsBanner
		qs.Filter("DisplayType",0).All(&textGoods)
		qs.Filter("DisplayType",1).All(&imageGoods)
		//fmt.Println("1111111111111111")
		//fmt.Println(v.Image)
		//定义行容器
		temp:=make(map[string]interface{})
		//行容器中有个类型对象v 有textGoods文章商品和imageGoods图片商品
		temp["GoodsTypes"]=v
		temp["textGoods"]=textGoods
		temp["imageGoods"]=imageGoods
		goods=append(goods,temp)

		}
		this.Data["goods"]=goods


	/*	o:=orm.NewOrm()
		var goodsTypes []models.GoodsType
		o.QueryTable("GoodsType").All(&goodsTypes)
		this.Data["goodsType"]=goodsTypes


		//获取轮播图
		var goodsBanners []models.IndexGoodsBanner
		o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&goodsBanners)
		this.Data["goodsBanners"]=goodsBanners

		//获取促销商品
		var promotionBanners []models.IndexPromotionBanner
		o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionBanners)
		this.Data["promotions"] = promotionBanners

		//获取首页商品展示
		var goods []map[string]interface{}

		for _,v := range goodsTypes{
			var textGoods []models.IndexTypeGoodsBanner
			var imageGoods []models.IndexTypeGoodsBanner
			qs:=o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType","GoodsSKU").Filter("GoodsType__Id",v.Id).OrderBy("Index")
			//获取文字商品
			qs.Filter("DisplayType",0).All(&textGoods)
			//获取图片商品
			qs.Filter("DisplayType",1).All(&imageGoods)

			//定义行容器
			temp := make(map[string]interface{})
			temp["goodsType"] = v
			temp["textGoods"] = textGoods
			temp["imageGoods"] = imageGoods

			//把行容器追加到总容器中
			goods = append(goods,temp)
		}
		this.Data["goods"] = goods*/
	
	
	this.TplName="index_sx.html"
}
//展示商品详情
func (this *Goodscontrollers)ShowDetail(){
	//获取数据
	id,err:=this.GetInt("Id")
	//校验数据
	if err!=nil{
		beego.Error("获取连接数据失败")
		this.Redirect("/indexsx",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var goodsSku models.GoodsSKU
	//goodsSku.Id=id
	////fmt.Println(id)
	//err=o.Read(&goodsSku)
	////fmt.Println(111111111)
	//if err!=nil{
	//	beego.Error("id不存在，获取id失败")
	//	this.Redirect("/indexsx",302)
	//	return
	//}
	//展示商品详情实现
	o.QueryTable("GoodsSKU").RelatedSel("Goods","GoodsType").Filter("Id",id).One(&goodsSku)


	//获取同一类型的新品推荐
	var newGoods []models.GoodsSKU
	qs:=o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",goodsSku.GoodsType.Name)
	qs.OrderBy("-Time").Limit(2,0).All(&newGoods)
	//返回数据
	this.Data["newGoods"]=newGoods
	this.Data["goodsSku"]=goodsSku
	this.TplName="detail.html"
}
//封装分页函数
func PageEdit(pageCount int,pageIndex int)[]int{
	//不足五页
	var pages []int
	if pageCount < 5{
		for i:=1;i<=pageCount;i++{
			pages = append(pages,i)
		}
	}else if pageIndex <= 3{
		for i:=1;i<=5;i++{
			pages = append(pages,i)
		}
	}else if pageIndex >= pageCount -2{
		for i:=pageCount - 4;i<=pageCount;i++{
			pages = append(pages,i)
		}
	}else {
		for i:=pageIndex - 2;i<=pageIndex + 2;i++{
			pages = append(pages,i)
		}
	}

	return pages
}
//展示商品列表页
func (this *Goodscontrollers)ShowList()  {
	//获取数据
	/*id,err:=this.GetInt("id")
	//校验数据
	//fmt.Println(111111111111)
	if err!=nil{
		beego.Error("获取商品类型连接失败")
		this.Redirect("/indexsx",302)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	//获取所有同一类型的商品
	var goods []models.GoodsSKU
	//或排序方式
	sort:=this.GetString("sort")
	if sort==""{
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).All(&goods)
	}else if sort=="price"{
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("Price").All(&goods)
	}else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).OrderBy("-Sales").All(&goods)
	}
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id).All(&goods)

	this.Data["id"]=id
	this.Data["sort"]=sort
	this.Data["goods"]=goods
	//fmt.Println(222222222222222,goods)
	this.TplName="list.html"*/
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("类型不存在")
		this.Redirect("/index_sx",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var goods []models.GoodsSKU
	//获取排序方式
	sort := this.GetString("sort")

	//实现分页

	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id)
	//获取总页码
	count,_ := qs.Count()
	pageSize := 1
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	pages := PageEdit(pageCount,pageIndex)
	this.Data["pages"] = pages
	//获取上一页，下一页的值
	var prePage,nextPage int
	//设置个范围
	if pageIndex -1 <= 0{
		prePage = 1
	}else {
		prePage = pageIndex - 1
	}


	if pageIndex +1 >= pageCount{
		nextPage = pageCount
	}else {
		nextPage = pageIndex + 1
	}


	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage

	qs = qs.Limit(pageSize,pageSize*(pageIndex - 1))

	//获取排序
	if sort == ""{
		qs.All(&goods)
	}else if sort == "price"{
		qs.OrderBy("Price").All(&goods)
	}else {
		qs.OrderBy("-Sales").All(&goods)
	}

	this.Data["sort"] = sort

this.Data["pageIndex"]=pageIndex
fmt.Println(pageIndex)

	//返回数据
	this.Data["id"] = id
	this.Data["goods"] = goods
	this.TplName = "list.html"
}







/*
//展示商品详情
func(this*Goodscontrollers)ShowDetail(){
	//获取数据
	id,err := this.GetInt("Id")
	//校验数据
	if err != nil {
		beego.Error("商品链接错误")
		this.Redirect("/index_sx",302)
		return
	}
	//处理数据
	//根据id获取商品有关数据
	o := orm.NewOrm()
	var goodsSku models.GoodsSKU
	//goodsSku.Id = id
	//o.Read(&goodsSku)
	o.QueryTable("GoodsSKU").RelatedSel("Goods","GoodsType").Filter("Id",id).One(&goodsSku)

	//获取同一类型的新品推荐
	var newGoods []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name",goodsSku.GoodsType.Name)
	qs.OrderBy("-Time").Limit(2,0).All(&newGoods)
	this.Data["newGoods"] = newGoods
	//传递数据
	this.Data["goodsSku"] = goodsSku
	this.TplName = "detail.html"
}

//分页函数
func PafeEdit(pageCount int,pageIdex int)[]int  {
	//不足五页
	var pages []int
	if pageCount<5{
		for i:=1;i<pageCount ;i++  {
			pages = append(pages,i)
		}
	}else if pageCount<=3{
		for i:=1;i<5 ;i++  {
			pages=append(pages,i)
		}

	}else if pageIdex>=pageCount-2{
		for i:=pageCount-4;i<=pageCount;i++{
			pages =append(pages,i)
		}
	}
	return pages
}

//展示详情页
func (this *Goodscontrollers)ShowList()  {
	//获取数据
	//id ,err:=this.GetInt("id")
	id,err := this.GetInt("id")
	//校验数据
	if err != nil {
		beego.Error("类型不存在")
		this.Redirect("/index_sx",302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var goods []models.GoodsSKU
	//获取排序方式
	sort := this.GetString("sort")

	//实现分页

	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",id)
	//获取总页码
	count,_ := qs.Count()
	pageSize := 1
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取当前页码
	pageIndex,err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	pages := PageEdit(pageCount,pageIndex)
	this.Data["pages"] = pages
	//获取上一页，下一页的值
	var prePage,nextPage int
	//设置个范围
	if pageIndex -1 <= 0{
		prePage = 1
	}else {
		prePage = pageIndex - 1
	}


	if pageIndex +1 >= pageCount{
		nextPage = pageCount
	}else {
		nextPage = pageIndex + 1
	}


	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage

	qs = qs.Limit(pageSize,pageSize*(pageIndex - 1))

	//获取排序
	if sort == ""{
		qs.All(&goods)
	}else if sort == "price"{
		qs.OrderBy("Price").All(&goods)
	}else {
		qs.OrderBy("-Sales").All(&goods)
	}

	this.Data["sort"] = sort



	//返回数据
	this.Data["id"] = id
	this.Data["goods"] = goods
	this.TplName = "list.html"
}*/

