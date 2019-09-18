/**
 * @Time : 2019-08-19 10:56
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package bank

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) Post(ctx context.Context, name string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Post",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name)
}
