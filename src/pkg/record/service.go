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
	"time"
)

var (
	ErrServiceCreate       = errors.New("记录创建错误")
	ErrServiceFindCard     = errors.New("获取卡信息错误")
	ErrServiceFindBusiness = errors.New("商户类型获取错误")
)

type Service interface {
	// 增加消费记录
	Post(ctx context.Context, cardId int64, businessType int64,
		businessName string, rate float64, amount float64, swipeTime *time.Time) (err error)

	// 消费列表
	List(ctx context.Context, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error)
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

func (c *service) List(ctx context.Context, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error) {
	// todo 应该会有很多条件 先从简单的开始

	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, 0, middleware.ErrCheckAuth
	}

	if page != 0 {
		page -= 1
	}
	_ = level.Debug(c.logger).Log("userId", userId)

	return c.repository.ExpenseRecord().List(userId, page, pageSize)
}

func (c *service) Post(ctx context.Context, cardId int64, businessType int64,
	businessName string, rate float64, amount float64, swipeTime *time.Time) (err error) {
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

	createdAt := time.Now()

	if swipeTime != nil {
		createdAt = *swipeTime
	}

	if err := c.repository.ExpenseRecord().Create(&types.ExpensesRecord{
		CardId:       card.Id,
		BusinessType: business.Id,
		BusinessName: businessName,
		Rate:         rate,
		Amount:       amount,
		Arrival:      amount - transform.Decimal(amount*rate),
		UserId:       userId,
		CreatedAt:    createdAt,
	}); err != nil {
		_ = level.Error(c.logger).Log("ExpenseRecord", "Create", "err", err.Error())
		return ErrServiceCreate
	}

	go func() {
		if err = c.repository.Merchant().FirstOrCreate(&types.Merchant{
			MerchantName: businessName,
			BusinessId:   business.Id,
			Business:     *business,
		}); err != nil {
			_ = level.Warn(c.logger).Log("Merchant", "FirstOrCreate", "err", err.Error())
		}
	}()

	return
}
