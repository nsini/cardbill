/**
 * @Time: 2019-08-18 11:35
 * @Author: solacowa@gmail.com
 * @File: respoinse
 * @Software: GoLand
 */

package encode

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/nsini/cardbill/src/middleware"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Err     error       `json:"error,omitempty"`
}

func (r Response) error() error { return r.Err }

var (
	ErrBadRoute      = errors.New("bad route")
	ErrToken         = errors.New("Token 失效, 请重新登陆")
	ErrParamsRefused = errors.New("参数校验未通过")
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		EncodeError(ctx, e.error(), w)
		return nil
	}
	resp := response.(Response)
	if resp.Err == nil {
		resp.Success = true
	}

	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("Context-Type", "application/json")
	return kithttp.EncodeJSONResponse(ctx, w, resp)
}

type errorer interface {
	error() error
}

func EncodeError(ctx context.Context, err error, w http.ResponseWriter) {
	switch err {
	case kitjwt.ErrTokenContextMissing, kitjwt.ErrTokenExpired:
		err = ErrToken
		w.WriteHeader(http.StatusUnauthorized)
	//case casbin.ErrUnauthorized:
	//	err = ErrCasbin
	//	w.WriteHeader(http.StatusForbidden)
	case middleware.ErrCheckAuth, jwt.ErrSignatureInvalid:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusOK)
	}
	headers, ok := ctx.Value("response-headers").(map[string]string)
	if ok {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
}
