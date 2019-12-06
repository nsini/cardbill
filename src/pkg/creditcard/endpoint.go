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

type getRequest struct {
	Id int64
}

type postRequest struct {
	Id          int64   `json:"id"`
	CardName    string  `json:"card_name"`
	BankId      int64   `json:"bank_id"`
	FixedAmount float64 `json:"fixed_amount"`
	MaxAmount   float64 `json:"max_amount"`
	BillingDay  int     `json:"billing_day"`
	Cardholder  int     `json:"cardholder"`
	Sate        int     `json:"sate"`
	CardNumber  int64   `json:"card_number"`
	TailNumber  int64   `json:"tail_number"`
}

type recordRequest struct {
	Id       int64
	Page     int
	PageSize int
}

type listRequest struct {
	BankId int64 `json:"bank_id"`
}

type StatisticsResponse struct {
	CreditAmount       float64 `json:"credit_amount"`       // 总额度
	CreditMaxAmount    float64 `json:"credit_max_amount"`   // 临时额度
	CreditNumber       int     `json:"credit_number"`       // 信用卡数量
	TotalConsumption   float64 `json:"total_consumption"`   // 总消费额度
	MonthlyConsumption float64 `json:"monthly_consumption"` // 当月消费额度
	InterestExpense    float64 `json:"interest_expense"`    // 利息支出
	CurrentInterest    float64 `json:"current_interest"`    // 当月利息出支
	UnpaidBill         float64 `json:"unpaid_bill"`         // 未还账单
	RepaidBill         float64 `json:"repaid_bill"`         // 已还账单
}

func makePostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postRequest)
		err := s.Post(ctx, req.CardName, req.BankId, req.FixedAmount, req.MaxAmount,
			req.BillingDay, req.Cardholder, req.CardNumber, req.TailNumber)
		return encode.Response{Err: err}, err
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		res, err := s.Get(ctx, req.Id)
		return encode.Response{Err: err, Data: res}, err
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
		err = s.Put(ctx, req.Id, req.CardName, req.BankId, req.FixedAmount, req.MaxAmount, req.BillingDay,
			req.Cardholder, req.Sate)
		return encode.Response{Err: err}, err
	}
}

func makeStatisticsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.Statistics(ctx)
		return encode.Response{Err: err, Data: res}, err
	}
}

func makeRecordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(recordRequest)
		res, count, err := s.Record(ctx, req.Id, req.Page, req.PageSize)
		return encode.Response{Err: err, Data: map[string]interface{}{
			"list": res,
			"pagination": map[string]interface{}{
				"total":    count,
				"current":  req.Page,
				"pageSize": req.PageSize,
			},
		}}, err
	}
}
