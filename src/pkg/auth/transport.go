/**
 * @Time: 2019-08-18 17:07
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package auth

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type endpoints struct {
}

func MakeHandler(svc Service, logger kitlog.Logger) http.Handler {
	//ctx := context.Background()
	r := mux.NewRouter()
	r.HandleFunc("/github/callback", svc.AuthLoginGithubCallback).Methods("GET")
	r.HandleFunc("/github/login", svc.AuthLoginGithub).Methods("GET")

	return r
}

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	res := response.(authResponse)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Authorization", res.Token)
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	case jwt.ErrSignatureInvalid:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
