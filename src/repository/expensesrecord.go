/**
 * @Time: 2019-08-18 00:32
 * @Author: solacowa@gmail.com
 * @File: expensesrecord
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type ExpenseRecordRepository interface {
	Create(record *types.ExpensesRecord) (err error)
	List(userId int64) (res []*types.ExpensesRecord, err error)
	RemainingAmount(cardId int64, billingDay time.Time, cardholder time.Time) (ra *RemainingAmount, err error)
}

type RemainingAmount struct {
	Amount float64
}

type expenseRecordRepository struct {
	db *gorm.DB
}

func NewExpenseRecordRepository(db *gorm.DB) ExpenseRecordRepository {
	return &expenseRecordRepository{db}
}

func (c *expenseRecordRepository) Create(record *types.ExpensesRecord) (err error) {
	return c.db.Save(record).Error
}

func (c *expenseRecordRepository) List(userId int64) (res []*types.ExpensesRecord, err error) {
	err = c.db.Where("user_id = ?", userId).Preload("CreditCard", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Bank")
	}).
		Preload("Business").
		Order("id DESC").Limit(20).Find(&res).Error
	return
}

func (c *expenseRecordRepository) RemainingAmount(cardId int64, billingDay time.Time, cardholder time.Time) (ra *RemainingAmount, err error) {
	var rs RemainingAmount
	err = c.db.Raw("SELECT SUM(amount) AS amount FROM expenses_records WHERE card_id = ? AND created_at > ? and created_at <= ?",
		cardId, billingDay.Format("2006-01-02 15:04:05"),
		cardholder.Format("2006-01-02 15:04:05")).Scan(&rs).Error
	return &rs, err
}
