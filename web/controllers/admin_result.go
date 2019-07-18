package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"lottery/models"
	"lottery/services"
)

type AdminResultController struct {
	Ctx           iris.Context
	ServiceResult services.ResultService
}

func (c *AdminResultController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	uid := c.Ctx.URLParamIntDefault("uid", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	pagePrev := ""
	pageNext := ""
	var datalist []models.LtResult
	var total int
	size := 100
	if giftId < 1 && uid < 1 {
		datalist = c.ServiceResult.GetAll(page, size)
	} else if giftId >= 1 && uid < 1 {
		datalist = c.ServiceResult.SearchByGift(giftId)
	} else if giftId < 1 && uid >= 1 {
		datalist = c.ServiceResult.SearchByUser(uid)
	} else {
		datalist = c.ServiceResult.SearchByGiftAndUser(uid, giftId)
	}
	total = (page-1)*size + len(datalist)
	if len(datalist) >= size {
		if giftId < 1 && uid < 1 {
			total = int(c.ServiceResult.CountAll())
		} else if giftId >= 1 && uid < 1 {
			total = int(c.ServiceResult.CountByGift(giftId))
		} else if giftId < 1 && uid >= 1 {
			total = int(c.ServiceResult.CountByUser(uid))
		} else {
			total = int(c.ServiceResult.CountByGiftAndUser(uid, giftId))
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/result.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "result",
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminResultController) GetDelete() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	if id < 0 {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("不存在的记录, <a href='' onclick='history.go(-1); return false;'>"),
		}
	}
	err := c.ServiceResult.Delete(id)
	if err != nil {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("Internal error:%s, <a href='' onclick='history.go(-1); return false;'>", err),
		}
	}

	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}

	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminResultController) GetCheat() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	if id < 0 {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("不存在的记录, <a href='' onclick='history.go(-1); return false;'>"),
		}
	}
	err := c.ServiceResult.Update(&models.LtResult{Id: id, SysStatus: 2}, []string{"sys_status"})
	if err != nil {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("Internal error:%s, <a href='' onclick='history.go(-1); return false;'>", err),
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminResultController) GetReset() mvc.Result {
	id := c.Ctx.URLParamIntDefault("id", 0)
	if id < 0 {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("不存在的记录, <a href='' onclick='history.go(-1); return false;'>"),
		}
	}
	err := c.ServiceResult.Update(&models.LtResult{Id: id, SysStatus: 0}, []string{"sys_status"})
	if err != nil {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("Internal error:%s, <a href='' onclick='history.go(-1); return false;'>", err),
		}
	}

	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/result"
	}

	return mvc.Response{
		Path: refer,
	}
}
