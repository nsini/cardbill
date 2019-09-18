/**
 * @Time : 2019-08-20 17:24
 * @Author : solacowa@gmail.com
 * @File : claims
 * @Software: GoLand
 */

package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

// ArithmeticCustomClaims 自定义声明
type ArithmeticCustomClaims struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}
