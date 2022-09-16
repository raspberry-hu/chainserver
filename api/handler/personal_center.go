package handler

import (
	"ChainServer/api/response"
	"ChainServer/common"
	"ChainServer/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func PersonalCenterSelling(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	donation := ctx.Query("donation")
	user_id, _ := strconv.Atoi(ctx.Query("id"))
	market_type := ctx.Query("market_type")
	state := ctx.Query("state")
	currency_name := ctx.Query("currency_name")
	market_creater := ctx.Query("market_creater")

	starting_price := ctx.Query("starting_price")
	end_time := ctx.Query("end_time")
	favorites_count := ctx.Query("favorites_count")
	create_time := ctx.Query("create_time")

	order_params_map := make(map[string]string)
	order_params_map["starting_price"] = starting_price
	order_params_map["end_time"] = end_time
	order_params_map["favorites_count"] = favorites_count
	order_params_map["create_time"] = create_time
	order_sql := MarketAllOrderByParams(order_params_map)
	fmt.Println("order ddddddddddddd", order_sql)

	var params_map = make(map[string]string)
	params_map["donation"] = donation
	params_map["market_type"] = market_type
	params_map["currency_name"] = currency_name
	params_map["state"] = state
	params_map["status"] = "1"
	params_map["creater"] = market_creater
	params_sql := ScreeningParams(params_map)

	limit_sql := fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	data_arry := models.MarketAllFind(user_id, params_sql, order_sql, limit_sql)
	count := models.MarketAllCount(params_sql)
	result := []map[string]interface{}{}
	for _, value := range data_arry {
		_, max_price := models.OrdersGetMaxPrice(fmt.Sprintf("%d", value.MarketId))
		var info = make(map[string]interface{})
		data, _ := json.Marshal(&value)
		json.Unmarshal(data, &info)
		info["end_time"] = time.Unix(int64(info["end_time"].(float64)), 0).Format("2006-01-02 15:04:05")
		info["nft_create_time"] = time.Unix(int64(info["nft_create_time"].(float64)), 0).Format("2006-01-02 15:04:05")
		info["highest_bid"] = max_price
		result = append(result, info)
	}
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result,
		Code:   200,
		Count:  count,
	}
	ctx.JSON(http.StatusOK, response)
	//response.ReturnSuccessResponseNew(ctx, result)
}

