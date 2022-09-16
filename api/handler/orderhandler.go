package handler

import (
	"ChainServer/api/request"
	"ChainServer/api/response"
	"ChainServer/common"
	"ChainServer/models"
	"ChainServer/models/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

// OrderNew
func OrderNew(ctx *gin.Context) {
	req := request.OrderNew{}
	ctx.ShouldBindJSON(&req)
	data, _ := json.Marshal(&req)
	logger.Info("order new ", string(data))
	market_id := req.MarketId
	nft_id := req.NftId
	market_type := req.MarketType
	token_id := req.TokenId
	tx_hash := req.TxHash
	create_time := int(time.Now().Unix()) //req.CreateTime
	buyer := req.Buyer
	seller := req.Seller
	price := req.Price
	currency_name := req.CurrencyName
	chain_name := req.ChainName
	lazy := req.Lazy
	donation := req.Donation
	var order_info models.TOrder
	order_info.MarketId = market_id
	order_info.NftId = nft_id
	order_info.MarketType = market_type
	order_info.TokenId = token_id
	order_info.CreateTime = create_time
	order_info.Buyer = buyer
	order_info.Seller = seller
	order_info.Price = price
	order_info.CurrencyName = currency_name
	order_info.Lazy = lazy
	order_info.ChainName = chain_name
	order_info.Donation = donation
	order_info.Status = 0
	order_info.TxHash = tx_hash
	order_id := models.OrderInsert(order_info)

	var collection_info models.TCollectionInfo
	nftInfo := models.NftFind(map[string]interface{}{"id": nft_id})
	collection_info.Id = nftInfo[0].CollectionId

	var result = make(map[string]interface{})

	//2-限价买卖成交时，查询当前collection中的信息
	if market_type == 2 {
		var amountCount float32
		var itemCount int

		collectionArr := models.CollectionFind(map[string]interface{}{"id": collection_info.Id})
		if len(collectionArr) == 0 {
			logger.Error("没找到Collection")
		} else {
			amountCount = collectionArr[0].Amount
			itemCount = collectionArr[0].Items
		}

		// 每完成一笔交易，对应的items减一
		if itemCount > 0 {
			itemCount = itemCount - 1
		}

		//字符串转成float类型
		priceNum, err := strconv.ParseFloat(price, 32)
		if err != nil {
			fmt.Sprintf("Error: %v", err)
		}
		amountCount = amountCount + float32(priceNum)
		//models.UpdateAmountByCollection(collection_info, map[string]interface{}{"amount": amountCount})
		//models.UpdateItemsByCollection(collection_info, map[string]interface{}{"items": itemCount})
		models.CollectionUpdate(collection_info, map[string]interface{}{"amount": amountCount, "items": itemCount})
		result["buyer"] = buyer
		result["max_price"] = price
	}
	//1-拍卖时，相继竞拍出价

	result["order_id"] = order_id
	result["collection_id"] = collection_info.Id
	response.ReturnSuccessResponseNew(ctx, result)
}

// OrderHashUpdate
func OrderHashUpdate(ctx *gin.Context) {
	req := request.UpdateHash{}
	ctx.ShouldBindJSON(&req)

	tx_hash := req.TxHash
	order_id := req.OrderId
	user_id := strings.ToLower(ctx.Request.Header.Get("user_id"))

	order_arry := models.OrderFind(map[string]interface{}{"id": order_id, "buyer": user_id})
	if len(order_arry) == 0 {
		logger.Error("找不到order")
	} else {
		if order_arry[0].TxHash != "" {
			logger.Error("order txhash 已存在")
		} else {
			models.OrderUpdate(order_arry[0], map[string]interface{}{"tx_hash": tx_hash})
		}
	}
	response.ReturnSuccessResponseNew(ctx, nil)
}

// OrderGetByOid
func OrderGetByOid(ctx *gin.Context) {
	oid := ctx.Query("oid")
	chain_type := ctx.Query("chain_type")
	res := model.OrderGetByOid(oid, chain_type)
	var recordList []response.OrderView
	for _, value := range res {
		var records response.OrderView

		records.Mid = value.Mid
		records.Tokenid = value.Tokenid
		records.Txhash = value.Txhash
		records.Buyer = value.Buyer
		records.TokenType = value.TokenType
		records.Seller = value.Seller
		records.Price = value.Price
		records.TokenType = value.TokenType
		records.Status = value.Status
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.OnchainTime = time.Unix(value.OnchainTime, 0).Format("2006-01-02 15:04:05")
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation

		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// OrderGetByTokenId
func OrderGetByMarketId(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	market_id := ctx.Query("market_id")
	order_arry := models.OrderFindLimit(map[string]interface{}{"market_id": market_id, "status": 1}, rows, offset)
	result := []map[string]interface{}{}
	for _, value := range order_arry {
		var info = make(map[string]interface{})
		info["price"] = value.Price
		info["time"] = time.Unix(int64(value.CreateTime), 0).Format("2006-01-02 15:04:05")
		info["bidders"] = value.Buyer
		info["status"] = "success"
		info["tx_hash"] = value.TxHash
		result = append(result, info)
	}
	count := models.OrderFindCount(map[string]interface{}{"market_id": market_id, "status": 1})
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result,
		Code:   200,
		Count:  count,
	}
	ctx.JSON(http.StatusOK, response)
	return
}

func OrderGetByTokenId(ctx *gin.Context) {
	tokenid := ctx.Query("tokenid")
	mid := ctx.Query("mid")
	ismid := ctx.Query("ismid")
	status := ctx.Query("status")
	chain_type := ctx.Query("chain_type")
	logger.Info("OrderGetByTokenId req:", "tokenid", tokenid, "mid", mid, "ismid", ismid, "status", status)
	res := model.OrdersByTokenId(tokenid, mid, ismid, status, chain_type)
	var recordList []response.OrderView
	for _, value := range res {
		var records response.OrderView

		records.Mid = value.Mid
		records.Tokenid = value.Tokenid
		records.Txhash = value.Txhash
		records.Buyer = value.Buyer
		records.TokenType = value.TokenType
		records.Seller = value.Seller
		records.Price = value.Price
		records.TokenType = value.TokenType
		records.Status = value.Status
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.OnchainTime = time.Unix(value.OnchainTime, 0).Format("2006-01-02 15:04:05")
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation

		recordList = append(recordList, records)
	}
	response.ReturnSuccessResponse(ctx, recordList)
}

// OrderGetMaxPrice
func OrderGetMaxPrice(ctx *gin.Context) {
	market_id := ctx.Query("market_id")
	buyer, max_price := models.OrdersGetMaxPrice(market_id)
	response.ReturnSuccessResponseNew(ctx, gin.H{
		"msg":       200,
		"buyer":     buyer,
		"max_price": max_price,
	})
}

// OrdersByBuyer buyer, mtype, status
func OrdersByBuyer(ctx *gin.Context) {
	buyer := ctx.Query("buyer")
	status := ctx.Query("status")
	mtype := ctx.Query("mtype")
	UserId := ctx.Query("user_id")
	chain_type := ctx.Query("chain_type")

	logger.Info("OrderGetByTokenId req:", "buyer", buyer, "mtype", mtype, "status", status)
	res := model.OrdersByBuyer(buyer, mtype, status, chain_type)
	var recordList []response.MyBids
	for _, value := range res {
		var records response.MyBids
		marketInfo := model.MarketByTokenIdMid(value.Tokenid, strconv.Itoa(value.Mid), "1", chain_type)
		if len(marketInfo) == 1 {
			nftInfo := model.NftGetByTokenId(value.Tokenid, chain_type)
			records.NftName = nftInfo[0].NftName
			records.NftDesc = nftInfo[0].NftDesc
			records.MediaUri = nftInfo[0].MediaUri
			records.CreateNumber = nftInfo[0].CreateNumber
			records.ChainType = nftInfo[0].ChainType
			var attention int
			if UserId != "" {
				id := model.MyAttentionFindNew(UserId, value.Id)
				if id != "" {
					attention = 1
				}
			}
			records.Attention = attention
			records.Mid = value.Mid
			records.Tokenid = value.Tokenid
			records.Txhash = value.Txhash
			records.Buyer = value.Buyer
			records.TokenType = value.TokenType
			records.Seller = value.Seller
			records.Price = marketInfo[0].StartingPrice // 起始拍卖价格
			records.TokenType = value.TokenType
			records.Status = value.Status
			records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
			records.OnchainTime = time.Unix(value.OnchainTime, 0).Format("2006-01-02 15:04:05")
			records.EndTime = time.Unix(marketInfo[0].EndTime, 0).Format("2006-01-02 15:04:05")
			records.OrderStatus = marketInfo[0].Status
			records.FreeGas = value.FreeGas
			records.Donation = value.Donation
			recordList = append(recordList, records)
		}
	}
	response.ReturnSuccessResponse(ctx, recordList)
}
