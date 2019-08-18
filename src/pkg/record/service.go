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
	Post(ctx context.Context, cardId int64, businessType int64,
		businessName string, rate float64, amount float64) (err error)
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

func (c *service) Post(ctx context.Context, cardId int64, businessType int64,
	businessName string, rate float64, amount float64) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	card, err := c.repository.CreditCard().FindById(cardId)
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
		Arrival:      transform.Decimal(amount * rate),
		UserId:       userId,
	}); err != nil {
		_ = level.Error(c.logger).Log("ExpenseRecord", "Create", "err", err.Error())
		return ErrServiceCreate
	}

	return
}
