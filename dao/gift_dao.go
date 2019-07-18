package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/comm"
	"lottery/models"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao {
	return &GiftDao{
		engine: engine,
	}
}

func (d *GiftDao) Get(id int) *models.LtGift {
	data := &models.LtGift{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *GiftDao) GetAll() []models.LtGift {
	datalist := make([]models.LtGift, 0)
	err := d.engine.
		Asc("sys_status").
		Asc("displayorder").
		Find(&datalist)
	if err != nil {
		log.Println("dao.gift_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *GiftDao) GetAllUse() []models.LtGift {
	datalist := make([]models.LtGift, 0)
	now := comm.NowUnix()
	err := d.engine.
		Cols(
			"id", "title", "prize_num", "left_num",
			"prize_code", "img", "displayorder", "gtype", "gdata").
		Desc("gtype").
		Asc("displayorder").
		Where("prize_num>=?", 0).
		Where("sys_status=?", 0).
		Where("time_begin<=?", now).
		Where("time_end>=?", now).
		Find(&datalist)
	if err != nil {
		log.Println("gift_dao.GetAllUse error:", err)
	}

	return datalist
}

func (d *GiftDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtGift{})
	if err != nil {
		return 0
	}
	return count
}

func (d *GiftDao) Delete(id int) error {
	data := &models.LtGift{
		Id:        id,
		SysStatus: 1,
	}
	_, err := d.engine.ID(data.Id).Update(data)
	return err
}

func (d *GiftDao) Update(gift *models.LtGift, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *GiftDao) Create(gift *models.LtGift) error {
	_, err := d.engine.Insert(gift)
	return err
}

func (d *GiftDao) DecrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.ID(id).
		Decr("left_num", num).
		Where("left_num >= ?", num). // 乐观锁
		Update(&models.LtGift{Id: id})
	return r, err
}

func (d *GiftDao) IncrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.ID(id).
		Incr("left_num", num).
		Update(&models.LtGift{Id: id})
	return r, err
}
