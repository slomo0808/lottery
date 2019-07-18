package services

import (
	"fmt"
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
	"strconv"
	"time"
)

type UserdayService interface {
	Get(id int) *models.LtUserday
	GetAll() []models.LtUserday
	CountAll() int64
	Update(Userday *models.LtUserday, cols []string) error
	Create(Userday *models.LtUserday) error
	Search(uid, day int) []models.LtUserday
	Count(uid, day int) int
	GetUserToday(uid int) *models.LtUserday
}

type userdayService struct {
	dao *dao.UserdayDao
}

func NewUserdayService() UserdayService {
	return &userdayService{dao.NewUserdayDao(datasource.InstanceDBMaster())}
}

func (s *userdayService) Get(id int) *models.LtUserday {
	return s.dao.Get(id)
}

func (s *userdayService) GetAll() []models.LtUserday {
	return s.dao.GetAll()
}

func (s *userdayService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userdayService) Update(Userday *models.LtUserday, cols []string) error {
	return s.dao.Update(Userday, cols)
}

func (s *userdayService) Create(Userday *models.LtUserday) error {
	return s.dao.Create(Userday)
}

func (s *userdayService) Search(uid, day int) []models.LtUserday {
	return s.dao.Search(uid, day)
}
func (s *userdayService) Count(uid, day int) int {
	return s.dao.Count(uid, day)
}

func (s *userdayService) GetUserToday(uid int) *models.LtUserday {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	list := s.dao.Search(uid, day)
	if list != nil && len(list) > 0 {
		return &list[0]
	} else {
		return nil
	}
}
