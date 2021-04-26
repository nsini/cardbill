/**
 * @Time : 2019-08-20 10:26
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

type logging struct {
	logger  log.Logger
	next    Service
	traceId string
}

func (s *logging) Info(ctx context.Context, userId int64) (res userInfoResult, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Info",
			"userId", userId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Info(ctx, userId)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	logger = log.With(logger, "user", "logging")
	return func(next Service) Service {
		return &logging{
			logger:  level.Info(logger),
			next:    next,
			traceId: traceId,
		}
	}
}

func (s *logging) Current(ctx context.Context) (user *types.User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Current",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Current(ctx)
}
