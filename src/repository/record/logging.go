/**
 * @Time : 4/1/21 6:04 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package record

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

func (l *loggingServer) List(ctx context.Context, userId int64, page, pageSize int, bankId int64, cardIds []int64, start, end *time.Time) (res []types.ExpensesRecord, total int, err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log(
			l.traceId, ctx.Value(l.traceId),
			"method", "List",
			"userId", userId,
			"page", page,
			"pageSize", pageSize,
			"bankId", bankId,
			"bankId", bankId,
			"total", total,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.List(ctx, userId, page, pageSize, bankId, cardIds, start, end)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "record", "logging")
	return func(next Service) Service {
		return &loggingServer{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}