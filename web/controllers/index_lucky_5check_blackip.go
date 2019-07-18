package controllers

import (
	"lottery/comm"
	"lottery/models"
)

func (c *IndexController) checkBlackip(ipStr string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIp(ipStr)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > comm.NowUnix() {
		// IP黑名单存在，并且还在黑名单有效期内
		return false, info
	}
	return true, info
}
