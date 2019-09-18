/**
 * @Time : 2019-09-18 18:17
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package bill

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

func (s loggingService) Repay(ctx context.Context, cardId int64, amount float64) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Repay",
			"cardId", cardId,
			"amount", amount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Repay(ctx, cardId, amount)
}
