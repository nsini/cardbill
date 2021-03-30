/**
 * @Time : 3/30/21 3:05 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package mp

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/repository"
	"time"
)

type Service interface {
	// 最近一周要还款的卡
	RecentRepay(ctx context.Context, userId int64, recent int) (res []recentRepayResult, err error)

	// 添加银行
	// bankName: 银行名称
	AddBank(ctx context.Context, bankName string) (err error)

	// 添加信用卡
	// userId: 用户ID
	// cardName: 卡笥名称
	// bankId: 银行ID
	// fixedAmount: 固定额
	// maxAmount: 最大金额
	// billingDay: 账单日
	// cardHolder: 每月几号或账单日后几天
	// holderType: 还款类型 0每月几号 1账单日后多少天
	// tailNumber: 卡片后四位
	AddCreditCard(ctx context.Context, userId int64, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder int, holderType int, tailNumber int64) (err error)

	// 银行列表
	// bankName: 银行名称
	BankList(ctx context.Context, bankName string) (res []bankResult, total int, err error)
}

type service struct {
	logger     log.Logger
	traceId    string
	repository repository.Repository
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
			BankName:   v.BankName,
			BankAvatar: fmt.Sprintf("./icons/banks/%s@3x.png", v.BankName),
		})
	}

	return
}

func (s *service) AddBank(ctx context.Context, bankName string) (err error) {
	panic("implement me")
}

func (s *service) AddCreditCard(ctx context.Context, userId int64, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder int, holderType int, tailNumber int64) (err error) {
	panic("implement me")
}

func (s *service) RecentRepay(ctx context.Context, userId int64, recent int) (res []recentRepayResult, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "RecentRepay")
	cards, err := s.repository.CreditCard().FindByUserId(userId, 0, -1)
	if err != nil {
		_ = level.Error(logger).Log("repository.CreditCard", "FindByUserId", "err", err.Error())
		return
	}

	var cardIds []int64
	for _, card := range cards {
		cardIds = append(cardIds, card.Id)
	}

	now := time.Now()
	t := now.AddDate(0, 0, +recent)

	list, err := s.repository.Bill().LastBill(cardIds, 10, &t)
	if err != nil {
		_ = level.Error(logger).Log("repository.Bill", "LastBill", "err", err.Error())
		return
	}

	for _, v := range list {
		res = append(res, recentRepayResult{
			CardName:     v.CreditCard.CardName,
			BankName:     v.CreditCard.Bank.BankName,
			BankAvatar:   fmt.Sprintf("./icons/banks/%s@3x.png", v.CreditCard.Bank.BankName),
			Amount:       v.Amount,
			RepaymentDay: v.RepaymentDay,
			TailNumber:   v.CreditCard.TailNumber,
		})
	}

	return
}

func New(logger log.Logger, traceId string, repository repository.Repository) Service {
	logger = log.With(logger, "mp", "service")
	return &service{
		logger:     logger,
		traceId:    traceId,
		repository: repository,
	}
}
