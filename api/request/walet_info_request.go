package request

type WalletInfo struct {
	WalletAddr string `json:"wallet_addr"` // 钱包地址
	UserId     int    `json:"user_id"`     // 钱包地址
	UserName   string `json:"user_name"`   // 用户名
	ImageUrl   string `json:"image_url"`   // 头像地址
	BannerUrl  string `json:"banner_url"`  // 横幅地址
	PassWord   string `json:"pass_word"`   // 密码
	ChainType  string `json:"chain_type"`  // 链类型
}

type MyAttentionRequest struct {
	WalletAddr string `json:"wallet_addr"` // 钱包地址
	UserId     int    `json:"user_id"`     // 钱包地址
	AddNftId   []int  `json:"add_nft_id"`  // 关注的token_id
	DelNftId   []int  `json:"del_nft_id"`  // 取消关注的token_id
	ChainName  string `json:"chain_name"`  // 链类型
}

type AuthRequest struct {
	WalletAddr string `json:"wallet_addr"` // 钱包地址
	UserId     int    `json:"user_id"`     // 钱包地址
	Sign       string `json:"sign"`        // 签名
}
