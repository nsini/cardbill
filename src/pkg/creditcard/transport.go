/**
 * @Time : 2019-08-19 11:17
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package creditcard

import (
	"context"
	"encoding/json"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/util/encode"
	"io/ioutil"
	"net/http"
	"strconv"
)

type endpoints struct {
	GetEndpoint        endpoint.Endpoint
	PostEndpoint       endpoint.Endpoint
	ListEndpoint       endpoint.Endpoint
	PutEndpoint        endpoint.Endpoint
	StatisticsEndpoint endpoint.Endpoint
	RecordEndpoint     endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		GetEndpoint:        makeGetEndpoint(svc),
		PostEndpoint:       makePostEndpoint(svc),
		ListEndpoint:       makeListEndpoint(svc),
		PutEndpoint:        makePutEndpoint(svc),
		StatisticsEndpoint: makeStatisticsEndpoint(svc),
		RecordEndpoint:     makeRecordEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Get":        ems,
		"Post":       ems,
		"List":       ems,
		"Put":        ems,
		"Statistics": ems,
		"Record":     ems,
	}

	for _, m := range mw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Put"] {
		eps.PutEndpoint = m(eps.PutEndpoint)
	}
	for _, m := range mw["Statistics"] {
		eps.StatisticsEndpoint = m(eps.StatisticsEndpoint)
	}
	for _, m := range mw["Record"] {
		eps.RecordEndpoint = m(eps.RecordEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/creditcard", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/creditcard", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			bankId, _ := strconv.ParseInt(r.URL.Query().Get("bank_id"), 10, 64)
			stateStr := r.URL.Query().Get("state")
			var state int
			if stateStr == "" {
				state = -1
			} else {
				state, _ = strconv.Atoi(stateStr)
			}
			return listRequest{bankId, state}, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/creditcard/{id:[0-9]+}", kithttp.NewServer(
		eps.GetEndpoint,
		decodeGetRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/creditcard/statistics", kithttp.NewServer(
		eps.StatisticsEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/creditcard/{id:[0-9]+}", kithttp.NewServer(
		eps.PutEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("PUT")

	r.Handle("/creditcard/{id:[0-9]}/record", kithttp.NewServer(
		eps.RecordEndpoint,
		decodeRecordRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeRecordRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	paramId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return nil, err
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if pageSize == 0 {
		pageSize = 20
	}

	return recordRequest{Id: id, Page: page, PageSize: pageSize}, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	paramId, ok := vars["id"]
	if !ok {
		return nil, encode.ErrBadRoute
	}

	id, _ := strconv.ParseInt(paramId, 10, 64)
	return getRequest{id}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	var req postRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}

	// todo 对参数进行处理

	return req, nil
}
