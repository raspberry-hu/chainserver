package handler

import (
	"ChainServer/api/response"
	"ChainServer/models"
	"ChainServer/models/model"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

// HomeLimiteds
func HomeLimiteds(ctx *gin.Context) {
	rows := ctx.Query("rows")
	chain_type := ctx.Query("chain_type")
	logger.Info("HomeLimiteds req:", "rows", rows)
	if rows == "" {
		rows = "8"
	}

	res := model.MarketHomeLimiteds("2", rows, chain_type)
	// 按CreateTime值进行排序算法
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res)-i-1; j++ {
			if res[j].CreateTime < res[j+1].CreateTime {
				res[j], res[j+1] = res[j+1], res[j]
			}
		}
	}
	// 补充nft的信息
	var recordList []response.MarketAndNftView
	for _, value := range res {
		var records response.MarketAndNftView
		nftInfo := model.NftGetByTokenId(value.Tokenid, value.ChainType)
		// market 数据
		records.Id = value.Id
		records.SN = value.SN
		records.Creater = value.Creater
		records.Tokenid = value.Tokenid
		records.MarketType = value.MarketType
		records.StartingPrice = value.StartingPrice
		records.TokenType = value.TokenType
		records.EndTime = time.Unix(value.EndTime, 0).Format("2006-01-02 15:04:05")
		records.Buyer = value.Buyer
		records.Bonus = value.Bonus
		records.Txhash = value.Txhash
		records.CancelHash = value.CancelHash
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.Sorting = value.Sorting
		records.Status = value.Status
		// nft 数据
		records.NftName = nftInfo[0].NftName
		records.NftDesc = nftInfo[0].NftDesc
		records.TokenId = nftInfo[0].TokenId
		records.TokenUri = nftInfo[0].TokenUri
		records.NftTxhash = nftInfo[0].Txhash
		records.NftCreater = nftInfo[0].Creater
		records.CreateNumber = nftInfo[0].CreateNumber
		records.Nft_CreateTime = nftInfo[0].CreateTime
		records.MediaUri = nftInfo[0].MediaUri
		records.CreateTax = nftInfo[0].CreateTax
		records.Owner = nftInfo[0].Owner
		records.NftType = nftInfo[0].NftType
		records.Approved = nftInfo[0].Approved

		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// HomeAuctions
func HomeAuctions(ctx *gin.Context) {
	rows := ctx.Query("rows")
	chain_type := ctx.Query("chain_type")
	logger.Info("HomeLimiteds req:", "rows", rows)
	if rows == "" {
		rows = "8"
	}

	res := model.MarketHomeLimiteds("1", rows, chain_type)
	// 补充nft的信息
	var recordList []response.MarketAndNftView
	for _, value := range res {
		var records response.MarketAndNftView
		nftInfo := model.NftGetByTokenId(value.Tokenid, value.ChainType)
		// market 数据
		records.Id = value.Id
		records.SN = value.SN
		records.Creater = value.Creater
		records.Tokenid = value.Tokenid
		records.MarketType = value.MarketType
		records.StartingPrice = value.StartingPrice
		records.TokenType = value.TokenType
		records.EndTime = time.Unix(value.EndTime, 0).Format("2006-01-02 15:04:05")
		records.Buyer = value.Buyer
		records.Bonus = value.Bonus
		records.Txhash = value.Txhash
		records.CancelHash = value.CancelHash
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.Sorting = value.Sorting
		records.Status = value.Status
		// nft 数据
		records.NftName = nftInfo[0].NftName
		records.NftDesc = nftInfo[0].NftDesc
		records.TokenId = nftInfo[0].TokenId
		records.TokenUri = nftInfo[0].TokenUri
		records.NftTxhash = nftInfo[0].Txhash
		records.NftCreater = nftInfo[0].Creater
		records.CreateNumber = nftInfo[0].CreateNumber
		records.Nft_CreateTime = nftInfo[0].CreateTime
		records.MediaUri = nftInfo[0].MediaUri
		records.CreateTax = nftInfo[0].CreateTax
		records.Owner = nftInfo[0].Owner
		records.NftType = nftInfo[0].NftType
		records.OfferCount = model.OrdersCountByMid(value.Id)
		records.Approved = nftInfo[0].Approved

		recordList = append(recordList, records)
	}
	// 按OfferCount值进行排序算法
	for i := 0; i < len(recordList); i++ {
		for j := 0; j < len(recordList)-i-1; j++ {
			if recordList[j].OfferCount < recordList[j+1].OfferCount {
				recordList[j], recordList[j+1] = recordList[j+1], recordList[j]
			}
		}
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// HomeBanners
func HomeBanners(ctx *gin.Context) {
	chain_type := ctx.Query("chain_type")
	res := model.BannerList(chain_type)
	response.ReturnSuccessResponse(ctx, res)
}

// HomeNotableCollections
func HomeNotable(ctx *gin.Context) {
	limits, _ := strconv.Atoi(ctx.Query("limit"))
	logger.Info("HomeTopCollections req: limit", limits)

	var res []models.TCollectionInfo
	var limitSql string
	if limits == 0 {
		limits = 10
	}
	limitSql = fmt.Sprintf("LIMIT %d", limits)
	res = models.RankingFindAll(limitSql)

	var recordList []response.HomeCollectionView
	for _, value := range res {
		var record response.HomeCollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"id": value.UserId})
		if len(ownerUserInfo) > 0 {
			record.OwnerName = ownerUserInfo[0].UserName
		}
		record.CollectionId = value.Id
		record.LogoImageURL = value.LogoImage
		record.BannerImageURL = value.BannerImageUrl
		record.CategoryName = value.Category
		record.CollectionName = value.CollectionName
		record.CollectionDesc = value.CollectionDesc
		record.Owner = value.UserId
		record.ChainName = value.ChainName
		recordList = append(recordList, record)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

func HomeTopCollections(ctx *gin.Context) {
	limits, _ := strconv.Atoi(ctx.Query("limit"))
	logger.Info("HomeTopCollections req: limit", limits)

	var res []models.TCollectionInfo
	var limitSql string
	if limits == 0 {
		limits = 15
	}
	limitSql = fmt.Sprintf("LIMIT %d", limits)
	res = models.RankingFindAll(limitSql)

	var recordList []response.HomeCollectionTop
	for _, value := range res {
		var record response.HomeCollectionTop
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"id": value.UserId})
		if len(ownerUserInfo) > 0 {
			record.OwnerName = ownerUserInfo[0].UserName
		}
		record.Owner = value.UserId
		record.CollectionId = value.Id
		record.LogoImageURL = value.LogoImage
		record.CollectionName = value.CollectionName
		record.CollectionDesc = value.CollectionDesc
		record.CreateTax = value.CreateTax
		record.Amount = value.Amount
		recordList = append(recordList, record)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// HomeTrending
func HomeTrending(ctx *gin.Context) {
	category := ctx.Query("category")
	limits, _ := strconv.Atoi(ctx.Query("limit"))
	var limitSql string
	if limits == 0 {
		limits = 10
	}
	limitSql = fmt.Sprintf("LIMIT %d", limits)
	logger.Info("HomeTrending req: limit", limits)

	var res []models.TCollectionInfo
	res = models.RankingByCondition(category, "", limitSql)

	var recordList []response.HomeCollectionView
	for _, value := range res {
		var record response.HomeCollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"id": value.UserId})
		if len(ownerUserInfo) > 0 {
			record.OwnerName = ownerUserInfo[0].UserName
		}
		record.CollectionId = value.Id
		record.LogoImageURL = value.LogoImage
		record.CategoryName = value.Category
		record.Owner = value.UserId
		record.ChainName = value.ChainName
		record.CollectionDesc = value.CollectionDesc
		record.BannerImageURL = value.BannerImageUrl
		recordList = append(recordList, record)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}
