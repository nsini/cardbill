/**
 * @Time: 2019-08-18 11:32
 * @Author: solacowa@gmail.com
 * @File: logging
 * @Software: GoLand
 */

package record

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

func (s loggingService) Post(ctx context.Context, cardId int64, businessType int64,
	businessName string, rate float64, amount float64, swipeTime *time.Time) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"cardId", cardId,
			"businessType", businessType,
			"businessName", businessName,
			"rate", rate,
			"amount", amount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, cardId, businessType, businessName, rate, amount, swipeTime)
}

func (s loggingService) List(ctx context.Context, page, pageSize int, bankId, cardId int64, start, end *time.Time) (res []*types.ExpensesRecord, count int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"page", page,
			"pageSize", pageSize,
			"bankId", bankId,
			"cardId", cardId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, pageSize, bankId, cardId, start, end)
}

func (s loggingService) Export(ctx context.Context, bankId, cardId int64, start, end *time.Time) (res []*types.ExpensesRecord, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"bankId", bankId,
			"cardId", cardId,
			//"start", start.String(),
			//"end", end.String(),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Export(ctx, bankId, cardId, start, end)
}
