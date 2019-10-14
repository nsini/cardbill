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

type ResData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func makeLastAmountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.LastAmount(ctx)
		var resp []ResData
		if err == nil && res != nil {
			for _, v := range res {
				resp = append(resp, ResData{
					Date:   v.Date.Format("2006-01-02"),
					Amount: v.Amount,
				})
			}
		}
		return encode.Response{Err: err, Data: resp}, err
	}
}

func makeMonthAmountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.MonthAmount(ctx)
		var resp []ResData
		if err == nil && res != nil {
			for _, v := range res {
				resp = append(resp, ResData{
					Date:   v.Date.Format("2006-01"),
					Amount: v.Amount,
				})
			}
		}
		return encode.Response{Err: err, Data: resp}, err
	}
}
