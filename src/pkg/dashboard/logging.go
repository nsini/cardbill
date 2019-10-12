/**
 * @Time : 2019-10-12 18:39
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package dashboard

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

func (s loggingService) LastAmount(ctx context.Context) (resp []LastAmountResponse, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "LastAmount",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.LastAmount(ctx)
}
