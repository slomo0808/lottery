package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"lottery/services"
)

type AdminController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceResult  services.ResultService
	ServiceBlackip services.BlackipService
	ServiceUserday services.UserdayService
	ServiceCode    services.CodeService
}

func (c *AdminController) Get() mvc.View {
	return mvc.View{
		Name: "admin/index.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "",
		},
		Layout: "admin/layout.html",
	}
}
