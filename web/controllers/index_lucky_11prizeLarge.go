package controllers

import (
	"lottery/comm"
	"lottery/models"
)

func (c *IndexController) prizeLarge(ip string,
	loginUser *models.ObjLoginUser,
	userInfo *models.LtUser,
	blackIpInfo *models.LtBlackip) {

	nowTime := comm.NowUnix()
	blackTime := 30 * 86400
	// 更新用户黑名单信息
	if userInfo == nil || userInfo.Id < 1 {
		userInfo = &models.LtUser{
			Username:   loginUser.Username,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
			SysIp:      ip,
		}
		c.ServiceUser.Create(userInfo)
	} else {
		userInfo = &models.LtUser{
			Id:         loginUser.Uid,
			Blacktime:  nowTime + blackTime,
			SysUpdated: nowTime,
		}
		c.ServiceUser.Update(userInfo, []string{"blacktime"})
	}
	// 更新ip黑名单信息
	if blackIpInfo == nil || blackIpInfo.Id < 1 {
		blackIpInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
		}
		c.ServiceBlackip.Create(blackIpInfo)
	} else {
		blackIpInfo.SysUpdated = nowTime
		blackIpInfo.Blacktime = nowTime + blackTime
		c.ServiceBlackip.Update(blackIpInfo, []string{"blacktime"})
	}
}
