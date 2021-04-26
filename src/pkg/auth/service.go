/**
 * @Time: 2019-08-18 16:56
 * @Author: solacowa@gmail.com
 * @File: service
 * @Software: GoLand
 */

package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/go-github/v26/github"
	"github.com/icowan/config"
	"github.com/jinzhu/gorm"
	"github.com/nsini/cardbill/src/encode"
	cbjwt "github.com/nsini/cardbill/src/jwt"
	jwt2 "github.com/nsini/cardbill/src/jwt"
	"github.com/nsini/cardbill/src/pkg/wechat"
	"github.com/nsini/cardbill/src/repository"
	"github.com/nsini/cardbill/src/repository/types"
	"golang.org/x/oauth2"
	oauthgithub "golang.org/x/oauth2/github"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrAuthLoginGitHubGetUser = errors.New("获取Github用户邮箱及名称失败.")
)

type Service interface {
	// github 授权登陆跳转
	AuthLoginGithub(w http.ResponseWriter, r *http.Request)

	// github 授权登陆回调
	AuthLoginGithubCallback(w http.ResponseWriter, r *http.Request)

	// weibo 授权登录跳转
	AuthLoginWeibo(w http.ResponseWriter, r *http.Request)

	// weibo 授权登陆回调
	AuthLoginWeiboCallback(w http.ResponseWriter, r *http.Request)

	// 微信授权登录
	AuthLoginWechat(w http.ResponseWriter, r *http.Request)

	// 微信授权登录回调
	AuthLoginWechatCallback(w http.ResponseWriter, r *http.Request)

	// 微信小程序授权登录
	AuthLoginMP(ctx context.Context, code, iv, rawData, signature, encryptedData, inviteCode string) (res loginResponse, err error)
}

type service struct {
	logger     log.Logger
	config     *config.Config
	repository repository.Repository
	traceId    string
	wechat     wechat.Service
}

type userInfo struct {
	AvatarURL string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	NickName  string `json:"nickName"`
	Province  string `json:"province"`
}

func (s *service) AuthLoginMP(ctx context.Context, code, iv, rawData, signature, encryptedData, inviteCode string) (res loginResponse, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "AuthLoginMP")
	var reqUserInfo userInfo
	if err = json.NewDecoder(bytes.NewBufferString(rawData)).Decode(&reqUserInfo); err != nil {
		_ = level.Warn(logger).Log("json", "NewDecoder", "userInfo", "Decode", "err", err.Error())
	}

	// todo: 校验 signature, encryptedData

	userInfo, sessionKey, err := s.wechat.MPLogin(ctx, code)
	if err != nil {
		_ = level.Error(logger).Log("wechat", "MPLogin", "err", err.Error())
		return
	}
	var user types.User
	if user, err = s.repository.Users().FindByUnionId(ctx, userInfo.UnionId); err != nil {
		if err != gorm.ErrRecordNotFound {
			_ = level.Error(logger).Log("gorm", "ErrRecordNotFound", "err", err.Error())
			err = encode.ErrAuthMPLogin.Error()
			return res, err
		}
		u := &types.User{
			OpenId:   userInfo.OpenId,
			UnionId:  userInfo.UnionId,
			Nickname: reqUserInfo.NickName,
			Sex:      reqUserInfo.Gender,
			City:     reqUserInfo.City,
			Province: reqUserInfo.Province,
			Country:  reqUserInfo.Country,
			Avatar:   reqUserInfo.AvatarURL,
			Remark:   "小程序登录",
		}

		if err = s.repository.Users().Save(ctx, u); err != nil {
			_ = level.Error(logger).Log("repository.User", "Save", "err", err.Error())
			err = encode.ErrAuthMPLogin.Error()
			return
		}
		user = *u
		_ = level.Info(logger).Log("repository.User", "FindByUnionId", "msg", "用户不存在,保存信息")
	}

	defer func() {
		user.Nickname = reqUserInfo.NickName
		user.Sex = reqUserInfo.Gender
		user.City = reqUserInfo.City
		user.Province = reqUserInfo.Province
		user.Country = reqUserInfo.Country
		user.Avatar = reqUserInfo.AvatarURL
		_ = s.repository.Users().Save(ctx, &user)
	}()

	sessionTimeout := 3600 * 24 * 31 * 12 * time.Second

	expAt := time.Now().Add(sessionTimeout).Unix()

	claims := jwt2.ArithmeticCustomClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(jwt2.GetJwtKey()))
	if err != nil {
		_ = level.Error(logger).Log("token", "SignedString", "err", err.Error())
	}

	//_ = s.cache.Set(ctx, fmt.Sprintf("login:%d:token", user.Id), tk, sessionTimeout)

	return loginResponse{
		Token:      tk,
		SessionKey: sessionKey,
		Avatar:     reqUserInfo.AvatarURL,
		Nickname:   reqUserInfo.NickName,
	}, nil
}

