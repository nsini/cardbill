/**
 * @Time : 2019-09-18 17:32
 * @Author : solacowa@gmail.com
 * @File : service_gen
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"github.com/nsini/cardbill/src/pkg/bill"
	"github.com/robfig/cron"
	"time"
)

func billCronJob(c *cron.Cron, service bill.Service) {
	spec := "0 0 1 * * ?" // 每天凌晨1点执行

	if err := c.AddFunc(spec, func() {
		if err := service.GenBill(context.Background(), time.Now().Day()); err != nil {
			_ = level.Error(logger).Log("service", "GenBill", "err", err.Error())
		}
	}); err != nil {
		_ = level.Error(logger).Log("c", "AddFunc", "err", err.Error())
	}
}
