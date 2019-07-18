package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"lottery/comm"
	"lottery/models"
	"lottery/services"
	"strconv"
)

type IndexController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceResult  services.ResultService
	ServiceBlackip services.BlackipService
	ServiceUserday services.UserdayService
	ServiceCode    services.CodeService
}

func (c *IndexController) Get() string {
	c.Ctx.Header("Content-Type", "text/html")
	return "welcome to Go抽奖系统, <a href='public/home.html'>开始抽奖</a>"
}

// localhost:8080/gifts
// 显示所有奖品信息
func (c *IndexController) GetGifts() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	datalist := c.ServiceGift.GetAll(true)
	list := make([]models.LtGift, 0)
	for _, data := range datalist {
		if data.SysStatus == 0 {
			list = append(list, data)
		}
	}
	rs["gift"] = list
	return rs
}

// localhost:8080/newprize
// 显示最新的获奖列表
func (c *IndexController) GetNewprize() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	rs["data"] = c.ServiceResult.GetNewPrize(10, []int{1, 2, 3})

	return rs
}

func (c *IndexController) GetLogin() {
	uid := comm.Random(100000)
	var loginuser = &models.ObjLoginUser{
		Uid:      uid,
		Username: fmt.Sprintf("user" + strconv.Itoa(uid)),
		Now:      comm.NowUnix(),
		Ip:       comm.ClientIp(c.Ctx.Request()),
	}
	// 生成cookie并放入浏览器
	comm.SetLoginUser(c.Ctx.ResponseWriter(), loginuser)

	// 跳转
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/home.html?from=login")
}

func (c *IndexController) GetLogout() {
	// 生成空cookie并放入浏览器
	comm.SetLoginUser(c.Ctx.ResponseWriter(), nil)

	// 跳转
	comm.Redirect(c.Ctx.ResponseWriter(), "/public/home.html?from=logout")
}
