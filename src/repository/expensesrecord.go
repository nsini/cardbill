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
)

type ExpenseRecordRepository interface {
	Create(record *types.ExpensesRecord) (err error)
	List(userId int64) (res []*types.ExpensesRecord, err error)
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
	err = c.db.Where("user_id = ?", userId).Preload("CreditCard").
		Preload("Business").Order("id DESC").Limit(20).Find(&res).Error
	return
}
