package handler

import (
	"ChainServer/api/request"
	"ChainServer/common"
	"ChainServer/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OrganizationPost(ctx *gin.Context) {
	var organization_info request.OrganizationInfo
	ctx.ShouldBindJSON(&organization_info)
	fmt.Println(organization_info)
}

func OrganizationGet(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	data_arry := models.OrganizationInfoFind(rows, offset, map[string]interface{}{})
	result_arry := []map[string]interface{}{}
	for _, data := range data_arry {
		result := make(map[string]interface{})
		result["org_name"] = data.OrgName
		result["image_url"] = data.ImageUrl
		result["user_id"] = data.UserId
		result_arry = append(result_arry, result)
	}
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result_arry,
		Code:   200,
	}
	ctx.JSON(http.StatusOK, response)
}
