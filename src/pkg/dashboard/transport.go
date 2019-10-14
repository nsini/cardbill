/**
 * @Time : 2019-10-12 18:42
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package dashboard

import (
	"context"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/util/encode"
	"net/http"
)

type endpoints struct {
	LastAmountEndpoint  endpoint.Endpoint
	MonthAmountEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		LastAmountEndpoint:  makeLastAmountEndpoint(svc),
		MonthAmountEndpoint: makeMonthAmountEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"LastAmount":  ems,
		"MonthAmount": ems,
	}

	for _, m := range mw["LastAmount"] {
		eps.LastAmountEndpoint = m(eps.LastAmountEndpoint)
	}
	for _, m := range mw["MonthAmount"] {
		eps.MonthAmountEndpoint = m(eps.MonthAmountEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/dashboard/last-amount", kithttp.NewServer(
		eps.LastAmountEndpoint,
		decodeLastAmountRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/dashboard/month-amount", kithttp.NewServer(
		eps.MonthAmountEndpoint,
		decodeLastAmountRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeLastAmountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return nil, nil
}
