/**
 * @Time: 2020/10/16 22:14
 * @Author: solacowa@gmail.com
 * @File: http
 * @Software: GoLand
 */

package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/cardbill/src/encode"
	"net/http"
)

func MakeHTTPHandler(pattern string, logger kitlog.Logger, s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"JsSDK":       ems,
		"Callback":    ems,
		"AuthCodeURL": ems,
	})

	r := mux.NewRouter()

	r.Handle(fmt.Sprintf("%s/jssdk", pattern), kithttp.NewServer(
		eps.JsSDKEndpoint,
		decodeJsSDKRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost, http.MethodGet)
	r.Handle(fmt.Sprintf("%s/auth/code-url", pattern), kithttp.NewServer(
		eps.AuthCodeURLEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle(fmt.Sprintf("%s/callback", pattern), kithttp.NewServer(
		eps.CallbackEndpoint,
		decodeCallbackRequest,
		encodeCallbackResponse,
		opts...,
	))
	r.Handle(fmt.Sprintf("%s/callback/mch", pattern), kithttp.NewServer(
		eps.CallbackMchEndpoint,
		decodeCallbackOrderNotifyRequest,
		encodeCallbackMchResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeCallbackOrderNotifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return callbackRequest{r: r}, nil
}

func decodeJsSDKRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req jsSDKRequest
	_ = r.ParseForm()
	link := r.FormValue("url")
	if link != "" {
		req.Link = link
		return req, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeCallbackRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return callbackRequest{r: r}, nil
}

func encodeCallbackResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(encode.Response)
	data := resp.Data.(callbackResponse)
	data.server.ServeHTTP(w, data.r, nil)
	return nil
}

func encodeCallbackMchResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(encode.Response)
	data := resp.Data.(callbackMchResponse)
	data.server.ServeHTTP(w, data.r, nil)
	return nil
}
