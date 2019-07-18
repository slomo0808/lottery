package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/datasource"
	"lottery/models"
	"lottery/services"
	"time"
)

func GetGiftPoolNum(id int) int {
	key := "gift_pool"
	rds := datasource.InstanceCache()
	reply, err := rds.Do("HGET", key, id)
	if err != nil {
		log.Println("prizedata.PrizeGift redis HGET error:", err, ", reply:", reply)
		return 0
	} else {
		num := comm.GetInt64(reply, 0)
		return int(num)
	}
}

func prizeServGift(id int) bool {
	key := "gift_pool"
	rds := datasource.InstanceCache()
	reply, err := rds.Do("HINCRBY", key, id, -1)
	if err != nil {
		log.Println("prizedata.PrizeGift redis HINCRBY error:", err, ", reply:", reply)
		return false
	}
	num := comm.GetInt64(reply, -1)
	return num >= 0
}

func PrizeGift(id, leftNum int) bool {
	ok := false
	ok = prizeServGift(id)
	if ok {
		giftService := services.NewGiftService()
		rows, err := giftService.DecrLeftNum(id, 1)
		if rows < 1 || err != nil {
			log.Println("prizedata.PrizeGift giftService.DecrLeftNum error:", err, ", rows:", rows)
			return false
		}
	}
	return ok
}

func PrizeCodeDiff(id int, codeService services.CodeService) string {
	return prizeServCodeDiff(id, codeService)

}

func prizeLocalCodeGift(id int, codeService services.CodeService) string {
	lockUid := 0 - id - 100000000
	LockLucky(lockUid)
	defer UnlockLucky(lockUid)

	codeId := 0
	codeInfo := codeService.NextUsingCode(id, codeId)
	if codeInfo.Id < 1 {
		log.Println("prizedata.prizeLocalCodeGift num codeInfo, gift_id = ", codeInfo.Id)
		return ""
	}
	if codeInfo == nil {
		log.Println("prizedata.prizeLocalCodeGift num codeInfo is nil")
		return ""
	}
	codeInfo.SysStatus = 2
	codeInfo.SysUpdated = comm.NowUnix()
	codeService.Update(codeInfo, nil)

	return codeInfo.Code
}

func ImportCacheCodes(id int, code string) bool {
	key := fmt.Sprintf("gift_code_%d", id)
	rds := datasource.InstanceCache()
	_, err := rds.Do("SADD", key, code)
	if err != nil {
		log.Println("prizedata.ImportCacheCodes redis SADD error:", err)
		return false
	}
	return true
}

func ReCacheCodes(id int, codeService services.CodeService) (sucNum, errNum int) {
	list := codeService.GetByGift(id)
	if list == nil || len(list) == 0 {
		return
	}
	key := fmt.Sprintf("gift_code_%d", id)
	tmpKey := "temp_" + key
	rds := datasource.InstanceCache()
	for _, data := range list {
		if data.SysStatus == 0 {
			_, err := rds.Do("SADD", tmpKey, data.Code)
			if err != nil {
				errNum++
				continue
			} else {
				sucNum++
			}
		}
	}
	_, err := rds.Do("RENAME", tmpKey, key)
	if err != nil {
		log.Println("prizedata.ImportCacheCodes redis RENAME error:", err)
		return 0, 0
	}
	return
}

func GetCacheCodeNum(id int, codeService services.CodeService) (int, int) {
	var num, cacheNum int
	list := codeService.GetByGift(id)
	if len(list) > 0 {
		for _, data := range list {
			if data.SysStatus == 0 {
				num++
			}
		}
	}

	// redis 中统计
	key := fmt.Sprintf("gift_code_%d", id)
	rds := datasource.InstanceCache()
	reply, err := rds.Do("SCARD", key)
	if err != nil {
		log.Println("prizedata.GetCacheCodeNum redis SCARD error:", err)
	}
	cacheNum = int(comm.GetInt64(reply, 0))
	return num, cacheNum
}

