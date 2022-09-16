package request

import uuid "github.com/satori/go.uuid"

// 登录ViewModel
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordStruct struct {
	UserName    string `json:"userName"`
	UserPass    string `json:"userPass"`
	NewPassword string `json:"newPassword"`
	Id          uint   `json:"id"`
}

type SetUserAuth struct {
	UUID        uuid.UUID `json:"uuid"`
	AuthorityId string    `json:"authorityId"`
}
