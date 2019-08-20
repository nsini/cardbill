/**
 * @Time: 2019-08-18 00:30
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package record

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/util/encode"
)

type tmePostRequest struct {
	CardId       string `json:"card_id"`
	BusinessType string `json:"business_type"`
	BusinessName string `json:"business_name"`
	Rate         string `json:"rate"`
	Amount       string `json:"amount"`
}

type postRequest struct {
	CardId       int64   `json:"card_id"`
	BusinessType int64   `json:"business_type"`
	BusinessName string  `json:"business_name"`
	Rate         float64 `json:"rate"`
	Amount       float64 `json:"amount"`
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.CardId, req.BusinessType, req.BusinessName, req.Rate, req.Amount)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.List(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}
