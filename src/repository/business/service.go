/**
 * @Time : 4/6/21 10:06 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package business

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type Middleware func(Service) Service

type Service interface {
	SaveMerchant(ctx context.Context, business *types.Merchant) (err error)
	Types(ctx context.Context) (res []types.Business, err error)
	FindByCode(ctx context.Context, code int64) (res types.Business, err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindByCode(ctx context.Context, code int64) (res types.Business, err error) {
	err = s.db.Model(&types.Business{}).Where("code = ?", code).First(&res).Error
	return
}

func (s *service) Types(ctx context.Context) (res []types.Business, err error) {
	err = s.db.Model(&types.Business{}).Order("code ASC").Find(&res).Error
	return
}

func (s *service) SaveMerchant(ctx context.Context, merchant *types.Merchant) (err error) {
	return s.db.Model(&types.Merchant{}).FirstOrCreate(merchant).
		Where("merchant_name = ?", merchant.MerchantName).Error
}

func NewService(db *gorm.DB) Service {
	return &service{db: db}
}
