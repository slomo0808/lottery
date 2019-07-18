package controllers

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
)

// localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""

	// 1.验证登录
	loginUser := comm.GetLoginUser(c.Ctx.Request())
	if loginUser == nil || loginUser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录"
		return rs
	}

	// 2.用户抽奖分布式锁
	ok := utils.LockLucky(loginUser.Uid)
	if ok {
		defer utils.UnlockLucky(loginUser.Uid)
	} else {
		rs["code"] = 102
		rs["msg"] = "正在抽奖，请稍后重试"
	}

	// 3.验证用户今日参与次数
	userDayNum := utils.IncrUserLuckyNum(loginUser.Uid)
	if userDayNum > conf.UserLimitMax {
		rs["code"] = 103
		rs["msg"] = "今日的抽奖次数已经用完，明天再来吧！"
		return rs
	} else {
		ok = c.CheckUserday(loginUser.Uid, userDayNum)
		if !ok {
			rs["code"] = 103
			rs["msg"] = "今日的抽奖次数已经用完，明天再来吧！"
			return rs
		}
	}

	// 4.验证Ip今日参与次数
	ip := comm.ClientIp(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "相同IP参与次数过多，明天再来吧！"
		return rs
	}

	limitBlack := false // 黑名单
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	// 5.验证Ip黑名单
	var blackIPInfo *models.LtBlackip
	if !limitBlack {
		ok, blackIPInfo = c.checkBlackip(ip)
		if !ok {
			fmt.Println("黑名单中的IP", ip, limitBlack)
			limitBlack = true
		}
	}
	// 6.验证用户黑名单
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginUser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginUser.Uid, limitBlack)
		}
	}

	// 7.获得0~9999的抽奖编码
	luckyNum := comm.Random(10000)

	// 8.匹配奖品
	prizeGift := c.prize(luckyNum, limitBlack)
	if prizeGift == nil || prizeGift.PrizeNum < 0 ||
		(prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		// 证明没有获奖或者奖品已经发放完
		rs["code"] = 205
		rs["msg"] = "很遗憾，没有中奖"
		return rs
	}

	// 9.有限制类的奖品发放
	if prizeGift.PrizeNum > 0 { // 有限量的奖品
		// 基于奖品池的验证
		if utils.GetGiftPoolNum(prizeGift.Id) <= 0 {
			rs["code"] = 206
			rs["msg"] = "很遗憾，没有中奖, 请下次再试"
			return rs
		}
		// 是否发奖成功
		ok = utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum)
		if !ok {
			rs["code"] = 207
			rs["msg"] = "很遗憾，没有中奖, 请下次再试"
			return rs
		}
	}
	// 10.不同编码优惠券的发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode)
		if code == "" {
			rs["code"] = 208
			rs["msg"] = "很遗憾，没有中奖，请下次再试"
			return rs
		}
		prizeGift.Gdata = code
	}
	// 11.记录中奖记录
	result := &models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginUser.Uid,
		Username:   loginUser.Username,
		PrizeCode:  luckyNum,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(result)
	if err != nil {
		log.Println("index_lucky GetLucky c.ServiceResult.Create result:", result, ",error:", err)
		rs["code"] = 209
		rs["msg"] = "很遗憾，没有中奖，请下次再试"
		return rs
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果获得了实物大奖，需要将用户、IP设置黑名单一段时间
		c.prizeLarge(ip, loginUser, userInfo, blackIPInfo)
	}
	// 12.返回抽奖结果
	rs["gift"] = prizeGift
	return rs
}
