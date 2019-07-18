package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/models"
)

type UserdayDao struct {
	engine *xorm.Engine
}

func NewUserdayDao(engine *xorm.Engine) *UserdayDao {
	return &UserdayDao{
		engine: engine,
	}
}

func (d *UserdayDao) Get(id int) *models.LtUserday {
	data := &models.LtUserday{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserdayDao) GetAll() []models.LtUserday {
	datalist := make([]models.LtUserday, 0)
	err := d.engine.
		Desc("id").
		Find(&datalist)
	if err != nil {
		log.Println("dao.userday_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *UserdayDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtUserday{})
	if err != nil {
		return 0
	}
	return count
}

func (d *UserdayDao) Search(uid, day int) []models.LtUserday {
	datalist := make([]models.LtUserday, 0)
	err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Desc("id").
		Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *UserdayDao) Count(uid, day int) int {
	info := &models.LtUserday{}
	ok, err := d.engine.
		Where("uid=?", uid).
		Where("day=?", day).
		Get(info)
	if !ok || err != nil {
		return 0
	} else {
		return info.Num
	}
}

func (d *UserdayDao) Update(gift *models.LtUserday, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *UserdayDao) Create(gift *models.LtUserday) error {
	_, err := d.engine.Insert(gift)
	return err
}
