package services

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"lottery/comm"
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

type BlackipService interface {
	Get(id int) *models.LtBlackip
	GetByIp(ip string) *models.LtBlackip
	Search(ip string) []models.LtBlackip
	GetAll(int, int) []models.LtBlackip
	CountAll() int64
	Update(blackip *models.LtBlackip, cols []string) error
	Create(blackip *models.LtBlackip) error
}

type blackipService struct {
	dao *dao.BlackipDao
}

func NewBlackipService() BlackipService {
	return &blackipService{dao.NewBlackipDao(datasource.InstanceDBMaster())}
}

func (s *blackipService) Get(id int) *models.LtBlackip {
	return s.dao.Get(id)
}

func (s *blackipService) GetAll(page, size int) []models.LtBlackip {
	return s.dao.GetAll(page, size)
}

func (s *blackipService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *blackipService) Update(data *models.LtBlackip, cols []string) error {
	s.updateByCache(data, cols)
	return s.dao.Update(data, cols)
}

func (s *blackipService) Create(data *models.LtBlackip) error {
	return s.dao.Create(data)
}

func (s *blackipService) GetByIp(ip string) *models.LtBlackip {
	data := s.getByCache(ip)
	if data == nil || data.Ip == "" {
		data = s.dao.GetByIp(ip)
		if data == nil || data.Ip == "" {
			data = &models.LtBlackip{Ip: ip}
		}
	}
	return data
}

func (s *blackipService) Search(ip string) []models.LtBlackip {
	return s.dao.Search(ip)
}

func (s *blackipService) getByCache(ip string) *models.LtBlackip {
	key := fmt.Sprintf("info_black_%s", ip)
	rds := datasource.InstanceCache()
	reply, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("blackip_service getByCache redis HGETALL reply:", reply, ", error:", err)
		return nil
	}
	ip = comm.GetStringFromStringMap(reply, "Ip", "")
	if ip == "" {
		return nil
	}
	return &models.LtBlackip{
		Id:         int(comm.GetInt64FromStringMap(reply, "Id", 0)),
		Ip:         ip,
		Blacktime:  int(comm.GetInt64FromStringMap(reply, "Blacktime", 0)),
		SysCreated: int(comm.GetInt64FromStringMap(reply, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(reply, "SysUpdated", 0)),
	}
}

func (s *blackipService) setByCache(data *models.LtBlackip) {
	if data == nil || data.Ip == "" {
		return
	}
	key := fmt.Sprintf("info_black_%s", data.Ip)
	rds := datasource.InstanceCache()
	params := redis.Args{key}
	params.Add("Id", data.Id, "Ip", data.Ip, "Blacktime",
		data.Blacktime, "SysCreated", data.SysCreated, "SysUpdated", data.SysUpdated)
	_, err := rds.Do("HMSET", params)
	if err != nil {
		log.Println("blackip_service getByCache redis HMSET params:", params, ", error:", err)
	}
}

func (s *blackipService) updateByCache(data *models.LtBlackip, cols []string) {
	if data == nil || data.Ip == "" {
		return
	}
	key := fmt.Sprintf("info_black_%s", data.Ip)
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
