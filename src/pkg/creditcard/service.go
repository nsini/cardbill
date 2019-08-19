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
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
)

type Service interface {
	// 增加信用卡
	Post(ctx context.Context, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder int) (err error)

	// 获取信用卡列表
	List(ctx context.Context) (res []*types.CreditCard, err error)

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

func (c *service) List(ctx context.Context) (res []*types.CreditCard, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}

	res, err = c.repository.CreditCard().FindByUserId(userId)

	return
}
