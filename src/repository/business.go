/**
 * @Time: 2019-08-18 11:06
 * @Author: solacowa@gmail.com
 * @File: business
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type BusinessRepository interface {
	FindById(id int64) (res *types.Business, err error)
	List(name string) (res []*types.Business, err error)
	Create(name string, code int64) (err error)
	FindByCode(code int64) (res *types.Business, err error)
}

type businessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) BusinessRepository {

	return &businessRepository{db: db}
}

func (c *businessRepository) FindById(id int64) (res *types.Business, err error) {
	var rs types.Business
	err = c.db.First(&rs, " id = ?", id).Error
	return &rs, err
}

func (c *businessRepository) List(name string) (res []*types.Business, err error) {
	query := c.db.Order("id DESC")
	if name != "" {
		query = query.Where("business_name like ?", "%"+name+"%").
			Or("code like ?", "%"+name+"%")
	}
	err = query.Find(&res).Error
	return
}

func (c *businessRepository) Create(name string, code int64) (err error) {
	return c.db.Save(&types.Business{
		BusinessName: name,
		Code:         code,
	}).Error
}

func (c *businessRepository) FindByCode(code int64) (res *types.Business, err error) {
	var rs types.Business
	err = c.db.First(&rs, "code = ?", code).Error
	return &rs, err
}
