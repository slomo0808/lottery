package controllers

import (
	"lottery/comm"
	"lottery/models"
)

func (c *IndexController) checkBlackUser(uid int) (bool, *models.LtUser) {
	info := c.ServiceUser.Get(uid)
	if info != nil && info.Blacktime > comm.NowUnix() {
		// 黑名单期内
		return false, info
	}
	return true, info
}
