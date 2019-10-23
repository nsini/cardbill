/**
 * @Time : 2019-09-18 18:20
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package bill

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
	"time"
)

type endpoints struct {
	RepayEndpoint      endpoint.Endpoint
	ListEndpoint       endpoint.Endpoint
	ListByCardEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		RepayEndpoint:      makeRepayEndpoint(svc),
		ListEndpoint:       makeListEndpoint(svc),
		ListByCardEndpoint: makeListByCardEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Repay":      ems,
		"List":       ems,
		"ListByCard": ems,
	}

	for _, m := range mw["Repay"] {
		eps.RepayEndpoint = m(eps.RepayEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["ListByCard"] {
		eps.ListByCardEndpoint = m(eps.ListByCardEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/bill/repay", kithttp.NewServer(
		eps.RepayEndpoint,
		decodeRepayRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/bill", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/bill/card/{cardId:[0-9]+}", kithttp.NewServer(
		eps.ListByCardEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}

	page -= 1

	if pageSize == 0 {
		pageSize = 10
	}

	var cardId int64

	vars := mux.Vars(r)
	id, ok := vars["cardId"]
	if ok {
		intId, _ := strconv.Atoi(id)
		cardId = int64(intId)
	}

	return listRequest{
		pageSize: pageSize,
		page:     page,
		cardId:   cardId,
	}, nil
}

func decodeRepayRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	var req repayRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}

	if req.Repayment != "" {
		t, err := time.Parse("2006-01-02", req.Repayment)
		if err != nil {
			return nil, err
		}
		req.RepaymentDay = &t
	}

	return req, nil
}
