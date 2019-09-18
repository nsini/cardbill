/**
 * @Time: 2019-08-18 17:09
 * @Author: solacowa@gmail.com
 * @File: user
 * @Software: GoLand
 */

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/repository/types"
)

type UserRepository interface {
	FindByEmail(email string) (res *types.User, err error)
	FindById(id int64) (res *types.User, err error)
	FindByAuthId(authId int64) (res *types.User, err error)
	Create(user *types.User) (err error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (c *userRepository) FindByEmail(email string) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "email = ?", email).Error
	return &rs, err
}

func (c *userRepository) FindByAuthId(authId int64) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "auth_id = ?", authId).Error
	return &rs, err
}

func (c *userRepository) Create(user *types.User) (err error) {
	err = c.db.Save(user).Error
	return err
}

func (c *userRepository) FindById(id int64) (res *types.User, err error) {
	var rs types.User
	err = c.db.Select("username").First(&rs, "id = ?", id).Error
	return &rs, err
}
