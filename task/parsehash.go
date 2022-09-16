package task

import (
	"ChainServer/config"
	"ChainServer/models"
	jsonrpc "ChainServer/task/rpcclient/jsonprc"
	"fmt"
	"strconv"
	"sync"

	"github.com/robfig/cron"
	"github.com/sea-project/go-logger"
)

var TaskCLientLock sync.Mutex
var load_client_once sync.Once

func CronTimer() {
	c := cron.New()                     // 新建一个定时任务对象
	c.AddFunc("0/5 * * * * *", func() { // 每5秒执行
		StartLoopParseHash() // todo 后续优化
	})
	c.Start()
}

var HttpClientMap = make(map[string]*jsonrpc.Http)

func InitHttpClientMap() {
	TaskCLientLock.Lock()
	defer TaskCLientLock.Unlock()
	for chain_name, v := range config.ChainTypeMap {
		rpc := v.(string)
		client := jsonrpc.NewETHTP(rpc)
		HttpClientMap[chain_name] = client
	}
}

func GetHttpClient(chain_name string) (client *jsonrpc.Http) {
	TaskCLientLock.Lock()
	defer TaskCLientLock.Unlock()
	v, ok := HttpClientMap[chain_name]
	if ok {
		client = v
	} else {
		rpc := config.ChainTypeMap[chain_name]
		client = jsonrpc.NewETHTP(rpc.(string))
		HttpClientMap[chain_name] = client
	}
	return client
}

func StartLoopParseHash() {
	load_client_once.Do(InitHttpClientMap)
	for chain_name, v := range config.ChainTypeMap {
		var rpc string
		rpc = v.(string)
		client := GetHttpClient(chain_name)
		go LoopParseHash(chain_name, rpc, client)
	}
}

// 遍历最新交易哈希解析
func LoopParseHash(chain_name, rpc string, client *jsonrpc.Http) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	logger.Info("chain_name ======= ", chain_name, "rpc ============ ", rpc)
	handleNft(chain_name, client)
	handleNftTransfer(chain_name, client)
	handleMarket(chain_name, client)
	handleOrder(chain_name, client)
	handleOrderLazy(chain_name, client)
}

// handleNft
func handleNft(chain_name string, client *jsonrpc.Http) {
	logger.Info("handleNft:", "status", "start")
	//recordList := model.GetNftHashList(chain_name)
	params_sql := fmt.Sprintf("chain_name = '%s' and tx_hash != '' and status = 0 and lazy = 0 ", chain_name)
	recordList := models.NftFind(params_sql)
	for _, value := range recordList {
		logger.Info("handleNft", "step", "for", "hash", value.TxHash)
		nftChainInfo, err := handleNftHash(client, value.TxHash)
		if err != nil {
			// 不建议自动判断失败更新字段
			logger.Warn("handleNft", "step", "handleNftHash", "hash", value.TxHash, "err", err)
			continue
			//err = model.NftStatusUpdate(value.Txhash, chain_name)
			//return
		} else {
			if value.BlockNumber > 0 && value.CreateTime > 0 {
				// 更新拥有者 token_id, lazy
				//err = model.NftOwnerUpdate(nftChainInfo.Owner, "0", "0", value.Id)
				models.NftUpdate(value, map[string]interface{}{"owner": nftChainInfo.Owner, "lazy": 0})
				if err != nil {
					logger.Warn("handleNft", "step", "NftOwnerUpdate", "token_id", nftChainInfo.TokenId, "owner", nftChainInfo.Owner, "err", err)
				}
			} else {
				// 全都更新 hash, token_id, creater, create_number, create_time, owner
				//err = model.NftChainUpdate(value.Txhash, nftChainInfo.TokenId, nftChainInfo.Owner,
				//	nftChainInfo.CreateNumber, nftChainInfo.CreateTime, nftChainInfo.Owner)
				// 全都更新  token_id, creater, create_number, create_time, owner
				update_map := make(map[string]interface{})
				update_map["token_id"] = nftChainInfo.TokenId
				update_map["owner"] = nftChainInfo.Owner
				update_map["block_number"] = nftChainInfo.CreateNumber
				create_time, _ := strconv.ParseInt(nftChainInfo.CreateTime, 10, 64)
				update_map["create_time"] = create_time
				update_map["owner"] = nftChainInfo.Owner
				update_map["status"] = 1
				models.NftUpdate(value, update_map)
				if err != nil {
					logger.Warn("handleNft", "step", "NftChainUpdate", "hash", value.TxHash,
						"token_id", nftChainInfo.TokenId, "creater", nftChainInfo.Owner,
						"create_number", nftChainInfo.CreateNumber, "create_time", nftChainInfo.CreateTime,
						"owner", nftChainInfo.Owner, "err", err)
				}
			}
		}

	}
	//logger.Info("handleNft:", "status", "end")
}

