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
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
)

type Service interface {
	Current(ctx context.Context) (user *types.User, err error)

	Info(ctx context.Context, userId int64) (res userInfoResult, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
	traceId    string
}

func (s *service) Info(ctx context.Context, userId int64) (res userInfoResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Info")
	userInfo, err := s.repository.User().FindById(userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.User", "FindById", "err", err.Error())
		return
	}

	res.Username = userInfo.Username
	return
}

func NewService(logger log.Logger, repository repository.Repository) Service {

	return &service{logger: logger, repository: repository}
}

func (s *service) Current(ctx context.Context) (user *types.User, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}
	return c.repository.User().FindById(userId)
}
