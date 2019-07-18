package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/models"
)

type ResultDao struct {
	engine *xorm.Engine
}

func NewResultDao(engine *xorm.Engine) *ResultDao {
	return &ResultDao{
		engine: engine,
	}
}

func (d *ResultDao) Get(id int) *models.LtResult {
	data := &models.LtResult{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *ResultDao) GetAll(page, size int) []models.LtResult {
	offset := (page - 1) * size
	datalist := make([]models.LtResult, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		log.Println("dao.result_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *ResultDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtResult{})
	if err != nil {
		return 0
	}
	return count
}

// 得到最新中奖记录
func (d *ResultDao) GetNewPrize(size int, giftIds []int) []models.LtResult {
	datalist := make([]models.LtResult, 0)
	err := d.engine.
		In("gift_id", giftIds).
		Desc("id").
		Limit(size).
		Find(&datalist)
	if err != nil {
		log.Println("dao.result_dao.GetNewPrize error:", err)
	}
	return datalist
}

func (d *ResultDao) SearchByGift(giftId int) []models.LtResult {
	datalist := make([]models.LtResult, 0)
	err := d.engine.
		Where("gift_id = ?", giftId).
		Desc("id").
		Find(&datalist)
	if err != nil {
		log.Println("dao.result_dao.SearchByGift error:", err)
	}
	return datalist
}

func (d *ResultDao) SearchByUser(uid int) []models.LtResult {
	datalist := make([]models.LtResult, 0)
	err := d.engine.
		Where("uid = ?", uid).
		Desc("id").
		Find(&datalist)
	if err != nil {
		log.Println("dao.result_dao.SearchByUser error:", err)
	}
	return datalist
}

func (d *ResultDao) SearchByGiftAndUser(uid, giftId int) []models.LtResult {
	datalist := make([]models.LtResult, 0)
	err := d.engine.
		Where("uid = ? and gift_id = ?", uid, giftId).
		Find(&datalist)
	if err != nil {
		log.Println("dao.result_dao.SearchByGiftAndUser error:", err)
	}
	return datalist
}

func (d *ResultDao) CountByGift(giftId int) int64 {
	count, err := d.engine.
		Where("gift_id = ?", giftId).
		Count(&models.LtResult{})
	if err != nil {
		return 0
	}
	return count
}

func (d *ResultDao) CountByUser(uid int) int64 {
	count, err := d.engine.
		Where("uid = ?", uid).
		Count(&models.LtResult{})
	if err != nil {
		return 0
	}
	return count
}

func (d *ResultDao) CountByGiftAndUser(uid, giftId int) int64 {
	count, err := d.engine.
		Where("uid = ? and gift_id = ?", uid, giftId).
		Count(&models.LtResult{})
	if err != nil {
		return 0
	}
	return count
}

func (d *ResultDao) Delete(id int) error {
	data := &models.LtResult{
		Id:        id,
		SysStatus: 1,
	}
	_, err := d.engine.ID(data.Id).Update(data)
	return err
}

func (d *ResultDao) Update(gift *models.LtResult, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *ResultDao) Create(gift *models.LtResult) error {
	_, err := d.engine.Insert(gift)
	return err
}
