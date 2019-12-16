/**
 * @Time : 2019-08-19 11:03
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package creditcard

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
	"github.com/nsini/cardbill/src/util/date"
	"strconv"
	"time"
)

type Service interface {
	// 增加信用卡
	Post(ctx context.Context, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder int, cardNumber, tailNumber int64) (err error)

	// 获取信用卡列表
	List(ctx context.Context, bankId int64) (res []*types.CreditCard, err error)

	// 更新信用卡信息
	Put(ctx context.Context, id int64, cardName string, bankId int64,
		fixedAmount, maxAmount float64, billingDay, cardHolder, state int) (err error)

	// 消费统计
	Statistics(ctx context.Context) (res *StatisticsResponse, err error)

	// 卡消费记录
	Record(ctx context.Context, id int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error)

	// 获取信用卡信息
	Get(ctx context.Context, id int64) (res *types.CreditCard, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repository repository.Repository) Service {
	return &service{logger: logger, repository: repository}
}

func (c *service) Get(ctx context.Context, id int64) (res *types.CreditCard, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	res, err = c.repository.CreditCard().FindById(id, userId, "Bank", "User")
	if err != nil {
		return
	}

	billingDay, _ := date.ParseCardBillAndHolderDay(res.BillingDay, res.Cardholder)

	cardHolder := billingDay.AddDate(0, -1, 0)

	ra, err := c.repository.ExpenseRecord().RemainingAmount(id, cardHolder, time.Now())
	if err != nil {
		return
	}

	remainingAmount := res.MaxAmount - ra.Amount
	res.RemainingAmount = remainingAmount

	return
}

func (c *service) Record(ctx context.Context, id int64, page, pageSize int) (res []*types.ExpensesRecord, count int64, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)
	return c.repository.ExpenseRecord().ListByCardId(userId, id, page, pageSize)
}

func (c *service) Statistics(ctx context.Context) (res *StatisticsResponse, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}

	var cardIds []int64
	if cards, err := c.repository.CreditCard().FindByUserId(userId, 0); err == nil {
		for _, v := range cards {
			cardIds = append(cardIds, v.Id)
		}
	}

	currentMonth := time.Now()

	creditTotalCh := make(chan int64)
	creditAmountCh := make(chan *repository.TotalAmount)
	sacCh := make(chan *repository.RemainingAmount)
	currSacCh := make(chan *repository.RemainingAmount)
	unpaidBillCh := make(chan *repository.BillAmount)
	repaidBillCh := make(chan *repository.BillAmount)

	go func() {
		// 信用卡数量
		creditTotal, err := c.repository.CreditCard().Count(userId)
		if err == nil {
			creditTotalCh <- creditTotal
		} else {
			creditTotalCh <- 0
			_ = level.Error(c.logger).Log("CreditCard", "Count", "err", err.Error())
		}
	}()

	go func() {
		// 信用卡总额度
		creditAmount, err := c.repository.CreditCard().Sum(userId)
		if err != nil {
			creditAmountCh <- nil
			_ = level.Error(c.logger).Log("CreditCard", "Sum", "err", err.Error())
		} else {
			creditAmountCh <- creditAmount
		}
	}()

	go func() {
		// 总消费
		sac, err := c.repository.ExpenseRecord().SumAmountCards(cardIds, nil)
		if err != nil {
			sacCh <- nil
			_ = level.Error(c.logger).Log("ExpenseRecord", "SumAmountCards", "err", err.Error())
		} else {
			sacCh <- sac
		}
	}()

	go func() {
		// 当月消费
		currSac, err := c.repository.ExpenseRecord().SumAmountCards(cardIds, &currentMonth)
		if err != nil {
			currSacCh <- nil
			_ = level.Error(c.logger).Log("ExpenseRecord", "SumAmountCards", "err", err.Error())
		} else {
			currSacCh <- currSac
		}
	}()

	go func() {
		// 账单
		unpaidBill, err := c.repository.Bill().SumByCards(cardIds, nil, repository.RepayFalse)
		if err != nil {
			unpaidBillCh <- nil
			_ = level.Error(c.logger).Log("Bill", "SumByCards", "err", err.Error())
		} else {
			unpaidBillCh <- unpaidBill
		}
	}()

	go func() {
		repaidBill, err := c.repository.Bill().SumByCards(cardIds, nil, repository.RepayTrue)
		if err != nil {
			repaidBillCh <- nil
			_ = level.Error(c.logger).Log("Bill", "SumByCards", "err", err.Error())
		} else {
			repaidBillCh <- repaidBill
		}
	}()

	totalAmount := <-creditAmountCh
	cardNumber := <-creditTotalCh
	sac := <-sacCh
	currSac := <-currSacCh
	unpaidBill := <-unpaidBillCh
	repaidBill := <-repaidBillCh

	close(creditTotalCh)
	close(creditAmountCh)
	close(sacCh)
	close(currSacCh)
	close(unpaidBillCh)
	close(repaidBillCh)

	interestExpense, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(2)+"f", sac.Amount-sac.Arrival), 64)
	currentInterest, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(2)+"f", currSac.Amount-currSac.Arrival), 64)

	return &StatisticsResponse{
		CreditAmount:       totalAmount.Amount,
		CreditMaxAmount:    totalAmount.MaxAmount,
		CreditNumber:       int(cardNumber),
		TotalConsumption:   sac.Amount,
		MonthlyConsumption: currSac.Amount,
		InterestExpense:    interestExpense,
		CurrentInterest:    currentInterest,
		UnpaidBill:         unpaidBill.Amount,
		RepaidBill:         repaidBill.Amount,
	}, nil
}

func (c *service) Post(ctx context.Context, cardName string, bankId int64, fixedAmount, maxAmount float64, billingDay, cardHolder int, cardNumber, tailNumber int64) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	return c.repository.CreditCard().Create(&types.CreditCard{
		CardName:    cardName,
		BankId:      bankId,
		FixedAmount: fixedAmount,
		MaxAmount:   maxAmount,
		BillingDay:  billingDay,
		Cardholder:  cardHolder,
		UserId:      userId,
		TailNumber:  tailNumber,
	})
}

func (c *service) Put(ctx context.Context, id int64, cardName string, bankId int64,
	fixedAmount, maxAmount float64, billingDay, cardHolder, state int) (err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return middleware.ErrCheckAuth
	}

	return c.repository.CreditCard().Update(&types.CreditCard{
		Id:          id,
		CardName:    cardName,
		BankId:      bankId,
		FixedAmount: fixedAmount,
		MaxAmount:   maxAmount,
		BillingDay:  billingDay,
		Cardholder:  cardHolder,
		UserId:      userId,
		State:       state,
	})

}

func (c *service) List(ctx context.Context, bankId int64) (res []*types.CreditCard, err error) {
	userId, ok := ctx.Value(middleware.UserIdContext).(int64)
	if !ok {
		return nil, middleware.ErrCheckAuth
	}

	return c.repository.CreditCard().FindByUserId(userId, bankId)
}
