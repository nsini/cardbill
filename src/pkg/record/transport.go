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
	"time"
)

type endpoints struct {
	PostEndpoint   endpoint.Endpoint
	ListEndpoint   endpoint.Endpoint
	ExportEndpoint endpoint.Endpoint
}

func MakeHandler(svc Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encode.EncodeError),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
	}

	eps := endpoints{
		PostEndpoint:   makePostEndpoint(svc),
		ListEndpoint:   makeListEndpoint(svc),
		ExportEndpoint: makeExportEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger), // 2
		//kitjwt.NewParser(kpljwt.JwtKeyFunc, jwt.SigningMethodHS256, kitjwt.StandardClaimsFactory), // 1
	}

	mw := map[string][]endpoint.Middleware{
		"Post":   ems,
		"List":   ems,
		"Export": ems,
	}

	for _, m := range mw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}
	for _, m := range mw["Export"] {
		eps.ExportEndpoint = m(eps.ExportEndpoint)
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
		decodeListRequest,
		encode.EncodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/record/export", kithttp.NewServer(
		eps.ExportEndpoint,
		decodeExportRequest,
		encodeExportResponse,
		opts...,
	)).Methods("GET")

	return r
}

func decodeExportRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	cardId, _ := strconv.ParseInt(r.URL.Query().Get("cardId"), 10, 64)
	bankId, _ := strconv.ParseInt(r.URL.Query().Get("bankId"), 10, 64)
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	var startTime, endTime *time.Time

	if t, err := time.Parse("2006-01-02", start); err == nil {
		startTime = &t
	}

	if t, err := time.Parse("2006-01-02", end); err == nil {
		endTime = &t
	}

	return listRequest{
		BankId: bankId,
		CardId: cardId,
		Start:  startTime,
		End:    endTime,
	}, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	cardId, _ := strconv.ParseInt(r.URL.Query().Get("cardId"), 10, 64)
	bankId, _ := strconv.ParseInt(r.URL.Query().Get("bankId"), 10, 64)
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	var startTime, endTime *time.Time

	if t, err := time.Parse("2006-01-02", start); err == nil {
		startTime = &t
	}

	if t, err := time.Parse("2006-01-02", end); err == nil {
		endTime = &t
	}

	if pageSize == 0 {
		pageSize = 10
	}

	return listRequest{
		Page: page, PageSize: pageSize,
		Start:  startTime,
		End:    endTime,
		BankId: bankId,
		CardId: cardId,
	}, nil
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {

	var req postRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	req.Rate /= 10000

	if req.TmpTime != "" {
		if t, err := time.Parse("2006-01-02T15:04:05Z", req.TmpTime); err == nil {
			tt := t.Local()
			req.SwipeTime = &tt
		} else {
			return nil, err
		}
	}

	return req, nil
}
