package request

type OrganizationInfo struct {
	OrgName     string `json:"org_name"`     // 组织名称
	OrgDesc     string `json:"org_desc"`     // 组织名称
	CreditCode  string `json:"credit_code"`  // 信用代码
	LegalPerson string `json:"legal_person"` // 头像地址
	FullName    string `json:"full_name"`    // 联系人
	Position    string `json:"position"`     // 职位
	PhoneNum    string `json:"phone_num"`    // 手机号
	Verified    string `json:"verified"`     //  验证码
	FileUrl1    string `json:"file_url1"`
	FileUrl2    string `json:"file_url2"`
	FileUrl3    string `json:"file_url3"`
}
