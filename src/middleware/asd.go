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
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/nsini/cardbill/src/encode"
	asdjwt "github.com/nsini/cardbill/src/jwt"
	cbjwt "github.com/nsini/cardbill/src/jwt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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

func CheckAuthMiddleware(logger log.Logger, tracer opentracing.Tracer) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if tracer != nil {
				var span opentracing.Span
				span, ctx = opentracing.StartSpanFromContextWithTracer(ctx, tracer, "CheckAuthMiddleware", opentracing.Tag{
					Key:   string(ext.Component),
					Value: "Middleware",
				})
				defer func() {
					span.LogKV("err", err)
					span.Finish()
				}()
			}

			token := ctx.Value(kithttp.ContextKeyRequestAuthorization).(string)
			if token == "" {
				_ = level.Warn(logger).Log("ctx", "Value", "err", encode.ErrAuthNotLogin.Error())
				return nil, encode.ErrAuthNotLogin.Error()
			}

			var clustom asdjwt.ArithmeticCustomClaims
			tk, err := jwt.ParseWithClaims(token, &clustom, asdjwt.JwtKeyFunc)
			if err != nil || tk == nil {
				_ = level.Error(logger).Log("jwt", "ParseWithClaims", "err", err)
				err = encode.ErrAuthNotLogin.Wrap(err)
				return
			}

			claim, ok := tk.Claims.(*asdjwt.ArithmeticCustomClaims)
			if !ok {
				_ = level.Error(logger).Log("tk", "Claims", "err", ok)
				err = encode.ErrAccountASD.Error()
				return
			}

			// 查询用户是否退出
			//if _, err = cache.Get(ctx, fmt.Sprintf("login:%d:token", claim.UserId), nil); err != nil {
			//	_ = level.Error(logger).Log("cache", "Get", "err", err)
			//	err = encode.ErrAuthNotLogin.Wrap(err)
			//	return
			//}

			// TODO 获取用户信息
			//if sysUser.Locked {
			//	err = encode.ErrAccountLocked.Error()
			//	_ = level.Error(logger).Log("sysUser", "Locked", "err", err)
			//	return
			//}

			ctx = context.WithValue(ctx, UserIdContext, claim.UserId)
			ctx = context.WithValue(ctx, "Authorization", token)
			return next(ctx, request)
		}
	}
}
