/**
 * @Time: 2019-08-18 10:40
 * @Author: solacowa@gmail.com
 * @File: creditcard
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type CreditCardRepository interface {
	FindById(id int64) (res *types.CreditCard, err error)
}

type creditCardRepository struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {

	return &creditCardRepository{db}
}

func (c *creditCardRepository) FindById(id int64) (res *types.CreditCard, err error) {
	var rs types.CreditCard
	err = c.db.First(&rs, "id = ? ", id).Error
	return &rs, err
}
