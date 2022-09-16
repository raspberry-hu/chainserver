package handler

import (
	"ChainServer/api/request"
	response2 "ChainServer/api/response"
	"ChainServer/common"
	"ChainServer/config"
	"ChainServer/models"
	"ChainServer/models/model"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.alipay-inc.com/antchain/restclient-go-sdk/client"
	"log"

	//"github.com/go-kit/kit/transport/http"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

const (
	RestBizTestBizID    = "a00e36c5"
	RestBizTestAccount  = "Zxi7788"
	RestBizTestKmsID    = "vfx371d2LLENUIUX1658393912735"
	RestBizTestTenantID = "LLENUIUX"
	RestContractName    = "TokenWithMarket"
)

var restClient *client.RestClient

func AntInit() {
	var err error
	configFilePath := "./conf/rest-config.json"
	restClient, err = client.NewRestClient(configFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to NewRestClient err:%+v", err))
	}
	if restClient.RestToken == "" {
		panic(fmt.Errorf("rest token:%+v is empty", restClient.RestToken))
	}
}

func NftNew(ctx *gin.Context) {
	req := request.NftNew{}
	logger.Info(req)
	ctx.ShouldBindJSON(&req)
	data, _ := json.Marshal(&req)
	logger.Info("nft new", string(data))
	//response.ReturnErrorResponseNew(ctx, 405, "蚂蚁链余额不足")
	//response.ReturnErrorResponse(ctx, 201, "余额不足", "余额不足")
	//response.ReturnErrorResponse(ctx, 201, "蚂蚁链余额不足")
	//response.ReturnErrorResponse(ctx, 201, "蚂蚁链余额不足", "蚂蚁链余额不足")
	//return
	logger.Info("nft ant count", req.AntCount)
	if req.AntCount == 0 {
		response2.ReturnErrorResponseNew(ctx, 202, "fail")
		return
	}
	//var status int
	nft_name := req.NftName
	nft_desc := req.NftDesc
	//rights_rules := req.RightsRules
	media_uri := req.MediaUri
	//explore_uri := req.ExploreUri //商品价格
	//create_tax := req.CreateTax
	//creater := req.Creater
	//chain_name := req.ChainName
	//media_ipfs_uri := req.MediaIpfsUri
	//meta_data_uri := req.MetaDataUri
	////要有一个collection选择为空的判断
	collection_id := req.CollectionId
	//categories_id := req.CategoriesId
	//tx_hash := req.TxHash
	//
	//creater = creater
	owner := req.Owner

	antTokenId := req.AntTokenId
	antCount := req.AntCount
	antNftOwner := req.AntNftOwner
	antTokenUrl := req.AntTokenUrl
	antTxHash := req.AntTxHash
	price := req.AntPrice
	//antPrice := req.AntPrice
	//
	//currency_name := req.CurrencyName

	//lazy := int64(req.Lazy)
	//var create_time int64 = 0

	//if lazy == 0 {
	//	status = 0
	//} else if lazy == 1 {
	//	status = 1
	//	create_time = time.Now().Unix()
	//}

	var nft_info models.TNft
	nft_info.NftName = nft_name
	nft_info.NftDesc = nft_desc
	//nft_info.RightsRules = rights_rules
	nft_info.MediaUri = media_uri
	//nft_info.CreateTax = create_tax
	//nft_info.Creater = creater
	nft_info.Owner = owner
	//nft_info.CreateTime = create_time
	//nft_info.ChainName = chain_name
	//nft_info.CurrencyName = currency_name
	//nft_info.Lazy = lazy
	//nft_info.Status = status
	//nft_info.MetaDataUri = meta_data_uri
	//nft_info.MediaIpfsUri = media_ipfs_uri
	nft_info.CollectionId = collection_id
	//nft_info.CategoriesId = categories_id
	//nft_info.ExploreUri = explore_uri
	//nft_info.TxHash = tx_hash
	nft_info.AntNftOwner = req.AntNftOwner
	nft_info.AntTxHash = req.AntTxHash
	nft_info.AntTokenUrl = req.AntTokenUrl
	nft_info.AntCount = req.AntCount
	nft_info.AntTokenId = req.AntTokenId
	nft_info.AntPrice = req.AntPrice
	//nft_info.AntCollection = req.AntCollection
	//nft_info.AntNftName = req.NftName

	var nftTemp models.TNft
	var antStartId = models.NftInsert(nftTemp) + 1

	//蚂蚁链铸币
	//
	//antTokenId := req.AntTokenId
	//antCount := req.AntCount
	//antNftOwner := req.AntNftOwner
	//antTokenUrl := req.AntTokenUrl
	//antTxHash := req.AntTxHash

	jsonArr := make([]interface{}, 0)
	jsonArr = append(jsonArr, antTokenId)
	jsonArr = append(jsonArr, antCount)
	jsonArr = append(jsonArr, antNftOwner)
	jsonArr = append(jsonArr, antTokenUrl)
	jsonArr = append(jsonArr, antStartId)
	jsonArr = append(jsonArr, antTxHash)
	jsonArr = append(jsonArr, price)

	//将输入参数转成byte
	inputparams, err := json.Marshal(jsonArr)
	if err != nil {
		panic(err)
	}

	u := uuid.New()
	orderId := fmt.Sprintf("callNewNft_%v", u.String())
	var gas int64 = int64(500000 * antCount)

	//调用合约
	baseResp, err := restClient.CallContract(RestBizTestBizID, orderId, RestBizTestAccount, RestBizTestTenantID, RestContractName, "NewNft(uint256,uint256,string,string,uint256,string,uint256)",
		string(inputparams), `["string"]`, RestBizTestKmsID, false, gas)
	if !(err == nil && baseResp.Success) || baseResp.Code == "0" {
		panic(fmt.Errorf("no succ resp baseResp:%+v err:%+v", baseResp, err))
		response2.ReturnErrorResponseNew(ctx, 405, "蚂蚁链余额不足")
		return
	}

	outputs := config.Output{}
	err = json.Unmarshal([]byte(baseResp.Data), &outputs)
	if err != nil {
		panic(err)
	}

	outputs1 := outputs.OutRes[0].(string)
	fmt.Printf("The owner of %d NFT is %+v\n", antNftOwner, outputs1)
	fmt.Printf("The owner of %d NFT is %+v, , the transaction hash is %s\n", antNftOwner, outputs1, outputs.Transaction.TxHash)

	//蚂蚁链铸币

	//增加nft资产
	logger.Info(antCount)
	for i := 0; i < antCount; i++ {
		//var collection_info models.TCollectionInfo
		//var nftCount int
		//collectionArr := models.CollectionFind(map[string]interface{}{"id": collection_id})
		//nftCount = collectionArr[0].Items
		//nftCount = nftCount + 1
		//models.CollectionUpdate(collection_info, map[string]interface{}{"items": nftCount})
		nft_info.AntNftId = i
		nft_info.CreateTime = int64(int(time.Now().Unix()))
		// 插入铸造交易时的hash
		nft_info.TransferHash = outputs.Transaction.TxHash
		nft_info.MarketType = 2
		nft_info.Creater = 1
		antStartId = models.NftInsert(nft_info)
	}
	collectionArray := models.CollectionFind(map[string]interface{}{"id": collection_id})
	if len(collectionArray) == 0 {
		response2.ReturnErrorResponseNew(ctx, 202, "collection fail")
		return
	} else {
		for i := 0; i < len(collectionArray); i++ {
			models.CollectionUpdate(collectionArray[i], map[string]interface{}{"items": collectionArray[i].Items + antCount})
		}
	}
	//同时进行挂单操作

	//
	//var collection_info models.TCollectionInfo
	//var nftCount int
	//collection_info.Id = collection_id
	//collectionArr := models.CollectionFind(map[string]interface{}{"id": collection_id})
	//if len(collectionArr) == 0 {
	//	logger.Error("没找到Collection")
	//} else {
	//	nftCount = collectionArr[0].Items
	//}
	//nftCount = nftCount + 1
	//models.UpdateItemsByCollection(collection_info, map[string]interface{}{"items": nftCount})
	//models.CollectionUpdate(collection_info, map[string]interface{}{"items": nftCount})
	var result = make(map[string]interface{})
	result["nft_id"] = antStartId
	if antStartId == 0 {
		response2.ReturnErrorResponseNew(ctx, 203, "fail")
		return
	} else {
		response2.ReturnSuccessResponseNew(ctx, result)
		return
	}
}

//获取资产的详情信息
func GetTokenInfo(ctx *gin.Context) {
	//根据tokenID\currentID锁定
	tokenId, _ := strconv.Atoi(ctx.Query("ant_token_id"))
	currentId, _ := strconv.Atoi(ctx.Query("ant_nft_id"))

	var recordlist []models.TNft
	var resp response2.TokenView
	recordlist = models.NftFind(map[string]interface{}{"ant_token_id": tokenId, "ant_nft_id": currentId})

	var response common.ResponseNew
	if len(recordlist) == 1 {
		resp.NftName = recordlist[0].NftName
		resp.NftDesc = recordlist[0].NftDesc
		resp.Creater = recordlist[0].Creater
		resp.CreateTime = recordlist[0].CreateTime
		resp.MediaUri = recordlist[0].MediaUri
		resp.Owner = recordlist[0].Owner
		resp.CollectionId = recordlist[0].CollectionId
		resp.ExploreUri = recordlist[0].ExploreUri
		resp.AntTokenId = recordlist[0].AntTokenId
		resp.AntNftId = recordlist[0].AntNftId
		resp.AntCount = recordlist[0].AntCount
		resp.AntPrice = recordlist[0].AntPrice
		resp.AntTokenUrl = recordlist[0].AntTokenUrl
		resp.AntTxHash = recordlist[0].AntTxHash
		resp.AntNftOwner = recordlist[0].AntNftOwner
		resp.MarketType = recordlist[0].MarketType
		resp.TransferHash = recordlist[0].TransferHash

		//根据ID查询相应的用户名
		var originator, current []models.UserInfo
		originator = models.UserInfoFind(map[string]interface{}{"id": recordlist[0].Creater})
		current = models.UserInfoFind(map[string]interface{}{"id": recordlist[0].Owner})
		var collectionInfo []models.TCollectionInfo
		collectionInfo = models.CollectionFind(map[string]interface{}{"id": recordlist[0].CollectionId})

		if len(originator) == 1 && len(current) == 1 && len(collectionInfo) == 1 {
			resp.ImageUrl = originator[0].ImageUrl
			resp.CreaterName = originator[0].UserName
			resp.OwnerImageUrl = current[0].ImageUrl
			resp.OwnerName = current[0].UserName
			resp.CollectionName = collectionInfo[0].CollectionName
		} else {
			log.Println("Not find account by id")
		}

		response = common.ResponseNew{
			Code:   200,
			Msg:    "success",
			Result: resp,
			Count:  1,
		}
	} else {
		response = common.ResponseNew{
			Msg:    "null",
			Result: recordlist,
			Count:  0,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// NftHashUpdate
// sn,hash
func NftHashUpdate(ctx *gin.Context) {
	req := request.UpdateHash{}
	ctx.ShouldBindJSON(&req)
	tx_hash := req.TxHash
	nft_id := req.NftId
	user_id := strings.ToLower(ctx.Request.Header.Get("user_id"))
	nft_arry := models.NftFind(map[string]interface{}{"id": nft_id, "owner": user_id})
	if len(nft_arry) == 0 {
		logger.Error("找不到nft")
	} else {
		if nft_arry[0].TxHash != "" {
			logger.Error("已经更新")
		} else {
			models.NftUpdate(nft_arry[0], map[string]interface{}{"tx_hash": tx_hash})
		}
	}
	response2.ReturnSuccessResponseNew(ctx, nil)
}

func NftTransfersUpdate(ctx *gin.Context) {
	req := request.UpdateHash{}
	ctx.ShouldBindJSON(&req)
	transfer_hash := req.TransferHash
	nft_id := req.NftId
	user_id := ctx.Request.Header.Get("user_id")
	nft_arry := models.NftFind(map[string]interface{}{"id": nft_id, "market_type": 0, "owner": user_id})
	if len(nft_arry) == 0 {
		logger.Error("找不到nft")
	} else {
		models.NftUpdate(nft_arry[0], map[string]interface{}{"transfer_hash": transfer_hash, "market_type": 3})
	}
	response2.ReturnSuccessResponseNew(ctx, nil)
}

func NftOwnerUpdate(ctx *gin.Context) {
	req := request.NftOwner{}
	var nft_info models.TNft
	ctx.ShouldBindJSON(&req)
	nft_id := req.NftId
	owner := req.Owner
	buyer := req.Buyer
	nft_info.Id = nft_id
	nft_info.Owner = owner
	if nft_id == 0 || owner != 0 || buyer != 0 {
		response2.ReturnErrorResponseNew(ctx, 201, "fail")
	} else {
		if buyer != owner {
			models.NftUpdate(nft_info, map[string]interface{}{"owner": buyer})
		}
		response2.ReturnSuccessResponseNew(ctx, nil)
	}
}

// GetNFRFromCollection 展示集合中的资产列表
func GetNFRFromCollection(ctx *gin.Context) {
	//进行分页展示
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows

	//获取当前的colelctionId
	collectionId, _ := strconv.Atoi(ctx.Query("collection_id"))

	/**
	根据collectionId查询当前未被购买的资产
	*/
	var total []models.TNft
	var recordlist []models.TNft
	sql := map[string]interface{}{"ant_nft_owner": AdminAccount, "collection_id": collectionId}
	recordlist = models.NftFindByLimit(sql, rows, offset)
	total = models.NftFind(map[string]interface{}{"ant_nft_owner": AdminAccount, "collection_id": collectionId})

	var response common.ResponseNew
	var resp []response2.TokenView

	if len(recordlist) > 0 {
		for _, value := range recordlist {
			var everresp response2.TokenView

			everresp.NftName = value.NftName
			everresp.NftDesc = value.NftDesc
			everresp.Creater = value.Creater
			everresp.MediaUri = value.MediaUri
			everresp.Owner = value.Owner
			everresp.CollectionId = value.CollectionId
			everresp.ExploreUri = value.ExploreUri
			everresp.AntTokenId = value.AntTokenId
			everresp.AntNftId = value.AntNftId
			everresp.AntCount = value.AntCount
			everresp.AntPrice = value.AntPrice
			everresp.AntTokenUrl = value.AntTokenUrl
			everresp.AntTxHash = value.AntTxHash
			everresp.AntNftOwner = value.AntNftOwner
			everresp.MarketType = value.MarketType

			//根据ID查询相应的用户名
			var originator, current []models.UserInfo
			originator = models.UserInfoFind(map[string]interface{}{"id": value.Creater})
			current = models.UserInfoFind(map[string]interface{}{"id": value.Owner})

			if len(originator) == 1 && len(current) == 1 {
				everresp.ImageUrl = originator[0].ImageUrl
				everresp.CreaterName = originator[0].UserName
				everresp.OwnerImageUrl = current[0].ImageUrl
				everresp.OwnerName = current[0].UserName
			} else {
				log.Println("Not find account by id")
			}

			resp = append(resp, everresp)
		}

		response = common.ResponseNew{
			Msg:    "success",
			Result: resp,
			Code:   200,
			Count:  len(total),
		}
	} else {
		response = common.ResponseNew{
			Msg:    "null",
			Result: resp,
			Count:  len(total),
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// GetNFTFromCollection 展示市场collection中的资产
func GetNFTFromCollection(ctx *gin.Context) {
	collectionId, _ := strconv.Atoi(ctx.Query("collection_id"))
	status, _ := strconv.Atoi(ctx.Query("status"))
	currencyName := ctx.Query("currency")
	//minPrice := ctx.Query("min_price")
	//maxPrice := ctx.Query("max_price")
	name := ctx.Query("name")

	var nftInfo []models.TNft
	//var martketInfo []models.TMarketList
	var recordList []response2.NftAndMarketView
	var collectionInfo []models.TCollectionInfo
	collectionInfo = models.CollectionFind(map[string]interface{}{"id": collectionId})

	if len(collectionInfo) == 0 {
		logger.Error("Not find the collection")
	}

	// 默认查找该collection下所有挂单的资产
	if collectionId != 0 {

		// 先查找该集合下创建并满足条件的nftstatus == 0 &&
		if status == 0 && currencyName == "" && name == "" {
			nftInfo = models.NftFind(map[string]interface{}{"collection_id": collectionId, "owner": collectionInfo[0].UserId})
		} else {
			nftInfo = models.CollectionNftFind(collectionId, status, currencyName, name, collectionInfo[0].UserId)
		}
		for _, temp := range nftInfo {
			//if temp.Status == 1 && temp.MarketType == 0 && temp.Owner == collectionInfo[0].UserId {
			if temp.Owner == collectionInfo[0].UserId {
				//查找初次挂单未成交的资产
				//if minPrice == "" && maxPrice == "" {
				//martketInfo = models.MarketListFind(map[string]interface{}{"nft_id": temp.Id, "market_type": temp.MarketType, "creater": temp.Owner})
				//} else {
				//	martketInfo = models.MarketCollectionListFind(temp.Id, temp.Owner, minPrice, maxPrice)
				//}

				var record response2.NftAndMarketView
				//if len(martketInfo) != 0 {
				//record.MId = martketInfo[0].Id
				//record.MId = nftInfo[0].Owner
				//record.MarketCreateTime = time.Unix(martketInfo[0].CreateTime, 0).Format("2006-01-02 15:04:05")
				record.MarketCreateTime = time.Unix(nftInfo[0].CreateTime, 0).Format("2006-01-02 15:04:05")
				//record.StartingPrice = martketInfo[0].StartingPrice
				record.StartingPrice = strconv.Itoa(temp.AntPrice)
				//record.TokenType = martketInfo[0].TokenId
				//record.EndTime = time.Unix(martketInfo[0].EndTime, 0).Format("2006-01-02 15:04:05")
				//record.Bonus = martketInfo[0].Reward
				//record.ChainType = martketInfo[0].ChainName
				record.ChainType = "ant chain"
				//record.OrderStatus = martketInfo[0].Status
				//record.FreeGas = strconv.Itoa(martketInfo[0].Lazy)
				//record.Donation = martketInfo[0].Donation
				// 获取创建者和拥有着的名字和地址
				//var creatorUserInfo, ownerUserInfo []models.UserInfo
				var creatorUserInfo []models.UserInfo
				creatorUserInfo = models.UserInfoFind(map[string]interface{}{"id": temp.AntNftOwner})
				//ownerUserInfo = models.UserInfoFind(map[string]interface{}{"id": temp.Owner})
				if len(creatorUserInfo) != 0 {
					record.UserName = creatorUserInfo[0].UserName
					record.ImageUrl = creatorUserInfo[0].ImageUrl
				}
				//if len(ownerUserInfo) != 0 {
				//	record.OwnerUserName = ownerUserInfo[0].UserName
				//	record.OwnerImageUrl = ownerUserInfo[0].ImageUrl
				//}

				record.NftId = temp.AntNftId
				record.NftName = temp.NftName
				record.NftDesc = temp.NftDesc
				record.TokenId = strconv.Itoa(temp.AntTokenId)
				record.MediaUri = temp.MediaUri
				record.MediaIpfsUri = temp.MediaIpfsUri
				record.ExploreUri = temp.ExploreUri
				record.Creater = temp.Creater
				record.Owner = temp.Owner
				record.CreateTax = temp.CreateTax
				record.CurrencyName = temp.CurrencyName
				// 将时间转换成字符串的形式
				record.CreateTime = time.Unix(temp.CreateTime, 0).Format("2006-01-02 15:04:05")
				record.CollectionId = temp.CollectionId
				record.Status = temp.Status
				record.MarKetStatus = temp.MarketType
				record.ChainType = temp.ChainName
				record.CreateNumber = temp.BlockNumber
				record.TokenUri = temp.AntTokenUrl
				record.Txhash = temp.AntTxHash

				// 查看该资产被收藏的个数
				var favorites []models.TFavorites
				favorites = models.FavoritesFind(map[string]interface{}{"nft_id": temp.Id})
				record.FavoritesCount = strconv.Itoa(len(favorites))
				recordList = append(recordList, record)
				//}
			}
		}
	}

	//根据条件查找

	response2.ReturnSuccessResponse(ctx, gin.H{
		"msg":    "success",
		"result": recordList,
	})
}

// NftGetByTokenId
func NftGetByTokenId(ctx *gin.Context) {
	tokenid := ctx.Query("tokenid")
	chain_type := ctx.Query("chain_type")
	logger.Info("NftGetByTokenId req:", "tokenid", tokenid)
	res := model.NftGetByTokenId(tokenid, chain_type)
	var recordList []response2.NftView
	for _, value := range res {
		var records response2.NftView
		records.NftId = value.Id
		records.NftName = value.NftName
		records.NftDesc = value.NftDesc
		records.RightsRules = value.RightsRules
		records.TokenId = value.TokenId
		records.TokenUri = value.TokenUri
		records.Txhash = value.Txhash
		records.Creater = value.Creater
		records.CreateNumber = value.CreateNumber
		records.MediaUri = value.MediaUri
		records.CreateTax = value.CreateTax
		records.Owner = value.Owner
		records.NftType = value.NftType
		records.Status = value.Status
		records.MarKetStatus = value.MarKetStatus
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		recordList = append(recordList, records)
	}
	response2.ReturnSuccessResponse(ctx, recordList)
}

// NftGetBySn
func NftGetBySn(ctx *gin.Context) {
	sn := ctx.Query("sn")
	chain_type := ctx.Query("chain_type")
	logger.Info("NftGetBySn req:", "sn", sn)
	res := model.NftGetBySn(sn, chain_type)
	var recordList []response2.NftView
	for _, value := range res {
		var records response2.NftView
		records.NftId = value.Id
		records.NftName = value.NftName
		records.NftDesc = value.NftDesc
		records.RightsRules = value.RightsRules
		records.TokenId = value.TokenId
		records.TokenUri = value.TokenUri
		records.Txhash = value.Txhash
		records.Creater = value.Creater
		records.CreateNumber = value.CreateNumber
		records.MediaUri = value.MediaUri
		records.CreateTax = value.CreateTax
		records.Owner = value.Owner
		records.NftType = value.NftType
		records.Status = value.Status
		records.MarKetStatus = value.MarKetStatus
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.ChainType = chain_type
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation
		recordList = append(recordList, records)
	}
	response2.ReturnSuccessResponse(ctx, recordList)
}

// NftMySalability
func NftMySalability(ctx *gin.Context) {
	address := ctx.Query("address")
	market_status := ctx.Query("market_status")
	user_id := ctx.Query("user_id")
	chain_type := ctx.Query("chain_type")

	logger.Info("NftMySalability req:", "address", address, "market_status", market_status)

	res := model.NftMySalability(address, market_status, chain_type)
	var recordList []response2.NftAndMarketView
	for _, value := range res {
		var records response2.NftAndMarketView
		if market_status != "0" {
			marketInfo := model.MarketByTokenId(value.Id, "1", value.ChainType)
			if len(marketInfo) == 1 {
				records.MId = marketInfo[0].Id
				records.StartingPrice = marketInfo[0].StartingPrice
				records.TokenType = marketInfo[0].TokenType
				records.MarketCreateTime = time.Unix(marketInfo[0].CreateTime, 0).Format("2006-01-02 15:04:05")
				records.EndTime = time.Unix(marketInfo[0].EndTime, 0).Format("2006-01-02 15:04:05")
				records.Bonus = marketInfo[0].Bonus
				records.OrderStatus = marketInfo[0].Status
			}
		}
		var attention = 0
		if user_id != "" {
			id := model.MyAttentionFindNew(user_id, value.Id)
			if id != "" {
				attention = 1
			}
		}
		records.Attention = attention

		records.SN = value.SN
		records.NftName = value.NftName
		records.NftDesc = value.NftDesc
		records.TokenId = value.TokenId
		records.TokenUri = value.TokenUri
		records.Txhash = value.Txhash
		records.Creater = value.Creater
		records.CreateNumber = value.CreateNumber
		records.MediaUri = value.MediaUri
		records.CreateTax = value.CreateTax
		records.Owner = value.Owner
		records.NftType = value.NftType
		records.Status = value.Status
		records.MarKetStatus = value.MarKetStatus
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.ChainType = value.ChainType
		records.TokenType = value.TokenType
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation
		records.NftId = value.Id
		// 创建者用户信息
		info_arry := models.UserInfoFind(map[string]interface{}{"user_id": value.Creater})
		if len(info_arry) > 0 {
			records.UserName = info_arry[0].UserName
			records.ImageUrl = info_arry[0].ImageUrl
		}
		// 拥有者信息
		owner_info_arry := models.UserInfoFind(map[string]interface{}{"user_id": value.Owner})
		if len(owner_info_arry) > 0 {
			records.OwnerUserName = owner_info_arry[0].UserName
			records.OwnerImageUrl = owner_info_arry[0].ImageUrl
		}

		recordList = append(recordList, records)
	}
	response2.ReturnSuccessResponse(ctx, recordList)
}

// NftType
func NftType(ctx *gin.Context) {
	logger.Info("NftType req:")
	chain_type := ctx.Query("chain_type")
	res := model.TypeList(chain_type)
	var recordList []response2.NftType
	for _, value := range res {
		var records response2.NftType

		records.Id = value.Id
		records.TypeName = value.TypeName

		recordList = append(recordList, records)
	}
	response2.ReturnSuccessResponse(ctx, recordList)
}
