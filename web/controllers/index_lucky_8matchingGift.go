package controllers

import (
	"lottery/conf"
	"lottery/models"
)

func (c *IndexController) prize(luckyNum int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := c.ServiceGift.GetAllUse(true)
	for _, gift := range giftList {
		if luckyNum >= gift.PrizeCodeA && luckyNum <= gift.PrizeCodeB {
			// 满足中奖编码区间条件，可以中奖了
			// 验证是否黑名单用户
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall {
				// 如果不是黑名单用户  或者 中的是虚拟奖品，则中奖
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}
