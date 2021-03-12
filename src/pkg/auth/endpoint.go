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

type (
	loginResponse struct {
		Token      string `json:"token"`
		OpenId     string `json:"openId"`
		SessionKey string `json:"sessionKey"`
		UnionId    string `json:"unionId"`
		Avatar     string `json:"avatar"`
		Nickname   string `json:"nickname"`
		ShareCode  string `json:"shareCode"`
	}
)
