/**
 * @Time : 2019-10-11 15:52
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package dashboard

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/middleware"
	"github.com/nsini/cardbill/src/repository"
	"time"
)

type Service interface {
	// 最近一个31天的消费记录
	LastAmount(ctx context.Context) (resp []LastAmountResponse, err error)

	// 近一年的消费记录
	MonthAmount(ctx context.Context) (resp []LastAmountResponse, err error)
}

type service struct {
	logger     log.Logger
	repository repository.Repository
}

func NewService(logger log.Logger, repo repository.Repository) Service {
	return &service{logger: logger, repository: repo}
}

type LastAmountResponse struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

func (c *service) LastAmount(ctx context.Context) (resp []LastAmountResponse, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	t := time.Now()

	days, err := c.repository.ExpenseRecord().SumDays(userId)
	if err != nil {
		return
	}

D:
	for i := 0; i <= 31; i++ {
		ttt := t.AddDate(0, 0, -i)
		for _, v := range days {
			tt, err := time.Parse("2006-01-02", v.Day)
			if err != nil {
				_ = level.Warn(c.logger).Log("time", "Parse", "err", err.Error())
				continue D
			}
			if tt.Month() == ttt.Month() && tt.Day() == ttt.Day() {
				resp = append(resp, LastAmountResponse{
					Date:   tt,
					Amount: v.Amount,
				})
				continue D
			}
		}
		resp = append(resp, LastAmountResponse{
			Date:   ttt,
			Amount: 0,
		})
	}

	return
}

func (c *service) MonthAmount(ctx context.Context) (resp []LastAmountResponse, err error) {
	userId := ctx.Value(middleware.UserIdContext).(int64)

	t := time.Now()

	days, err := c.repository.ExpenseRecord().SumMonth(userId)
	if err != nil {
		return
	}

D:
	for i := 0; i <= 12; i++ {
		ttt := t.AddDate(0, -i, 0)
		for _, v := range days {
			tt, err := time.Parse("2006-01", v.Day)
			if err != nil {
				_ = level.Warn(c.logger).Log("time", "Parse", "err", err.Error())
				continue D
			}
			if tt.Year() == ttt.Year() && tt.Month() == ttt.Month() {
				resp = append(resp, LastAmountResponse{
					Date:   tt,
					Amount: v.Amount,
				})
				continue D
			}
		}
		resp = append(resp, LastAmountResponse{
			Date:   ttt,
			Amount: 0,
		})
	}

	return
}

func (c *service) EveryMont(ctx context.Context) {

}
