/**
 * @Time : 3/30/21 5:27 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package mp

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/cardbill/src/encode"
	"net/http"
	"strconv"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"RecentRepay": ems,
		"BankList":    ems,
	})

	r := mux.NewRouter()

	r.Handle("/recent-repay", kithttp.NewServer(
		eps.RecentRepayEndpoint,
		decodeRecentRepayRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/banks", kithttp.NewServer(
		eps.BankListEndpoint,
		decodeBankListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/login", kithttp.NewServer(
		eps.LoginEndpoint,
		decodeMpLoginRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodeMpLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req mpLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if req.Code == "" {
		return nil, encode.ErrAuthMPLoginCode.Error()
	}
	return req, nil
}

func decodeBankListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req bankRequest
	req.bankName = r.URL.Query().Get("bankName")

	return req, nil
}

func decodeRecentRepayRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req recentRepayRequest
	recent, _ := strconv.Atoi(r.URL.Query().Get("recent"))
	if recent <= 0 {
		recent = 10
	}
	req.recent = recent

	return req, nil
}