//
func handleNftTransfer(chain_name string, client *jsonrpc.Http) {
	logger.Info("transfer_hash:", "status", "start")
	params_sql := fmt.Sprintf("chain_name = '%s' and transfer_hash != '' and market_type = 3 and lazy = 0 ", chain_name)
	recordList := models.NftFind(params_sql)
	for _, value := range recordList {
		nftChainInfo, err := handleNftHash(client, value.TransferHash)
		update_map := make(map[string]interface{})
		if err != nil {
			//update_map["market_type"] = 0
			//models.NftUpdate(value, update_map)
			// TODO  区块链上交易未被确认时也可能会 err 暂时不做处理 后期加上时间判断超过一天的在做处理
			logger.Warn("transfer_hash", "step", "transfer_hash", "hash", value.TransferHash, "err", err)
		} else {
			update_map["owner"] = nftChainInfo.Owner
			update_map["market_type"] = 0
			if nftChainInfo.TokenId == value.TokenId {
				models.NftUpdate(value, update_map)
			}
			if err != nil {
				logger.Warn("handleNft", "step", "NftChainUpdate", "hash", value.TxHash,
					"token_id", nftChainInfo.TokenId, "creater", nftChainInfo.Owner,
					"create_number", nftChainInfo.CreateNumber, "create_time", nftChainInfo.CreateTime,
					"owner", nftChainInfo.Owner, "err", err)
			}
		}
	}
}

