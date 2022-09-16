package request

type CollectionRankNew struct {
	CollectionId   	int     `json:"collection_id"`   	// 集合ID
	LogoImageURL    string  `json:"logo_image"`      	// 集合标志图片
	CollectionName 	string  `json:"collection_name"` 	// 集合名称
	BannerImageURL 	string	`json:"banner_image_url"`	// 背景图片
	CreateTax	   	int 	`json:"create_tax"`			// 铸造税
	CategoryName  	string  `json:"category_name`  	 	// 类别名称
	CollectionDesc 	string	`json:"collection_desc"` 	// 集合描述
	ChainName      	string  `json:"chain_name"`     	// 区块链名称
	Owner          	string	`json:"owner"`				// 创建者
	Items          	int     `json:"items"`          	// 资产总量
	Favorites 	   	int 	`json:"favorites"`			// 收藏资产数额
}
