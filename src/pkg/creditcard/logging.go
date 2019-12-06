/**
 * @Time : 2019-08-19 11:13
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package creditcard

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) Post(ctx context.Context, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder int, cardNumber, tailNumber int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"cardName", cardName,
			"bankId", bankId,
			"fixedAmount", fixedAmount,
			"maxAmount", maxAmount,
			"billingDay", billingDay,
			"cardHolder", cardHolder,
			"cardNumber", cardNumber,
			"tailNumber", tailNumber,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, cardName, bankId, fixedAmount, maxAmount, billingDay,
		cardHolder, cardNumber, tailNumber)
}

func (s loggingService) List(ctx context.Context, bankId int64) (res []*types.CreditCard, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"bankId", bankId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, bankId)
}

func (s loggingService) Get(ctx context.Context, id int64) (res *types.CreditCard, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Get",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s loggingService) Statistics(ctx context.Context) (res *StatisticsResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Statistics",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Statistics(ctx)
}

func (s loggingService) Record(ctx context.Context, id int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Record",
			"id", id,
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Record(ctx, id, page, pageSize)
}
