package models

import (
	. "ChainServer/models/postgresql"
)

type TOrganizationInfo struct {
	ID       int
	OrgName  string
	ImageUrl string
	UserId   int
}

func OrganizationInfoFind(rows, offset int, sql interface{}) (t_organization_info []TOrganizationInfo) {
	Db.Limit(rows).Offset(offset).Where(sql).Find(&t_organization_info)
	return t_organization_info
}
