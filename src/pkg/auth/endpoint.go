/**
 * @Time: 2019-08-18 17:06
 * @Author: solacowa@gmail.com
 * @File: endpoint
 * @Software: GoLand
 */

package auth

type authResponse struct {
	//Success bool   `json:"success"`
	Code  int                    `json:"code"`
	Token string                 `json:"token"`
	Data  map[string]interface{} `json:"data"`
	Err   error                  `json:"error"`
}
