/**
 * @Time : 2019-08-20 10:20
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package user

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
)

type Service interface {
	Current(ctx context.Context) (user *types.User, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {

	return &service{logger: logger, repository: repository}
}

func (c *service) Current(ctx context.Context) (user *types.User, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}
	return c.repository.User().FindById(userId)
}
