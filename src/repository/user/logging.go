/**
 * @Time : 3/31/21 1:40 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package user

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

func (l *loggingServer) Save(ctx context.Context, user *types.User) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "Save",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.Save(ctx, user)
}

func (l *loggingServer) FindByUnionId(ctx context.Context, unionId string) (user types.User, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "FindByUnionId",
			"unionId", unionId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.FindByUnionId(ctx, unionId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "users", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}