func PersonalCenterOwner(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows

	//owner_id, _ := strconv.Atoi(ctx.Query("owner"))
	owner_name := ctx.Query("owner_name")
	//根据userId查找资产信息
	var nft_info []models.TNft
	var personal_total []models.TNft
	nft_info = models.NftFindByLimit(map[string]interface{}{"ant_nft_owner": owner_name}, rows, offset)
	personal_total = models.NftFind(map[string]interface{}{"ant_nft_owner": owner_name})

	var resplist []response.TokenView
	var resp response.TokenView

	var response common.ResponseNew
	if len(nft_info) > 0 {
		for _, value := range nft_info {
			resp.NftName = value.NftName
			resp.NftDesc = value.NftDesc
			resp.Creater = value.Creater
			resp.CreateTime = value.CreateTime
			resp.MediaUri = value.MediaUri
			resp.Owner = value.Owner
			resp.CollectionId = value.CollectionId
			resp.ExploreUri = value.ExploreUri
			resp.AntTokenId = value.AntTokenId
			resp.AntNftId = value.AntNftId
			resp.AntCount = value.AntCount
			resp.AntPrice = value.AntPrice
			resp.AntTokenUrl = value.AntTokenUrl
			resp.AntTxHash = value.AntTxHash
			resp.AntNftOwner = value.AntNftOwner
			resp.MarketType = value.MarketType
			resp.TransferHash = value.TransferHash

			//根据ID查询相应的用户名
			var originator, current []models.UserInfo
			originator = models.UserInfoFind(map[string]interface{}{"id": value.Creater})
			current = models.UserInfoFind(map[string]interface{}{"id": value.Owner})
			var collectionInfo []models.TCollectionInfo
			collectionInfo = models.CollectionFind(map[string]interface{}{"id": value.CollectionId})

			if len(originator) == 1 && len(current) == 1 && len(collectionInfo) == 1 {
				resp.ImageUrl = originator[0].ImageUrl
				resp.CreaterName = originator[0].UserName
				resp.OwnerImageUrl = current[0].ImageUrl
				resp.OwnerName = current[0].UserName
				resp.CollectionName = collectionInfo[0].CollectionName
			} else {
				log.Println("Not find account by id")
			}

			resplist = append(resplist, resp)
		}

		response = common.ResponseNew{
			Code:   200,
			Msg:    "success",
			Result: resplist,
			Count:  len(personal_total),
		}
	} else {
		response = common.ResponseNew{
			Msg:    "null",
			Result: resplist,
			Count:  len(personal_total),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func PersonalCenterFavorites(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	selling_state := ctx.Query("selling_state")
	donation := ctx.Query("donation")
	wallet_addr := ctx.Query("wallet_addr")
	favorites_addr := ctx.Query("favorites_addr")
	market_type := ctx.Query("market_type")
	state := ctx.Query("state")
	currency_name := ctx.Query("currency_name")
	market_creater := ctx.Query("market_creater")
	nft_creater := ctx.Query("nft_creater")
	nft_owner := ctx.Query("nft_owner")

	starting_price := ctx.Query("starting_price")
	end_time := ctx.Query("end_time")
	favorites_count := ctx.Query("favorites_count")
	create_time := ctx.Query("create_time")
	nft_create_time := ctx.Query("nft_create_time")

	order_params_map := make(map[string]string)
	order_params_map["starting_price"] = starting_price
	order_params_map["end_time"] = end_time
	order_params_map["favorites_count"] = favorites_count
	order_params_map["create_time"] = create_time
	order_params_map["nft_create_time"] = nft_create_time

	order_sql := MarketAllOrderByParams(order_params_map)

	var params_map = make(map[string]interface{})
	params_map["donation"] = donation
	params_map["market_type"] = market_type
	params_map["currency_name"] = currency_name
	params_map["state"] = state
	params_map["status"] = "1"
	params_map["creater"] = market_creater
	params_map["selling_state"] = selling_state
	params_map["nft_creater"] = nft_creater
	params_map["nft_owner"] = nft_owner

	params_map["f.wallet_addr"] = favorites_addr

	params_sql := PersonalCenterScreeningParams(params_map)

	limit_sql := fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	data_arry := models.FavoritesLeftNftLeftMarketList(wallet_addr, params_sql, order_sql, limit_sql)
	count := models.FavoritesLeftNftLeftMarketListCount(params_sql)
	result := []map[string]interface{}{}
	for _, value := range data_arry {
		_, max_price := models.OrdersGetMaxPrice(fmt.Sprintf("%d", value.MarketId))
		var info = make(map[string]interface{})
		data, _ := json.Marshal(&value)
		json.Unmarshal(data, &info)
		info["end_time"] = time.Unix(int64(info["end_time"].(float64)), 0).Format("2006-01-02 15:04:05")
		info["nft_create_time"] = time.Unix(int64(info["nft_create_time"].(float64)), 0).Format("2006-01-02 15:04:05")
		info["highest_bid"] = max_price
		WallectInfoAdd(value, info)
		collection_info := CollectionInfoAdd(value.CollectionId)
		info["collection_info"] = collection_info
		result = append(result, info)
	}
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result,
		Code:   200,
		Count:  count,
	}
	ctx.JSON(http.StatusOK, response)
}

func PersonalCenterScreeningParams(params_map map[string]interface{}) string {
	var params_sql = "where "
	for k, v := range params_map {
		if v == "" {
			continue
		}
		if k == "state" {
			now := time.Now().Unix()
			if v == "1" {
				params_sql += fmt.Sprintf("m.end_time > %d and ", now)
			} else if v == "0" {
				params_sql += fmt.Sprintf("m.end_time < %d and ", now)
			}
		} else if k == "currency_name" || k == "creater" {
			params_sql += fmt.Sprintf("n.%s = '%s' and ", k, v)
		} else if k == "selling_state" {
			if v == "0" {
				params_sql += "m.market_type is not NUll and "
			} else if v == "1" {
				params_sql += "m.market_type is NUll and "
			}
		} else if k == "nft_owner" {
			params_sql += fmt.Sprintf("n.owner = '%d' and ", v)
		} else if k == "nft_creater" {
			params_sql += fmt.Sprintf("n.creater = '%d' and ", v)
		} else if k == "donation" {
			params_sql += fmt.Sprintf("m.%s = '%s' and ", k, v)
		} else if k == "common_search" {
			params_sql += fmt.Sprintf("n.nft_name like '%%%s%%' and ", v)
		} else if k == "f.wallet_addr" {
			params_sql += fmt.Sprintf("%s = '%s' and ", k, v)
		} else {
			params_sql += fmt.Sprintf("n.%s = '%s' and ", k, v)
		}
	}
	if params_sql == "where " {
		params_sql = ""
	} else {
		params_sql = strings.TrimRight(params_sql, "and ")
	}
	return params_sql
}
