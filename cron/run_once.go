package cron

import (
	"log"
	"lottery/comm"
	"lottery/services"
	"lottery/web/utils"
	"time"
)

func ConfigureAppOneCron() {
	go resetAllGiftPrizeData()
	go distributiontAllGiftPool()
}

func resetAllGiftPrizeData() {
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	nowTIme := comm.NowUnix()
	for _, gift := range list {
		if gift.PrizeTime > 0 && (gift.PrizeData == "" || gift.PrizeEnd < nowTIme) {
			log.Println("crontab start utils.ResetGiftPrizeData giftInfo=", gift)
			utils.ResetGiftPrizeData(&gift, giftService)
			giftService.GetAll(true)
			log.Println("crontab end utils.ResetGiftPrizeData giftInfo=", gift)
		}
	}
	// 没5分钟执行一次
	time.AfterFunc(5*time.Minute, resetAllGiftPrizeData)
}

func distributiontAllGiftPool() {
	log.Println("crontab start utils.DistributionGiftPool")
	num := utils.DistributionGiftPool()
	log.Println("crontab end utils.DistributionGiftPool, num = ", num)

	// 每分钟执行一次
	time.AfterFunc(time.Minute, distributiontAllGiftPool)
}
