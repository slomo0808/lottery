package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"lottery/comm"
	"lottery/models"
	"lottery/services"
)

type AdminUserController struct {
	Ctx         iris.Context
	ServiceUser services.UserService
}

func (c *AdminUserController) Get() mvc.View {
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	datalist := c.ServiceUser.GetAll(page, size)
	var total int
	var pagePrev, pageNext string

	total = (page-1)*size + len(datalist)
	if len(datalist) >= size {
		total = int(c.ServiceUser.CountAll())
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/user.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "user",
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
			"Now":      comm.NowUnix(),
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminUserController) GetBlack() mvc.Result {
	uid, err := c.Ctx.URLParamInt("id")
	time, err2 := c.Ctx.URLParamInt("time")
	if err == nil && err2 == nil {
		timeBlack := comm.NowUnix() + time*86400
		c.ServiceUser.Update(
			&models.LtUser{Id: uid, Blacktime: timeBlack, SysUpdated: comm.NowUnix()},
			[]string{"blacktime"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "admin/user"
	}
	return mvc.Response{
		Path: refer,
	}
}
