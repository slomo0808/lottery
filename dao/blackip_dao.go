package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/models"
)

type BlackipDao struct {
	engine *xorm.Engine
}

func NewBlackipDao(engine *xorm.Engine) *BlackipDao {
	return &BlackipDao{
		engine: engine,
	}
}

func (d *BlackipDao) Get(id int) *models.LtBlackip {
	data := &models.LtBlackip{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *BlackipDao) GetByIp(ip string) *models.LtBlackip {
	data := &models.LtBlackip{}
	ok, err := d.engine.Where("ip = ?", ip).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Ip = ""
		return data
	}
}

func (d *BlackipDao) Search(ip string) []models.LtBlackip {
	datalist := make([]models.LtBlackip, 0)
	err := d.engine.Where("ip = ?", ip).Find(&datalist)
	if err != nil {
		log.Println("dao.blackip_dao.Search error:", err)
		return datalist
	} else {
		return datalist
	}
}

func (d *BlackipDao) GetAll(page, size int) []models.LtBlackip {
	offset := (page - 1) * size
	datalist := make([]models.LtBlackip, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		log.Println("dao.blackip_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *BlackipDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtBlackip{})
	if err != nil {
		return 0
	}
	return count
}

func (d *BlackipDao) Update(gift *models.LtBlackip, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *BlackipDao) Create(gift *models.LtBlackip) error {
	_, err := d.engine.Insert(gift)
	return err
}
