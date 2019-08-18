/**
 * @Time: 2019-08-17 20:32
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package record

import (
	"context"
	"github.com/nsini/cardbill/src/middleware"
)

type Service interface {
	Post(ctx context.Context, cardId int64, businessType int,
		businessName string, rateId int64, amount float64) (err error)
}

type service struct {
}

func NewService() Service {

	return &service{}
}

func (c *service) Post(ctx context.Context, cardId int64, businessType int,
	businessName string, rateId int64, amount float64) (err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	return
}
