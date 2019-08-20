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
	List() (res []*types.Business, err error)
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

func (c *businessRepository) List() (res []*types.Business, err error) {
	err = c.db.Order("id DESC").Find(&res).Error
	return
}
