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
		Id           int64     `json:"id"`
		CardName     string    `json:"cardName"`     // 卡名
		BankName     string    `json:"bankName"`     // 银行名称
		BankAvatar   string    `json:"bankAvatar"`   // 银行头像
		Amount       float64   `json:"amount"`       // 还款金额
		RepaymentDay time.Time `json:"repaymentDay"` // 最后还款日
		TailNumber   int64     `json:"tailNumber"`   // 卡片尾号
	}
	bankResult struct {
		Id         int64  `json:"id"`
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
		CardAvatar   string    `json:"cardAvatar,omitempty"`
		Id           int64     `json:"id,omitempty"`
		CardName     string    `json:"cardName,omitempty"`     // 卡名
		BankName     string    `json:"bankName,omitempty"`     // 银行名称
		BankAvatar   string    `json:"bankAvatar,omitempty"`   // 银行头像
		Amount       float64   `json:"amount,omitempty"`       // 金额
		TailNumber   int64     `json:"tailNumber,omitempty"`   // 卡片尾号
		CreatedAt    time.Time `json:"createdAt,omitempty"`    // 创建时间
		BusinessType string    `json:"businessType,omitempty"` // 渠道类型
		BusinessName string    `json:"businessName,omitempty"` // 渠道名称
		BusinessCode int64     `json:"businessCode,omitempty"` // 渠道号
		Rate         float64   `json:"rate,omitempty"`         // 费率
		Arrival      float64   `json:"arrival,omitempty"`      // 到账金额
	}

	cardsResult struct {
		Id         int64  `json:"id"`
		CardName   string `json:"cardName"`   // 卡名
		BankName   string `json:"bankName"`   // 银行名称
		BankAvatar string `json:"bankAvatar"` // 银行头像
		CardAvatar string `json:"cardAvatar"`
		TailNumber int64  `json:"tailNumber"` // 卡片尾号
	}

	recordAddRequest struct {
		CardId       int64   `json:"cardId"`
		BusinessType int64   `json:"businessType"`
		BusinessName string  `json:"businessName"`
		Rate         float64 `json:"rate"`
		Amount       float64 `json:"amount"`
		TmpTime      string  `json:"swipeTime"`
		SwipeTime    *time.Time
	}

	businessTypesResult struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
	}

	statisticResult struct {
		CreditAmount       float64 `json:"creditAmount"`
		CreditMaxAmount    float64 `json:"creditMaxAmount"`
		CreditNumber       int     `json:"creditNumber"`
		TotalConsumption   float64 `json:"totalConsumption"`
		MonthlyConsumption float64 `json:"monthlyConsumption"`
		InterestExpense    float64 `json:"interestExpense"`
		CurrentInterest    float64 `json:"currentInterest"`
		UnpaidBill         float64 `json:"unpaidBill"`
		RepaidBill         float64 `json:"repaidBill"`
	}

	recordDetailRequest struct {
		Id int64
	}
	recordDetailResult struct {
		CardId       int64     `json:"cardId"`       // 你的银行卡id
		BusinessType int64     `json:"businessType"` // 商户类型Code 对应businesses表
		BusinessName string    `json:"businessName"` // 商户名称 对应用merchant的名称
		Rate         float64   `json:"rate"`         // 费率
		Amount       float64   `json:"amount"`       // 消费金额
		Arrival      float64   `json:"arrival"`      // 实际到账
		CreatedAt    time.Time `json:"createdAt"`
		CardName     string    `json:"cardName"`
		BankName     string    `json:"bankName"`
		Merchant     string    `json:"merchant"`   // 商户名称
		BankAvatar   string    `json:"bankAvatar"` // 银行头像
		TailNumber   int64     `json:"tailNumber"` // 卡片尾号
		CardAvatar   string    `json:"cardAvatar"`
	}

	billResult struct {
		Id           int64      `json:"id"`
		CardName     string     `json:"cardName"`   // 卡名
		BankName     string     `json:"bankName"`   // 银行名称
		BankAvatar   string     `json:"bankAvatar"` // 银行头像
		CardAvatar   string     `json:"cardAvatar"`
		TailNumber   int64      `json:"tailNumber"` // 卡片尾号
		Amount       float64    `json:"amount"`
		IsRepay      bool       `json:"isRepay"`
		RepayTime    *time.Time `json:"repayTime"`
		RepaymentDay time.Time  `json:"repaymentDay"`
		CreatedAt    time.Time  `json:"createdAt"`

		Records []recordResult `json:"records"`
	}

	creditCardRequest struct {
		BankId      int64   `json:"bankId"`
		CardName    string  `json:"cardName"`
		TailNumber  int64   `json:"tailNumber"`
		FixedAmount float64 `json:"fixedAmount"`
		BillingDay  int     `json:"billingDay"`
		Cardholder  int     `json:"cardholder"`
	}

	cardResult struct {
		CardName         string    `json:"cardName"`         // 卡名
		BankName         string    `json:"bankName"`         // 银行名称
		BankAvatar       string    `json:"bankAvatar"`       // 银行头像
		CardAvatar       string    `json:"cardAvatar"`       // 卡图片
		TailNumber       int64     `json:"tailNumber"`       // 卡片尾号
		Amount           float64   `json:"amount"`           // 固定金额
		BillingDay       int       `json:"billingDay"`       // 账单日
		Cardholder       int       `json:"cardholder"`       // 还款日
		CreatedAt        time.Time `json:"createdAt"`        // 创建时间
		TotalConsumption float64   `json:"totalConsumption"` // 总消费
	}
)

