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
	FindById(id, userId int64, column ...string) (res *types.CreditCard, err error)
	FindByUserId(userId, bankId int64) (res []*types.CreditCard, err error)
	Create(card *types.CreditCard) error
	Update(card *types.CreditCard) error
	FindByBillDay(day int) (res []*types.CreditCard, err error)
	Count(userId int64) (total int64, err error)
	Sum(userId int64) (total *TotalAmount, err error)
}

type TotalAmount struct {
	Amount    float64
	MaxAmount float64
}

type creditCardRepository struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {
	return &creditCardRepository{db}
}

func (c *creditCardRepository) Count(userId int64) (total int64, err error) {
	err = c.db.Model(&types.CreditCard{}).Where("user_id = ?", userId).Count(&total).Error
	return
}

func (c *creditCardRepository) Sum(userId int64) (total *TotalAmount, err error) {
	var totalAmount TotalAmount
	err = c.db.Model(&types.CreditCard{}).Select("SUM(fixed_amount) AS amount, SUM(max_amount) as max_amount").
		Where("user_id = ?", userId).Scan(&totalAmount).Error

	return &totalAmount, err
}

func (c *creditCardRepository) FindByBillDay(day int) (res []*types.CreditCard, err error) {
	err = c.db.Where("billing_day = ?", day).Find(&res).Error
	return
}

func (c *creditCardRepository) FindById(id, userId int64, column ...string) (res *types.CreditCard, err error) {
	var rs types.CreditCard
	query := c.db.Model(&rs)
	if len(column) > 0 {
		for _, v := range column {
			query = query.Preload(v)
		}
	}
	err = query.First(&rs, "id = ? AND user_id = ?", id, userId).Error
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
