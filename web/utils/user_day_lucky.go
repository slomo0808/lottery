package utils

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/datasource"
	"math"
	"time"
)

const userFrameSize = 2

func init() {
	resetGroupUserList()
}

func resetGroupUserList() {
	log.Println("user_day_luck.resetGroupIPList start")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < iPFrameSize; i++ {
		key := fmt.Sprintf("day_user_num_%d", i)
		cacheObj.Do("DEL", key)
	}
	log.Println("user_day_luck.resetGroupUserList finished")

	// IP当天的统计数，0点时归零
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList)
}

func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	key := fmt.Sprintf("day_user_num_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1)
	if err != nil {
		log.Println("user_day_lucky.IncrUserLuckyNum redis HINCRBY error:", err)
		return math.MaxInt32
	}
	return rs.(int64)
}

func InitUserLuckyNum(uid int, num int64) {
	if num < 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf("day_user_num_%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, num)
	if err != nil {
		log.Println("user_day_lucky.InitUserLuckyNum redis HSET key:", key, ", nid:", uid, ", error:", err)
	}
}
