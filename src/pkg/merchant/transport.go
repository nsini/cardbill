/**
 * @Time: 2019-10-01 09:38
 * @Author: solacowa@gmail.com
 * @File: transport
 * @Software: GoLand
 */

package merchant

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
		ListEndpoint: makeListEndpoint(svc),
	}

	ems := []endpoint.Middleware{
		middleware.CheckLogin(logger),
	}

	mw := map[string][]endpoint.Middleware{
		"List": ems,
	}

	for _, m := range mw["List"] {
		eps.ListEndpoint = m(eps.ListEndpoint)
	}

	r := mux.NewRouter()

	r.Handle("/merchant", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.EncodeResponse,
		opts...)).Methods(http.MethodGet)

	return r
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	val := r.URL.Query().Get("q")
	return listRequest{name: val}, nil
}
