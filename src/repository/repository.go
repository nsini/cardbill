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
	"github.com/nsini/cardbill/src/repository/business"
	"github.com/nsini/cardbill/src/repository/card"
	"github.com/nsini/cardbill/src/repository/cardbill"
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
	Card() card.Service
	BusinessSvc() business.Service
	CardBill() cardbill.Service
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
	card          card.Service
	businessSvc   business.Service
	cardBill      cardbill.Service
}

func (c *repository) CardBill() cardbill.Service {
	return c.cardBill
}

func (c *repository) BusinessSvc() business.Service {
	return c.businessSvc
}

func (c *repository) Card() card.Service {
	return c.card
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

	cardSvc := card.NewService(db)
	cardSvc = card.NewLogging(logger, traceId)(cardSvc)

	businessSvc := business.NewService(db)
	businessSvc = business.NewLogging(logger, traceId)(businessSvc)

	cardBillSvc := cardbill.NewService(db)
	cardBillSvc = cardbill.NewLogging(logger, traceId)(cardBillSvc)

	return &repository{
		cardBill:      cardBillSvc,
		businessSvc:   businessSvc,
		card:          cardSvc,
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
