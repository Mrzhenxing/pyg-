package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"regexp"
	"time"
	"fmt"
	"math/rand"
	"github.com/astaxie/beego/orm"
	"pyg/pyg/models"
	"github.com/astaxie/beego/utils"
	//_ "github.com/weilaihui/fdfs_client"
)

type Usercontrollers struct {
	beego.Controller
}

//展示注册页面
func (this *Usercontrollers) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册页面
func RespFunc(this *beego.Controller, resp map[string]interface{}) {
	//3.把容器传递给前段
	this.Data["json"] = resp
	//4.指定传递方式
	this.ServeJSON()
}

//定义信息结构体
type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}

//发送短信
func (this *Usercontrollers) HandleSendMsg() {
	//接受数据
	phone := this.GetString("phone")
	resp := make(map[string]interface{})

	defer RespFunc(&this.Controller, resp)
	//返回json格式数据
	//校验数据
	if phone == "" {
		beego.Error("获取电话号码失败")
		//2.给容器赋值
		resp["errno"] = 1
		resp["errmsg"] = "获取电话号码错误"
		return
	}
	//检查电话号码格式是否正确
	reg, _ := regexp.Compile(`^1[3-9][0-9]{9}$`)
	result := reg.FindString(phone)
	if result == "" {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 2
		resp["errmsg"] = "电话号码格式错误"
		return
	}
	//发送短信   SDK调用
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "LTAIu4sh9mfgqjjr", "sTPSi0Ybj0oFyqDTjQyQNqdq9I9akE")
	if err != nil {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 3
		resp["errmsg"] = "初始化短信错误"
		return
	}

	rand.Seed(time.Now().UnixNano())
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06d", rand.Int31n(1000000))
	fmt.Println("---------", vcode)

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = "品优购"
	request.QueryParams["TemplateCode"] = "SMS_164275022"
	request.QueryParams["TemplateParam"] = `{"code":` + vcode + `}`

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		beego.Error("电话号码格式错误", err)
		//2.给容器赋值
		resp["errno"] = 4
		resp["errmsg"] = "短信发送失败"
		return
	}
	//json数据解析
	var message Message
	json.Unmarshal(response.GetHttpContentBytes(), &message)
	if message.Message != "OK" {
		beego.Error("电话号码格式错误")
		//2.给容器赋值
		resp["errno"] = 6
		resp["errmsg"] = message.Message
		return
	}

	resp["errno"] = 5
	resp["errmsg"] = "发送成功"
	resp["code"] = vcode

}

//处理提交出册页面
func (this *Usercontrollers) HandleRegister() {
	//
	password := this.GetString("password")
	repassword := this.GetString("repassword")
	phone := this.GetString("phone")

	//校验数据
	if phone == "" || password == "" || repassword == "" {
		beego.Error("密码不能为空")
		this.TplName = "register.html"
		return
	}
	if password != repassword {
		beego.Error("两次密码输入不一致")
		this.TplName = "register.html"
		return
	}

	//处理数据
	//创建orm对象
	o := orm.NewOrm()
	//创建处理对象
	var user models.User

	user.Name = phone
	user.Pwd = password
	user.Phone = phone
	_, err := o.Insert(&user)
	fmt.Println(1111111111111111)
	fmt.Println(err)
	fmt.Println(user.Name)
	this.Ctx.SetCookie("userName", user.Name, 60*5)
	//返回数据
	this.Redirect("/register-email", 302)
}

//展示邮箱
func (this *Usercontrollers) ShowEmail() {
	this.TplName = "register-email.html"

}

