package conf

import "time"

const GtypeVirtual = 0   // 虚拟币
const GtypeCodeSame = 1  // 虚拟券，相同的码
const GtypeCodeDiff = 2  // 虚拟券，不同的码
const GtypeGiftSmall = 3 // 实物小奖
const GtypeGiftLarge = 4 // 实物大奖

// 超过次数直接不允许抽奖
const IpLimitMax = 500

// 超过次数可以抽奖，但不能抽中实物奖品
const IpPrizeMax = 100
const UserLimitMax = 3000

var SignSecret = []byte("0123456789abcdef")

var CookieSecret = "hellokitty"

const SysTimeForm = "2006-01-02 15:04:05"
const SysTimeFormShort = "2006-01-02"

var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")

// 每天发奖计划中，不同时间对应的概率
var PrizeDataRandomDayTIme = [100]int{
	// 24小时， 每小时平均分3%的机会
	// 24 * 3 = 72
	// 剩余的28， 分别分到其他高峰期的7个小时
	// 增加%4的概率
	0, 0, 0,
	1, 1, 1,
	2, 2, 2,
	3, 3, 3,
	4, 4, 4,
	5, 5, 5,
	6, 6, 6,
	7, 7, 7,
	8, 8, 8,
	9, 9, 9, 9, 9, 9, 9,
	10, 10, 10, 10, 10, 10, 10,
	11, 11, 11,
	12, 12, 12,
	13, 13, 13,
	14, 14, 14,
	15, 15, 15, 15, 15, 15, 15,
	16, 16, 16, 16, 16, 16, 16,
	17, 17, 17, 17, 17, 17, 17,
	18, 18, 18,
	19, 19, 19,
	20, 20, 20, 20, 20, 20, 20,
	21, 21, 21, 21, 21, 21, 21,
	22, 22, 22,
	23, 23, 23,
}
