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
	FindById(id, userId int64) (res *types.CreditCard, err error)
	FindByUserId(userId, bankId int64) (res []*types.CreditCard, err error)
	Create(card *types.CreditCard) error
	Update(card *types.CreditCard) error
	FindByBillDay(day int) (res []*types.CreditCard, err error)
}

type creditCardRepository struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {
	return &creditCardRepository{db}
}

func (c *creditCardRepository) FindByBillDay(day int) (res []*types.CreditCard, err error) {
	err = c.db.Where("billing_day = ?", day).Find(&res).Error
	return
}

func (c *creditCardRepository) FindById(id, userId int64) (res *types.CreditCard, err error) {
	var rs types.CreditCard
	err = c.db.First(&rs, "id = ? AND user_id = ?", id, userId).Error
	return &rs, err
}

func (c *creditCardRepository) Create(card *types.CreditCard) error {
	return c.db.Save(card).Error
}

func (c *creditCardRepository) FindByUserId(userId, bankId int64) (res []*types.CreditCard, err error) {
	query := c.db.Where("user_id = ?", userId)
	if bankId != 0 {
		query = query.Where("bank_id = ?", bankId)
	}
	err = query.Order("bank_id DESC").Preload("Bank").Find(&res).Error
	return
}

func (c *creditCardRepository) Update(card *types.CreditCard) error {
	tx := c.db.Begin()
	err := c.db.Model(&card).Where("id = ? AND user_id = ?", card.Id, card.UserId).Update(card).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
