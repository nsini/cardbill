/**
 * @Time : 2019-08-19 14:02
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package business

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
)

var (
	ErrBusinessExists = errors.New("商户类型已存在")
	ErrBusinessName   = errors.New("商户类型名称不能为空")
	ErrBusinessCode   = errors.New("商户类型MCC码不能正确")
)

type Service interface {
	// 商户类型列表
	List(ctx context.Context, name string) (res []*types.Business, err error)

	// 创建商户类型
	Post(ctx context.Context, name string, code int64) (err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{logger: logger, repository: repository}
}

func (c *service) List(ctx context.Context, name string) (res []*types.Business, err error) {
	return c.repository.Business().List(name)
}

func (c *service) Post(ctx context.Context, name string, code int64) (err error) {
	if _, err := c.repository.Business().FindByCode(code); err == nil {
		return ErrBusinessExists
	}
	return c.repository.Business().Create(name, code)
}
