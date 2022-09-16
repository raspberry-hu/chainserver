package types

type HomeNftMarket struct {
	Id            string `json:"id"`
	MediaUri      string `json:"media_uri"`
	NftName       string `json:"nft_name"`
	TokenId       string `json:"tokenid"`
	StartingPrice string `json:"starting_price"`
	TokenType     string `json:"token_type"`
}

type NftMarket struct {
	Id            string `json:"id"`
	MediaUri      string `json:"media_uri"`
	NftName       string `json:"nft_name"`
	TokenId       string `json:"tokenid"`
	StartingPrice string `json:"starting_price"`
	TokenType     string `json:"token_type"`
	TokenAddr     string `json:"token_addr"`
	CreateNumber  string `json:"create_number"`
	CreateTime    string `json:"create_time"`
	Owner         string `json:"owner"`
}
