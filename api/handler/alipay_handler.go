package handler

import (
	"ChainServer/api/response"
	"ChainServer/config"
	"ChainServer/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sea-project/go-logger"
	"github.com/smartwalle/alipay/v3"
	"github.com/smartwalle/xid"
	"log"
	"net/http"
	"strconv"
	"time"
)

var aliClient *alipay.Client
var alipayTemp *temp

type temp struct {
	userId    string
	price     string
	tokenId   string
	currentId string
	ownerName string
}

const (
	kAppId        = "2021000121627867"
	kPrivateKey   = "MIIEogIBAAKCAQEAi6Dr0ZNfc2qhDkoYrB23V2ve2Rruv2LBhNOuzNYi0T1Mo50Q2AHBR2XfGbYFpSLwDketOcnb7UtFILTDqw7PI6Yi6Vj+/ZYCSZ7luDvFhcut+hbkdEV/tFv7pjUjsizUjv51vy7LWvuUa6UWuWgT9XbXurqDGnsfJekM8WTOq2ApuluESvy3AcGveO8gv/wGCkSLgHArveTGIw85MJrJMj0Y9AX2koQT5miVC6IPP74Y5HYxwgTo6/ADSEBKNn2NoGoJ+L8oDyK/M0z3RiP7tuDdiqrIhVgcuKifyBd1g4p7B8P/Aq8a8SN06zclJ1AMSv5QVeZGfhtPqEj1Vo5CLwIDAQABAoIBAChTocg1uCutcDagT9/l9T3aedJaZPoBm8KxIJsofYXRHoFiq6q3Vws38OeMGrVHEe4N5Yn7Mvml86EulBSjgk/Ze9vJSFwVJzP8IHzFRpcN7IF+exzZtbhxmIy4bEbZi8qA06ET8sekQYmVdKq31Ivgdw4HMDZFuQlJ9eMCKm50NyxrzX52ET1ShKO/PGpcfoLiyMbZfylDblo35WuyKua54MIP7C6Hbf9goQebQjLb8ai/+vVaN6ARTd4bTHBQITQjkUvJsUZFhHhM6lOAqRtlr6R96VMHHda+tmzzrPBrebKTCkvxapD4Os2FmDvYyMgFVLrm+79XrXJdjmOvrgECgYEA8NX6sYERNLAsesCwrrGXrstUNUmDNVq3aFJofmVymcpCOJbujsKFfWOyPDNv6zZtt7hujZAUi3oapIqpWmT/c0yC3W1QQDSgx5wakFfN2BYtWjBmKDtoJIfsjnGY1LYnSO+IT36c87azZGAT9uGdP4DsUqzEaP3qMXeJwoiqep8CgYEAlGuWWzKb4NbVcyydOeELR+PRO+Zk+4stDAxK+ztdi+eXO05d+a+zwsfDG+zN0L+MC2odOytggCVS6M+kpMmIQMn+Ce32aUK4WrsZsnLXKKiCBlpNGGCnP2ogg98fR1/GmjV9ushLaTKcDQDaZJGkBbD+V63zd0Zd9Ljt7jBgnnECgYAO4vsrC7JXkmg9cjm5oqqgmFrtLE0a+C+MGEPzRCwQS4tKWjIGywlbVdVHmVparLOdfp3+zCAo+vQ4pYWQW9vacalJLJ+gSGCD1idiMrs7lytYftNhu0JVt70slOMAiv3kqUHAwC/NdMaj9rhlM074BO0Wsy003DUkt6HhT3dSOwKBgHFY4NDxC35gU07MKZ7EMRtL7sTyJPi9xz9GPBU1tzFbQnG2XaqL9pqweF7hMCVVw5wMBBrl+6Kh3nmR6kk25+mi2XG329FzdNtFvFA9x/dzCSnU2L/fQJr7b62GPpsBl+i5JTX6NS03y1la526ak0sNapCHdkIG6UY13O9k3sThAoGAHUJSMS0/DkGomnPyxmLvO5OietFJTEalBLl0ieZxNsiYyw5bl9rPwEmltw5oQL4Q6YKG8X/nsX13zRdH/vSvrVdHF7H3fSnAi8g2UargW6eTJPnhBxPLj1T6y1vhNljytougwexr/r2tH66THM9z6NFAQfqO7kJjadaw53vqfOg="
	kServerPort   = "9204"
	kServerDomain = "http://8.130.97.237" + ":" + kServerPort
)

