/**
 * @Time : 2019-09-18 16:45
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package bill

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
	"github.com/nsini/cardbill/src/util/date"
	"time"
)

type Service interface {
	// 生成账单
	GenBill(ctx context.Context, day int) (err error)

	// 还款
	Repay(ctx context.Context, cardId int64, amount float64, repaymentDay *time.Time) (err error)

	// 账单列表
	List(ctx context.Context, page, pageSize int) (res []*types.Bill, count int64, err error)

	// 信用卡账单列表
	ListByCard(ctx context.Context, cardId int64, page, pageSize int) (res []*types.Bill, count int64, err error)

	// 分期
	// id: 账单ID
	// period: 分期数量
	// installmentAmount: 分期金额
	// monthlyRepayment: 月还款金额
	Installment(ctx context.Context, int int64, period int, installmentAmount, monthlyRepayment float64) (err error)

	// 最近一周要还款的卡
	RecentRepay(ctx context.Context, recent int) (res []*types.Bill, err error)
}

var (
	ErrNotPermission = errors.New("您没有权限修改别人的账单")
	ErrNoBill        = errors.New("没有账单")
)

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{logger, repository}
}

func (c *service) RecentRepay(ctx context.Context, recent int) (res []*types.Bill, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	cards, err := c.repository.CreditCard().FindByUserId(userId, 0, -1)
	if err != nil {
		return
	}

	var cardIds []int64

	for _, card := range cards {
		cardIds = append(cardIds, card.Id)
	}

	now := time.Now()

	t := now.AddDate(0, 0, +recent)

	return c.repository.Bill().LastBill(cardIds, 10, &t)
}

func (c *service) Installment(ctx context.Context, int int64, period int, installmentAmount, monthlyRepayment float64) (err error) {
	// todo 得查一下该账单是否属于本人

	return
}

func (c *service) ListByCard(ctx context.Context, cardId int64, page, pageSize int) (res []*types.Bill, count int64, err error) {
	// userId := ctx.Value(middleware.UserIdContext).(int64)

	return c.repository.Bill().FindByCardIds([]int64{cardId}, page, pageSize)
}

func (c *service) List(ctx context.Context, page, pageSize int) (res []*types.Bill, count int64, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	cards, err := c.repository.CreditCard().FindByUserId(userId, 0, -1)
	if err != nil {
		return
	}

	var cardIds []int64

	for _, v := range cards {
		cardIds = append(cardIds, v.Id)
	}

	if len(cardIds) < 1 {
		// ErrNoBill
		return
	}

	return c.repository.Bill().FindByCardIds(cardIds, page, pageSize)

}

func (c *service) Repay(ctx context.Context, cardId int64, amount float64, repaymentDay *time.Time) (err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	var card *types.CreditCard
	if card, err = c.repository.CreditCard().FindById(cardId, userId); err != nil {
		return
	}

	if card.Id == 0 {
		return ErrNotPermission
	}

	var y, d int
	var m time.Month

	if repaymentDay != nil {
		y, m, d = repaymentDay.Date()
	} else {
		y, m, _ = time.Now().Date()
		d = card.Cardholder
	}

	t := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

	return c.repository.Bill().Repay(cardId, amount, t)
}

func (c *service) GenBill(ctx context.Context, day int) (err error) {
	cards, err := c.repository.CreditCard().FindByBillDay(day)
	if err != nil {
		_ = level.Error(c.logger).Log("CreditCard", "FindByBillDay", "err", err.Error())
		return
	}

	curr := time.Now()
	year, month, _ := curr.Date()

	for _, card := range cards {
		endTime := time.Date(year, month, card.BillingDay, 0, 0, 0, 1, time.Local)

		// 还款日计算会有问题 可以考虑使用util.date.ParseCardBillAndHolderDay方法生成
		billing, holder := date.ParseCardBillAndHolderDay(card.BillingDay, card.Cardholder)

		billAmount, err := c.repository.ExpenseRecord().RemainingAmount(card.Id, billing, endTime)
		if err != nil {
			_ = level.Error(c.logger).Log("cardId", card.Id, "startTime", billing.String(), "endTime", endTime, "ExpenseRecord", "RemainingAmount", "err", err.Error())
			continue
		}

		if err = c.repository.Bill().Create(card.Id, billAmount.Amount, holder); err != nil {
			_ = level.Error(c.logger).Log("cardId", card.Id, "amount", billAmount.Amount, "Bill", "Create", "err", err.Error())
		}
	}

	return
}
