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
	"net/http"
)

func encodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	res := response.(authResponse)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Authorization", res.Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    res.Token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7200})
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
	default:
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  -1,
		"error": err.Error(),
	})
}
