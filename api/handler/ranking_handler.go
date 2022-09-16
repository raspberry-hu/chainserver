package handler

import (
	"ChainServer/api/response"
	"ChainServer/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

// 获取collection列表
func RankingGetByLimit(ctx *gin.Context) {
	// 当前的页码
	page, _ := strconv.Atoi(ctx.Query("page"))
	// 每页所要展示的条数 默认为10
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	if rows == 0 {
		rows = 10
	}
	// 数据库查询的起始下标位置
	offset := (page - 1) * rows
	category := ctx.Query("category")
	chainType := ctx.Query("chain")

	var limitSql string
	if page != 0 {
		limitSql = fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	}

	logger.Info("RankingGetByLimit req: category, chain", category, chainType)
	var res []models.TCollectionInfo
	if category == "" && chainType == "" {
		res = models.RankingFindAll(limitSql)
	} else {
		res = models.RankingByCondition(category, chainType, limitSql)
	}
	var recordList []response.CollectionView
	for _, value := range res {
		var records response.CollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"id": value.UserId})
		if len(ownerUserInfo) > 0 {
			records.OwnerName = ownerUserInfo[0].UserName
		}
		records.CollectionId = value.Id
		records.LogoImageURL = value.LogoImage
		records.CollectionName = value.CollectionName
		records.BannerImageURL = value.BannerImageUrl
		records.CollectionDesc = value.CollectionDesc
		records.CategoryName = value.Category
		records.ChainName = value.ChainName
		records.CreateTax = value.CreateTax
		records.Owner = value.UserId
		records.Favorites = value.Favorites
		records.Items = value.Items
		records.Amount = value.Amount
		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}
