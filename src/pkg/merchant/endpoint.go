/**
 * @Time: 2019-10-01 09:36
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package merchant

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type listRequest struct {
	page, pageSize int
	name           string
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.name, req.page, req.pageSize)
		return encode.Response{Data: res, Err: err}, err
	}
}
