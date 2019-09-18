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

func (s loggingService) Current(ctx context.Context) (user *types.User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Current",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Current(ctx)
}
