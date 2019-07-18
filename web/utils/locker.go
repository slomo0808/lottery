package utils

import (
	"fmt"
	"lottery/datasource"
)

func getLuckyLockKey(uid int) string {
	return fmt.Sprintf("lucky_lock_%d", uid)
}

func LockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("set", key, 1, "EX", 3, "NX")
	if rs.(string) == "OK" {
		return true
	} else {
		return false
	}
}

func UnlockLucky(uid int) bool {
	key := getLuckyLockKey(uid)
	cacheObj := datasource.InstanceCache()
	rs, _ := cacheObj.Do("DEL", key)
	if rs.(int64) == 0 || rs.(int64) == 1 {
		return true
	} else {
		return false
	}
}
