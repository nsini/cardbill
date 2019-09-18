/**
 * @Time : 2019-08-19 10:54
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package bank

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type postRequest struct {
	Name string `json:"name"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.Name)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.List(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
