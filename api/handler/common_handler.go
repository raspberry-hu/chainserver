package handler

import (
	myjwt "ChainServer/api/jwt"
	"ChainServer/api/request"
	"ChainServer/common"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

var WalletAddrNonce sync.Map
var ExpiresAt int64 = 86400 // 过期时间
var NotBefore int64 = 1000  // 签名生效时间

// 生成令牌
func GenerateToken(ctx *gin.Context) {
	var auth_req request.AuthRequest
	ctx.ShouldBindJSON(&auth_req)
	wallet_addr := strings.ToLower(auth_req.WalletAddr)
	sign := auth_req.Sign
	user_name := ""
	// 验证签名

	value, ok := WalletAddrNonce.Load(wallet_addr)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  "找不到nonce",
		})
		return
	}
	msg := value.(string)
	verify, err := common.ServiceVerify(wallet_addr, msg, sign)
	if ! verify || err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  "签名验证失败",
		})
		return
	}

	j := &myjwt.JWT{
		[]byte(wallet_addr),
	}
	claims := myjwt.CustomClaims{
		wallet_addr,
		user_name,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - NotBefore), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + ExpiresAt), // 过期时间
		},
	}
	token, err := j.CreateToken(claims)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 203,
			"msg":  err.Error(),
		})
		return
	}
	nonce := fmt.Sprintf("%v", time.Now().Unix())
	WalletAddrNonce.Store(wallet_addr, nonce)
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "登录成功！",
		"token": token,
	})
	return
}

func GetNonce(ctx *gin.Context) {
	wallet_addr := ctx.Query("wallet_addr")
	value, ok := WalletAddrNonce.Load(wallet_addr)
	var nonce string
	if ok {
		nonce = value.(string)
	} else {
		nonce = fmt.Sprintf("%v", time.Now().Unix())
		WalletAddrNonce.Store(wallet_addr, nonce)
	}
	response := &common.ResponseCommon{
		Msg:    "success",
		Result: nonce,
		Code:   200,
	}
	ctx.JSON(http.StatusOK, response)
}

func Expiration(ctx *gin.Context) {
	var expiration_time int64
	wallet_addr := strings.ToLower(ctx.Request.Header.Get("wallet"))
	token := ctx.Request.Header.Get("token")
	j := myjwt.NewJWT(wallet_addr)
	// parseToken 解析token包含的信息
	claims, err := j.ParseToken(token)
	if err != nil {
		expiration_time = time.Now().Unix() - 1000
	} else {
		expiration_time = claims.ExpiresAt
	}

	response := &common.ResponseCommon{
		Msg:    "success",
		Result: expiration_time,
		Code:   200,
	}
	ctx.JSON(http.StatusOK, response)
}

