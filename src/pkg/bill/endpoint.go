/**
 * @Time : 2019-09-18 18:17
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package bill

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type repayRequest struct {
	CardId int64   `json:"card_id"`
	Amount float64 `json:"amount"`
}

func makeRepayEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(repayRequest)
		err := s.Repay(ctx, req.CardId, req.Amount)
		return encode.Response{Err: err}, err
	}
}
