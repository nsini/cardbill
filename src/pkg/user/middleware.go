/**
 * @Time : 3/26/21 5:42 PM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package user

type Middleware func(Service) Service
