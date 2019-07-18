package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"math"
	"time"
)

const iPFrameSize = 2

func init() {
	resetGroupIPList()
}

func resetGroupIPList() {
	log.Println("ip_day_luck.resetGroupIPList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < iPFrameSize; i++ {
		key := fmt.Sprintf("day_ip_num_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("ip_day_luck.resetGroupIPList finished")

	// IP当天的统计数，0点时归零
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIPList)
}

func IncrIpLuckyNum(ipStr string) int64 {
	ip := comm.Ip4ToInt(ipStr)
	frame := ip % iPFrameSize
	key := fmt.Sprintf("day_ip_num_%d", frame)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1)
	if err != nil {
		log.Println("ip_day_lucky.IncreIpLuckyNum.Do error:", err)
		return math.MaxInt32
	}
	return rs.(int64)
}
