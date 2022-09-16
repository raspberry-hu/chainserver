package handler

import (
	"ChainServer/api/request"
	"ChainServer/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func WalletInfoPost(ctx *gin.Context) {
	var wallerInfoRequest request.WalletInfo
	ctx.ShouldBindJSON(&wallerInfoRequest)
	wallet_addr := wallerInfoRequest.WalletAddr
	user_name := wallerInfoRequest.UserName
	image_url := wallerInfoRequest.ImageUrl
	banner_url := wallerInfoRequest.BannerUrl
	info_arry := models.WalletInfoFind(map[string]interface{}{"wallet_addr": wallet_addr})

	if strings.ToLower(ctx.Request.Header.Get("wallet")) != wallet_addr {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 201,
			"msg":  "fail",
		})
		return
	}

	var err error
	var t_wallet_info models.TWalletInfo
	t_wallet_info.WalletAddr = wallet_addr
	t_wallet_info.UserName = user_name
	t_wallet_info.ImageUrl = image_url
	t_wallet_info.BannerUrl = banner_url
	if len(info_arry) == 1 {
		models.WalletInfoUpdate(info_arry[0], t_wallet_info)
	} else {
		models.WalletInfoInsert(t_wallet_info)
		//第一次登录，设定该账户的permission
		var accountPermission models.TPermissions
		permissionInfo := models.TPermissionsFind(map[string]interface{}{"wallet_addr": wallet_addr})
		if len(permissionInfo) == 0 {
			accountPermission.WalletAddr = wallet_addr
			id := models.TPermissionsInsert(accountPermission)
			accountPermission.Id = id
		}
	}
	var code int
	var msg string
	if err != nil {
		code = 201
		msg = "fail"
	} else {
		code = 200
		msg = "success"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func WalletInfoGet(ctx *gin.Context) {
	wallet_addr := ctx.Query("wallet_addr")
	var t_wallet_info []models.TWalletInfo
	t_wallet_info = models.WalletInfoFind(map[string]interface{}{"wallet_addr": wallet_addr})

	var result = make(map[string]interface{})

	result["permissions"] = 0
	if len(t_wallet_info) > 0 {
		result["image_url"] = t_wallet_info[0].ImageUrl
		result["banner_url"] = t_wallet_info[0].BannerUrl
		result["user_name"] = t_wallet_info[0].UserName
		result["user_desc"] = t_wallet_info[0].UserDesc
		result["wallet_addr"] = t_wallet_info[0].WalletAddr
	}
	list := models.TPermissionsFind(map[string]interface{}{"wallet_addr": wallet_addr})
	if len(list) > 0 {
		result["permissions"] = 1
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "success",
		"result": result,
	})
}

func WalletInfoList(ctx *gin.Context) {
	wallet_addr := ctx.Query("wallet_addr")
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	limit_sql := fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	data_arry := models.WalletInfoList(wallet_addr, limit_sql)

	results := []map[string]interface{}{}
	for _, value := range data_arry {
		var result = make(map[string]interface{})
		result["image_url"] = value.ImageUrl
		result["banner_url"] = value.BannerUrl
		result["user_name"] = value.UserName
		result["user_desc"] = value.UserDesc
		result["wallet_addr"] = value.WalletAddr
		results = append(results, result)
	}
	count := models.WalletInfoLikeCount(wallet_addr)
	ctx.JSON(http.StatusOK, gin.H{
		"code":   200,
		"msg":    "success",
		"result": results,
		"count":  count,
	})
}
