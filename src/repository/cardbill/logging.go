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
	"time"
)

type loggingServer struct {
	logger  log.Logger
	next    Service
	traceId string
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
