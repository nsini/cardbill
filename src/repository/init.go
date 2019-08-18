/**
 * @Time: 2019-08-17 20:32
 * @Author: solacowa@gmail.com
 * @File: init
 * @Software: GoLand
 */

package repository

import "github.com/jinzhu/gorm"

type Repository interface {
	Bank() BankRepository
	ExpenseRecord() ExpenseRecordRepository
}

type repository struct {
	bankRepository          BankRepository
	expenseRecordRepository ExpenseRecordRepository
}

func NewRepository(db *gorm.DB) Repository {

	return &repository{
		bankRepository:          NewBankRepository(db),
		expenseRecordRepository: NewExpenseRecordRepository(db),
	}
}

func (c *repository) Bank() BankRepository {
	return c.bankRepository
}

func (c *repository) ExpenseRecord() ExpenseRecordRepository {
	return c.expenseRecordRepository
}
