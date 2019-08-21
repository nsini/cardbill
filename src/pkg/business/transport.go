/**
 * @Time : 2019-08-19 14:10
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package business

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
	"strings"
)

type endpoints struct {
	ListEndpoint endpoint.Endpoint
	PostEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		ListEndpoint: makeListEndpoint(svc),
		PostEndpoint: makePostEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"List": ems,
		"Post": ems,
	}

	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/business", kithttp.NewServer(
		eps.ListEndpoint,
		func(ctx context.Context, r *http.Request) (request interface{}, err error) {
			return listRequest{}, nil
		},
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/business", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	var req listRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrBusinessName
	}

	if req.Code < 1000 && req.Code > 9999 {
		return nil, ErrBusinessCode
	}

	return req, nil
}
