package services

import (
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

type CodeService interface {
	Get(id int) *models.LtCode
	GetAll(page, size int) []models.LtCode
	CountAll() int64
	CountByGift(giftId int) int64
	GetByGift(giftId int) []models.LtCode
	NextUsingCode(giftId, codeId int) *models.LtCode
	UpdateByCode(data *models.LtCode, cols []string) error
	Delete(id int) error
	Update(code *models.LtCode, cols []string) error
	Create(code *models.LtCode) error
}

type codeService struct {
	dao *dao.CodeDao
}

func NewCodeService() CodeService {
	return &codeService{dao.NewCodeDao(datasource.InstanceDBMaster())}
}

func (s *codeService) Get(id int) *models.LtCode {
	return s.dao.Get(id)
}

func (s *codeService) GetAll(page, size int) []models.LtCode {
	return s.dao.GetAll(page, size)
}

func (s *codeService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *codeService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *codeService) Update(code *models.LtCode, cols []string) error {
	return s.dao.Update(code, cols)
}

func (s *codeService) Create(code *models.LtCode) error {
	return s.dao.Create(code)
}

func (s *codeService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *codeService) GetByGift(giftId int) []models.LtCode {
	return s.dao.GetByGift(giftId)
}

func (s *codeService) NextUsingCode(giftId, codeId int) *models.LtCode {
	return s.dao.NextUsingCode(giftId, codeId)
}

func (s *codeService) UpdateByCode(data *models.LtCode, cols []string) error {
	return s.dao.UpdateByCode(data, cols)
}