// handleNft
func handleMarket(chain_name string, client *jsonrpc.Http) {
	logger.Info("handleMarket:", "status", "start")
	//client := jsonrpc.NewETHTP(rpc)
	// 先处理新增挂单记录
	//free_gas := "0"
	//recordList := model.GetMarketHashList(chain_name)
	params_sql := fmt.Sprintf("chain_name = '%s' and tx_hash != '' and status = 0 and lazy = 0 ", chain_name)
	recordList := models.MarketListFind(params_sql)
	for _, value := range recordList {
		logger.Info("handleMarket", "step", "for recordList", "hash", value.TxHash)
		chainInfo, err := handleNftHash(client, value.TxHash)
		// 如果失败 除了更新时间还要更新状态为2(失败不能自动处理，建议手动处理)
		if err != nil {
			//model.MarketChainUpdate(value.Txhash,chainInfo.CreateTime,"2",chain_name)
			logger.Warn("handleMarket", "step", "handleNftHash", "hash", value.TxHash, "err", err)
		} else {
			// 如果成功 除了更新时间还要更新状态为1
			//err = model.MarketChainUpdate(value.Txhash, chainInfo.CreateTime, "1", chain_name)
			models.MarketListUpdate(value, map[string]interface{}{"create_time": chainInfo.CreateTime, "status": 1})
			if err != nil {
				logger.Warn("handleMarket", "step", "MarketChainUpdate", "hash", value.TxHash, "err", err)
			}
			// 还要更新nft表的market_status状态
			//err = model.NftMarketStatusUpdate(value.Tokenid, value.MarketType, chain_name)
			var nft_info models.TNft
			nft_info.Id = value.NftId
			models.NftUpdate(nft_info, map[string]interface{}{"market_type": value.MarketType})
			if err != nil {
				logger.Warn("handleMarket", "step", "NftMarketStatusUpdate", "tokenid", value.TokenId, "market_status", value.MarketType, "err", err)
			}
		}
	}
	// 再处理取消订单记录
	//cancelRecordList := model.GetMarketCancelHashList(chain_name, free_gas)
	cancel_sql := fmt.Sprintf("chain_name = '%s' and cancel_hash != '' and status = 1 and lazy = 0 ", chain_name)
	cancelRecordList := models.MarketListFind(cancel_sql)
	for _, value := range cancelRecordList {
		logger.Info("handleMarket", "step", "for cancelRecordList", "hash", value.TxHash)
		chainInfo, err := handleNftHash(client, value.CancelHash)
		// 如果失败 什么也不做
		if err != nil {
			logger.Warn("handleNft", "step", "handleNftHash", "hash", value.TxHash, "err", err)
		} else {
			// 如果成功 ,并且tokenid相等 更新状态为3
			if value.TokenId == chainInfo.TokenId {
				//err = model.MarketStatusUpdate(value.Txhash, "3", chain_name)
				models.MarketListUpdate(value, map[string]interface{}{"status": 3})
				if err != nil {
					logger.Warn("handleMarket", "step", "MarketStatusUpdate", "hash", value.TxHash, "err", err)
				}
			}
			// 还要更新nft表的market_status状态
			//err = model.NftMarketStatusUpdate(value.Tokenid, 0, chain_name)
			var nft_info models.TNft
			nft_info.Id = value.NftId
			models.NftUpdate(nft_info, map[string]interface{}{"market_type": 0})
			if err != nil {
				logger.Warn("handleMarket", "step", "NftMarketStatusUpdate", "tokenid", value.TokenId, "market_status", 0, "err", err)
			}
		}

	}
	// 最后处理deal订单记录：限价出售和拍卖出售
	//dealRecordList := model.GetMarketDealHashList(chain_name, free_gas)
	deal_sql := fmt.Sprintf("chain_name = '%s' and deal_hash != '' and status = 1 and lazy = 0 ", chain_name)
	dealRecordList := models.MarketListFind(deal_sql)
	for _, value := range dealRecordList {
		logger.Info("handleMarket", "step", "for dealRecordList", "hash", value.DealHash)
		chainInfo, err := handleNftHash(client, value.DealHash)
		// 如果失败 什么也不做
		if err != nil { // TODO 失败将market status 更为 2
			logger.Warn("handleNft", "step", "handleNftHash", "hash", value.TxHash, "err", err)
			//err = model.MarketUpdateStatusByTokenid(value.Tokenid, "2", chain_name)
		} else {
			// 如果成功 ,并且tokenid相等 更新状态为4
			if value.TokenId == chainInfo.TokenId {
				//err = model.MarketUpdateStatusByTokenid(value.Tokenid, "4", chain_name)
				models.MarketListUpdate(value, map[string]interface{}{"status": 4})
				if err != nil {
					logger.Warn("handleMarket", "step", "MarketStatusUpdate", "hash", value.TxHash, "err", err)
				}
				// 还要更新nft表的market_status状态
				//err = model.NftMarketStatusUpdate(value.Tokenid, 0, chain_name)
				var nft_info models.TNft
				nft_info.Id = value.NftId
				//err = model.NftMarketStatusUpdate(value.Tokenid, 0, chain_name)
				models.NftUpdate(nft_info, map[string]interface{}{"market_type": 0})
				if err != nil {
					logger.Warn("handleMarket", "step", "NftMarketStatusUpdate", "tokenid", value.TokenId, "market_status", 0, "err", err)
				}
				if chainInfo.Owner != "" {
					// 还要新拥有者 token_id, owner
					//err = model.NftOwnerUpdate(chainInfo.Owner, "0", "0", value.NftId)
					models.NftUpdate(nft_info, map[string]interface{}{"lazy": 0, "owner": chainInfo.Owner})
					if err != nil {
						logger.Warn("handleNft", "step", "NftOwnerUpdate", "token_id", chainInfo.TokenId, "owner", chainInfo.Owner, "err", err)
					}
				}
			}
		}

	}
	//logger.Info("handleMarket:", "status", "end")
}

