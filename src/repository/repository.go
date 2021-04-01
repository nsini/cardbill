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
	"github.com/nsini/cardbill/src/repository/record"
	"github.com/nsini/cardbill/src/repository/user"
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
	Users() user.Service
	Record() record.Service
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
	users         user.Service
	record        record.Service
}

func (c *repository) Record() record.Service {
	return c.record
}

func (c *repository) Users() user.Service {
	return c.users
}

func (c *repository) ChinaBank() bank.Service {
	return c.chinaBank
}

func NewRepository(db *gorm.DB, logger log.Logger, traceId string) Repository {

	bankSvc := bank.NewService(db)
	bankSvc = bank.NewLogging(logger, traceId)(bankSvc)

	userSvc := user.NewService(db)
	userSvc = user.NewLogging(logger, traceId)(userSvc)

	recordSvc := record.NewService(db)
	recordSvc = record.NewLogging(logger, traceId)(recordSvc)

	return &repository{
		record:        recordSvc,
		users:         userSvc,
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
