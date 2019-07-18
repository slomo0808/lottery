package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/models"
)

type CodeDao struct {
	engine *xorm.Engine
}

func NewCodeDao(engine *xorm.Engine) *CodeDao {
	return &CodeDao{
		engine: engine,
	}
}

func (d *CodeDao) Get(id int) *models.LtCode {
	data := &models.LtCode{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *CodeDao) GetAll(page, size int) []models.LtCode {
	datalist := make([]models.LtCode, 0)
	offset := (page - 1) * size
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		log.Println("dao.code_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *CodeDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtCode{})
	if err != nil {
		return 0
	}
	return count
}

func (d *CodeDao) CountByGift(giftId int) int64 {
	count, err := d.engine.
		Where("gift_id = ?", giftId).
		Count(&models.LtCode{})
	if err == nil {
		return count
	} else {
		return 0
	}
}

func (d *CodeDao) GetByGift(giftId int) []models.LtCode {
	datalist := make([]models.LtCode, 0)
	err := d.engine.
		Where("gift_id = ?", giftId).
		Find(&datalist)
	if err != nil {
		log.Println("dao.code_dao.GetByGift error:", err)
	}
	return datalist
}

func (d *CodeDao) Delete(id int) error {
	data := &models.LtCode{
		Id:        id,
		SysStatus: 1,
	}
	_, err := d.engine.ID(data.Id).Update(data)
	return err
}

func (d *CodeDao) Update(gift *models.LtCode, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *CodeDao) Create(gift *models.LtCode) error {
	_, err := d.engine.Insert(gift)
	return err
}

// 找到下一个可用的最小的优惠券
func (d *CodeDao) NextUsingCode(giftId, codeId int) *models.LtCode {
	datalist := make([]models.LtCode, 0)
	err := d.engine.
		Where("gift_id = ?", giftId).
		Where("sys_status = ?", 0).
		Where("id > ?", codeId).
		Find(&datalist)
	if err != nil || len(datalist) == 0 {
		log.Println("dao.code_dao.NextUsingCode len(datalist)=", len(datalist), "error:", err)
		return &models.LtCode{Id: 0}
	}
	return &datalist[0]
}

// 根据唯一的code来更新
func (d *CodeDao) UpdateByCode(data *models.LtCode, cols []string) error {
	_, err := d.engine.Where("code = ?", data.Code).MustCols(cols...).Update(data)
	return err
}