func prizeServCodeDiff(id int, codeService services.CodeService) string {
	key := fmt.Sprintf("gift_code_%d", id)
	rds := datasource.InstanceCache()
	reply, err := rds.Do("SPOP", key)
	if err != nil {
		log.Println("prizedata.GetCacheCodeNum redis SPOP error:", err)
		return ""
	}
	code := comm.GetString(reply, "")
	if code == "" {
		log.Println("prizedata.GetCacheCodeNum redis SPOP reply:", reply)
		return ""
	}
	codeService.UpdateByCode(&models.LtCode{Code: code, SysStatus: 2, SysUpdated: comm.NowUnix()}, nil)
	return code
}

func ResetGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	if giftInfo == nil || giftInfo.Id < 1 {
		return
	}
	id := giftInfo.Id
	now := comm.NowUnix()
	if giftInfo.SysStatus == 1 || giftInfo.TimeBegin > now || giftInfo.TimeEnd < now || giftInfo.LeftNum <= 0 ||
		giftInfo.PrizeNum <= 0 {
		if giftInfo.PrizeData != "" {
			clearGiftPrizeData(giftInfo, giftService)
		}
		return
	}
	// 发奖周期
	dayNum := giftInfo.PrizeTime
	if dayNum <= 0 {
		// 如果没有设置发奖周期，则将所有奖品都放入奖品池中
		setGiftPool(id, giftInfo.LeftNum)
		return
	}

	// 重置发奖计划数据
	// 首先清空奖品池
	setGiftPool(id, 0)

	// 实际奖品计划分布运算
	prizeNum := giftInfo.PrizeNum
	avgNum := prizeNum / dayNum
	// 每天可以分配的奖品数量
	dayPrizeNum := make(map[int]int)
	if avgNum >= 1 {
		for day := 0; day < dayNum; day++ {
			dayPrizeNum[day] = avgNum
		}
	}
	// 剩下的随机分配到任意哪天
	prizeNum -= avgNum * dayNum
	for prizeNum > 0 {
		prizeNum--
		day := comm.Random(dayNum)
		_, ok := dayPrizeNum[day]
		if !ok {
			dayPrizeNum[day] = 1
		} else {
			dayPrizeNum[day] += 1
		}
	}

	// 每天的map，每小时的map， 60分钟的数组， 奖品数量
	prizeData := make(map[int]map[int][60]int)
	for day, num := range dayPrizeNum {
		//计算出这一天的发奖计划
		dayPrizeData := getGiftPrizeDataOneDay(num)
		prizeData[day] = dayPrizeData
	}

	// 将周期内每天，每小时，每分钟的数据prizeData格式化为([时间:数量])
	datalist := formatGiftPrizeData(now, dayNum, prizeData)
	bytes, err := json.Marshal(datalist)

	if err != nil {
		log.Println("prizedata.ResetGiftPrizeData json.Marshal error:", err)
	} else {
		info := &models.LtGift{
			Id:         giftInfo.Id,
			LeftNum:    giftInfo.LeftNum,
			PrizeData:  string(bytes),
			PrizeBegin: now,
			PrizeEnd:   now + dayNum*86400,
			SysUpdated: now,
		}
		err := giftService.Update(info, nil)
		if err != nil {
			log.Println("prizedata.ResetGiftPrizeData giftService.Update error:", err)
		}
	}
}

// 清空原发奖计划
func clearGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	info := &models.LtGift{
		Id:        giftInfo.Id,
		PrizeData: "",
	}
	err := giftService.Update(info, []string{"prize_data"})
	if err != nil {
		log.Println("prizedata.clearGiftPrizeData giftService.Update info:", info, ". error:", err)
	}
	setGiftPool(giftInfo.Id, 0)
}

// 设置奖品池的库存数量
func setGiftPool(id, num int) {
	key := "gift_pool"
	cache := datasource.InstanceCache()
	_, err := cache.Do("HSET", key, id, num)
	if err != nil {
		log.Println("prizedata.setGiftPool redis HSET error:", err)
	}
}