// handleOrder
func handleOrder(chain_name string, client *jsonrpc.Http) {
	// 处理新增购买记录
	params_sql := fmt.Sprintf("chain_name = '%s' and tx_hash != '' and status = 0 and lazy = 0 ", chain_name)
	recordList := models.OrderFind(params_sql)
	for _, value := range recordList {
		logger.Info("handleOrder", "step", "for recordList", "hash", value.TxHash)
		chainInfo, err := handleNftHash(client, value.TxHash)
		// 如果失败 除了更新时间还要更新状态为2(失败不能自动处理，建议手动处理)
		// TODO  区块链上交易未被确认时也可能会 err 暂时不做处理 后期加上时间判断超过一天的在做处理
		if err != nil {
			//model.OrderChainUpdate(value.Txhash, "0", "2", chain_name)
			logger.Warn("handleOrder", "step", "handleNftHash", "hash", value.TxHash, "err", err)
		} else {
			// 如果成功 除了更新时间还要更新状态为1
			//err = model.OrderChainUpdate(value.Txhash, chainInfo.CreateTime, "1", chain_name)
			models.OrderUpdate(value, map[string]interface{}{"onchain_time": chainInfo.CreateTime, "status": 1})
			if err != nil {
				logger.Warn("handleOrder", "step", "MarketChainUpdate", "hash", value.TxHash, "err", err)
			}
			// 如果是限价出售，则更新market和nft的状态
			if value.MarketType == 2 && value.TokenId == chainInfo.TokenId {
				{
					var market_info models.TMarketList
					market_info.Id = value.MarketId
					models.MarketListUpdate(market_info, map[string]interface{}{"status": 4})
					if chainInfo.Owner != "" {
						// 还要新拥有者 token_id, owner
						var nft_info models.TNft
						nft_info.Id = value.NftId
						models.NftUpdate(nft_info, map[string]interface{}{"owner": chainInfo.Owner, "lazy": 0,
							"market_type": 0})
						if err != nil {
							logger.Warn("handleNft", "step", "NftOwnerUpdate", "token_id", chainInfo.TokenId, "owner", chainInfo.Owner, "err", err)
						}
					}
				}
			}

			//如果是竞价拍卖，则不需要再进行额外的更新操作
		}
	}
}

//// handleBidDeal 处理拍卖成交订单
//func handleBidDeal(chain_name string, client *jsonrpc.Http)  {
//
//}

