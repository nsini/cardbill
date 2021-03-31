/**
 * @Time : 3/30/21 3:05 PM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package mp

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/nsini/cardbill/src/encode"
	"github.com/nsini/cardbill/src/middleware"
	"time"
)

type (
	recentRepayRequest struct {
		recent int
	}
	recentRepayResult struct {
		CardName     string    `json:"cardName"`     // 卡名
		BankName     string    `json:"bankName"`     // 银行名称
		BankAvatar   string    `json:"bankAvatar"`   // 银行头像
		Amount       float64   `son:"amount"`        // 还款金额
		RepaymentDay time.Time `json:"repaymentDay"` // 最后还款日
		TailNumber   int64     `json:"tailNumber"`   // 卡片尾号
	}
	bankResult struct {
		BankName   string `json:"bankName"`
		BankAvatar string `json:"bankAvatar"`
	}
	bankRequest struct {
		bankName string
	}

	loginResult struct {
		Token      string `json:"token"`
		OpenId     string `json:"openId"`
		SessionKey string `json:"sessionKey"`
		UnionId    string `json:"unionId"`
		Avatar     string `json:"avatar"`
		Nickname   string `json:"nickname"`
		ShareCode  string `json:"shareCode"`
	}

	mpLoginRequest struct {
		Code          string `json:"code"`
		Iv            string `json:"iv"`
		RawData       string `json:"rawData"`
		Signature     string `json:"signature"`
		EncryptedData string `json:"encryptedData"`
		InviteCode    string `json:"inviteCode"`
	}
)

type Endpoints struct {
	RecentRepayEndpoint endpoint.Endpoint
	BankListEndpoint    endpoint.Endpoint
	LoginEndpoint       endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		RecentRepayEndpoint: makeRecentRepayEndpoint(s),
		BankListEndpoint:    makeBankListEndpoint(s),
		LoginEndpoint:       makeLoginEndpoint(s),
	}

	for _, m := range dmw["RecentRepay"] {
		eps.RecentRepayEndpoint = m(eps.RecentRepayEndpoint)
	}
	for _, m := range dmw["BankList"] {
		eps.BankListEndpoint = m(eps.BankListEndpoint)
	}
	return eps
}

func makeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(mpLoginRequest)
		res, err := s.Login(ctx, req.Code, req.Iv, req.RawData, req.Signature, req.EncryptedData, req.InviteCode)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeBankListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bankRequest)
		res, total, err := s.BankList(ctx, req.bankName)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeRecentRepayEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recentRepayRequest)
		res, err := s.RecentRepay(ctx, userId, req.recent)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
