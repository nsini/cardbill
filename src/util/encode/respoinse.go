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
	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/nsini/cardbill/src/middleware"
	"net/http"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}

type errorer interface {
	error() error
}

func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
}
