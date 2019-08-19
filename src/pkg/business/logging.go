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

func (s loggingService) List(ctx context.Context) (res []*types.Merchant, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.List(ctx)
}
