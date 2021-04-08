/**
 * @Time : 3/30/21 5:08 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package bank

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

func (l *loggingServer) Find(ctx context.Context, bankId int64) (res types.Bank, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Find",
			"bankId", bankId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Find(ctx, bankId)
}

func (l *loggingServer) List(ctx context.Context, bankName string) (res []types.Bank, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"bankName", bankName,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, bankName)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "bank", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
