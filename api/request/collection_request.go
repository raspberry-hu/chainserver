package request

type CollectionCreate struct {
	UserId           int    `json:"user_id"`
	CollectionName   string `json:"collection_name"`
	ChainName        string `json:"chain_name"`
	LogoImage        string `json:"logo_image"`
	FeaturedImageUrl string `json:"featured_image_url"`
	BannerImageUrl   string `json:"banner_image_url"`
	CollectionDesc   string `json:"collection_desc"`
	CurrencyName     string `json:"currency_name"`
	CreateTax        int    `json:"create_tax"`
	Category         string `json:"category"`
}
