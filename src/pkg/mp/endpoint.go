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
	"time"
)

type (
	recentRepayRequest struct {
		recent int
	}
	recentRepayResult struct {
		Id           int64     `json:"id"`
		CardName     string    `json:"cardName"`     // 卡名
		BankName     string    `json:"bankName"`     // 银行名称
		BankAvatar   string    `json:"bankAvatar"`   // 银行头像
		Amount       float64   `json:"amount"`       // 还款金额
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

	makeTokenRequest struct {
		AppKey string `json:"appKey"`
	}

	recordRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"pageSize"`
		BankId   int64
		CardId   int64
		Start    *time.Time
		End      *time.Time
	}
	recordResult struct {
		CardAvatar   string    `json:"cardAvatar"`
		Id           int64     `json:"id"`
		CardName     string    `json:"cardName"`     // 卡名
		BankName     string    `json:"bankName"`     // 银行名称
		BankAvatar   string    `json:"bankAvatar"`   // 银行头像
		Amount       float64   `json:"amount"`       // 金额
		TailNumber   int64     `json:"tailNumber"`   // 卡片尾号
		CreatedAt    time.Time `json:"createdAt"`    // 创建时间
		BusinessType string    `json:"businessType"` // 渠道类型
		BusinessName string    `json:"businessName"` // 渠道名称
		BusinessCode int64     `json:"businessCode"` // 渠道号
		Rate         float64   `json:"rate"`         // 费率
		Arrival      float64   `json:"arrival"`      // 到账金额
	}

	cardsResult struct {
		Id         int64  `json:"id"`
		CardName   string `json:"cardName"`   // 卡名
		BankName   string `json:"bankName"`   // 银行名称
		BankAvatar string `json:"bankAvatar"` // 银行头像
		TailNumber int64  `json:"tailNumber"` // 卡片尾号
	}
)

type Endpoints struct {
	RecentRepayEndpoint endpoint.Endpoint
	BankListEndpoint    endpoint.Endpoint
	LoginEndpoint       endpoint.Endpoint
	MakeTokenEndpoint   endpoint.Endpoint
	RecordEndpoint      endpoint.Endpoint
	CreditCardsEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		RecentRepayEndpoint: makeRecentRepayEndpoint(s),
		BankListEndpoint:    makeBankListEndpoint(s),
		LoginEndpoint:       makeLoginEndpoint(s),
		MakeTokenEndpoint:   makeMakeTokenEndpoint(s),
		RecordEndpoint:      makeRecordEndpoint(s),
		CreditCardsEndpoint: makeCreditCardsEndpoint(s),
	}

	for _, m := range dmw["RecentRepay"] {
		eps.RecentRepayEndpoint = m(eps.RecentRepayEndpoint)
	}
	for _, m := range dmw["BankList"] {
		eps.BankListEndpoint = m(eps.BankListEndpoint)
	}
	for _, m := range dmw["Record"] {
		eps.RecordEndpoint = m(eps.RecordEndpoint)
	}
	for _, m := range dmw["CreditCards"] {
		eps.CreditCardsEndpoint = m(eps.CreditCardsEndpoint)
	}
	return eps
}

func makeCreditCardsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		//if !ok {
		//	err = encode.ErrAuthNotLogin.Error()
		//	return
		//}
		res, err := s.CreditCards(ctx, 2)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeRecordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		//if !ok {
		//	err = encode.ErrAuthNotLogin.Error()
		//	return
		//}
		req := request.(recordRequest)
		res, total, err := s.Record(ctx, 2, req.BankId, req.CardId, req.Start, req.End, req.Page, req.PageSize)
		return encode.Response{
			Data: map[string]interface{}{
				"list":  res,
				"total": total,
			},
			Error: err,
		}, err
	}
}

func makeMakeTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(makeTokenRequest)
		res, err := s.MakeToken(ctx, req.AppKey)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
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
		//userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		//if !ok {
		//	err = encode.ErrAuthNotLogin.Error()
		//	return
		//}
		req := request.(recentRepayRequest)
		res, err := s.RecentRepay(ctx, 2, req.recent)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
