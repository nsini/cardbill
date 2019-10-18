/**
 * @Time: 2019-10-01 09:33
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package merchant

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/cardbill/src/repository/types"
)

type Service interface {
	List(ctx context.Context, name string, page, pageSize int) (res []*types.Merchant, err error)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{logger: logger}
}

func (c *service) List(ctx context.Context, name string, page, pageSize int) (res []*types.Merchant, err error) {

	return
}
