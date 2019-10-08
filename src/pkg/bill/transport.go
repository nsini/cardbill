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
	"time"
)

type endpoints struct {
	RepayEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		RepayEndpoint: makeRepayEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Repay": ems,
	}

	for _, m := range mw["Repay"] {
		eps.RepayEndpoint = m(eps.RepayEndpoint)
	}

	r := mux.NewRouter()
	r.Handle("/bill/repay", kithttp.NewServer(
		eps.RepayEndpoint,
		decodeRepayRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	return r
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