func AliInit() {
	var err error

	if aliClient, err = alipay.New(kAppId, kPrivateKey, false); err != nil {
		log.Println("初始化支付宝失败", err)
		return
	}

	// 使用支付宝证书
	if err = aliClient.LoadAppPublicCertFromFile("./conf/appCertPublicKey_2021000121627867.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}

	if err = aliClient.LoadAliPayRootCertFromFile("./conf/alipayRootCert.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}
	if err = aliClient.LoadAliPayPublicCertFromFile("./conf/alipayCertPublicKey_RSA2.crt"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}
	log.Println("加载证书成功")
	//var s = gin.Default()
	//s.GET("/alipay", pay)
	//s.GET("/callback", callback)
	//s.POST("/notify", notify)
	//s.Run(":" + kServerPort)
}

func Pay(c *gin.Context) {
	var tradeNo = fmt.Sprintf("%d", xid.Next())
	price := c.Query("ant_price")
	tokenId := c.Query("ant_tokenId")
	currentId := c.Query("ant_nftId")
	userId := c.Query("ant_userId")
	ownerName := c.Query("ant_nft_owner")
	alipayTemp = &temp{userId, price, tokenId, currentId, ownerName}
	log.Println(alipayTemp.userId, alipayTemp.price, alipayTemp.tokenId, alipayTemp.currentId, alipayTemp.ownerName)
	var p = alipay.TradePagePay{}
	p.NotifyURL = kServerDomain + "/api/v1/notify"
	p.ReturnURL = kServerDomain + "/api/v1/callback"
	p.Subject = "NFT购买:" + tokenId + currentId + tradeNo
	p.OutTradeNo = tradeNo
	p.TotalAmount = price
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, _ := aliClient.TradePagePay(p)

	//c.Redirect(http.StatusTemporaryRedirect, url.String())
	response.ReturnSuccessResponseNew(c, url.String())
}

func PurchaseStatus(ctx *gin.Context) {
	tokenId := ctx.Query("ant_tokenId")
	currentId := ctx.Query("ant_nftId")
	currentIdNow, _ := strconv.Atoi(currentId)
	nftTemp := models.NftFind(map[string]interface{}{"ant_token_id": tokenId})
	if len(nftTemp) == 0 {
		response.ReturnErrorResponseNew(ctx, 201, "fail")
		return
	} else {
		for i := 0; i < len(nftTemp); i++ {
			if nftTemp[i].AntNftId == currentIdNow {
				if nftTemp[i].AntNFTBuyer != "" && "admin" != nftTemp[i].AntNftOwner {
					response.ReturnSuccessResponseNew(ctx, "success")
					return
				}
			}
		}
	}
	response.ReturnErrorResponseNew(ctx, 201, "fail")
	return
}

func Callback(c *gin.Context) {
	c.Request.ParseForm()

	ok, err := aliClient.VerifySign(c.Request.Form)
	if err != nil {
		log.Println("回调验证签名发生错误", err)
		return
	}

	if ok == false {
		log.Println("回调验证签名未通过")
		return
	}

	log.Println(alipayTemp.userId, alipayTemp.price, alipayTemp.tokenId, alipayTemp.currentId, alipayTemp.ownerName)

	var outTradeNo = c.Request.Form.Get("out_trade_no")
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := aliClient.TradeQuery(p)
	if err != nil {
		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())
		return
	}
	if rsp.IsSuccess() == false {
		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
		return
	}

	c.String(http.StatusOK, "订单 %s 支付成功", outTradeNo)
	log.Println("订单支付成功")

	tokenId, _ := strconv.Atoi(alipayTemp.tokenId)
	subId, _ := strconv.Atoi(alipayTemp.currentId)
	count, _ := strconv.Atoi(alipayTemp.price)
	buyer := alipayTemp.ownerName

	var recordlist []models.TNft
	recordlist = models.NftFind(map[string]interface{}{"ant_token_id": tokenId, "ant_nft_id": subId})
	if len(recordlist) == 0 {
		response.ReturnErrorResponseNew(c, 209, "未找到该数字资产")
		return
	}
	log.Println(recordlist)
	jsonArr := make([]interface{}, 0)
	jsonArr = append(jsonArr, tokenId)
	jsonArr = append(jsonArr, recordlist[0].Id)
	jsonArr = append(jsonArr, count)
	jsonArr = append(jsonArr, buyer)

	//将输入参数转成byte
	inputparams, err := json.Marshal(jsonArr)
	if err != nil {
		panic(err)
	}

	u := uuid.New()
	orderId := fmt.Sprintf("callBuy_%v", u.String())
	var gas int64 = 500000

	//调用合约
	baseResp, err := restClient.CallContract(RestBizTestBizID, orderId, RestBizTestAccount, RestBizTestTenantID, RestContractName, "buy(uint256,uint256,uint256,string)",
		string(inputparams), `["string"]`, RestBizTestKmsID, false, gas)
	if !(err == nil && baseResp.Success) || baseResp.Code == "0" {
		panic(fmt.Errorf("no succ resp baseResp:%+v err:%+v", baseResp, err))
	}

	outputs := config.Output{}
	err = json.Unmarshal([]byte(baseResp.Data), &outputs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The current owner of %d token is %+v", tokenId, outputs.OutRes[0].(string))
	fmt.Printf("The current owner of %d token is %+v, the transaction hash is %s\n", tokenId, outputs.OutRes[0].(string), outputs.Transaction.TxHash)

	nft_arry := models.NftFind(map[string]interface{}{"ant_token_id": tokenId})
	nftTemp := models.NftFind(map[string]interface{}{"ant_token_id": tokenId})
	antCount := 0
	if len(nftTemp) == 0 {
		logger.Error("找不到nft")
	} else {
		antCount = nftTemp[0].AntCount - 1
	}
	if len(nft_arry) == 0 {
		logger.Error("找不到nft")
	} else {
		//if nft_arry[0].TxHash != "" {
		//	logger.Error("已经更新")
		//} else {
		//	models.NftUpdate(nft_arry[0], map[string]interface{}{"tx_hash": tx_hash})
		//}
		for i := 0; i < len(nft_arry); i++ {
			models.NftUpdate(nft_arry[i], map[string]interface{}{"ant_count": antCount})
			if nft_arry[i].AntNftId == subId {
				models.NftUpdate(nft_arry[i], map[string]interface{}{"ant_count": antCount, "ant_nft_owner": alipayTemp.ownerName,
					"buy_time": int(time.Now().Unix()), "transfer_hash": outputs.Transaction.TxHash, "market_type": 3,
					"ant_nft_buyer": alipayTemp.ownerName, "owner": alipayTemp.userId})
				collectionArray := models.CollectionFind(map[string]interface{}{"id": nft_arry[i].CollectionId})
				if len(collectionArray) == 0 {
					response.ReturnErrorResponseNew(c, 202, "collection not find")
					return
				} else {
					for j := 0; j < len(collectionArray); j++ {
						models.CollectionUpdate(collectionArray[j], map[string]interface{}{"items": collectionArray[j].Items - 1, "amount": collectionArray[j].Amount + float32(nft_arry[i].AntPrice)})
					}
				}
			}
		}
	}

	logger.Info("nft buy", alipayTemp.currentId)
}

func Notify(c *gin.Context) {
	c.Request.ParseForm()

	ok, err := aliClient.VerifySign(c.Request.Form)
	if err != nil {
		log.Println("异步通知验证签名发生错误", err)
		return
	}

	if ok == false {
		log.Println("异步通知验证签名未通过")
		return
	}

	log.Println("异步通知验证签名通过")

	var outTradeNo = c.Request.Form.Get("out_trade_no")
	var p = alipay.TradeQuery{}
	p.OutTradeNo = outTradeNo
	rsp, err := aliClient.TradeQuery(p)
	if err != nil {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s \n", outTradeNo, err.Error())
		return
	}
	if rsp.IsSuccess() == false {
		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s \n", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
		return
	}

	log.Printf("订单 %s 支付成功 \n", outTradeNo)
}