//处理邮箱
func (this *Usercontrollers) HandleEmail() {
	//获取数据
	email := this.GetString("email")
	password := this.GetString("password")
	repassword := this.GetString("repassword")

	//校验数据

	if email == "" || password == "" || repassword == "" {
		beego.Error("输入不争取，请重新输入，注：密码或者邮箱不能为空")
		this.Data["errmsg"] = "输入不争取，请重新输入，注：密码或者邮箱不能为空"
		this.TplName = "register-email.html"
		return
	}
	if password != repassword {
		beego.Error("两次密码输入不一致")
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register-email.html"
		return
	}
	//校验邮箱
	regx, err := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	if err != nil {
		beego.Error("匹配错误", err)
		return
	}
	result := regx.FindString(email)
	if result == "" {
		beego.Error("邮箱输入格式不正确")
		this.Data["errmsg"] = "邮箱输入格式不正确"
		this.TplName = "register-email.html"
		return
	}

	//处理数据
	//处理数据
	//发送邮件
	//utils     全局通用接口  工具类  邮箱配置
	config := `{"username":"632703728@qq.com","password":"rwyfzgkshdfjbfcc","host":"smtp.qq.com","port":587}`
	emailReg := utils.NewEMail(config)

	//内容配置
	emailReg.Subject = "品优购用户激活"
	emailReg.From = "632703728@qq.com"
	emailReg.To = []string{email}
	userName := this.Ctx.GetCookie("userName")
	emailReg.HTML = `<a href="http://192.168.91.88:8080/active?userName=` + userName + `"> 点击激活该用户</a>`

	//发送
	emailReg.Send()

	//插入邮箱
	o := orm.NewOrm()
	var user models.User
	//由于更新 需要先查询 在赋值更新
	user.Name = userName
	fmt.Println(222222222222222)
	fmt.Println(userName)
	err = o.Read(&user, "Name")
	fmt.Println(33333333333333)
	fmt.Println(err)
	if err != nil {
		beego.Error("用户不存在")
		this.Redirect("/register-email", 302)
		return
	}
	user.Email = email
	o.Update(&user, "Email")
	//处理数据

	//返回数据

	//返回数据

	this.Ctx.WriteString("邮件已发送，请去目标邮箱激活用户！")

	//返回数据

}

//激活
func (this *Usercontrollers) Active() {
	//获取数据
	userName := this.GetString("userName")
	//校验数据
	if userName == "" {
		beego.Error("用户名错误")
		this.Redirect("/register-email", 302)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户不存在")
		this.Redirect("/register-email", 302)
		return
	}
	user.Active = true
	o.Update(&user)

	//返回数据
	this.Redirect("/login", 302)
}

//展示登录页面
func (this *Usercontrollers) ShowLogin() {
	name := this.Ctx.GetCookie("Loginname")
	if name == "" {
		this.Data["checked"] = ""
	} else {
		this.Data["checked"] = "checked"
	}
	this.Data["name"] = name

	this.TplName = "login.html"
}

//处理登录页面
func (this *Usercontrollers) HandleLogin() {
	//获取数据
	name := this.GetString("name")
	pwd := this.GetString("pwd")

	//判读数据
	if name == "" || pwd == "" {
		this.Data["errmsg"] = "获取数据错误"
		this.TplName = "login.html"
		return
	}
	//处理数据
	o := orm.NewOrm()
	var user models.User

	regx, err := regexp.Compile(`^\w[\w\.-]*@[0-9a-z][0-9a-z-]*(\.[a-z]+)*\.[a-z]{2,6}$`)
	if err != nil {
		beego.Error("正则转化错误")
		return
	}
	result := regx.FindString(name)
	//如果匹配到了
	if result != "" {
		user.Email = name
		err := o.Read(&user, "Email")
		if err != nil {
			this.Data["errmsg"] = "邮箱未注册"
			this.TplName = "login.html"
			return
		}
		if user.Pwd != pwd {
			this.Data["errmsg"] = "密码不正确"
			this.TplName = "login.html"
			return

		}
		this.Redirect("/index", 302)

	} else {
		user.Name = name
		err := o.Read(&user, "Name")
		if err != nil {
			this.Data["errmsg"] = "用户不存在"
			this.TplName = "login.html"
			return
		}
		if user.Pwd != pwd {
			this.Data["errmsg"] = "密码不正确"
			this.TplName = "login.html"
			return
		}
		if user.Active == false {
			this.Data["errmsg"] = "用户未激活"
			this.TplName = "login.html"
			return
		}
	}

	//返回数据
	m1 := this.GetString("m1")
	if m1 == "2" {
		this.Ctx.SetCookie("Loginname", user.Name, 60*60)
	} else {
		this.Ctx.SetCookie("Loginname", user.Name, -1)
	}
	this.SetSession("name", user.Name)
	this.Redirect("/index", 302)

}

