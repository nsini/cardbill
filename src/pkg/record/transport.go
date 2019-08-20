/**
 * @Time: 2019-08-18 11:34
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package record

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
	r.Handle("/record", kithttp.NewServer(
		eps.PostEndpoint,
		decodePostRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/record", kithttp.NewServer(
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

	var req tmePostRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(body), &req); err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(req.Amount, 10)
	businessType, err := strconv.ParseInt(req.BusinessType, 10, 64)
	rate, err := strconv.ParseFloat(req.Rate, 10)
	cardId, err := strconv.ParseInt(req.CardId, 10, 64)
	if err != nil {
		return nil, err
	}
	return postRequest{
		Amount:       amount,
		BusinessName: req.BusinessName,
		BusinessType: businessType,
		Rate:         rate / 10000,
		CardId:       cardId,
	}, nil
}
