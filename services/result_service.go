package services

import (
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

type ResultService interface {
	Get(id int) *models.LtResult
	GetAll(int, int) []models.LtResult
	CountAll() int64
	GetNewPrize(size int, giftIds []int) []models.LtResult
	SearchByGift(giftId int) []models.LtResult
	SearchByUser(uid int) []models.LtResult
	SearchByGiftAndUser(uid, giftId int) []models.LtResult
	CountByGift(giftId int) int64
	CountByUser(uid int) int64
	CountByGiftAndUser(uid, giftId int) int64
	Delete(id int) error
	Update(Result *models.LtResult, cols []string) error
	Create(Result *models.LtResult) error
}

type resultService struct {
	dao *dao.ResultDao
}

func NewResultService() ResultService {
	return &resultService{dao.NewResultDao(datasource.InstanceDBMaster())}
}

func (s *resultService) Get(id int) *models.LtResult {
	return s.dao.Get(id)
}

func (s *resultService) GetAll(page, size int) []models.LtResult {
	return s.dao.GetAll(page, size)
}

func (s *resultService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *resultService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *resultService) Update(Result *models.LtResult, cols []string) error {
	return s.dao.Update(Result, cols)
}

func (s *resultService) Create(Result *models.LtResult) error {
	return s.dao.Create(Result)
}

func (s *resultService) GetNewPrize(size int, giftIds []int) []models.LtResult {
	return s.dao.GetNewPrize(size, giftIds)
}
func (s *resultService) SearchByGift(giftId int) []models.LtResult {
	return s.dao.SearchByGift(giftId)
}
func (s *resultService) SearchByUser(uid int) []models.LtResult {
	return s.dao.SearchByUser(uid)
}
func (s *resultService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}
func (s *resultService) CountByUser(uid int) int64 {
	return s.dao.CountByUser(uid)
}

func (s *resultService) SearchByGiftAndUser(uid, giftId int) []models.LtResult {
	return s.dao.SearchByGiftAndUser(uid, giftId)
}

func (s *resultService) CountByGiftAndUser(uid, giftId int) int64 {
	return s.dao.CountByGiftAndUser(uid, giftId)
}