//退出处理
func (this *Usercontrollers) LoginOut() {
	this.DelSession("name")
	this.Redirect("/index", 302)
}

//展示用户中心
func (this *Usercontrollers) ShowUserCenterInfo() {
	o:=orm.NewOrm()
	name:=this.GetSession("name")
	var user models.User
	//获取用户名
	user.Name=name.(string)
	//查询用户名
	o.Read(&user,"Name")
	//查询默认收货地址
	var address models.Address
	err:=o.QueryTable("Address").RelatedSel("User").Filter("User__Name",name.(string)).Filter("IsDefault",true).One(&address)
	if err!=nil{
		beego.Error("获取默认地址失败")
	}
	this.Data["show"]="用户中心"
	this.Layout="layout_center_info.html"
	this.Data["address"]=address
	this.Data["user"]=user
	this.TplName = "user_center_info.html"
}

//展示收货地址
func (this *Usercontrollers) ShowSite() {
	///要展示当前地址 需要获取当前用户的默认地址
	o:=orm.NewOrm()
	name:=this.GetSession("name")
	//查询默认地址
	var address models.Address
	err:=o.QueryTable("Address").RelatedSel("User").Filter("User__Name",name.(string)).Filter("IsDefault",true).One(&address)
	if err!=nil{
		beego.Error("获取默认地址失败")
	}
	qian:=address.Phone[:3]
	hou:=address.Phone[7:]
	
	this.Layout="layout_center_info.html"
	address.Phone=qian+"****"+hou
	this.Data["show"]="收货地址"
	this.Data["address"]=address
	this.TplName = "user_center_site.html"
}

//处理收货地址
func (this *Usercontrollers) HandleSite() {
	//获取数据
	receive := this.GetString("receive")
	address := this.GetString("address")
	postCode := this.GetString("postCode")
	phone := this.GetString("phone")
	//校验数据
	if receive == "" || address == "" || postCode == "" || phone == "" {
		beego.Error("获取数据失败")
		this.TplName = "user_center_site.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	var addr models.Address

	addr.Receiver = receive
	addr.Addr = address
	addr.PostCode = postCode
	addr.Phone = phone
	name := this.GetSession("name")
	fmt.Println("11111111111111", name)
	fmt.Println("11111111111111", phone)

	var user models.User
	user.Name = name.(string)
	//先查询在赋值
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户不存在")
		this.TplName = "user_center_site.html"
		return
	}
	addr.User = &user
	//查询用户是否有默认地址，如果没有，直接插入，如果有，把默认地址修改为非默认地址
	//查询当前用户是否有默认地址
	var oldaddress models.Address
	err = o.QueryTable("Address").RelatedSel("User").Filter("User__Name", name.(string)).Filter("IsDefault", true).One(&oldaddress)

	//if err!=nil{
	//	addr.IsDefault=true
	//
	//}else {
	//	addr.IsDefault=false
	//	o.Update(&oldaddress})
	//}
	//代码优化

	if err == nil {
		oldaddress.IsDefault=false
		o.Update(&oldaddress)
	}
		addr.IsDefault=true

	fmt.Println("222222222222", addr)
	_, err = o.Insert(&addr)
	if err != nil {
		beego.Error("插入数据失败", err)
		this.TplName = "user_center_site.html"
		return
	}
	//返回数据
	this.Redirect("/user/user_center_site", 302)

}

func (this *Usercontrollers)ShowUserOrder()  {
	
}
