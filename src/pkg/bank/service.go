/**
 * @Time: 2019-08-18 17:18
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package bank

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/cardbill/src/repository"
)

type Service interface {
	// 增加银行
	Post(ctx context.Context, bankName string) (err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{logger: logger, repository: repository}
}

func (c *service) Post(ctx context.Context, bankName string) (err error) {
	return c.repository.Bank().Create(bankName)
}
