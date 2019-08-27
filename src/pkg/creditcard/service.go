/**
 * @Time : 2019-08-19 11:03
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package creditcard

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type Service interface {
	// 增加信用卡
	Post(ctx context.Context, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder int) (err error)

	// 获取信用卡列表
	List(ctx context.Context, bankId int64) (res []*types.CreditCard, err error)

	// 更新信用卡信息
	Put(ctx context.Context, id int64, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder, state int) (err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{logger: logger, repository: repository}
}

func (c *service) Post(ctx context.Context, cardName string, bankId int64, fixedAmount, maxAmount float64, billingDay, cardHolder int) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	return c.repository.CreditCard().Create(&types.CreditCard{
		CardName:    cardName,
		BankId:      bankId,
		FixedAmount: fixedAmount,
		MaxAmount:   maxAmount,
		BillingDay:  billingDay,
		Cardholder:  cardHolder,
		UserId:      userId,
	})
}

func (c *service) Put(ctx context.Context, id int64, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder, state int) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	return c.repository.CreditCard().Update(&types.CreditCard{
		Id:          id,
		CardName:    cardName,
		BankId:      bankId,
		FixedAmount: fixedAmount,
		MaxAmount:   maxAmount,
		BillingDay:  billingDay,
		Cardholder:  cardHolder,
		UserId:      userId,
		State:       state,
	})

}

func (c *service) List(ctx context.Context, bankId int64) (res []*types.CreditCard, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}

	res, err = c.repository.CreditCard().FindByUserId(userId, bankId)
	if err != nil {
		return
	}

	for key, card := range res {
		curr := time.Now()
		year, month, _ := curr.Date()

		billingMonth := month - 1

		// 账单日
		billingDay := time.Date(year, billingMonth, card.BillingDay, 0, 0, 0, 1, &time.Location{}) // .Format("2006-01-02 15:04:05")
		// 还款日
		cardholder := time.Date(year, month, card.Cardholder, 23, 59, 59, 59, &time.Location{}) // .Format("2006-01-02 15:04:05")

		ra, err := c.repository.ExpenseRecord().RemainingAmount(card.Id, billingDay, cardholder)
		if err != nil {
			_ = level.Error(c.logger).Log("ExpenseRecord", "RemainingAmount", "err", err.Error())
			continue
		}
		res[key].BillingAmount = ra.Amount
	}

	return
}
