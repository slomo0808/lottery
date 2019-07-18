package services

import (
	"encoding/json"
	"log"
	"lottery/comm"
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
	"strconv"
	"strings"
)

type GiftService interface {
	Get(id int, useCache bool) *models.LtGift
	GetAll(bool) []models.LtGift
	GetAllUse(useCache bool) []models.ObjGiftPrize
	CountAll() int64
	Delete(id int) error
	Update(gift *models.LtGift, cols []string) error
	Create(gift *models.LtGift) error
	DecrLeftNum(id, num int) (int64, error)
	IncrLeftNum(id, num int) (int64, error)
}

type giftService struct {
	dao *dao.GiftDao
}

func NewGiftService() GiftService {
	return &giftService{dao.NewGiftDao(datasource.InstanceDBMaster())}
}

func (s *giftService) Get(id int, useCache bool) *models.LtGift {
	if !useCache {
		return s.dao.Get(id)
	}
	gifts := s.GetAll(true)
	for _, gift := range gifts {
		if gift.Id == id {
			return &gift
		}
	}
	return nil
}

// 设置是否使用缓存， 如果使用缓存
// 先从缓存内读取，如果读不到，再从数据库读取，同时设置缓存的值
func (s *giftService) GetAll(useCache bool) []models.LtGift {
	if !useCache {
		return s.dao.GetAll()
	}
	gifts := s.getAllByCache()
	if len(gifts) < 1 {
		gifts = s.dao.GetAll()
		s.setAllByCache(gifts)
	}
	return gifts
}

func (s *giftService) CountAll() int64 {
	// return s.dao.CountAll()
	gifts := s.GetAll(true)
	return int64(len(gifts))
}

func (s *giftService) Delete(id int) error {
	// 改变了数据，记得一定要更新缓存
	data := &models.LtGift{Id: id}
	s.updateByCache(data, nil)
	return s.dao.Delete(id)
}

func (s *giftService) Update(gift *models.LtGift, cols []string) error {
	// 改变了数据，记得一定要更新缓存
	s.updateByCache(gift, cols)
	return s.dao.Update(gift, cols)
}

func (s *giftService) Create(gift *models.LtGift) error {
	// 改变了数据，记得一定要更新缓存
	s.updateByCache(gift, nil)
	return s.dao.Create(gift)
}

func (s *giftService) GetAllUse(useCache bool) []models.ObjGiftPrize {
	datalist := make([]models.LtGift, 0)
	if !useCache {
		datalist = s.dao.GetAllUse()
	} else {
		now := comm.NowUnix()
		gifts := s.GetAll(true)
		for _, gift := range gifts {
			if gift.Id >= 1 && gift.SysStatus == 0 &&
				gift.TimeBegin <= now &&
				gift.TimeEnd >= now &&
				gift.PrizeNum >= 0 {
				datalist = append(datalist, gift)
			}
		}
	}

	objDatalist := make([]models.ObjGiftPrize, 0)

	if datalist != nil {
		for _, data := range datalist {
			sArr := strings.Split(data.PrizeCode, "-")
			if len(sArr) == 2 {
				prizeCodeA, err1 := strconv.Atoi(sArr[0])
				prizeCodeB, err2 := strconv.Atoi(sArr[1])
				if err1 == nil && err2 == nil && prizeCodeA <= prizeCodeB && prizeCodeA >= 0 && prizeCodeB <= 10000 {
					objData := models.ObjGiftPrize{
						Id:           data.Id,
						Title:        data.Title,
						PrizeNum:     data.PrizeNum,
						LeftNum:      data.LeftNum,
						PrizeCodeA:   prizeCodeA,
						PrizeCodeB:   prizeCodeB,
						PrizeTime:    data.PrizeTime,
						Img:          data.Img,
						Displayorder: data.Displayorder,
						Gtype:        data.Gtype,
						Gdata:        data.Gdata,
					}
					objDatalist = append(objDatalist, objData)
				}
			}
		}
	}
	return objDatalist
}

func (s *giftService) DecrLeftNum(id, num int) (int64, error) {
	return s.dao.DecrLeftNum(id, num)
}

func (s *giftService) IncrLeftNum(id, num int) (int64, error) {
	return s.dao.IncrLeftNum(id, num)
}

