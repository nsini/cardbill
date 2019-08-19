/**
 * @Time : 2019-08-19 10:00
 * @Author : solacowa@gmail.com
 * @File : merchant
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type MerchantRepository interface {
	Create(merchant *types.Merchant) error
	FindByName(name string) (res []*types.Merchant, err error)
	FirstOrCreate(merchant *types.Merchant) error
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (c *merchantRepository) Create(merchant *types.Merchant) error {
	return c.db.Save(merchant).Error
}

func (c *merchantRepository) FindByName(name string) (res []*types.Merchant, err error) {
	err = c.db.Where("merchant_name like ?", "%"+name+"%").Find(&res).Error
	return
}

func (c *merchantRepository) FirstOrCreate(merchant *types.Merchant) error {
	return c.db.FirstOrCreate(merchant, "merchant_name = ?", merchant.MerchantName).Error
}
