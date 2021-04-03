/**
 * @Time: 2021/4/3 下午11:13
 * @Author: solacowa@gmail.com
 * @File: logging
 * @Software: GoLand
 */

package card

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

func (l *loggingServer) FindByUserId(ctx context.Context, userId int64) (res []types.CreditCard, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByUserId",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByUserId(ctx, userId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "card", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
