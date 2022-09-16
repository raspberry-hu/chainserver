package handler

import (
	"ChainServer/api/request"
	response2 "ChainServer/api/response"
	"ChainServer/common"
	"ChainServer/models"
	"ChainServer/models/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

const AdminAccount = "admin"

// MarketNew
// sn,creater,tokenid,market_type,starting_price,token_type,end_time,bonus
func MarketNew(ctx *gin.Context) {
	req := request.MarketNew{}
	ctx.ShouldBindJSON(&req)
	data, _ := json.Marshal(&req)
	logger.Info("market new ", string(data))

	var status int
	nft_id := req.NftId
	token_id := req.TokenId
	creater := req.Creater
	//token_id := req.TokenId
	market_type := req.MarketType
	starting_price := req.StartingPrice
	currency_name := req.CurrencyName
	end_time := req.EndTime
	reward := req.Reward
	chain_name := req.ChainName
	lazy := req.Lazy
	donation := req.Donation
	create_time := 0
	donation_userid := req.DonationUserId
	tx_hash := req.TxHash

	if lazy == 0 {
		status = 0
	} else if lazy == 1 {
		status = 1
		create_time = int(time.Now().Unix())
	}
	////不考虑伪铸币的情况
	//status = 1
	//create_time = int(time.Now().Unix())

	var market_info models.TMarketList
	market_info.Creater = creater
	market_info.TokenId = token_id
	market_info.MarketType = market_type
	market_info.StartingPrice = starting_price
	market_info.CurrencyName = currency_name
	market_info.Reward = reward
	market_info.ChainName = chain_name
	market_info.EndTime = int64(end_time)
	market_info.Lazy = lazy
	market_info.CreateTime = int64(create_time)
	market_info.Donation = donation
	market_info.NftId = nft_id
	market_info.Status = status
	market_info.DonationUserId = donation_userid
	market_info.TxHash = tx_hash
	if lazy == 1 {
		market_info.MarketType = market_type
	}
	market_id := models.MarketListInsert(market_info)
	var nft_info models.TNft
	if market_id > 0 && lazy == 1 {
		nft_info.Id = nft_id
		models.NftUpdate(nft_info, map[string]interface{}{"market_type": market_type})
	}
	////不考虑伪铸币的情况时
	//nft_info.Id = nft_id
	//models.NftUpdate(nft_info, map[string]interface{}{"market_type": market_type})
	var result = make(map[string]interface{})
	result["market_id"] = market_id
	response2.ReturnSuccessResponseNew(ctx, result)
}

/**
用来在市场中展示，相关的资产信息进行重复展示。
*/
func TokenMarket(ctx *gin.Context) {
	//进行分页展示
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	//account := ctx.Query("admin")
	offset := (page - 1) * rows

	var recordlist []models.TNft
	var total []models.TNft

	/**
	根据资产的owner进行查询,即市场中仅展示未被购买的资产
	*/
	sql := map[string]interface{}{"ant_nft_owner": AdminAccount}
	recordlist = model.MarketInfoFindLimit(sql, rows, offset)
	total = models.NftFind(map[string]interface{}{"ant_nft_owner": AdminAccount})
	//recordlist = models.NftFind(sql)

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

// MarketHashUpdate
func MarketHashUpdate(ctx *gin.Context) {
	req := request.UpdateHash{}
	ctx.ShouldBindJSON(&req)
	market_id := req.MarketId
	tx_hash := req.TxHash
	// wallet_addr := strings.ToLower(ctx.Request.Header.Get("wallet"))
	user_id := strings.ToLower(ctx.Request.Header.Get("user_id"))
	market_arry := models.MarketListFind(map[string]interface{}{"id": market_id, "creater": user_id})
	if len(market_arry) == 0 {
		logger.Error("market 找不到数据")
	} else {
		if market_arry[0].TxHash != "" {
			logger.Error("已经更新")
		} else {
			models.MarketListUpdate(market_arry[0], map[string]interface{}{"tx_hash": tx_hash})
			logger.Info("success")
		}
	}
	response2.ReturnSuccessResponseNew(ctx, nil)
}

// MarketCancelHashUpdate
func MarketCancelHashUpdate(ctx *gin.Context) {
	req := request.UpdateCancelHash{}
	ctx.ShouldBindJSON(&req)
	market_id := req.MarketId
	nft_id := req.NftId
	lazy := req.Lazy
	cancel_hash := req.CancelHash
	user_id := strings.ToLower(ctx.Request.Header.Get("user_id"))

	market_arry := models.MarketListFind(map[string]interface{}{"id": market_id, "creater": user_id})
	if len(market_arry) == 0 {
		logger.Error("market 找不到数据")
	} else {
		if lazy == 1 {
			models.MarketListUpdate(market_arry[0], map[string]interface{}{"status": 3})
			var nft_info models.TNft
			nft_info.Id = nft_id
			models.NftUpdate(nft_info, map[string]interface{}{"market_type": 0})
		} else {
			if market_arry[0].CancelHash != "" {
				logger.Error("已经更新")
			} else {
				models.MarketListUpdate(market_arry[0], map[string]interface{}{"cancel_hash": cancel_hash})
				logger.Info("success")
			}
		}
	}
	response2.ReturnSuccessResponseNew(ctx, nil)
}

// MarketDealHashUpdate
func MarketDealHashUpdate(ctx *gin.Context) {
	req := request.UpdateCancelHash{}
	ctx.ShouldBindJSON(&req)
	token_id := req.TokenId
	DealHash := req.DealHash
	nft_id := req.NftId
	lazy := req.Lazy
	market_id := req.MarketId
	//wallet_addr代表seller
	// wallet_addr := strings.ToLower(ctx.Request.Header.Get("wallet"))
	user_id, _ := strconv.Atoi(ctx.Request.Header.Get("user_id"))
	market_arry := models.MarketListFind(map[string]interface{}{"id": market_id, "creater": user_id})
	if len(market_arry) == 0 {
		logger.Error("找不到market")
	}

	//拍卖交易结束，将collection对应的交易额进行累加
	var collection_info models.TCollectionInfo
	var market_info models.TMarketList
	market_info.Id = market_id
	var buyer string
	var maxPrice float32
	nftInfo := models.NftFind(map[string]interface{}{"id": nft_id})
	if len(nftInfo) == 0 {
		logger.Error("没找到nft")
	} else {
		// 错误
		//collection_info.Id = nftInfo[0].CollectionId
	}

	//不收取账户gas费用
	if lazy == 1 && token_id == "" {
		models.MarketListUpdate(market_arry[0], map[string]interface{}{"status": 3})
		var nft_info models.TNft
		nft_info.Id = nft_id
		models.NftUpdate(nft_info, map[string]interface{}{"market_status": 0})
	} else if lazy == 0 && token_id != "" {
		models.MarketListUpdate(market_arry[0], map[string]interface{}{"deal_hash": DealHash})
		var amountCount float32
		var itemCount int

		//查找对应的collection中当前的
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

		//查询当前所有竞拍者中出价最高者
		buyer, maxPrice = models.OrdersGetMaxPrice(strconv.Itoa(market_id))
		//再次更新market中最终的buyer
		models.MarketListUpdate(market_info, map[string]interface{}{"buyer": buyer})

		amountCount = amountCount + maxPrice
		//models.UpdateAmountByCollection(collection_info, map[string]interface{}{"amount": amountCount})
		//models.UpdateItemsByCollection(collection_info, map[string]interface{}{"items": itemCount})
		models.CollectionUpdate(collection_info, map[string]interface{}{"amount": amountCount, "items": itemCount})
	}

	response2.ReturnSuccessResponseNew(ctx, gin.H{
		"msg":       200,
		"buyer":     buyer,
		"max_price": maxPrice,
	})
}

// MarketGetByMid
func MarketGetByMid(ctx *gin.Context) {
	mid := ctx.Query("mid")
	chain_type := ctx.Query("chain_type")
	logger.Info("MarketGetByMid req:", "mid", mid)
	res := model.MarketGetByMid(mid, chain_type)
	var recordList []response2.MarketView
	for _, value := range res {
		var records response2.MarketView

		records.Creater = value.Creater
		records.Tokenid = value.Tokenid
		records.MarketType = value.MarketType
		records.StartingPrice = value.StartingPrice
		records.TokenType = value.TokenType
		records.Buyer = value.Buyer
		records.Bonus = value.Bonus
		records.Txhash = value.Txhash
		records.Sorting = value.Sorting
		records.Status = value.Status
		records.CreateTime = time.Unix(value.CreateTime, 0).Format("2006-01-02 15:04:05")
		records.EndTime = time.Unix(value.EndTime, 0).Format("2006-01-02 15:04:05")
		records.FreeGas = value.FreeGas
		records.Donation = value.Donation
		records.NftId = value.NftId
		recordList = append(recordList, records)
	}
	response2.ReturnSuccessResponse(ctx, recordList)
}

// MarketAllHandle page,rows,type,status
func MarketAllHandle(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	donation := ctx.Query("donation")
	user_id, _ := strconv.Atoi(ctx.Query("user_id"))
	market_type := ctx.Query("market_type")
	state := ctx.Query("state")
	currency_name := ctx.Query("currency_name")
	market_creater := ctx.Query("market_creater")
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
	//order_sql := ""

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

// 列表详情页
func GDNfttDetailsHandle(ctx *gin.Context) {
	n_id := ctx.Query("nft_id")
	nft_id, err := strconv.Atoi(n_id)
	user_id, _ := strconv.Atoi(ctx.Query("user_id"))

	result := []map[string]interface{}{}
	if err != nil {
		response := &common.ResponseCommon{
			Msg:    "not found",
			Code:   404,
			Result: result,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	limit_sql := ""
	params_sql := fmt.Sprintf("where n.id = %d", nft_id)
	order_sql := ""
	data_arry := models.MarketALeftNftFind(user_id, params_sql, order_sql, limit_sql)
	msg := "success"
	code := 200
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
		organization_info := OrgAdd(value.NftOwner) //不懂这里啥意思
		info["organization_info"] = organization_info
		info["collection_info"] = collection_info
		result = append(result, info)
	}
	if len(result) == 0 {
		msg = "not found"
		code = 404
	}

	response := &common.ResponseCommon{
		Msg:    msg,
		Code:   code,
		Result: result,
	}
	ctx.JSON(http.StatusOK, response)
}

// 列表详情页
func GDMarketDetailsHandle(ctx *gin.Context) {
	m_id := ctx.Query("market_id")
	user_id, _ := strconv.Atoi(ctx.Query("user_id"))
	market_id, err := strconv.Atoi(m_id)
	result := []map[string]interface{}{}
	if err != nil {
		response := &common.ResponseCommon{
			Msg:    "not found",
			Code:   404,
			Result: result,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}
	limit_sql := ""
	params_sql := fmt.Sprintf("where m.id = %d and m.status = 1 ", market_id)
	order_sql := ""
	data_arry := models.MarketAllFind(user_id, params_sql, order_sql, limit_sql)
	msg := "success"
	code := 200
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
		organization_info := OrgAdd(value.NftOwner) //不懂啥意思
		info["organization_info"] = organization_info
		info["collection_info"] = collection_info
		result = append(result, info)
	}
	if len(result) == 0 {
		msg = "not found"
		code = 404
	}

	response := &common.ResponseCommon{
		Msg:    msg,
		Code:   code,
		Result: result,
	}
	ctx.JSON(http.StatusOK, response)
}

// 列表页
func GDMarketListHandle(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	offset := (page - 1) * rows
	limit_sql := fmt.Sprintf("LIMIT %d OFFSET %d", rows, offset)
	var params_map = make(map[string]string)
	params_map["status"] = "1"
	params_sql := ScreeningParams(params_map)
	order_sql := ""
	user_id, _ := strconv.Atoi("")
	data_arry := models.MarketAllFind(user_id, params_sql, order_sql, limit_sql)
	result := []map[string]interface{}{}
	for _, value := range data_arry {
		_, max_price := models.OrdersGetMaxPrice(fmt.Sprintf("%d", value.MarketId))
		var info = make(map[string]interface{})
		info["nft_name"] = value.NftName
		info["price"] = max_price
		info["media_uri"] = value.MediaUri
		// info["donation"] = value.Donation
		info["market_id"] = value.MarketId
		result = append(result, info)
	}
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result,
		Code:   200,
		Count:  1,
	}
	ctx.JSON(http.StatusOK, response)
}

func MarketGetHistory(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	rows, _ := strconv.Atoi(ctx.Query("rows"))
	market_id, _ := strconv.Atoi(ctx.Query("market_id"))
	nft_id, _ := strconv.Atoi(ctx.Query("nft_id"))
	offset := (page - 1) * rows
	var sql string
	if market_id != 0 {
		sql = fmt.Sprintf("market_id != %d and status = 1 and nft_id = %d", market_id, nft_id)
	} else {
		sql = fmt.Sprintf("status = 1 and nft_id = %d", nft_id)
	}
	order_arry := models.OrderFindLimit(sql, rows, offset)
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
	count := models.OrderFindCount(sql)
	response := &common.ResponseNew{
		Msg:    "success",
		Result: result,
		Code:   200,
		Count:  count,
	}
	ctx.JSON(http.StatusOK, response)
	return
}

func ScreeningParams(params_map map[string]string) string {
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
		} else if k == "currency_name" {
			params_sql += fmt.Sprintf("m.%s = '%s' and ", k, v)
		} else if k == "creater" {
			params_sql += fmt.Sprintf("m.%s = '%s' and ", k, v)
		} else {
			params_sql += fmt.Sprintf("m.%s = %s and ", k, v)
		}
	}
	if params_sql == "where " {
		params_sql = ""
	} else {
		params_sql = strings.TrimRight(params_sql, "and ")
	}
	return params_sql
}

func MarketAllOrderByParams(order_params_map map[string]string) string {
	var order_by_sql = ""
	for k, v := range order_params_map {
		if v != "" {
			if k == "nft_create_time" {
				order_by_sql = fmt.Sprintf("order by nft_create_time %s", v)
			} else if k == "favorites_count" {
				order_by_sql = fmt.Sprintf("order by %s %s", k, v)
			} else {
				order_by_sql = fmt.Sprintf("order by m.%s %s nulls last", k, v)
			}
		}
	}
	return order_by_sql
}

func WallectInfoAdd(value models.MarketAllNft, info map[string]interface{}) {
	creater_user_info_arry := models.UserInfoFind(map[string]interface{}{"id": value.NftCreater})
	creater_user_info := make(map[string]interface{})
	owner_user_info := make(map[string]interface{})
	if len(creater_user_info_arry) > 0 {
		t_user_info := creater_user_info_arry[0]
		creater_user_info["image_url"] = t_user_info.ImageUrl
		creater_user_info["banner_url"] = t_user_info.BannerUrl
		creater_user_info["user_name"] = t_user_info.UserName
		creater_user_info["user_desc"] = t_user_info.UserDesc
		creater_user_info["user_id"] = t_user_info.Id
	} else {
		creater_user_info["user_name"] = value.NftCreater
		creater_user_info["user_id"] = value.NftCreater
		creater_user_info["image_url"] = ""
		creater_user_info["banner_url"] = ""
		creater_user_info["user_desc"] = ""
	}
	info["creater_user_info"] = creater_user_info

	if value.NftCreater == value.NftOwner {
		//owner_wallet_info =
		info["owner_user_info"] = creater_user_info
	} else {
		owner_user_info_arry := models.UserInfoFind(map[string]interface{}{"id": value.NftOwner})
		if len(owner_user_info_arry) > 0 {
			t_user_info := owner_user_info_arry[0]
			owner_user_info["image_url"] = t_user_info.ImageUrl
			owner_user_info["banner_url"] = t_user_info.BannerUrl
			owner_user_info["user_name"] = t_user_info.UserName
			owner_user_info["user_desc"] = t_user_info.UserDesc
			owner_user_info["user_id"] = t_user_info.Id
		} else {
			owner_user_info["user_name"] = value.NftOwner
			owner_user_info["user_id"] = value.NftOwner
			owner_user_info["image_url"] = ""
			owner_user_info["banner_url"] = ""
			owner_user_info["user_desc"] = ""
		}
		info["owner_user_info"] = owner_user_info
	}
}

func CollectionInfoAdd(collection_id int) map[string]interface{} {
	var t_collection_info models.TCollectionInfo
	t_collection_info.Id = collection_id
	var collection_info = make(map[string]interface{})
	if collection_id == 0 {
		return collection_info
	}
	collection_info_arry := models.CollectionFind(t_collection_info)
	if len(collection_info_arry) > 0 {
		collection_info["user_id"] = collection_info_arry[0].UserId
		collection_info["collection_name"] = collection_info_arry[0].CollectionName
		collection_info["chain_name"] = collection_info_arry[0].ChainName
		collection_info["currency_name"] = collection_info_arry[0].CurrencyName
		collection_info["logo_image"] = collection_info_arry[0].LogoImage
		collection_info["featured_image_url"] = collection_info_arry[0].FeaturedImageUrl
		collection_info["banner_imageUrl"] = collection_info_arry[0].BannerImageUrl
		collection_info["collection_desc"] = collection_info_arry[0].CollectionDesc
		collection_info["category"] = collection_info_arry[0].Category
	}
	return collection_info
}

func OrgAdd(user_id int) (organization_info map[string]interface{}) {
	organization_info = make(map[string]interface{})
	organization_info_arry := models.OrganizationInfoFind(1, 0, map[string]interface{}{"user_id": user_id})
	if len(organization_info_arry) > 0 {
		organization_info["image_url"] = organization_info_arry[0].ImageUrl
		organization_info["user_id"] = organization_info_arry[0].ID
		organization_info["org_name"] = organization_info_arry[0].OrgName
	}
	return
}
