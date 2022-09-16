package handler

import (
	"ChainServer/api/request"
	"ChainServer/api/response"
	"ChainServer/common"
	"ChainServer/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

func CollectionInfoCreate(ctx *gin.Context) {
	req := request.CollectionCreate{}
	ctx.ShouldBindJSON(&req)
	var collection_info models.TCollectionInfo
	collection_info.UserId = req.UserId
	collection_info.CollectionName = req.CollectionName
	collection_info.ChainName = req.ChainName
	collection_info.LogoImage = req.LogoImage
	collection_info.FeaturedImageUrl = req.FeaturedImageUrl
	// 判定banner是否为空
	if req.BannerImageUrl == "" {
		collection_info.BannerImageUrl = "http://101.43.151.77:8989/download/jupiter/img/4_1650101172_440619.png"
	} else {
		collection_info.BannerImageUrl = req.BannerImageUrl
	}
	collection_info.CollectionDesc = req.CollectionDesc
	collection_info.CurrencyName = req.CurrencyName
	collection_info.Status = 1
	collection_info.Category = req.Category
	collection_info.CreateTax = req.CreateTax
	collection_id := models.CollectionInsert(collection_info)

	var categories_info models.TCategoriesInfo
	categories_info.CategoriesName = collection_info.Category
	categories_info.CollectionId = collection_id
	categories_info.UserId = req.UserId
	models.CategoriesInsert(categories_info)

	response.ReturnSuccessResponseNew(ctx, gin.H{
		"msg":          "created successfully",
		"collectionId": collection_id,
	})
}

func CollectionListGet(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	collection_name := ctx.Query("common_search")
	// wallet_addr_list := models.CollectionLikeFind(collection_name)
	// UserId_list := models.CollectionLikeFind(collection_name)
	var contains_sql string
	var limit_sql string
	if collection_name != "" {
		contains_sql += fmt.Sprintf("where C.collection_name like '%%%s%%'", collection_name)
	}

	if page != 0 {
		limit_sql = fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	}
	//contains_sql, limit_sql :=  JoiningSql(page,rows,wallet_addr_list)
	collection_list_arry := models.CollectionListFind(contains_sql, limit_sql)
	// count
	count := models.CollectionCount(contains_sql)
	response := &common.ResponseNew{
		Msg:    "success",
		Result: collection_list_arry,
		Code:   200,
		Count:  count,
	}
	ctx.JSON(http.StatusOK, response)
}

func JoiningSql(page, rows int, user_id_list []models.UserIdList) (string, string) {
	var limit_sql string
	var contains_sql string
	offset := (page - 1) * rows
	if len(user_id_list) == 0 {
		page = 0
	} else {
		for _, user_id := range user_id_list {
			contains_sql += fmt.Sprintf(`'%d', `, user_id.UserId)
		}

		contains_sql = strings.TrimRight(contains_sql, ", ")
		contains_sql = fmt.Sprintf(`where c.user_id in (%s)`, contains_sql)
	}
	if page != 0 {
		limit_sql = fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	}
	return contains_sql, limit_sql
}

// ExploreCollections
func ExploreCollections(ctx *gin.Context) {
	category := ctx.Query("tab")
	var limit_sql string
	// 当前的页码
	page, _ := strconv.Atoi(ctx.Query("page"))
	// 每页所要展示的条数 默认为10
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	if rows == 0 {
		rows = 10
	}
	// 数据库查询的起始下标位置
	offset := (page - 1) * rows
	if page != 0 {
		limit_sql = fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	}

	var res []models.TCollectionInfo
	if category != "" {
		res = models.RankingByCondition(category, "", limit_sql)
	} else {
		res = models.RankingFindAll(limit_sql)
	}

	var recordList []response.HomeCollectionView
	for _, value := range res {
		var records response.HomeCollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"user_id": value.UserId})
		if len(ownerUserInfo) > 0 {
			records.OwnerName = ownerUserInfo[0].UserName
		}
		records.CollectionId = value.Id
		records.LogoImageURL = value.LogoImage
		records.CollectionName = value.CollectionName
		records.BannerImageURL = value.BannerImageUrl
		records.CollectionDesc = value.CollectionDesc
		records.Owner = value.UserId
		records.CategoryName = value.Category
		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

func GeneralCollectionInfo(ctx *gin.Context) {
	collectionId := ctx.Query("collection_id")
	var recordList []response.CollectionView
	var res []models.TCollectionInfo
	res = models.CollectionFind(map[string]interface{}{"id": collectionId})
	for _, value := range res {
		var record response.CollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"user_id": value.UserId})
		if len(ownerUserInfo) > 0 {
			record.OwnerName = ownerUserInfo[0].UserName
		}
		record.CollectionId = value.Id
		record.CollectionName = value.CollectionName
		record.CollectionDesc = value.CollectionDesc
		record.Owner = value.UserId
		record.LogoImageURL = value.LogoImage
		record.BannerImageURL = value.BannerImageUrl
		record.CurrencyName = value.CurrencyName
		record.CreateTax = value.CreateTax
		record.Items = value.Items
		record.Favorites = value.Favorites
		record.Amount = value.Amount
		record.ChainName = value.ChainName
		record.CategoryName = value.Category
		recordList = append(recordList, record)
	}
	response.ReturnSuccessResponse(ctx, gin.H{
		"msg":    "get successfully",
		"result": recordList,
	})
}

