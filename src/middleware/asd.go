/**
 * @Time: 2019-08-18 09:43
 * @Author: solacowa@gmail.com
 * @File: asd
 * @Software: GoLand
 */

package middleware

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type ASDContext string

const (
	UserIdContext ASDContext = "user-id-context"
)

var (
	ErrCheckAuth = errors.New("请进行授权登陆")
)

func CheckLogin(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			ctx = context.WithValue(ctx, UserIdContext, 1)
			return next(ctx, request)
		}
	}
}
