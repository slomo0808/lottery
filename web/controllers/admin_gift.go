package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"log"
	"lottery/comm"
	"lottery/models"
	"lottery/services"
	"lottery/web/utils"
	"lottery/web/viewmodels"
)

type AdminGiftController struct {
	Ctx         iris.Context
	ServiceGift services.GiftService
}

func (c *AdminGiftController) Get() mvc.View {
	datalist := c.ServiceGift.GetAll(false)
	total := len(datalist)
	for i, data := range datalist {
		prizedata := make([][2]int, 0)
		err := json.Unmarshal([]byte(data.PrizeData), &prizedata)
		if err != nil || len(prizedata) < 1 {
			datalist[i].PrizeData = "[]"
		} else {
			newpd := make([]string, len(prizedata))
			for i, pd := range prizedata {
				ct := comm.FormatFromUnixTime(int64(pd[0]))
				newpd[i] = fmt.Sprintf("【%s】:%d", ct, pd[1])
			}
			str, err := json.Marshal(newpd)
			if err == nil && len(str) > 0 {
				datalist[i].PrizeData = string(str)
			} else {
				datalist[i].PrizeData = "[]"
			}
		}
		num := utils.GetGiftPoolNum(data.Id)
		datalist[i].Title = fmt.Sprintf("【%d】%s", num, datalist[i].Title)
	}
	return mvc.View{
		Name: "admin/gift.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "gift",
			"Datalist": datalist,
			"Total":    total,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) GetEdit() mvc.View {
	id := c.Ctx.URLParamIntDefault("id", 0)
	giftInfo := viewmodels.ViewGift{}
	if id > 0 {
		data := c.ServiceGift.Get(id, true)
		giftInfo.Id = data.Id
		giftInfo.Title = data.Title
		giftInfo.PrizeNum = data.PrizeNum
		giftInfo.PrizeCode = data.PrizeCode
		giftInfo.PrizeTime = data.PrizeTime
		giftInfo.Img = data.Img
		giftInfo.Displayorder = data.Displayorder
		giftInfo.Gtype = data.Gtype
		giftInfo.Gdata = data.Gdata
		giftInfo.TimeBegin = comm.FormatFromUnixTime(int64(data.TimeBegin))
		giftInfo.TimeEnd = comm.FormatFromUnixTime(int64(data.TimeEnd))
	}
	return mvc.View{
		Name: "admin/giftEdit.html",
		Data: iris.Map{
			"Title":   "管理后台",
			"Channel": "gift",
			"info":    giftInfo,
		},
		Layout: "admin/layout.html",
	}
}

func (c *AdminGiftController) PostSave() mvc.Result {
	data := viewmodels.ViewGift{}
	err := c.Ctx.ReadForm(&data)
	if err != nil {
		log.Println("controller admin_gift.PostSave.ReadForm error:", err)
		return mvc.Response{
			Text: fmt.Sprintf("ReadForm转换异常:error = %s", err),
		}
	}
	giftInfo := models.LtGift{}
	giftInfo.Id = data.Id
	giftInfo.Title = data.Title
	giftInfo.PrizeNum = data.PrizeNum
	giftInfo.PrizeCode = data.PrizeCode
	giftInfo.PrizeTime = data.PrizeTime
	giftInfo.Img = data.Img
	giftInfo.Displayorder = data.Displayorder
	giftInfo.Gtype = data.Gtype
	giftInfo.Gdata = data.Gdata
	t1, err1 := comm.ParseTime(data.TimeBegin)
	t2, err2 := comm.ParseTime(data.TimeEnd)
	if err1 != nil || err2 != nil {
		return mvc.Response{
			Text: fmt.Sprintf("开始时间、结束时间格式不正确, err1=%s, err2=%s", err1, err2),
		}
	}
	giftInfo.TimeBegin = int(t1.Unix())
	giftInfo.TimeEnd = int(t2.Unix())
	if giftInfo.Id > 0 {
		// 数据更新
		dataBefore := c.ServiceGift.Get(giftInfo.Id, false)
		if dataBefore.Id > 0 {
			if dataBefore.PrizeNum != giftInfo.PrizeNum {
				giftInfo.LeftNum = (giftInfo.PrizeNum - dataBefore.PrizeNum) + dataBefore.LeftNum
				if giftInfo.LeftNum < 0 || giftInfo.PrizeNum <= 0 {
					giftInfo.LeftNum = 0
				}
				// 总数发生变化, 重新计算发奖计划
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			if giftInfo.PrizeTime != dataBefore.PrizeTime {
				log.Println("PrizeTime changed ResetGiftPrizeData")
				// 抽奖周期发生变化， 重新计算发奖计划
				utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
			}
			giftInfo.SysUpdated = comm.NowUnix()
			c.ServiceGift.Update(&giftInfo, nil)
		} else {
			giftInfo.Id = 0
		}
	}
	if giftInfo.Id == 0 {
		giftInfo.LeftNum = giftInfo.PrizeNum
		giftInfo.SysIp = comm.ClientIp(c.Ctx.Request())
		giftInfo.SysCreated = comm.NowUnix()
		giftInfo.SysUpdated = comm.NowUnix()
		c.ServiceGift.Create(&giftInfo)
		// 新的奖品，新的发奖计划
		utils.ResetGiftPrizeData(&giftInfo, c.ServiceGift)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Delete(id)
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}

func (c *AdminGiftController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceGift.Update(&models.LtGift{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	return mvc.Response{
		Path: "/admin/gift",
	}
}