func (s *service) AuthLoginWechat(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s *service) AuthLoginWechatCallback(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s *service) AuthLoginWeibo(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s *service) AuthLoginWeiboCallback(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (s *service) AuthLoginGithub(w http.ResponseWriter, r *http.Request) {
	githubOauthConfig := s.auth2Config()

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	u := githubOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (s *service) auth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.config.GetString("github", "client_id"),
		ClientSecret: s.config.GetString("github", "client_secret"),
		Scopes:       []string{"SCOPE1", "SCOPE2", "user:email"},
		Endpoint:     oauthgithub.Endpoint,
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	// todo 用jwt生成 然后 jwt 解析出来
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (s *service) AuthLoginGithubCallback(w http.ResponseWriter, r *http.Request) {
	var resp authResponse

	ctx := context.Background()
	// state := r.URL.Query().Get("state") // todo 它需要验证一下可以考虑使用jwt生成  先用cookie 简单处理一下吧...

	if httpProxy := s.config.GetString("server", "http_proxy"); httpProxy != "" {
		_ = level.Debug(s.logger).Log("use-proxy", httpProxy)
		dialer := &net.Dialer{
			Timeout:   time.Duration(5 * int64(time.Second)),
			KeepAlive: time.Duration(5 * int64(time.Second)),
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {
					return url.Parse(httpProxy)
				},
				DialContext: dialer.DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		})
	}

	githubOauthConfig := s.auth2Config()

	if r.URL.Query().Get("error") != "" {
		resp.Err = errors.New(r.URL.Query().Get("error") + ": " + r.URL.Query().Get("error_description"))
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	token, err := githubOauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		if strings.Contains(err.Error(), "server response missing access_token") {
			http.Redirect(w, r, s.config.GetString("server", "domain")+"/#/user/login", http.StatusPermanentRedirect)
		}
		_ = s.logger.Log("githubOauthConfig", "Exchange", "err", err.Error())
		resp.Err = err
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	if token == nil || !token.Valid() {
		_ = s.logger.Log("token", "nil", "or", "token.valid is false")
		resp.Err = errors.New("token is nil or token.valid is false")
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	client := github.NewClient(githubOauthConfig.Client(ctx, token))
	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		_ = s.logger.Log("client.users", "Get", "err", err.Error())
		resp.Err = err
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	if user == nil {
		resp.Err = ErrAuthLoginGitHubGetUser
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	authId := user.GetID()

	username := user.GetName()
	if username == "" {
		username = user.GetLogin()
	}

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		_ = s.logger.Log("invalid", "oauth github state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rs, member, err := s.AuthLogin(authId, user.GetEmail(), username)
	if err != nil {
		resp = authResponse{Err: err}
		_ = encodeLoginResponse(ctx, w, resp)
		return
	}

	//_ = c.casbin.GetEnforcer().LoadPolicy()

	params := url.Values{}
	params.Add("token", strings.Replace(rs, "Bearer ", "", -1))
	params.Add("username", member.Username)

	http.Redirect(w, r, s.config.GetString("server", "domain")+"/?"+params.Encode(), http.StatusPermanentRedirect)

}

func (s *service) AuthLogin(authId int64, email, username string) (rs string, member *types.User, err error) {
	member, err = s.repository.User().FindByAuthId(authId)

	if member == nil || err != nil {
		member = &types.User{
			Username: username,
			Email:    email,
			AuthId:   authId,
		}
		if err = s.repository.User().Create(member); err != nil {
			_ = s.logger.Log("User", "Create", "err", err.Error())
			return "", nil, err
		}
	}

	rs, err = s.sign(strconv.Itoa(int(authId)), member.Id)
	rs = "Bearer " + rs
	return
}

func (s *service) sign(authId string, uid int64) (string, error) {
	sessionTimeout, err := s.config.Int64("server", "session_timeout")
	if err != nil {
		sessionTimeout = 3600
	}
	expAt := time.Now().Add(time.Duration(sessionTimeout) * time.Second).Unix()

	// 创建声明
	claims := cbjwt.ArithmeticCustomClaims{
		UserId:   uid,
		Username: authId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成token
	return token.SignedString([]byte(cbjwt.GetJwtKey()))
}

func NewService(logger log.Logger, cf *config.Config, store repository.Repository) Service {
	logger = log.With(logger, "auth", "service")
	return &service{
		logger:     logger,
		config:     cf,
		repository: store,
	}
}
