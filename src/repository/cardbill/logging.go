/**
 * @Time : 4/6/21 3:26 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package cardbill

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/repository/types"
	"time"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (l *loggingServer) FindById(ctx context.Context, id int64) (res types.Bill, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindById",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindById(ctx, id)
}

func (l *loggingServer) LastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res []types.Bill, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "LastBill",
			"cardIds", cardIds,
			"t", t,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.LastBill(ctx, cardIds, limit, t)
}

func (l *loggingServer) CountLastBill(ctx context.Context, cardIds []int64, limit int, t *time.Time) (res int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "CountLastBill",
			"cardIds", cardIds,
			"t", t,
			"total", res,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CountLastBill(ctx, cardIds, limit, t)
}

func (l *loggingServer) SumByCards(ctx context.Context, cardIds []int64, t *time.Time, repay Repay) (res BillAmount, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SumByCards",
			"cardIds", cardIds,
			"t", t,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SumByCards(ctx, cardIds, t, repay)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "cardbill", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
