/**
 * @Time : 2019-10-12 18:40
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package dashboard

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

func makeLastAmountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.LastAmount(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
