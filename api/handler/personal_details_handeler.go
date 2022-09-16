package handler

import (
	"ChainServer/api/response"
	"ChainServer/models"
	"ChainServer/models/model"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PersonalDetails(ctx *gin.Context) {
	//address := ctx.Query("address")
	//market_status := ctx.Query("market_status")
	user_id, _ := strconv.Atoi(ctx.Query("user_id"))
	//page, _ := strconv.Atoi(ctx.Query("page"))
	//offset, _ := strconv.Atoi(ctx.Query("offset"))
	//limit_sql := ""
	//chain_type := ctx.Query("chain_type")
	// 已经售卖

	var screen_sql string
	screen_sql = fmt.Sprintf(` where n.owner = '%d' and n.market_status != 0 and n.status = 1`, user_id)
	selling_arry := model.Test(screen_sql, "", "")
	var selling_arry_result []model.MarketAndNftAll
	for _, selling := range *selling_arry {
		// 创建者用户信息
		info_arry := models.UserInfoFind(map[string]interface{}{"id": selling.Creater})
		if len(info_arry) > 0 {
			selling.UserName = info_arry[0].UserName
			selling.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := models.UserInfoFind(map[string]interface{}{"id": selling.Owner})
		if len(owner_info_arry) > 0 {
			selling.OwnerUserName = owner_info_arry[0].UserName
			selling.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}
		selling_arry_result = append(selling_arry_result, selling)
	}

	//selling_arry :=  model.MarketListCommonFind(page, offset, fmt.Sprintf(`owner = '%s' and market_status != 0 and status = 1`,wallet_addr))

	//fmt.Println(selling_arry)
	// 拥有的
	nft_owner_arry := model.TestRight(fmt.Sprintf("where n.owner = '%d' and n.status = 1", user_id), "", "")
	var nft_owner_result []model.MarketAndNftAll
	for _, nft_owner := range *nft_owner_arry {
		// 创建者用户信息
		info_arry := models.UserInfoFind(map[string]interface{}{"id": nft_owner.Creater})
		if len(info_arry) > 0 {
			nft_owner.UserName = info_arry[0].UserName
			nft_owner.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := models.UserInfoFind(map[string]interface{}{"id": nft_owner.Owner})
		if len(owner_info_arry) > 0 {
			nft_owner.OwnerUserName = owner_info_arry[0].UserName
			nft_owner.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}
		nft_owner_result = append(nft_owner_result, nft_owner)
	}

	// 已创建的
	nft_creater_arry := model.TestRight(fmt.Sprintf("where n.creater = '%d' and n.status = 1", user_id), "", "")
	var nft_creater_result []model.MarketAndNftAll
	for _, nft_create := range *nft_creater_arry {
		// 创建者用户信息
		info_arry := models.UserInfoFind(map[string]interface{}{"id": nft_create.Creater})
		if len(info_arry) > 0 {
			nft_create.UserName = info_arry[0].UserName
			nft_create.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := models.UserInfoFind(map[string]interface{}{"id": nft_create.Owner})
		if len(owner_info_arry) > 0 {
			nft_create.OwnerUserName = owner_info_arry[0].UserName
			nft_create.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}
		nft_creater_result = append(nft_owner_result, nft_create)
	}

	fmt.Println(nft_creater_arry)
	var result = make(map[string]interface{})
	result["selling_arry"] = selling_arry_result
	result["owner_arrt"] = nft_owner_result
	result["creater_arry"] = nft_creater_result
	response.ReturnSuccessResponse(ctx, result)
}

func PersonalDetailsNew(ctx *gin.Context) {
	//address := ctx.Query("address")
	//market_status := ctx.Query("market_status")
	wallet_addr := ctx.Query("walletAddr")
	page, _ := strconv.Atoi(ctx.Query("page"))
	offset, _ := strconv.Atoi(ctx.Query("rows"))
	end_time := ctx.Query("end_time")
	starting_price := ctx.Query("starting_price")
	create_time := ctx.Query("create_time")
	var order_params string
	if end_time != "" {
		order_params = fmt.Sprintf("end_time %s", end_time)
	} else if starting_price != "" {
		order_params = fmt.Sprintf("starting_price %s", starting_price)
	} else if create_time != "" {
		order_params = fmt.Sprintf("create_time %s", create_time)
	}
	//var
	//chain_type := ctx.Query("chain_type")
	// 已经售卖
	// MarketAndNftAll
	market_sql := fmt.Sprintf(`creater = '%s' and status = 1`, wallet_addr)
	//market_sql += order_sql
	market_list := model.MarketListCommonFind(page, offset, market_sql, order_params)
	var selling_arry []model.MarketAndNftAll
	for _, market := range market_list {
		var market_and_nft_all model.MarketAndNftAll
		market_and_nft_all.Id = market.Id
		market_and_nft_all.SN = market.SN
		market_and_nft_all.Creater = market.Creater
		market_and_nft_all.Tokenid = market.Tokenid
		market_and_nft_all.MarketType = market.MarketType
		market_and_nft_all.StartingPrice = market.StartingPrice
		market_and_nft_all.TokenType = market.TokenType
		market_and_nft_all.EndTime = time.Unix(market.EndTime, 0).Format("2006-01-02 15:04:05")
		market_and_nft_all.CreateTime = time.Unix(market.CreateTime, 0).Format("2006-01-02 15:04:05")
		market_and_nft_all.Buyer = market.Buyer
		market_and_nft_all.Bonus = market.Bonus
		market_and_nft_all.Txhash = market.Txhash
		market_and_nft_all.CancelHash = market.CancelHash
		market_and_nft_all.DealHash = market.DealHash
		market_and_nft_all.Sorting = market.Sorting
		market_and_nft_all.Status = market.Status
		market_and_nft_all.ChainType = market.ChainType
		market_and_nft_all.FreeGas = market.FreeGas
		market_and_nft_all.Donation = market.Donation
		market_and_nft_all.NftId = market.NftId

		nft_list := model.NftCommonFind(page, offset, map[string]interface{}{"id": market.NftId})
		if len(nft_list) > 0 {
			nft := nft_list[0]
			market_and_nft_all.NftName = nft.NftName
			market_and_nft_all.NftDesc = nft.NftDesc
			market_and_nft_all.TokenId = nft.TokenId
			market_and_nft_all.TokenUri = nft.TokenUri
			market_and_nft_all.NftTxhash = nft.TxHash
			market_and_nft_all.NftCreater = nft.Creater
			market_and_nft_all.CreateNumber = nft.CreateNumber
			market_and_nft_all.Nft_CreateTime = time.Unix(int64(nft.CreateTime), 0).Format("2006-01-02 15:04:05")
			market_and_nft_all.MediaUri = nft.MediaUri
			market_and_nft_all.CreateTax = nft.CreateTax
			market_and_nft_all.Owner = nft.Owner
			market_and_nft_all.NftType = nft.NftType
			market_and_nft_all.Approved = nft.Approved
			info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft.Creater})
			if len(info_arry) > 0 {
				market_and_nft_all.UserName = info_arry[0].UserName
				market_and_nft_all.ImageUrl = info_arry[0].ImageUrl
			}
			// 拥有者信息
			owner_info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft.Owner})
			if len(owner_info_arry) > 0 {
				market_and_nft_all.OwnerUserName = owner_info_arry[0].UserName
				market_and_nft_all.OwnerImageUrl = owner_info_arry[0].ImageUrl
			}
		}
		selling_arry = append(selling_arry, market_and_nft_all)
	}
	fmt.Println(selling_arry)
	// 拥有的
	nft_owner_arry := model.NftCommonFind(page, offset, map[string]interface{}{"owner": wallet_addr})
	var nft_owner_result []model.MarketAndNftAll
	fmt.Println(nft_owner_arry)
	for _, nft_owner := range nft_owner_arry {
		var market_and_nft_all model.MarketAndNftAll
		market_and_nft_all.NftName = nft_owner.NftName
		market_and_nft_all.NftDesc = nft_owner.NftDesc
		market_and_nft_all.TokenId = nft_owner.TokenId
		market_and_nft_all.TokenUri = nft_owner.TokenUri
		market_and_nft_all.NftTxhash = nft_owner.TxHash
		market_and_nft_all.NftCreater = nft_owner.Creater
		market_and_nft_all.CreateNumber = nft_owner.CreateNumber
		market_and_nft_all.Nft_CreateTime = time.Unix(int64(nft_owner.CreateTime), 0).Format("2006-01-02 15:04:05")
		//int64(nft_owner.CreateTime)
		market_and_nft_all.MediaUri = nft_owner.MediaUri
		market_and_nft_all.CreateTax = nft_owner.CreateTax
		market_and_nft_all.Owner = nft_owner.Owner
		market_and_nft_all.NftType = nft_owner.NftType
		market_and_nft_all.ChainType = nft_owner.ChainType
		market_and_nft_all.TokenType = nft_owner.TokenType
		market_and_nft_all.FreeGas = nft_owner.FreeGas

		info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft_owner.Creater})
		if len(info_arry) > 0 {
			market_and_nft_all.UserName = info_arry[0].UserName
			market_and_nft_all.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft_owner.Owner})
		if len(owner_info_arry) > 0 {
			market_and_nft_all.OwnerUserName = owner_info_arry[0].UserName
			market_and_nft_all.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}
		market_owner_list := model.MarketListCommonFind(page, offset, fmt.Sprintf(`nft_id = %d and status = 1 `, nft_owner.ID), "")
		if len(market_owner_list) > 0 {
			market := market_owner_list[0]
			market_and_nft_all.Id = market.Id
			market_and_nft_all.SN = market.SN
			market_and_nft_all.Creater = market.Creater
			market_and_nft_all.Tokenid = market.Tokenid
			market_and_nft_all.MarketType = market.MarketType
			market_and_nft_all.StartingPrice = market.StartingPrice
			market_and_nft_all.EndTime = time.Unix(market.EndTime, 0).Format("2006-01-02 15:04:05")
			market_and_nft_all.CreateTime = time.Unix(market.CreateTime, 0).Format("2006-01-02 15:04:05")
			market_and_nft_all.Buyer = market.Buyer
			market_and_nft_all.Bonus = market.Bonus
			market_and_nft_all.Txhash = market.Txhash
			market_and_nft_all.CancelHash = market.CancelHash
			market_and_nft_all.DealHash = market.DealHash
			market_and_nft_all.Sorting = market.Sorting
			market_and_nft_all.Status = market.Status
			market_and_nft_all.Donation = market.Donation
			market_and_nft_all.NftId = market.NftId

		}
		nft_owner_result = append(nft_owner_result, market_and_nft_all)
	}

	// 已创建的
	nft_creater_arry := model.NftCommonFind(page, offset, map[string]interface{}{"creater": wallet_addr})
	var nft_creater_result []model.MarketAndNftAll
	for _, nft_owner := range nft_creater_arry {
		var market_and_nft_all model.MarketAndNftAll
		market_and_nft_all.NftName = nft_owner.NftName
		market_and_nft_all.NftDesc = nft_owner.NftDesc
		market_and_nft_all.TokenId = nft_owner.TokenId
		market_and_nft_all.TokenUri = nft_owner.TokenUri
		market_and_nft_all.NftTxhash = nft_owner.TxHash
		market_and_nft_all.NftCreater = nft_owner.Creater
		market_and_nft_all.CreateNumber = nft_owner.CreateNumber
		market_and_nft_all.Nft_CreateTime = time.Unix(int64(nft_owner.CreateTime), 0).Format("2006-01-02 15:04:05")
		market_and_nft_all.MediaUri = nft_owner.MediaUri
		market_and_nft_all.CreateTax = nft_owner.CreateTax
		market_and_nft_all.Owner = nft_owner.Owner
		market_and_nft_all.NftType = nft_owner.NftType
		market_and_nft_all.Approved = nft_owner.Approved
		market_and_nft_all.FreeGas = nft_owner.FreeGas
		market_and_nft_all.ChainType = nft_owner.ChainType
		market_and_nft_all.TokenType = nft_owner.TokenType
		info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft_owner.Creater})
		if len(info_arry) > 0 {
			market_and_nft_all.UserName = info_arry[0].UserName
			market_and_nft_all.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := model.TWalletInfoFind(map[string]interface{}{"wallet_addr": nft_owner.Owner})
		if len(owner_info_arry) > 0 {
			market_and_nft_all.OwnerUserName = owner_info_arry[0].UserName
			market_and_nft_all.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}
		market_owner_list := model.MarketListCommonFind(page, offset, fmt.Sprintf(`nft_id = %d and status = 1 `, nft_owner.ID), "")
		if len(market_owner_list) > 0 {
			market := market_owner_list[0]
			market_and_nft_all.Id = market.Id
			market_and_nft_all.SN = market.SN
			market_and_nft_all.Creater = market.Creater
			market_and_nft_all.Tokenid = market.Tokenid
			market_and_nft_all.MarketType = market.MarketType
			market_and_nft_all.StartingPrice = market.StartingPrice
			market_and_nft_all.EndTime = time.Unix(market.EndTime, 0).Format("2006-01-02 15:04:05")
			market_and_nft_all.CreateTime = time.Unix(market.CreateTime, 0).Format("2006-01-02 15:04:05")
			market_and_nft_all.Buyer = market.Buyer
			market_and_nft_all.Bonus = market.Bonus
			market_and_nft_all.Txhash = market.Txhash
			market_and_nft_all.CancelHash = market.CancelHash
			market_and_nft_all.DealHash = market.DealHash
			market_and_nft_all.Sorting = market.Sorting
			market_and_nft_all.Status = market.Status
			market_and_nft_all.Donation = market.Donation
			market_and_nft_all.NftId = market.NftId

		}
		nft_creater_result = append(nft_creater_result, market_and_nft_all)
	}
	var result = make(map[string]interface{})
	result["selling_arry"] = selling_arry
	result["owner_arrt"] = nft_owner_result
	result["creater_arry"] = nft_creater_result
	response.ReturnSuccessResponse(ctx, result)
}
