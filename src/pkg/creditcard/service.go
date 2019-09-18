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

	// todo 考虑写个定时任务生成账单

	for key, card := range res {
		curr := time.Now()
		year, month, _ := curr.Date()

		// todo 如果没有生成的话再走它 应该叫 "预计本期账单" 下次再更新吧
		// 上期账单
		startBillingDay := time.Date(year, month-1, card.BillingDay, 0, 0, 0, 1, &time.Location{})
		endBillingDay := time.Date(year, month, card.BillingDay+1, 0, 0, 0, 1, &time.Location{})

		ra, err := c.repository.ExpenseRecord().RemainingAmount(card.Id, startBillingDay, endBillingDay)
		if err != nil {
			_ = level.Error(c.logger).Log("ExpenseRecord", "RemainingAmount", "err", err.Error())
			continue
		}

		// todo 预计下期账单
		currStartBilling := time.Date(year, month, card.BillingDay+1, 0, 0, 0, 1, &time.Location{})
		currEndBilling := time.Date(year, month+1, card.BillingDay, 0, 0, 0, 1, &time.Location{})
		nextRes, err := c.repository.ExpenseRecord().RemainingAmount(card.Id, currStartBilling, currEndBilling)
		if err != nil {
			_ = level.Error(c.logger).Log("ExpenseRecord", "RemainingAmount", "err", err.Error())
			continue
		}

		res[key].BillingAmount = ra.Amount
		res[key].NextBillingAmount = nextRes.Amount
	}

	return
}
