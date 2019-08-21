/**
 * @Time : 2019-08-19 14:09
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package business

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type listRequest struct {
	Name string `json:"name"`
	Code int64  `json:"code"`
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.Name)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		err := s.Post(ctx, req.Name, req.Code)
		return encode.Response{Err: err, Data: nil}, err
	}
}
