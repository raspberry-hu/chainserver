package models

import (
	"ChainServer/api/request"
	. "ChainServer/models/postgresql"
	myUtils "ChainServer/utils"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sea-project/go-logger"
)

type UserInfo struct {
	Id         int       `json:"id"`
	UserPass   string    `json:"userPass"`
	UserName   string    `json:"userName"`
	UserDesc   string    `json:"userDesc"`
	ImageUrl   string    `json:"imageUrl"`
	BannerUrl  string    `json:"bannerUrl"`
	EmailAddr  string    `json:"emailAddr"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Permission int       `json:"permission"`
}

// TableName 会将 UserInfo 的表名重写为 `t_user_info`
func (UserInfo) TableName() string {
	return "t_user_info"
}

// 根据用户Id查询用户信息
func UserFindByID(id int) (user UserInfo, err error) {
	err = Db.Where("id = $1 ", id).Find(&user).Error
	return user, err
}

// UserInfoInsert 新增用户
func UserInfoInsert(user UserInfo) bool {
	var usr UserInfo                                    //用以判断用户名是否已存在
	user.UserPass = myUtils.MD5V([]byte(user.UserPass)) //变为md5摘要
	err := Db.Where("user_name = ?", user.UserName).First(&usr).Error
	if gorm.IsRecordNotFoundError(err) { //如果用户名不存在，可进行注册
		er := Db.Create(&user)
		flag := Db.NewRecord(user) // 主键是否为空
		if er.Error != nil || flag {
			logger.Error("insert user err failed", er.Error)
			return false
		}
		return true
	} else {
		return false
	}
}

// UserInfoFind使用sql根据条件查询用户切片
func UserInfoFind(sql interface{}) (user_info_arry []UserInfo) {
	Db.Where(sql).Find(&user_info_arry)
	return user_info_arry
}

// // UserInfoUpdate 更新用户信息
// func UserInfoUpdate(user_info UserInfo, update_info interface{}) bool {
// 	Db.Model(&user_info).Updates(update_info)
// 	return true
// }

// UserInfoList 分页查询用户信息 ,
func UserInfoList(user_name, limit_sql string) (list []TWalletInfo) {
	sql := fmt.Sprintf("SELECT * FROM mintklub.t_user_info where   user_name like '%%%s%%' %s",
		user_name, limit_sql)
	Db.Raw(sql).Find(&list)
	return list
}

//
func UserInfoLikeCount(user_name string) (total int) {
	sql := fmt.Sprintf("SELECT count(*) FROM mintklub."+
		"t_user_info where  user_name like '%%%s%%'",
		user_name)
	rows, err := Db.Raw(sql).Rows()
	if err != nil {
		logger.Error("UserInfoLikeCount error", "err", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			logger.Error("UserInfoLikeCount error", "err", err)
		}
	}
	return total
}

// 登录
func Login(u *UserInfo) (user *UserInfo, err error) {
	user = &UserInfo{}
	pass := myUtils.MD5V([]byte(u.UserPass))
	err = Db.Where("user_name = $1 AND user_pass = $2", u.UserName, pass).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

// UpdateUser 更新
func UpdateUser(newPerson UserInfo) (UserInfo, error) {
	var u UserInfo
	db := Db
	err := db.First(&u, newPerson.Id).Error
	if err != nil {
		return u, err
	}
	newPerson.UpdateTime = time.Now()
	newPerson.UserPass = u.UserPass
	err = db.Save(&newPerson).Error
	return newPerson, err
}

// ChangePssword 更改密码
func ChangePssword(model request.ChangePasswordStruct) (person UserInfo, err error) {
	db := Db
	person = UserInfo{}
	err = db.First(&person, model.Id).Error
	if err != nil {
		return person, err
	}
	err = Db.Model(&person).Update("user_pass", myUtils.MD5V([]byte(model.NewPassword))).Error
	return person, nil
}

// DeleteUser 删除
func DeleteUser(id int) (err error) {
	db := Db
	err = db.Delete(&UserInfo{}, id).Error
	return err
}