func PersonalCollection(ctx *gin.Context) {
	UserId, _ := strconv.Atoi(ctx.Query("user_id"))
	var limitSql string
	// 当前的页码
	page, _ := strconv.Atoi(ctx.Query("page"))
	// 每页所要展示的条数 默认为10
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	if rows == 0 {
		rows = 15
	}
	// 数据库查询的起始下标位置
	offset := (page - 1) * rows
	if page != 0 {
		limitSql = fmt.Sprintf("WHERE user_id='%d' LIMIT %d OFFSET %d", UserId, rows, offset)
	}

	res := models.RankingFindAll(limitSql)
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
		records.Owner = value.UserId
		records.Items = value.Items
		records.CreateTax = value.CreateTax
		records.Favorites = value.Favorites
		records.Amount = value.Amount
		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// ExistedCollection铸币页面获取账户已创建的collection
func ExistedCollection(ctx *gin.Context) {
	UserId := ctx.Query("user_id")
	var recordList []response.CollectionView
	var res []models.TCollectionInfo
	res = models.CollectionFind(map[string]interface{}{"user_id": UserId})
	for _, value := range res {
		var record response.CollectionView
		// 获取拥有着的名字和地址
		var ownerUserInfo []models.UserInfo
		ownerUserInfo = models.UserInfoFind(map[string]interface{}{"user_id": UserId})
		if len(ownerUserInfo) > 0 {
			record.OwnerName = ownerUserInfo[0].UserName
		}
		record.CollectionId = value.Id
		record.CollectionName = value.CollectionName
		record.CollectionDesc = value.CollectionDesc
		record.Owner = value.UserId
		record.LogoImageURL = value.LogoImage
		record.BannerImageURL = value.BannerImageUrl
		record.CurrencyName = value.CurrencyName
		record.CreateTax = value.CreateTax
		record.Items = value.Items
		record.Favorites = value.Favorites
		record.Amount = value.Amount
		record.ChainName = value.ChainName
		record.CategoryName = value.Category
		recordList = append(recordList, record)
	}
	response.ReturnSuccessResponse(ctx, gin.H{
		"msg":    "get successfully",
		"result": recordList,
	})
}

// PersonalNFTFromCollection 展示个人collection中的资产
func PersonalNFTFromCollection(ctx *gin.Context) {
	collectionId, _ := strconv.Atoi(ctx.Query("collection_id"))
	status, _ := strconv.Atoi(ctx.Query("status"))
	currencyName := ctx.Query("currency")
	minPrice := ctx.Query("min_price")
	maxPrice := ctx.Query("max_price")
	name := ctx.Query("name")

	var nftInfo []models.TNft
	var martketInfo []models.TMarketList
	var recordList []response.NftAndMarketView
	var collectionInfo []models.TCollectionInfo
	collectionInfo = models.CollectionFind(map[string]interface{}{"id": collectionId})
	if len(collectionInfo) == 0 {
		logger.Error("Not find the collection")
	}

	// 默认查找该collection下所有挂单和未挂单的资产
	if collectionId != 0 {

		// 先查找该集合下创建并满足条件的nft
		if status == 0 && currencyName == "" && name == "" {
			nftInfo = models.NftFind(map[string]interface{}{"collection_id": collectionId, "owner": collectionInfo[0].UserId})
		} else {
			nftInfo = models.CollectionNftFind(collectionId, status, currencyName, name, collectionInfo[0].UserId)
		}
		for _, temp := range nftInfo {
			if temp.Status == 1 && temp.Owner == collectionInfo[0].UserId {
				//查找初次挂单未成交的资产
				if minPrice == "" && maxPrice == "" {
					martketInfo = models.MarketListFind(map[string]interface{}{"nft_id": temp.Id, "market_type": temp.MarketType, "creater": temp.Owner})
				} else {
					martketInfo = models.MarketCollectionListFind(temp.Id, temp.Owner, minPrice, maxPrice)
				}

				var marketRecord response.NftAndMarketView
				var nftRecord response.NftAndMarketView
				if len(martketInfo) != 0 {
					marketRecord.MId = martketInfo[0].Id
					marketRecord.MarketCreateTime = time.Unix(martketInfo[0].CreateTime, 0).Format("2006-01-02 15:04:05")
					marketRecord.StartingPrice = martketInfo[0].StartingPrice
					marketRecord.TokenType = martketInfo[0].TokenId
					marketRecord.EndTime = time.Unix(martketInfo[0].EndTime, 0).Format("2006-01-02 15:04:05")
					marketRecord.Bonus = martketInfo[0].Reward
					marketRecord.ChainType = martketInfo[0].ChainName
					marketRecord.OrderStatus = martketInfo[0].Status
					marketRecord.FreeGas = strconv.Itoa(martketInfo[0].Lazy)
					marketRecord.Donation = martketInfo[0].Donation
					// 获取创建者和拥有着的名字和地址
					var creatorUserInfo, ownerUserInfo []models.UserInfo
					creatorUserInfo = models.UserInfoFind(map[string]interface{}{"ID": temp.Creater})
					ownerUserInfo = models.UserInfoFind(map[string]interface{}{"ID": temp.Owner})
					fmt.Println(creatorUserInfo)
					if len(creatorUserInfo) != 0 {
						marketRecord.UserName = creatorUserInfo[0].UserName
						marketRecord.ImageUrl = creatorUserInfo[0].ImageUrl
					}
					if len(ownerUserInfo) != 0 {
						marketRecord.OwnerUserName = ownerUserInfo[0].UserName
						marketRecord.OwnerImageUrl = ownerUserInfo[0].ImageUrl
					}
					marketRecord.NftId = temp.Id
					marketRecord.NftName = temp.NftName
					marketRecord.NftDesc = temp.NftDesc
					marketRecord.TokenId = temp.TokenId
					marketRecord.MediaUri = temp.MediaUri
					marketRecord.ExploreUri = temp.ExploreUri
					marketRecord.MediaIpfsUri = temp.MediaIpfsUri
					marketRecord.Creater = temp.Creater
					marketRecord.Owner = temp.Owner
					marketRecord.CreateTax = temp.CreateTax
					marketRecord.CurrencyName = temp.CurrencyName
					// 将时间转换成字符串的形式
					marketRecord.CreateTime = time.Unix(temp.CreateTime, 0).Format("2006-01-02 15:04:05")
					marketRecord.CollectionId = temp.CollectionId
					marketRecord.Status = temp.Status
					marketRecord.MarKetStatus = temp.MarketType
					marketRecord.ChainType = temp.ChainName
					marketRecord.CreateNumber = temp.BlockNumber
					marketRecord.TokenUri = temp.MetaDataUri
					marketRecord.Txhash = temp.TxHash

					// 查看该资产被收藏的个数
					var favorites []models.TFavorites
					favorites = models.FavoritesFind(map[string]interface{}{"nft_id": temp.Id})
					marketRecord.FavoritesCount = strconv.Itoa(len(favorites))
					recordList = append(recordList, marketRecord)
				} else if minPrice == "" && maxPrice == "" {
					// 查看未挂单的资产
					nftRecord.NftId = temp.Id
					nftRecord.NftName = temp.NftName
					nftRecord.NftDesc = temp.NftDesc
					nftRecord.TokenId = temp.TokenId
					nftRecord.MediaUri = temp.MediaUri
					nftRecord.MediaIpfsUri = temp.MediaIpfsUri
					nftRecord.ExploreUri = temp.ExploreUri
					nftRecord.Creater = temp.Creater
					nftRecord.Owner = temp.Owner
					nftRecord.CreateTax = temp.CreateTax
					nftRecord.CurrencyName = temp.CurrencyName
					// 将时间转换成字符串的形式
					nftRecord.CreateTime = time.Unix(temp.CreateTime, 0).Format("2006-01-02 15:04:05")
					nftRecord.CollectionId = temp.CollectionId
					nftRecord.Status = temp.Status
					nftRecord.MarKetStatus = temp.MarketType
					nftRecord.ChainType = temp.ChainName
					nftRecord.CreateNumber = temp.BlockNumber
					nftRecord.TokenUri = temp.MetaDataUri
					nftRecord.Txhash = temp.TxHash
					recordList = append(recordList, nftRecord)
				}

			}
		}
	}

	response.ReturnSuccessResponse(ctx, gin.H{
		"msg":    "success",
		"result": recordList,
	})
}
