/**
 * @Time: 2019-08-17 20:32
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package record

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
	"github.com/nsini/cardbill/src/util/transform"
)

var (
	ErrServiceCreate       = errors.New("记录创建错误")
	ErrServiceFindCard     = errors.New("获取卡信息错误")
	ErrServiceFindBusiness = errors.New("商户类型获取错误")
)

type Service interface {
	// 增加消费记录
	Post(ctx context.Context, cardId int64, businessType int64,
		businessName string, rate float64, amount float64) (err error)

	// 消费列表
	List(ctx context.Context) (res []*types.ExpensesRecord, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{
		logger:     logger,
		repository: repository,
	}
}

func (c *service) List(ctx context.Context) (res []*types.ExpensesRecord, err error) {
	// todo 应该会有很多条件 先从简单的开始

	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}

	_ = level.Debug(c.logger).Log("userId", userId)

	return c.repository.ExpenseRecord().List(userId)
}

func (c *service) Post(ctx context.Context, cardId int64, businessType int64,
	businessName string, rate float64, amount float64) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	card, err := c.repository.CreditCard().FindById(cardId, userId)
	if err != nil {
		_ = level.Warn(c.logger).Log("CreditCard", "FindById", "err", err.Error())
		return ErrServiceFindCard
	}

	business, err := c.repository.Business().FindById(businessType)
	if err != nil {
		_ = level.Warn(c.logger).Log("Business", "FindById", "err", err.Error())
		return ErrServiceFindBusiness
	}

	if err := c.repository.ExpenseRecord().Create(&types.ExpensesRecord{
		CardId:       card.Id,
		BusinessType: business.Id,
		BusinessName: businessName,
		Rate:         rate,
		Amount:       amount,
		Arrival:      amount - transform.Decimal(amount*rate),
		UserId:       userId,
	}); err != nil {
		_ = level.Error(c.logger).Log("ExpenseRecord", "Create", "err", err.Error())
		return ErrServiceCreate
	}

	go func() {
		if err = c.repository.Merchant().FirstOrCreate(&types.Merchant{
			MerchantName: businessName,
			BusinessId:   business.Id,
		}); err != nil {
			_ = level.Warn(c.logger).Log("Merchant", "FirstOrCreate", "err", err.Error())
		}
	}()

	return
}
