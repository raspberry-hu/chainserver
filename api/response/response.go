package response

import (
	"ChainServer/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReturnSuccessResponse 接口返回数据信息
func ReturnSuccessResponse(ctx *gin.Context, data interface{}) {
	msgId := common.GetRandString(16)
	response := &common.Response{
		ID:      msgId,
		Version: "2.0",
		Error:   nil,
		Result:  data,
		Code:    200,
	}
	ctx.JSON(http.StatusOK, response)
	return
}

func ReturnSuccessResponseNew(ctx *gin.Context, data interface{}) {
	response := &common.ResponseNew{
		Msg:    "success",
		Result: data,
		Code:   200,
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// ReturnErrorResponse 返回错误
func ReturnErrorResponse(ctx *gin.Context, code int64, err string, errmsg string) {
	msgId := common.GetRandString(16)
	response := &common.Response{
		ID:      msgId,
		Version: "2.0",
		Error: &common.Error{
			Code: code,
			Msg:  errmsg,
			Data: err,
		},
	}
	ctx.JSON(http.StatusOK, response)
	return
}

func ReturnErrorResponseNew(ctx *gin.Context, code int, errmsg string) {
	response := &common.ResponseNew{
		Msg:  errmsg,
		Code: code,
	}
	ctx.JSON(http.StatusOK, response)
	return
}
