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
	"time"
)

type repayRequest struct {
	CardId       int64      `json:"card_id"`
	Amount       float64    `json:"amount"`
	Repayment    string     `json:"repayment"`
	RepaymentDay *time.Time `json:"repayment_day"`
}

type listRequest struct {
	page, pageSize int
}

func makeRepayEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(repayRequest)
		err := s.Repay(ctx, req.CardId, req.Amount, req.RepaymentDay)
		return encode.Response{Err: err}, err
	}
}

func makeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listRequest)
		list, count, err := s.List(ctx, req.page, req.pageSize)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"count": count,
			"list":  list,
		}}, err
	}
}
