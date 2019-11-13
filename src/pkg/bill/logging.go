/**
 * @Time : 2019-09-18 18:17
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package bill

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

func (s loggingService) Repay(ctx context.Context, cardId int64, amount float64, repaymentDay *time.Time) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Repay",
			"cardId", cardId,
			"amount", amount,
			"repaymentDay", repaymentDay,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Repay(ctx, cardId, amount, repaymentDay)
}

func (s loggingService) List(ctx context.Context, page, pageSize int) (res []*types.Bill, count int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"page", page,
			"pageSize", pageSize,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, page, pageSize)
}

func (s loggingService) ListByCard(ctx context.Context, cardId int64, page, pageSize int) (res []*types.Bill, count int64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "ListByCard",
			"page", page,
			"pageSize", pageSize,
			"cardId", cardId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.ListByCard(ctx, cardId, page, pageSize)
}

func (s loggingService) RecentRepay(ctx context.Context, recent int) (res []*types.Bill, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "RecentRepay",
			"recent", recent,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.RecentRepay(ctx, recent)
}
