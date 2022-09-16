package api

import (
	"ChainServer/api/handler"
	"ChainServer/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

func RouterStart() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 允许使用跨域请求  全局中间件
	router.Use(Cors())
	// 传参 设定路由组 允许路由组使用路由
	//router.Router(engine)

	chain := router.Group("/api/v1")
	auth := router.Group("/api/v1")

	//jwt.JWTAuth()
	auth.Use()
	{
		/**
		 * NFT 路由
		 */
		// 新增NFT √
		auth.POST("/nft/new", handler.NftNew)
		// 判断购买状态
		auth.GET("/nft/purchasestatus", handler.PurchaseStatus)
		// 新增查看蚂蚁链中创建成功的资产详情
		auth.GET("token/detail", handler.GetTokenInfo)
		// 查询当前账户下已创建的collection
		chain.GET("/existed/collection", handler.ExistedCollection)
		// 根据nft_id 更新 txhash √
		auth.POST("/nft/hash/update", handler.NftHashUpdate)
		// 根据nft_id 更新 transfer_hash √
		auth.POST("/nft/transfers/update", handler.NftTransfersUpdate)
		// 通过支付宝购买nft
		auth.GET("/alipay", handler.Pay)
		// 支付宝回调函数
		auth.GET("/callback", handler.Callback)
		// 支付宝通知函数
		auth.POST("/notify", handler.Notify)
		// 新增Market
		auth.POST("/market/new", handler.MarketNew)
		// 交易hash更新 √
		auth.POST("/market/hash/update", handler.MarketHashUpdate)
		// 交易cancel hash更新 √
		auth.POST("/market/cancelHash/update", handler.MarketCancelHashUpdate)
		// 交易deal hash更新 √
		auth.POST("/market/dealHash/update", handler.MarketDealHashUpdate)

		// 新增Order √
		auth.POST("/order/new", handler.OrderNew)
		// 交易hash更新 √
		auth.POST("/order/hash/update", handler.OrderHashUpdate)
		// 关注信息添加删除
		auth.POST("/my_attention/post", handler.AttentionPost)
		chain.GET("/order/max_price/get", handler.OrderGetMaxPrice)

		// 市场列表页
		chain.GET("/market/all", handler.MarketAllHandle)
		// 蚂蚁链中查看市场中的资产列表
		auth.GET("/market/token/list", handler.TokenMarket)
		// 列表页
		chain.GET("/market/list", handler.GDMarketListHandle)
		// 市场详情页
		chain.GET("/market/detail", handler.GDMarketDetailsHandle)
		// nft详情页
		chain.GET("/nft/detail", handler.GDNfttDetailsHandle)

		/**
		钱包信息
		*/
		//新增 修改钱包信息
		auth.POST("/wallet_info/post", handler.WalletInfoPost)
		//查询钱包信息
		chain.GET("/wallet_info/get", handler.WalletInfoGet)

		// 钱包用户模糊搜索
		chain.GET("/wallet_info/list", handler.WalletInfoList)

		// 获取竞拍者的出价信息
		chain.GET("/bidding/discipline/get", handler.OrderGetByMarketId)

		chain.GET("/transaction/history/get", handler.MarketGetHistory)
		// 个人中心正在售卖
		chain.GET("/personal_center/selling/get", handler.MarketAllHandle)
		// 个人中心 创建 or 拥有
		chain.GET("/personal_center/owner/get", handler.PersonalCenterOwner)
		// 个人中心 关注
		chain.GET("/personal_center/favorites/get", handler.PersonalCenterFavorites)
		// 个人collection列表展示
		chain.GET("/personal_center/collections", handler.PersonalCollection)
		// 个人collection下资产的展示
		chain.GET("/personal_center/collection/details", handler.PersonalNFTFromCollection)
		//捐赠对象列表
		//chain.GET("/organisation/audit/get", handler.OrganizationGet)

		/**
		 * Home 首页定制查询
		 */
		// 查询限价购买列表 √
		chain.GET("/home/limiteds", handler.HomeLimiteds)
		// 查询拍卖列表 √
		chain.GET("/home/auctions", handler.HomeAuctions)
		// 查询banner列表 √
		chain.GET("/home/banners", handler.HomeBanners)
		// 查询notable drops
		chain.GET("/home/notable", handler.HomeNotable)
		// 查询top collection列表
		chain.GET("/home/top", handler.HomeTopCollections)
		// 查询trend categories
		chain.GET("/home/trending", handler.HomeTrending)

		/**
		statics collection ranking排列页
		*/
		chain.GET("/rankings", handler.RankingGetByLimit)

		/***
		collections展示页
		*/
		chain.GET("/explore/collections", handler.ExploreCollections)
		// 查询资产信息页中顶部collection信息
		chain.GET("/explore/collectioninfo", handler.GeneralCollectionInfo)
		// 查询collection中的资产信息
		chain.GET("/collection/details", handler.GetNFRFromCollection)

		/**
		  个人详情页
		*/
		chain.GET("/personal/details/get", handler.PersonalDetailsNew)
		//auth.POST("/organisation/audit/post", handler.OrganizationPost)
		auth.POST("/change/owner/post", handler.NftOwnerUpdate)

		/**
		  集合
		*/
		auth.POST("/create/collection/post", handler.CollectionInfoCreate)
		chain.GET("/collection/list/get", handler.CollectionListGet)

		chain.POST("/auth/token", handler.GenerateToken)
		chain.GET("/nonce/get", handler.GetNonce)
		chain.POST("/expiration/time", handler.Expiration)

		// ---------------------------------------------------------------------------------------------------------- 华丽的分割线
		/*
			// 根据tokenid获取nft √
			chain.GET("/nft/getByTokenId", handler.NftGetByTokenId) // 根据tokenid获取nft √
			// 根据sn获取nft √
			chain.GET("/nft/getBySn", handler.NftGetBySn)
			// 根据address获取可销售的nft 发布的拍卖status=1 发布的限价销售status=2 √
			chain.GET("/nft/mySalability", handler.NftMySalability)
			// 获取NFT分类 √
			chain.GET("/nft/type", handler.NftType) //
			// 根据mid获取Market √
			chain.GET("/market/getByMid", handler.MarketGetByMid)
			chain.GET("/order/getByOid", handler.OrderGetByOid)
			chain.GET("/order/getByBuyer", handler.OrdersByBuyer)
		*/

		// 用户相关路由
		chain.POST("/login", handler.Login)
		chain.POST("/register", handler.Register)
		chain.POST("/user/changePass", handler.ChangePassword)
		chain.DELETE("/user/delete/:id", handler.DeleteUser)
		chain.POST("/user/update", handler.UpdateUser)
		chain.POST("/user/find", handler.FindUserByID)
	}

	logger.Info("API server start", "listen", config.Conf.Console.Port)
	router.Run(config.Conf.Console.Port)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,"+
				"X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, "+
				"X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, "+
				"Content-Type, Pragma, Wallet")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
