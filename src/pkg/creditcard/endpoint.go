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
	Id          int64   `json:"id"`
	CardName    string  `json:"card_name"`
	BankId      int64   `json:"bank_id"`
	FixedAmount float64 `json:"fixed_amount"`
	MaxAmount   float64 `json:"max_amount"`
	BillingDay  int     `json:"billing_day"`
	Cardholder  int     `json:"cardholder"`
	Sate        int     `json:"sate"`
}

type listRequest struct {
	BankId int64 `json:"bank_id"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.CardName, req.BankId, req.FixedAmount, req.MaxAmount, req.BillingDay, req.Cardholder)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		res, err := s.List(ctx, req.BankId)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makePutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postRequest)
		err = s.Put(ctx, req.Id, req.CardName, req.BankId, req.FixedAmount, req.MaxAmount, req.BillingDay, req.Cardholder, req.Sate)
		return encode.Response{Err: err}, err
	}
}
