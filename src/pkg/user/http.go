/**
 * @Time : 2019-08-20 10:28
 * @Author : solacowa@gmail.com
 * @File : transport
 * @Software: GoLand
 */

package user

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/nsini/cardbill/src/encode"
	"net/http"
)

func MakeHTTPHandler(svc Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	eps := NewEndpoint(svc, map[string][]endpoint.Middleware{
		"Info":    dmw,
		"Current": dmw,
	})

	r := mux.NewRouter()

	r.Handle("/current", kithttp.NewServer(
		eps.CurrentEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")
	r.Handle("/info", kithttp.NewServer(
		eps.CurrentEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods("GET")

	return r
}
