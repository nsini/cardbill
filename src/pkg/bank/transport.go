/**
 * @Time : 2019-08-19 10:55
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package bank

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
)

type endpoints struct {
	PostEndpoint endpoint.Endpoint
	ListEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		PostEndpoint: makePostEndpoint(svc),
		ListEndpoint: makeListEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Post": ems,
		"List": ems,
	}

	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/bank", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/bank", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return nil, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
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

	if req.Name == "" {
		return nil, ErrBankNameNotNull
	}

	// todo 对参数进行处理

	return req, nil
}
