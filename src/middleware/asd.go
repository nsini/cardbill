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
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	cbjwt "github.com/nsini/cardbill/src/jwt"
	"strings"
)

type ASDContext string

const (
	UserIdContext ASDContext = "user-id-context"
)

var (
	ErrCheckAuth = errors.New("请进行授权登陆")
	//ErrorASD     = errors.New("权限验证失败！")
)

func CheckLogin(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			if token == "" {
				return nil, ErrCheckAuth
			}
			if strings.Contains(token, "Bearer ") {
				token = strings.Split(token, "Bearer ")[1]
			}

			var clustom cbjwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &clustom, cbjwt.JwtKeyFunc)
			if err != nil || tk == nil {
				_ = logger.Log("jwt", "ParseWithClaims", "err", err)
				return nil, ErrCheckAuth
			}

			claim, ok := tk.Claims.(*cbjwt.ArithmeticCustomClaims)
			if !ok {
				_ = logger.Log("tk", "Claims", "err", ok)
				err = ErrCheckAuth
				return
			}

			ctx = context.WithValue(ctx, UserIdContext, claim.UserId)
			return next(ctx, request)
		}
	}
}
