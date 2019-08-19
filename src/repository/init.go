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
	CreditCard() CreditCardRepository
	Business() BusinessRepository
	User() UserRepository
	Merchant() MerchantRepository
}

type repository struct {
	bank          BankRepository
	expenseRecord ExpenseRecordRepository
	creditCard    CreditCardRepository
	business      BusinessRepository
	user          UserRepository
	merchant      MerchantRepository
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		bank:          NewBankRepository(db),
		expenseRecord: NewExpenseRecordRepository(db),
		creditCard:    NewCreditCardRepository(db),
		business:      NewBusinessRepository(db),
		user:          NewUserRepository(db),
		merchant:      NewMerchantRepository(db),
	}
}

func (c *repository) Bank() BankRepository {
	return c.bank
}

func (c *repository) ExpenseRecord() ExpenseRecordRepository {
	return c.expenseRecord
}

func (c *repository) CreditCard() CreditCardRepository {
	return c.creditCard
}

func (c *repository) Business() BusinessRepository {
	return c.business
}

func (c *repository) User() UserRepository {
	return c.user
}

func (c *repository) Merchant() MerchantRepository {
	return c.merchant
}
