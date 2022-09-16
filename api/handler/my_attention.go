package handler

import (
	"ChainServer/api/request"
	"ChainServer/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sea-project/go-logger"
)

func AttentionPost(ctx *gin.Context) {
	var myattention request.MyAttentionRequest
	ctx.ShouldBindJSON(&myattention)
	// WalletAddr := myattention.WalletAddr
	UserId := myattention.UserId
	add_nft_id_arry := myattention.AddNftId

	del_nft_id_arry := myattention.DelNftId
	chain_name := myattention.ChainName
	fmt.Println(chain_name)

	if UserId != 0 {
		for _, add_nft_id := range add_nft_id_arry {
			favorites_arry := models.FavoritesFind(map[string]interface{}{"user_id": UserId, "nft_id": add_nft_id})
			if len(favorites_arry) > 0 {
				continue
			} else {
				var favorites_info models.TFavorites
				favorites_info.NftId = add_nft_id
				favorites_info.UserId = UserId
				favorites_info.CreateTime = time.Now()
				models.FavoritesInsert(favorites_info)
				//增加收藏个数
				var collection_info models.TCollectionInfo
				var favoriteCount int
				nftInfo := models.NftFind(map[string]interface{}{"id": add_nft_id})
				collection_info.Id = nftInfo[0].CollectionId
				collectionArr := models.CollectionFind(map[string]interface{}{"id": collection_info.Id})
				if len(collectionArr) == 0 {
					logger.Error("没找到Collection")
				} else {
					favoriteCount = collectionArr[0].Favorites
				}
				favoriteCount = favoriteCount + 1

				models.UpdateFavoriteByCollection(collection_info, map[string]interface{}{"favorites": favoriteCount})
			}
		}
		for _, del_nft_id := range del_nft_id_arry {
			favorites_arry := models.FavoritesFind(map[string]interface{}{"nft_id": del_nft_id, "user_id": UserId})
			if len(favorites_arry) > 0 {
				models.FavoritesDel(favorites_arry[0])
				//	减少收藏个数
				var collection_info models.TCollectionInfo
				var favoriteCount int
				nftInfo := models.NftFind(map[string]interface{}{"id": del_nft_id})
				collection_info.Id = nftInfo[0].CollectionId
				collectionArr := models.CollectionFind(map[string]interface{}{"id": collection_info.Id})
				if len(collectionArr) == 0 {
					logger.Error("没找到Collection")
				} else {
					favoriteCount = collectionArr[0].Favorites
				}
				favoriteCount = favoriteCount - 1

				models.UpdateFavoriteByCollection(collection_info, map[string]interface{}{"favorites": favoriteCount})
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