func handleOrderLazy(chain_name string, client *jsonrpc.Http) {
	params_sql := fmt.Sprintf("chain_name = '%s' and tx_hash != '' and status = 0 and lazy = 1 ", chain_name)
	recordList := models.OrderFind(params_sql)
	for _, value := range recordList {
		logger.Info("handleOrder", "step", "for recordList", "hash", value.TxHash)
		chainInfo, err := handleNftHash(client, value.TxHash)
		if chainInfo.TokenId == "" {
			order_list := models.OrderFind(map[string]interface{}{"market_id": value.MarketId}) // token_id != 0
			if len(order_list) > 0 {
				chainInfo.TokenId = order_list[0].TokenId
			} else {
				continue
			}
		}
		var nft_info models.TNft
		nft_info.Id = value.NftId
		models.NftUpdate(nft_info, map[string]interface{}{"token_id": chainInfo.TokenId, "block_number": chainInfo.CreateNumber})

		// 更新 market
		var market_info models.TMarketList
		market_info.Id = value.MarketId
		if chainInfo.TokenId != "" {
			models.MarketListUpdate(market_info, map[string]interface{}{"token_id": chainInfo.TokenId})
		} else {
			logger.Warn("tokenid err ", chainInfo)
		}

		// 如果失败 除了更新时间还要更新状态为2(失败不能自动处理，建议手动处理)
		if err != nil {
			logger.Warn("handleOrder", "step", "handleNftHash", "hash", value.TxHash, "err", err)
		} else {
			// 如果成功 除了更新时间还要更新状态为1
			models.OrderUpdate(value, map[string]interface{}{"onchain_time": chainInfo.CreateTime, "status": 1})
			models.OrderUpdate(value, map[string]interface{}{"token_id": chainInfo.TokenId})
			value.TokenId = chainInfo.TokenId
			// 如果是限价拍卖，则更新market和nft的状态
			if value.MarketType == 2 && value.TokenId == chainInfo.TokenId {
				{
					models.MarketListUpdate(market_info, map[string]interface{}{"status": 4})
					if err != nil {
						logger.Warn("handleOrder", "step", "MarketUpdateStatusByTokenid", "hash", value.TxHash, "tokenid", value.TokenId, "status", "4", "err", err)
					}
					models.NftUpdate(nft_info, map[string]interface{}{"market_type": 0})
					if err != nil {
						logger.Warn("handleOrder", "step", "NftUpdateMarketStatusByTokenid", "hash", value.TxHash, "tokenid", value.TokenId, "status", "0", "err", err)
					}
					if chainInfo.Owner != "" {
						// 还要新拥有者 token_id, owner
						models.NftUpdate(nft_info, map[string]interface{}{"owner": chainInfo.Owner, "lazy": 0})
						if err != nil {
							logger.Warn("handleNft", "step", "NftOwnerUpdate", "token_id", chainInfo.TokenId, "owner", chainInfo.Owner, "err", err)
						}
					}
				}
			}
		}
	}
	// 处理del
	del_flow(chain_name, client)
}

func del_flow(chain_name string, client *jsonrpc.Http) {
	// 最后处理deal订单记录
	deal_sql := fmt.Sprintf("chain_name = '%s' and deal_hash != '' and status = 1 and lazy = 1 ", chain_name)
	dealRecordList := models.MarketListFind(deal_sql)
	for _, value := range dealRecordList {
		logger.Info("handleMarket", "step", "for dealRecordList", "hash", value.DealHash)
		chainInfo, err := handleNftHash(client, value.DealHash)
		// 如果失败 什么也不做
		if err != nil { // TODO 失败将market status 更为 2
			logger.Warn("handleNft", "step", "handleNftHash", "hash", value.TxHash, "err", err)
		} else {
			// 如果成功 ,并且tokenid相等 更新状态为4
			if value.TokenId == chainInfo.TokenId {
				//err = model.MarketUpdateStatusByTokenid(value.TokenId, "4", chain_name)
				models.MarketListUpdate(value, map[string]interface{}{"status": 4})
				if err != nil {
					logger.Warn("handleMarket", "step", "MarketStatusUpdate", "hash", value.TxHash, "err", err)
				}
				// 还要更新nft表的market_status状态
				var nft_info models.TNft
				nft_info.Id = value.NftId
				models.NftUpdate(nft_info, map[string]interface{}{"market_type": 0})
				if err != nil {
					logger.Warn("handleMarket", "step", "NftMarketStatusUpdate", "tokenid", value.TokenId, "market_status", 0, "err", err)
				}
				if chainInfo.Owner != "" {
					// 还要新拥有者 token_id, owner
					models.NftUpdate(nft_info, map[string]interface{}{"owner": chainInfo.Owner, "lazy": 0})
					if err != nil {
						logger.Warn("handleNft", "step", "NftOwnerUpdate", "token_id", chainInfo.TokenId, "owner", chainInfo.Owner, "err", err)
					}
				}
			}
		}
	}
}
