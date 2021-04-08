/**
 * @Time : 3/30/21 5:04 PM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package bank

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
	"strings"
)

type Middleware func(Service) Service

type Service interface {
	Find(ctx context.Context, bankId int64) (res types.Bank, err error)
	List(ctx context.Context, bankName string) (res []types.Bank, total int, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Find(ctx context.Context, bankId int64) (res types.Bank, err error) {
	query := s.db.Model(&types.Bank{})
	query = query.Where("id = ?", bankId)
	err = query.First(&res).Error
	return
}

func (s *service) List(ctx context.Context, bankName string) (res []types.Bank, total int, err error) {
	query := s.db.Model(&types.Bank{})

	if !strings.EqualFold(bankName, "") {
		query = query.Where("bank_name LIKE '%?%'", bankName)
	}
	err = query.Count(&total).Order("id DESC").Find(&res).Error
	return
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
