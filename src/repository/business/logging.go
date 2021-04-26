/**
 * @Time : 4/6/21 10:09 AM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package business

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

func (l *loggingServer) FindByCode(ctx context.Context, code int64) (res types.Business, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByCode",
			"code", code,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByCode(ctx, code)
}

func (l *loggingServer) Types(ctx context.Context) (res []types.Business, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Types",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Types(ctx)
}

func (l *loggingServer) SaveMerchant(ctx context.Context, business *types.Merchant) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "SaveMerchant",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.SaveMerchant(ctx, business)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "business", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
