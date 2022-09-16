package response

type PageRows struct {
	Total  int         `json:"total"`  // 记录总数
	Result interface{} `json:"result"` // 记录列表
}
