package dao

import (
	"github.com/go-xorm/xorm"
	"log"
	"lottery/models"
)

type UserDao struct {
	engine *xorm.Engine
}

func NewUserDao(engine *xorm.Engine) *UserDao {
	return &UserDao{
		engine: engine,
	}
}

func (d *UserDao) Get(id int) *models.LtUser {
	data := &models.LtUser{}
	ok, err := d.engine.ID(id).Get(data)
	if ok && err == nil {
		return data
	} else {
		data.Id = 0
		return data
	}
}

func (d *UserDao) GetAll(page, size int) []models.LtUser {
	offset := (page - 1) * size
	datalist := make([]models.LtUser, 0)
	err := d.engine.
		Desc("id").
		Limit(size, offset).
		Find(&datalist)
	if err != nil {
		log.Println("dao.user_dao.GetAll error:", err)
		return datalist
	}

	return datalist
}

func (d *UserDao) CountAll() int64 {
	count, err := d.engine.Count(&models.LtUser{})
	if err != nil {
		return 0
	}
	return count
}

func (d *UserDao) Update(gift *models.LtUser, cols []string) error {
	_, err := d.engine.ID(gift.Id).MustCols(cols...).Update(gift)
	return err
}

func (d *UserDao) Create(gift *models.LtUser) error {
	_, err := d.engine.Insert(gift)
	return err
}
