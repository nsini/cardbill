/**
 * @Time : 2019-08-19 11:15
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package creditcard

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type postRequest struct {
	CardName    string  `json:"card_name"`
	BankId      int64   `json:"bank_id"`
	FixedAmount float64 `json:"fixed_amount"`
	MaxAmount   float64 `json:"max_amount"`
	BillingDay  int     `json:"billing_day"`
	CardHolder  int     `json:"card_holder"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.CardName, req.BankId, req.FixedAmount, req.MaxAmount, req.BillingDay, req.CardHolder)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.List(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return encode.Response{Err: err}
	}
}
