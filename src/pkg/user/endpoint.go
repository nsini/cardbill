/**
 * @Time : 2019-08-20 10:26
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/encode"
	"github.com/nsini/cardbill/src/middleware"
)

type (
	userInfoResult struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
	}
)

type Endpoints struct {
	InfoEndpoint    endpoint.Endpoint
	CurrentEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		InfoEndpoint:    makeInfoEndpoint(s),
		CurrentEndpoint: makeCurrentEndpoint(s),
	}

	for _, m := range dmw["Info"] {
		eps.InfoEndpoint = m(eps.InfoEndpoint)
	}
	for _, m := range dmw["Current"] {
		eps.CurrentEndpoint = m(eps.CurrentEndpoint)
	}

	return eps
}

func makeInfoEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			return nil, encode.ErrAccountNotLogin.Error()
		}
		res, err := s.Info(ctx, userId)
		return encode.Response{Error: err, Data: res}, err
	}
}

func makeCurrentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Current(ctx)
		return encode.Response{Error: err, Data: res}, err
	}
}