func (s *giftService) getAllByCache() []models.LtGift {
	key := "allGifts"
	rds := datasource.InstanceCache()
	res, err := rds.Do("get", key)
	if err != nil {
		log.Println("gift_service.getAllByCache rds.Do result:", res, ",error:", err)
		return nil
	}
	str := comm.GetString(res, "")
	if str == "" {
		return nil
	}
	datalist := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(str), &datalist); err != nil {
		log.Println("gift_service.getAllByCache.Unmarshal error:", err)
		return nil
	}
	gifts := make([]models.LtGift, len(datalist))
	for i, data := range datalist {
		id := comm.GetInt64FromMap(data, "Id", 0)
		if id < 1 {
			gifts[i] = models.LtGift{}
		} else {
			gift := models.LtGift{
				Id:           int(id),
				Title:        comm.GetStringFromMap(data, "Title", ""),
				PrizeNum:     int(comm.GetInt64FromMap(data, "PrizeNum", 0)),
				LeftNum:      int(comm.GetInt64FromMap(data, "LeftNum", 0)),
				PrizeCode:    comm.GetStringFromMap(data, "PrizeCode", ""),
				PrizeTime:    int(comm.GetInt64FromMap(data, "PrizeTime", 0)),
				Img:          comm.GetStringFromMap(data, "Img", ""),
				Displayorder: int(comm.GetInt64FromMap(data, "Displayorder", 0)),
				Gtype:        int(comm.GetInt64FromMap(data, "Gtype", 0)),
				Gdata:        comm.GetStringFromMap(data, "Gdata", ""),
				TimeBegin:    int(comm.GetInt64FromMap(data, "TimeBegin", 0)),
				TimeEnd:      int(comm.GetInt64FromMap(data, "TimeEnd", 0)),
				//PrizeData:    comm.GetStringFromMap(data, "PrizeData", ""),
				PrizeBegin: int(comm.GetInt64FromMap(data, "PrizeBegin", 0)),
				PrizeEnd:   int(comm.GetInt64FromMap(data, "PrizeEnd", 0)),
				SysStatus:  int(comm.GetInt64FromMap(data, "SysStatus", 0)),
				SysCreated: int(comm.GetInt64FromMap(data, "SysCreated", 0)),
				SysUpdated: int(comm.GetInt64FromMap(data, "SysUpdated", 0)),
				SysIp:      comm.GetStringFromMap(data, "SysIp", ""),
			}
			gifts[i] = gift
		}

	}
	return gifts
}

func (s *giftService) setAllByCache(gifts []models.LtGift) {
	strVal := ""
	if len(gifts) > 0 {
		datalist := make([]map[string]interface{}, len(gifts))
		// 格式转换
		for i, gift := range gifts {
			data := make(map[string]interface{})
			data["Id"] = gift.Id
			data["Title"] = gift.Title
			data["PrizeNum"] = gift.PrizeNum
			data["LeftNum"] = gift.LeftNum
			data["PrizeCode"] = gift.PrizeCode
			data["PrizeTime"] = gift.PrizeTime
			data["Img"] = gift.Img
			data["Displayorder"] = gift.Displayorder
			data["Gtype"] = gift.Gtype
			data["Gdata"] = gift.Gdata
			data["TimeBegin"] = gift.TimeBegin
			data["TimeEnd"] = gift.TimeEnd
			data["PrizeBegin"] = gift.PrizeBegin
			data["PrizeEnd"] = gift.PrizeEnd
			data["SysStatus"] = gift.SysStatus
			data["SysUpdated"] = gift.SysUpdated
			data["SysCreated"] = gift.SysCreated
			data["SysIp"] = gift.SysIp
			datalist[i] = data
		}
		bytes, err := json.Marshal(datalist)
		if err != nil {
			log.Println("gift_service.setAllByCache json.Marshal error:", err)
		}
		strVal = string(bytes)
	}
	key := "allGifts"
	rds := datasource.InstanceCache()
	_, err := rds.Do("set", key, strVal)
	if err != nil {
		log.Println("gift_service.setAllByCache redis set key:", key, ", error:", err)
	}
}

func (s *giftService) updateByCache(data *models.LtGift, cols []string) {
	if data == nil || data.Id < 1 {
		return
	}
	key := "allGifts"
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