// 计算出一天的发奖计划
func getGiftPrizeDataOneDay(num int) map[int][60]int {
	result := make(map[int][60]int)
	hourData := [24]int{}
	if num > 100 {
		hourNum := 0
		for _, h := range conf.PrizeDataRandomDayTIme {
			hourData[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourData[h]
			n := num * (d / 100)
			hourData[h] = n
			hourNum += n
		}
		num -= hourNum
	}
	for num > 0 {
		num--
		hourIndex := comm.Random(100)
		h := conf.PrizeDataRandomDayTIme[hourIndex]
		hourData[h]++
	}
	// 将每个小时n内的奖品数量分配到60分钟
	for h, hnum := range hourData {
		if hnum <= 0 {
			continue
		}
		minuteData := [60]int{}
		if hnum > 60 {
			avgMinute := hnum / 60
			for i := 0; i < 60; i++ {
				minuteData[i] = avgMinute
			}
			hnum -= avgMinute * 60
		}
		for hnum > 0 {
			hnum--
			minuteIndex := comm.Random(60)
			minuteData[minuteIndex]++
		}
		result[h] = minuteData
	}
	return result
}

// 将周期内每天，每小时，每分钟的数据prizeData格式化为([时间:数量])
// map[int]map[int][60]minute ---> [][int,int]
func formatGiftPrizeData(now, dayNum int, prizeData map[int]map[int][60]int) [][2]int {
	result := make([][2]int, 0)
	nowHour := time.Now().Hour()
	// 处理日期的数据
	for dn := 0; dn < dayNum; dn++ {
		dayData, ok := prizeData[dn]
		if !ok {
			continue
		}
		// 计算发奖计划中的 天
		dayTime := now + dn*86400
		// 处理小时的数据
		for hn := 0; hn < 24; hn++ {
			hourData, ok := dayData[(hn+nowHour)%24]
			if !ok {
				continue
			}
			// 计算发奖计划中的 小时
			hourTime := dayTime + hn*3600
			// 处理分钟的数据
			for mn := 0; mn < 60; mn++ {
				num := hourData[mn]
				if num <= 0 {
					continue
				}
				minuteTIme := hourTime + mn*60
				result = append(result, [2]int{minuteTIme, num})
			}
		}
	}
	return result
}

// 填充奖品池
func DistributionGiftPool() int {
	totalNum := 0
	nowTime := comm.NowUnix()
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	if list != nil && len(list) > 0 {
		for _, gift := range list {
			if gift.SysStatus != 0 {
				continue
			}
			if gift.PrizeNum < 1 {
				// 不限量产品
				continue
			}
			if gift.TimeBegin > nowTime || gift.TimeEnd < nowTime {
				continue
			}
			if len(gift.PrizeData) <= 7 {
				continue
			}
			var cronData [][2]int
			err := json.Unmarshal([]byte(gift.PrizeData), &cronData)
			if err != nil {
				log.Println("prizedata.DistributionGiftPool json.Unmarshal error:", err)
			}
			index := 0
			giftNum := 0
			for i, data := range cronData {
				ct := data[0]
				num := data[1]
				if ct <= nowTime {
					giftNum += num
					index = i + 1
				} else {
					break
				}
			}
			// 更新奖品池
			if giftNum > 0 {
				incrGiftPool(gift.Id, giftNum)
				totalNum += giftNum
			}
			// 更新奖品的计划
			if index > 0 {
				if index >= len(cronData)-1 {
					cronData = make([][2]int, 0)
				} else {
					cronData = cronData[index:]
				}
				bytes, err := json.Marshal(cronData)
				if err != nil {
					log.Println(
						"prizedata.DistributionGiftPool json.Marshal cronData:", cronData, ", error:", err)
				}
				cols := []string{"prize_data"}
				err = giftService.Update(&models.LtGift{
					Id:         gift.Id,
					PrizeData:  string(bytes),
					SysUpdated: nowTime}, cols)
				if err != nil {
					log.Println(
						"prizedata.DistributionGiftPool Update error:", err)
				}
			}
		}
		if totalNum > 0 {
			giftService.GetAll(true)
		}
	}
	return totalNum
}

func incrGiftPool(giftId, num int) int {
	key := "gift_pool"
	rds := datasource.InstanceCache()
	resultNum, err := redis.Int64(rds.Do("HINCRBY", key, giftId, num))
	if err != nil {
		log.Println("prizedata.increGiftPool redis HINCRBY error:", err)
		return 0
	}
	if int(resultNum) < num {
		// 递增小于预期值, 补偿y一次
		num2 := num - int(resultNum)
		resultNum, err = redis.Int64(rds.Do("HINCRBY", key, giftId, num2))
		if err != nil {
			log.Println("prizedata.increGiftPool redis HINCRBY error:", err)
			return 0
		}
	}
	return int(resultNum)
}
