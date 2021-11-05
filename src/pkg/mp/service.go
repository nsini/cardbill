/**
 * @Time : 3/30/21 3:05 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package mp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/encode"
	jwt2 "github.com/nsini/cardbill/src/jwt"
	"github.com/nsini/cardbill/src/pkg/wechat"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/card"
	"github.com/nsini/cardbill/src/repository/cardbill"
	"github.com/nsini/cardbill/src/repository/record"
	"github.com/nsini/cardbill/src/repository/types"
	"github.com/nsini/cardbill/src/util/transform"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	// MakeToken 生成TOKEN
	MakeToken(ctx context.Context, appKey string) (token string, err error)
	// UserInfo 用户基本信息
	UserInfo(ctx context.Context, userId int64) (res userInfo, err error)
	// Login 微信小程序授权登录
	Login(ctx context.Context, code, iv, rawData, signature, encryptedData, inviteCode string) (res loginResult, err error)
	// RecentRepay 最近一周要还款的卡
	RecentRepay(ctx context.Context, userId int64, recent int) (res []recentRepayResult, err error)
	// RecentRepayCount 统计最近要还的款数量
	RecentRepayCount(ctx context.Context, userId int64, recent int) (res int, err error)
	// BillDetail 账单详情
	BillDetail(ctx context.Context, userId, billId int64) (res billResult, err error)
	// BillRepay 还款
	BillRepay(ctx context.Context, userId, billId int64) (err error)
	// CardBill 卡账单
	CardBill(ctx context.Context, userId, cardId int64) (res []billResult, err error)
	// AddCreditCard 添加信用卡
	// userId: 用户ID, cardName: 卡名称, bankId: 银行ID
	// fixedAmount: 固定额, maxAmount: 最大金额
	// billingDay: 账单日, cardHolder: 每月几号或账单日后几天
	// holderType: 还款类型 0每月几号 1账单日后多少天
	// tailNumber: 卡片后四位
	AddCreditCard(ctx context.Context, userId int64, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder int, holderType int, tailNumber int64) (err error)
	// CreditCards 信用卡列表
	CreditCards(ctx context.Context, userId int64) (res []cardsResult, err error)
	// CreditCard 信用卡详情
	CreditCard(ctx context.Context, userId, cardId int64) (res cardResult, err error)
	// BankList 银行列表
	// bankName: 银行名称
	BankList(ctx context.Context, bankName string) (res []bankResult, total int, err error)
	// Record 刷卡记录
	Record(ctx context.Context, userId int64, bankId, cardId int64, start, end *time.Time, page, pageSize int) (res []recordResult, total int, err error)
	// RecordAdd 添加刷卡记录
	RecordAdd(ctx context.Context, userId, cardId int64, amount, rate float64, businessType int64, businessName string, swipeTime *time.Time) (err error)
	// RecordDetail 记录详情
	RecordDetail(ctx context.Context, userId, recordId int64) (res recordDetailResult, err error)
	// BusinessTypes 商户类开列表
	BusinessTypes(ctx context.Context) (res []businessTypesResult, err error)
	// Statistics 统计数据
	Statistics(ctx context.Context, userId int64) (res statisticResult, err error)
	// CreditCardNames 获取信用卡名称
	CreditCardNames(ctx context.Context, bankId int64) (res []cardsResult, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
	wechat     wechat.Service
	host       string
}

func (s *service) CardBill(ctx context.Context, userId, cardId int64) (res []billResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CardBill")

	list, err := s.repository.CardBill().LastBill(ctx, []int64{cardId}, 12, nil)
	if err != nil {
		_ = level.Error(logger).Log("repository.CardBill", "LastBill", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, billResult{
			Id:           v.Id,
			CardName:     v.CreditCard.CardName,
			Amount:       v.Amount,
			IsRepay:      v.IsRepay,
			RepayTime:    v.RepayTime,
			RepaymentDay: v.RepaymentDay,
			CreatedAt:    v.CreatedAt,
			Records:      nil,
		})
	}

	return
}

func (s *service) CreditCard(ctx context.Context, userId, cardId int64) (res cardResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreditCard")

	cd, err := s.repository.Card().FindById(ctx, userId, cardId)
	if err != nil {
		_ = level.Error(logger).Log("repository.Card", "FindById", "err", err.Error())
		return
	}

	return cardResult{
		CardName:         cd.CardName,
		BankName:         cd.Bank.BankName,
		BankAvatar:       fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, cd.Bank.BankName),
		CardAvatar:       fmt.Sprintf("%s/icons/cards/%s-%s.png", s.host, cd.Bank.BankName, cd.CardName),
		TailNumber:       cd.TailNumber,
		Amount:           cd.FixedAmount,
		BillingDay:       cd.BillingDay,
		Cardholder:       cd.Cardholder,
		CreatedAt:        cd.CreatedAt,
		TotalConsumption: 0,
	}, nil
}

func (s *service) CreditCardNames(ctx context.Context, bankId int64) (res []cardsResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreditCardNames")

	bankInfo, err := s.repository.ChinaBank().Find(ctx, bankId)
	if err != nil {
		_ = level.Error(logger).Log("repository.ChinaBank", "Find", "err", err.Error())
		return
	}

	list, err := s.repository.Card().FindByBankId(ctx, bankId)
	if err != nil {
		_ = level.Error(logger).Log("repository.Card", "FindByUserId", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, cardsResult{
			Id:         v.Id,
			CardName:   v.CardName,
			BankName:   v.Bank.BankName,
			BankAvatar: fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, bankInfo.BankName),
			TailNumber: v.TailNumber,
			CardAvatar: fmt.Sprintf("%s/icons/cards/%s-%s.png", s.host, bankInfo.BankName, v.CardName),
		})
	}

	return
}

func (s *service) BillRepay(ctx context.Context, userId, billId int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "BillRepay")

	bill, err := s.repository.CardBill().FindById(ctx, billId)
	if err != nil {
		_ = level.Error(logger).Log("repository.CardBill", "FindById", "err", err.Error())
		return
	}

	if bill.CreditCard.UserId != userId {
		_ = level.Warn(logger).Log("userId", userId, "CreditCard.UserId", bill.CreditCard.UserId)
		return encode.ErrMpNotPermission.Error()
	}

	t := time.Now()
	bill.RepayTime = &t
	bill.IsRepay = true
	bill.CreditCard = types.CreditCard{}

	return s.repository.CardBill().Save(ctx, &bill)
}

func (s *service) BillDetail(ctx context.Context, userId, billId int64) (res billResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "BillDetail")
	bill, err := s.repository.CardBill().FindById(ctx, billId)
	if err != nil {
		_ = level.Error(logger).Log("repository.CardBill", "FindById", "err", err.Error())
		return
	}

	if userId != bill.CreditCard.UserId {
		return
	}

	records, err := s.repository.Record().FindByBill(ctx, userId, bill.CreditCard.Id, bill.CreatedAt.AddDate(0, -1, 0), bill.CreatedAt)
	if err != nil {
		_ = level.Error(logger).Log("repository.Record", "FindByBill", "err", err.Error())
		return
	}

	res.CardName = bill.CreditCard.CardName
	res.BankName = bill.CreditCard.Bank.BankName
	res.CardAvatar = fmt.Sprintf("%s/icons/cards/%s-%s.png", s.host, bill.CreditCard.Bank.BankName, bill.CreditCard.CardName)
	res.TailNumber = bill.CreditCard.TailNumber
	res.BankAvatar = fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, bill.CreditCard.Bank.BankName)
	res.Amount = transform.Decimal(bill.Amount)
	res.IsRepay = bill.IsRepay
	res.RepayTime = bill.RepayTime
	res.RepaymentDay = bill.RepaymentDay

	var list []recordResult

	for _, v := range records {
		list = append(list, recordResult{
			Id:           v.Id,
			Amount:       v.Amount,
			CreatedAt:    v.CreatedAt,
			BusinessType: v.Business.BusinessName,
			BusinessName: v.BusinessName,
			BusinessCode: v.Business.Code,
			Rate:         transform.Decimal(v.Rate * 100),
			Arrival:      transform.Decimal(v.Arrival),
		})
	}
	res.Records = list
	return
}

func (s *service) RecentRepayCount(ctx context.Context, userId int64, recent int) (res int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "RecentRepayCount")
	cards, err := s.repository.Card().FindByUserId(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.CreditCard", "FindByUserId", "err", err.Error())
		return
	}

	var cardIds []int64
	for _, cd := range cards {
		cardIds = append(cardIds, cd.Id)
	}

	now := time.Now()
	t := now.AddDate(0, 0, +recent)

	res, err = s.repository.CardBill().CountLastBill(ctx, cardIds, recent, &t)
	if err != nil {
		_ = level.Error(logger).Log("repository.Bill", "CountLastBill", "err", err.Error())
		return
	}
	return
}

func (s *service) UserInfo(ctx context.Context, userId int64) (res userInfo, err error) {
	panic("implement me")
}

func (s *service) RecordDetail(ctx context.Context, userId, recordId int64) (res recordDetailResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Statistics")
	rd, err := s.repository.Record().FindById(ctx, userId, recordId)
	if err != nil {
		_ = level.Error(logger).Log("repository.Record", "FindById", "err", err.Error())
		return
	}

	return recordDetailResult{
		CardId:       rd.CardId,
		BusinessType: rd.Business.Code,
		BusinessName: rd.Business.BusinessName,
		Rate:         transform.Decimal(rd.Rate * 100),
		Amount:       transform.Decimal(rd.Amount),
		Arrival:      transform.Decimal(rd.Arrival),
		CreatedAt:    rd.CreatedAt,
		CardName:     rd.CreditCard.CardName,
		BankName:     rd.CreditCard.Bank.BankName,
		Merchant:     rd.BusinessName,
		BankAvatar:   fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, rd.CreditCard.Bank.BankName),
		CardAvatar:   fmt.Sprintf("%s/icons/cards/%s-%s.png", s.host, rd.CreditCard.Bank.BankName, rd.CreditCard.CardName),
		TailNumber:   rd.CreditCard.TailNumber,
	}, nil
}

func (s *service) Statistics(ctx context.Context, userId int64) (res statisticResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Statistics")

	var cardIds []int64
	cards, err := s.repository.Card().FindAllByUserId(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("Card", "Count", "FindByUserId", err.Error())
		return
	}
	for _, v := range cards {
		cardIds = append(cardIds, v.Id)
	}

	currentMonth := time.Now()

	creditTotalCh := make(chan int)
	creditAmountCh := make(chan card.TotalAmount)
	sacCh := make(chan record.RemainingAmount)
	currSacCh := make(chan record.RemainingAmount)
	unpaidBillCh := make(chan cardbill.BillAmount)

	go func() {
		// 信用卡数量
		total, err := s.repository.Card().Count(ctx, userId)
		if err != nil {
			_ = level.Error(logger).Log("Card", "Count", "err", err.Error())
		}
		creditTotalCh <- total
	}()

	go func() {
		// 信用卡总额度
		creditAmount, err := s.repository.Card().Sum(ctx, userId, 0)
		if err != nil {
			_ = level.Error(logger).Log("Card", "Sum", "err", err.Error())
		}
		creditAmountCh <- creditAmount
	}()

	go func() {
		// 总消费
		sac, err := s.repository.Record().SumAmountCards(ctx, cardIds, nil)
		if err != nil {
			_ = level.Error(logger).Log("Record", "SumAmountCards", "err", err.Error())
		} else {
			sacCh <- sac
		}
	}()

	go func() {
		// 当月消费
		currSac, err := s.repository.Record().SumAmountCards(ctx, cardIds, &currentMonth)
		if err != nil {
			_ = level.Error(logger).Log("Record", "SumAmountCards", "err", err.Error())
		}
		currSacCh <- currSac
	}()

	go func() {
		// 账单
		unpaidBill, err := s.repository.CardBill().SumByCards(ctx, cardIds, nil, cardbill.RepayFalse)
		if err != nil {
			_ = level.Error(logger).Log("Bill", "SumByCards", "err", err.Error())
		}
		unpaidBillCh <- unpaidBill
	}()

	totalAmount := <-creditAmountCh
	cardNumber := <-creditTotalCh
	sac := <-sacCh
	currSac := <-currSacCh
	unpaidBill := <-unpaidBillCh

	close(creditTotalCh)
	close(creditAmountCh)
	close(sacCh)
	close(currSacCh)
	close(unpaidBillCh)

	interestExpense, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(2)+"f", sac.Amount-sac.Arrival), 64)
	currentInterest, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(2)+"f", currSac.Amount-currSac.Arrival), 64)

	return statisticResult{
		CreditAmount:       totalAmount.Amount,
		CreditMaxAmount:    totalAmount.MaxAmount,
		CreditNumber:       cardNumber,
		TotalConsumption:   sac.Amount,
		MonthlyConsumption: currSac.Amount,
		InterestExpense:    interestExpense,
		CurrentInterest:    currentInterest,
		UnpaidBill:         unpaidBill.Amount,
		RepaidBill:         0,
	}, nil

}

func (s *service) BusinessTypes(ctx context.Context) (res []businessTypesResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "BusinessTypes")
	list, err := s.repository.BusinessSvc().Types(ctx)
	if err != nil {
		_ = level.Error(logger).Log("repository.BusinessSvc", "Types", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, businessTypesResult{
			Code: v.Code,
			Name: v.BusinessName,
		})
	}
	return
}

func (s *service) RecordAdd(ctx context.Context, userId, cardId int64, amount, rate float64, businessType int64, businessName string, swipeTime *time.Time) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "RecordAdd")
	crd, err := s.repository.Card().FindById(ctx, userId, cardId)
	if err != nil {
		_ = level.Warn(logger).Log("Card", "FindById", "err", err.Error())
		return
	}

	business, err := s.repository.BusinessSvc().FindByCode(ctx, businessType)
	if err != nil {
		_ = level.Warn(logger).Log("BusinessSvc", "FindByCode", "err", err.Error())
		return
	}

	if swipeTime == nil {
		t := time.Now()
		swipeTime = &t
	}

	if err := s.repository.Record().Save(ctx, &types.ExpensesRecord{
		CardId:       crd.Id,
		BusinessType: business.Id,
		BusinessName: businessName,
		Rate:         rate,
		Amount:       amount,
		Arrival:      amount - transform.Decimal(amount*rate),
		UserId:       userId,
		CreatedAt:    *swipeTime,
	}); err != nil {
		_ = level.Error(logger).Log("ExpenseRecord", "Create", "err", err.Error())
		return encode.ErrMpRecordAdd.Error()
	}

	go func() {
		if err = s.repository.BusinessSvc().SaveMerchant(ctx, &types.Merchant{
			MerchantName: businessName,
			BusinessId:   business.Id,
			Business:     business,
		}); err != nil {
			_ = level.Warn(logger).Log("BusinessSvc", "SaveMerchant", "err", err.Error())
		}
	}()

	return
}

func (s *service) CreditCards(ctx context.Context, userId int64) (res []cardsResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Record")
	list, err := s.repository.Card().FindByUserId(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.Card", "FindByUserId", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, cardsResult{
			Id:         v.Id,
			CardName:   v.CardName,
			BankName:   v.Bank.BankName,
			BankAvatar: fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, v.Bank.BankName),
			TailNumber: v.TailNumber,
			CardAvatar: fmt.Sprintf("%s/icons/cards/%s-%s.png", s.host, v.Bank.BankName, v.CardName),
		})
	}
	return
}

func (s *service) Record(ctx context.Context, userId int64, bankId, cardId int64, start, end *time.Time, page, pageSize int) (res []recordResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Record")
	var ids []int64
	ids = append(ids, cardId)

	if cardId < 1 {
		if cards, err := s.repository.CreditCard().FindByUserId(userId, bankId, -1); err == nil {
			for _, v := range cards {
				ids = append(ids, v.Id)
			}
		}
	}
	list, total, err := s.repository.Record().List(ctx, userId, page, pageSize, bankId, ids, start, end)
	if err != nil {
		_ = level.Error(logger).Log("repository.Record", "List", "err", err.Error())
		return
	}
	for _, v := range list {
		res = append(res, recordResult{
			CardAvatar:   "",
			Id:           v.Id,
			CardName:     v.CreditCard.CardName,
			BankName:     v.CreditCard.Bank.BankName,
			BankAvatar:   fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, v.CreditCard.Bank.BankName),
			Amount:       v.Amount,
			TailNumber:   v.CreditCard.TailNumber,
			CreatedAt:    v.CreatedAt,
			BusinessType: v.Business.BusinessName,
			BusinessName: v.BusinessName,
			BusinessCode: v.Business.Code,
			Rate:         v.Rate,
			Arrival:      v.Arrival,
		})
	}
	return
}

func (s *service) MakeToken(ctx context.Context, appKey string) (token string, err error) {
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "MakeToken")
	token = appKey
	return
}

type userInfo struct {
	AvatarURL string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	NickName  string `json:"nickName"`
	Province  string `json:"province"`
}

func (s *service) Login(ctx context.Context, code, iv, rawData, signature, encryptedData, inviteCode string) (res loginResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "Login")
	var reqUserInfo userInfo
	if err = json.NewDecoder(bytes.NewBufferString(rawData)).Decode(&reqUserInfo); err != nil {
		_ = level.Warn(logger).Log("json", "NewDecoder", "userInfo", "Decode", "err", err.Error())
	}

	// todo: 校验 signature, encryptedData

	userInfo, sessionKey, err := s.wechat.MPLogin(ctx, code)
	if err != nil {
		_ = level.Error(logger).Log("wechat", "MPLogin", "err", err.Error())
		return
	}
	var user types.User
	openId := userInfo.UnionId
	if strings.EqualFold(openId, "") {
		openId = userInfo.OpenId
	}
	if user, err = s.repository.Users().FindByUnionId(ctx, openId); err != nil {
		if err != gorm.ErrRecordNotFound {
			_ = level.Error(logger).Log("gorm", "ErrRecordNotFound", "err", err.Error())
			err = encode.ErrAuthMPLogin.Error()
			return res, err
		}
		u := &types.User{
			OpenId:   userInfo.OpenId,
			UnionId:  userInfo.UnionId,
			Nickname: userInfo.Nickname,
			Username: userInfo.Nickname,
			Sex:      reqUserInfo.Gender,
			City:     reqUserInfo.City,
			Province: reqUserInfo.Province,
			Country:  reqUserInfo.Country,
			Avatar:   reqUserInfo.AvatarURL,
			Remark:   "小程序登录",
		}

		if err = s.repository.Users().Save(ctx, u); err != nil {
			_ = level.Error(logger).Log("repository.User", "Save", "err", err.Error())
			err = encode.ErrAuthMPLogin.Error()
			return
		}
		user = *u
		_ = level.Info(logger).Log("repository.User", "FindByUnionId", "msg", "用户不存在,保存信息")
	}

	defer func() {
		user.Nickname = reqUserInfo.NickName
		user.Sex = reqUserInfo.Gender
		user.City = reqUserInfo.City
		user.Province = reqUserInfo.Province
		user.Country = reqUserInfo.Country
		user.Avatar = reqUserInfo.AvatarURL
		_ = s.repository.Users().Save(ctx, &user)
	}()

	sessionTimeout := 3600 * 24 * 31 * 12 * time.Second

	expAt := time.Now().Add(sessionTimeout).Unix()

	claims := jwt2.ArithmeticCustomClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(jwt2.GetJwtKey()))
	if err != nil {
		_ = level.Error(logger).Log("token", "SignedString", "err", err.Error())
	}

	//_ = s.cache.Set(ctx, fmt.Sprintf("login:%d:token", user.Id), tk, sessionTimeout)

	return loginResult{
		Token:      tk,
		SessionKey: sessionKey,
		Avatar:     reqUserInfo.AvatarURL,
		Nickname:   reqUserInfo.NickName,
	}, nil
}

func (s *service) BankList(ctx context.Context, bankName string) (res []bankResult, total int, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "BankList")

	list, total, err := s.repository.ChinaBank().List(ctx, bankName)
	if err != nil {
		_ = level.Error(logger).Log("repository.ChinaBank", "List", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, bankResult{
			Id:         v.Id,
			BankName:   v.BankName,
			BankAvatar: fmt.Sprintf("./icons/banks/%s@3x.png", v.BankName),
		})
	}

	return
}

func (s *service) AddCreditCard(ctx context.Context, userId int64, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder int, holderType int, tailNumber int64) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "AddCreditCard")
	if err = s.repository.Card().Save(ctx, &types.CreditCard{
		CardName:    cardName,
		BankId:      bankId,
		FixedAmount: fixedAmount,
		MaxAmount:   maxAmount,
		BillingDay:  billingDay,
		Cardholder:  cardHolder,
		UserId:      userId,
		TailNumber:  tailNumber,
	}); err != nil {
		_ = level.Error(logger).Log("repository.CreditCard", "Save", "err", err.Error())
	}
	return
}

func (s *service) RecentRepay(ctx context.Context, userId int64, recent int) (res []recentRepayResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "RecentRepay")
	cards, err := s.repository.Card().FindByUserId(ctx, userId)
	if err != nil {
		_ = level.Error(logger).Log("repository.CreditCard", "FindByUserId", "err", err.Error())
		return
	}

	var cardIds []int64
	for _, cd := range cards {
		cardIds = append(cardIds, cd.Id)
	}

	now := time.Now()
	t := now.AddDate(0, 0, +recent)

	list, err := s.repository.CardBill().LastBill(ctx, cardIds, recent, &t)
	if err != nil {
		_ = level.Error(logger).Log("repository.Bill", "LastBill", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, recentRepayResult{
			Id:           v.Id,
			CardName:     v.CreditCard.CardName,
			BankName:     v.CreditCard.Bank.BankName,
			BankAvatar:   fmt.Sprintf("%s/icons/banks/%s@3x.png", s.host, v.CreditCard.Bank.BankName),
			Amount:       v.Amount,
			RepaymentDay: v.RepaymentDay,
			TailNumber:   v.CreditCard.TailNumber,
		})
	}

	return
}

func New(logger log.Logger, traceId, host string, repository repository.Repository, wechatSvc wechat.Service) Service {
	logger = log.With(logger, "mp", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
		host:       host,
		wechat:     wechatSvc,
	}
}
