/**
 * @Time: 2019-08-18 00:24
 * @Author: solacowa@gmail.com
 * @File: bank
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type BankRepository interface {
	Create(name string) error
	List() (res []*types.Bank, err error)
}

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &bankRepository{db}
}

func (c *bankRepository) Create(name string) error {
	return c.db.Save(&types.Bank{BankName: name}).Error
}

func (c *bankRepository) List() (res []*types.Bank, err error) {
	err = c.db.Order("id DESC").Find(&res).Error
	return
}
