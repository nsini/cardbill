/**
 * @Time: 2019-10-01 09:34
 * @Author: solacowa@gmail.com
 * @File: logger
 * @Software: GoLand
 */

package merchant

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

func (s loggingService) List(ctx context.Context, name string, page, pageSize int) (res []*types.Merchant, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, name, page, pageSize)
}
