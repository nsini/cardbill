/**
 * @Time : 2019-08-20 10:26
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

func makeCurrentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Current(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
