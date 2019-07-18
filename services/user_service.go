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

type UserService interface {
	Get(id int) *models.LtUser
	GetAll(int, int) []models.LtUser
	CountAll() int64
	Update(User *models.LtUser, cols []string) error
	Create(User *models.LtUser) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{dao.NewUserDao(datasource.InstanceDBMaster())}
}

func (s *userService) Get(id int) *models.LtUser {
	data := s.getByCache(id)
	if data == nil || data.Id < 1 {
		data = s.dao.Get(id)
		if data == nil || data.Id < 1 {
			data = &models.LtUser{Id: id}
		}
		s.setByCache(data)
	}
	return data
}

func (s *userService) GetAll(page, size int) []models.LtUser {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userService) Update(data *models.LtUser, cols []string) error {
	s.updateByCache(data, cols)
	return s.dao.Update(data, cols)
}

func (s *userService) Create(data *models.LtUser) error {
	return s.dao.Create(data)
}

func (s *userService) getByCache(id int) *models.LtUser {
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	// 返回 map[string][string], error
	reply, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service getByCache redis HGETALL reply:", reply, ", error:", err)
		return nil
	}
	dataId := comm.GetInt64FromStringMap(reply, "Id", 0)
	if dataId < 1 {
		return nil
	}
	return &models.LtUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(reply, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(reply, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(reply, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(reply, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(reply, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(reply, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(reply, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(reply, "SysIp", ""),
	}
}

func (s *userService) setByCache(data *models.LtUser) {
	if data == nil || data.Id < 1 {
		return
	}
	id := data.Id
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	params := redis.Args{key}
	params = params.Add("Id", id)
	if data.Username != "" {
		params = params.Add("Username", data.Username)
		params = params.Add("Blacktime", data.Blacktime)
		params = params.Add("Realname", data.Realname)
		params = params.Add("Mobile", data.Mobile)
		params = params.Add("Address", data.Address)
		params = params.Add("SysCreated", data.SysCreated)
		params = params.Add("SysUpdated", data.SysUpdated)
		params = params.Add("SysIp", data.SysIp)
	}
	_, err := rds.Do("HMSET", key, params)
	if err != nil {
		log.Println("user_service getByCache redis HMSET params:", params, ", error:", err)
	}
}

func (s *userService) updateByCache(data *models.LtUser, cols []string) {
	if data == nil || data.Id < 1 {
		return
	}
	key := fmt.Sprintf("info_user_%d", data.Id)
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
