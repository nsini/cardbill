/**
 * @Time : 2019-08-19 14:07
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package business

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

func (s loggingService) List(ctx context.Context, name string) (res []*types.Business, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, name)
}

func (s loggingService) Post(ctx context.Context, name string, code int64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"code", code,
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, code)
}
