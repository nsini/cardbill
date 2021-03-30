/**
 * @Time: 2019-08-17 20:32
 * @Author: solacowa@gmail.com
 * @File: init
 * @Software: GoLand
 */

package repository

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/bank"
)

type Repository interface {
	Bank() BankRepository
	ExpenseRecord() ExpenseRecordRepository
	CreditCard() CreditCardRepository
	Business() BusinessRepository
	User() UserRepository
	Merchant() MerchantRepository
	Bill() BillRepository
	ChinaBank() bank.Service
}

type repository struct {
	bank          BankRepository
	expenseRecord ExpenseRecordRepository
	creditCard    CreditCardRepository
	business      BusinessRepository
	user          UserRepository
	merchant      MerchantRepository
	bill          BillRepository
	chinaBank     bank.Service
}

func (c *repository) ChinaBank() bank.Service {
	return c.chinaBank
}

func NewRepository(db *gorm.DB, logger log.Logger, traceId string) Repository {

	bankSvc := bank.NewService(db)
	bankSvc = bank.NewLogging(logger, traceId)(bankSvc)

	return &repository{
		chinaBank:     bankSvc,
		bank:          NewBankRepository(db),
		expenseRecord: NewExpenseRecordRepository(db),
		creditCard:    NewCreditCardRepository(db),
		business:      NewBusinessRepository(db),
		user:          NewUserRepository(db),
		merchant:      NewMerchantRepository(db),
		bill:          NewBillRepository(db),
	}
}

func (c *repository) Bill() BillRepository {
	return c.bill
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