type Endpoints struct {
	RecentRepayEndpoint      endpoint.Endpoint
	RecentRepayCountEndpoint endpoint.Endpoint
	BankListEndpoint         endpoint.Endpoint
	LoginEndpoint            endpoint.Endpoint
	MakeTokenEndpoint        endpoint.Endpoint
	RecordEndpoint           endpoint.Endpoint
	RecordAddEndpoint        endpoint.Endpoint
	CreditCardsEndpoint      endpoint.Endpoint
	BusinessTypesEndpoint    endpoint.Endpoint
	StatisticsEndpoint       endpoint.Endpoint
	RecordDetailEndpoint     endpoint.Endpoint
	BillDetailEndpoint       endpoint.Endpoint
	BillRepayEndpoint        endpoint.Endpoint
	CreditCardNamesEndpoint  endpoint.Endpoint
	CreditCardAddEndpoint    endpoint.Endpoint
	CreditCardEndpoint       endpoint.Endpoint
	CardBillEndpoint         endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		RecentRepayEndpoint:      makeRecentRepayEndpoint(s),
		BankListEndpoint:         makeBankListEndpoint(s),
		LoginEndpoint:            makeLoginEndpoint(s),
		MakeTokenEndpoint:        makeMakeTokenEndpoint(s),
		RecordEndpoint:           makeRecordEndpoint(s),
		CreditCardsEndpoint:      makeCreditCardsEndpoint(s),
		RecordAddEndpoint:        makeRecordAddEndpoint(s),
		BusinessTypesEndpoint:    makeBusinessTypesEndpoint(s),
		StatisticsEndpoint:       makeStatisticsEndpoint(s),
		RecordDetailEndpoint:     makeRecordDetailEndpoint(s),
		RecentRepayCountEndpoint: makeRecentRepayCountEndpoint(s),
		BillDetailEndpoint:       makeBillDetailEndpoint(s),
		BillRepayEndpoint:        makeBillRepayEndpoint(s),
		CreditCardNamesEndpoint:  makeCreditCardNamesEndpoint(s),
		CreditCardAddEndpoint:    makeCreditCardAddEndpoint(s),
		CreditCardEndpoint:       makeCreditCardEndpoint(s),
		CardBillEndpoint:         makeCardBillEndpoint(s),
	}

	for _, m := range dmw["RecentRepay"] {
		eps.RecentRepayEndpoint = m(eps.RecentRepayEndpoint)
	}
	for _, m := range dmw["RecentRepayCount"] {
		eps.RecentRepayCountEndpoint = m(eps.RecentRepayCountEndpoint)
	}
	for _, m := range dmw["BankList"] {
		eps.BankListEndpoint = m(eps.BankListEndpoint)
	}
	for _, m := range dmw["Record"] {
		eps.RecordEndpoint = m(eps.RecordEndpoint)
	}
	for _, m := range dmw["RecordAdd"] {
		eps.RecordAddEndpoint = m(eps.RecordAddEndpoint)
	}
	for _, m := range dmw["BusinessTypes"] {
		eps.BusinessTypesEndpoint = m(eps.BusinessTypesEndpoint)
	}
	for _, m := range dmw["Statistics"] {
		eps.StatisticsEndpoint = m(eps.StatisticsEndpoint)
	}
	for _, m := range dmw["RecordDetail"] {
		eps.RecordDetailEndpoint = m(eps.RecordDetailEndpoint)
	}
	for _, m := range dmw["BillDetail"] {
		eps.BillDetailEndpoint = m(eps.BillDetailEndpoint)
	}
	for _, m := range dmw["BillRepay"] {
		eps.BillRepayEndpoint = m(eps.BillRepayEndpoint)
	}
	for _, m := range dmw["CreditCardNames"] {
		eps.CreditCardNamesEndpoint = m(eps.CreditCardNamesEndpoint)
	}
	for _, m := range dmw["CreditCardAdd"] {
		eps.CreditCardAddEndpoint = m(eps.CreditCardAddEndpoint)
	}
	for _, m := range dmw["CreditCard"] {
		eps.CreditCardEndpoint = m(eps.CreditCardEndpoint)
	}
	for _, m := range dmw["CardBill"] {
		eps.CardBillEndpoint = m(eps.CardBillEndpoint)
	}
	for _, m := range dmw["CreditCards"] {
		eps.CreditCardsEndpoint = m(eps.CreditCardsEndpoint)
	}
	return eps
}

func makeCardBillEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(recordDetailRequest)
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		res, err := s.CardBill(ctx, userId, req.Id)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCreditCardEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(recordDetailRequest)
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		res, err := s.CreditCard(ctx, userId, req.Id)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeCreditCardAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(creditCardRequest)
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		err = s.AddCreditCard(ctx, userId, req.CardName, req.BankId, req.FixedAmount, req.FixedAmount, req.BillingDay, req.Cardholder, 1, req.TailNumber)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCreditCardNamesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(creditCardRequest)
		res, err := s.CreditCardNames(ctx, req.BankId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeBillRepayEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recordDetailRequest)
		err = s.BillRepay(ctx, userId, req.Id)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeBillDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recordDetailRequest)
		res, err := s.BillDetail(ctx, userId, req.Id)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeRecordDetailEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recordDetailRequest)
		res, err := s.RecordDetail(ctx, userId, req.Id)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeStatisticsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		res, err := s.Statistics(ctx, userId)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeBusinessTypesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		res, err := s.BusinessTypes(ctx)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeRecordAddEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recordAddRequest)
		err = s.RecordAdd(ctx, userId, req.CardId, req.Amount, req.Rate, req.BusinessType, req.BusinessName, req.SwipeTime)
		return encode.Response{
			Error: err,
		}, err
	}
}

func makeCreditCardsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		res, err := s.CreditCards(ctx, userId)
		if res == nil {
			return encode.Response{
				Data:  true,
				Error: err,
			}, err
		}
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}

func makeRecordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recordRequest)
		res, total, err := s.Record(ctx, userId, req.BankId, req.CardId, req.Start, req.End, req.Page, req.PageSize)
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

func makeRecentRepayCountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userId, ok := ctx.Value(middleware.UserIdContext).(int64)
		if !ok {
			err = encode.ErrAuthNotLogin.Error()
			return
		}
		req := request.(recentRepayRequest)
		res, err := s.RecentRepayCount(ctx, userId, req.recent)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
