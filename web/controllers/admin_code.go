package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"strings"
)

type AdminCodeController struct {
	Ctx         iris.Context
	ServiceCode services.CodeService
	ServiceGift services.GiftService
}

func (c *AdminCodeController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	var num, cacheNum int
	var datalist []models.LtCode
	var total int
	if giftId > 0 {
		datalist = c.ServiceCode.GetByGift(giftId)
		num, cacheNum = utils.GetCacheCodeNum(giftId, c.ServiceCode)

	} else {
		datalist = c.ServiceCode.GetAll(page, size)

	}
	total = (page-1)*size + len(datalist)
	if len(datalist) >= size {
		if giftId > 0 {
			total = int(c.ServiceCode.CountByGift(giftId))
		} else {
			total = int(c.ServiceCode.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/code.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "code",
			"Datalist": datalist,
			"GiftId":   giftId,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
			"Total":    total,
			"CodeNum":  num,
			"CacheNum": cacheNum,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminCodeController) PostImport() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	if giftId < 1 {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("没有指定的奖品ID,无法进行导入<a href='' onclick='history.go(-10;return false;'>返回</a> "),
		}
	}
	gift := c.ServiceGift.Get(giftId, false)
	if gift == nil || gift.Id < 1 || gift.Gtype != conf.GtypeCodeDiff {
		return mvc.Response{
			ContentType: "text/html",
			Text:        fmt.Sprintf("奖品信息不存在或奖品类型不是差异化虚拟券，无法导入<a href='' onclick='history.go(-10;return false;'>返回</a> "),
		}
	}
	codes := c.Ctx.FormValue("codes")
	codeArr := strings.Split(codes, "\n")
	var sucNum, errNum int
	for _, code := range codeArr {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		err := c.ServiceCode.Create(&models.LtCode{
			GiftId:     giftId,
			Code:       code,
			SysCreated: comm.NowUnix(),
			SysUpdated: comm.NowUnix(),
			SysStatus:  0,
		})
		if err == nil {
			// 成功导入数据库，还需导入缓存
			ok := utils.ImportCacheCodes(giftId, code)
			if ok {
				sucNum++
			} else {
				errNum++
			}
		} else {
			log.Println("admin_code.PostImport.Create error:", err)
			errNum += 1
		}
	}
	return mvc.Response{
		ContentType: "text/html",
		Text: fmt.Sprintf("成功导入%d条，导入失败%d条, <a href='' onclick='history.go(-1);return false;'>返回</a>",
			sucNum, errNum),
	}
}

func (c *AdminCodeController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		return mvc.Response{
			Text: fmt.Sprintf("没有这个优惠券, error:%s", err),
		}
	}
	err = c.ServiceCode.Delete(id)
	if err != nil {
		return mvc.Response{
			Text: fmt.Sprintf("Internal Error:%s", err),
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminCodeController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err != nil {
		return mvc.Response{
			Text: fmt.Sprintf("没有这个优惠券, error:%s", err),
		}
	}
	err = c.ServiceCode.Update(&models.LtCode{
		Id:         id,
		SysUpdated: comm.NowUnix(),
		SysStatus:  0,
	}, []string{"sys_status"})
	if err != nil {
		return mvc.Response{
			Text: fmt.Sprintf("Internal Error:%s", err),
		}
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

func (c *AdminCodeController) GetRecache() {
	refer := c.Ctx.GetHeader("referer")
	if refer == "" {
		refer = "/admin/code"
	}
	id, err := c.Ctx.URLParamInt("id")
	if id < 1 || err != nil {
		rs := fmt.Sprintf("没有指定优惠券所属的奖品ID, <a href='' onclick='history.go(-1); return false'>返回</a>")
		c.Ctx.HTML(rs)
		return
	}
	sucNum, errNum := utils.ReCacheCodes(id, c.ServiceCode)
	rs := fmt.Sprintf(
		"sucNum=%d, errNum=%d, <a href='' onclick='history.go(-1); return false'>返回</a>", sucNum, errNum)
	c.Ctx.HTML(rs)
}
