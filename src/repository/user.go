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
	FindByGithubId(githubId int64) (res *types.User, err error)
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

func (c *userRepository) FindByGithubId(githubId int64) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "github_id = ?", githubId).Error
	return &rs, err
}

func (c *userRepository) Create(user *types.User) (err error) {
	err = c.db.Save(user).Error
	return err
}
