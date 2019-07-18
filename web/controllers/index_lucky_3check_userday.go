package controllers

import (
	"fmt"
	"log"
	"lottery/comm"
	"lottery/conf"
	"lottery/models"
	"lottery/web/utils"
	"strconv"
	"time"
)

func (c *IndexController) CheckUserday(uid int, userDayNum int64) bool {
	userdayInfo := c.ServiceUserday.GetUserToday(uid)
	if userdayInfo != nil && userdayInfo.Uid == uid {
		// 存在抽奖记录
		if userdayInfo.Num > conf.UserLimitMax {
			if int(userDayNum) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		} else {
			userdayInfo.Num++
			if int(userDayNum) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			if err103 := c.ServiceUserday.Update(userdayInfo, nil); err103 != nil {
				log.Println("index_lucky_3check_userday ServiceUserday.Update error103:", err103)
			}
		}
	} else {
		// 创建今天的抽奖记录
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		data := &models.LtUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: comm.NowUnix(),
		}
		if err103 := c.ServiceUserday.Create(data); err103 != nil {
			log.Println("index_lucky_3check_userday ServiceUserday.Create error103:", err103)
		}

		utils.InitUserLuckyNum(uid, 1)

	}
	return true
}
