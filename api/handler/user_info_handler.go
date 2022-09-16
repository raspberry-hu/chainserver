package handler

import (
	"ChainServer/api/request"
	"ChainServer/middleware"
	"ChainServer/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"ChainServer/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 注册
func Register(ctx *gin.Context) {
	var request models.UserInfo
	ctx.ShouldBindJSON(&request)
	code := -1
	msg := "注册失败"
	//request.CreateTime = time.Now() //给创建时间赋值为当前时间
	flag := models.UserInfoInsert(request)
	if flag {
		code = 200
		msg = "注册成功"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})

}

// 登录
func Login(ctx *gin.Context) {
	code := -1
	msg := "登录失败"
	var m request.LoginRequest
	_ = ctx.ShouldBindJSON(&m)

	u := &models.UserInfo{UserName: m.Username, UserPass: m.Password}
	if user, err := models.Login(u); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":   code,
			"msg":    fmt.Sprintf("%s %s", msg, err.Error()),
			"result": nil,
		})
	} else {
		tokenNext(ctx, *user)
	}

}

//登录以后签发jwt
func tokenNext(c *gin.Context, user models.UserInfo) {
	j := &middleware.JWT{
		SigningKey: []byte(config.Conf.Jwt.SigningKey), //jwt签名
	}
	clams := request.UserClaims{

		ID:       uint(user.Id),
		TrueName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),       //签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60*60*24*7), //过期时间 7天
			Issuer:    "qmPlus",                              //签名发行者
		},
	}
	token, err := j.CreateToken(clams)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "登录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"msg":       "登录成功",
		"user":      user,
		"token":     token,
		"expiresAt": clams.StandardClaims.ExpiresAt * 1000,
	})

}

// 更新人员信息（也可用于更新密码）
func UpdateUser(c *gin.Context) {
	parameter := models.UserInfo{}
	_ = c.ShouldBindJSON(&parameter)
	user, err := models.UpdateUser(parameter)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": user,
	})
}

//根据id查询用户信息
func FindUserByID(c *gin.Context) {
	parameter := models.UserInfo{}
	_ = c.ShouldBindJSON(&parameter)
	user, err := models.UserFindByID(parameter.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询失败" + err.Error(),
		})
		return

	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": user,
	})

}

// 删除用户
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败，参数错误 " + err.Error(),
		})

		return
	}
	err = models.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败 " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// 修改密码
func ChangePassword(c *gin.Context) {
	var params request.ChangePasswordStruct
	_ = c.ShouldBindJSON(&params)
	//u := &models.UserInfo{UserName: params.Username, UserPass: params.Password}
	if _, err := models.ChangePssword(params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改失败，请检查用户名密码 " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
