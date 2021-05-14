package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Token    string `gorm:"not null" json:"token"`
	Email    string `gorm:"" json:"email"`
	Avatar   string `gorm:"not null" json:"avatar"`
	IsAdmin  bool   `gorm:"default:1" json:"isadmin"`
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	err := DB.Model(&User{}).Where("username=?", username).Update("token", token).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func CreateUser(user *User) error {
	return DB.Save(user).Error
}

func (user *User) UpdateUser() (err error) {
	return DB.Save(&user).Error
}

func UpdateUserByMaps(maps interface{}, items map[string]interface{}) (err error) {
	err = DB.Model(&User{}).Where(maps).Updates(items).Error
	return
}

func UpdateUserNameByID(newName string, id int) (err error) {
	var m = make(map[string]interface{})
	m["id"] = id
	err = UpdateUserByMaps(m, map[string]interface{}{"username": newName})
	return
}

func UpdateUserAvatar(user *User, avatar string) error {
	return DB.Model(&user).Update(User{Avatar: avatar}).Error
}

func UpdateUserAvatarByID(newAvatar string, id int) (err error) {
	var m = make(map[string]interface{})
	m["id"] = id
	err = UpdateUserByMaps(m, map[string]interface{}{"avatar": "/" + newAvatar})
	return
}

func UpdateUserEmail(user *User, email string) error {
	if len(email) > 0 {
		return DB.Model(user).Update("email", email).Error
	}
	return DB.Model(user).Update("email", gorm.Expr("NULL")).Error
}

func (user *User) DeleteUserByID(id int) error {
	user.ID = uint(id)
	// Unscoped 永久删除
	return DB.Unscoped().Delete(&user).Error
}

func DeleteUserByMaps(maps interface{}) (err error) {
	err = DB.Unscoped().Where(maps).Delete(&User{}).Error
	return
}

func IfUsernameExist(username string) bool {
	var user User
	DB.Model(&User{}).Select("id").Where("username=?", username).First(&user)
	return user.ID > 0
}

func GetUserByID(id int) (user *User, err error) {
	err = DB.First(&user, id).Error
	return
}

func GetUserObjectByID(id int) (user User, err error) {
	err = DB.Model(&User{}).Where("id=?", id).First(&user).Error
	return
}

func GetUserObjectByName(name string) (user User, err error) {
	err = DB.Model(&User{}).Where("username=?", name).Find(&user).Error
	return
}

func GetUserObjectByMaps(maps interface{}) (user User, err error) {
	err = DB.Model(&User{}).Where(maps).First(&user).Error
	return
}

func GetUserByUsername(name string) (user *User, err error) {
	err = DB.First(&user, "username=?", name).Error
	return
}

func GetUserByEmail(email string) (user *User, err error) {
	// err=DB.Where("email=?",email).First(&user).Error
	err = DB.First(&user, "email=?", email).Error
	return
}

func GetUsers(limit int, order string, maps interface{}) (user []User, err error) {
	err = DB.Model(&User{}).Order(order).Limit(limit).Find(&user).Error
	return
}

func ListUsers() (users []*User, err error) {
	err = DB.Order("id").Find(&users).Error
	return
}

func CountUsers() (count int) {
	DB.Model(&User{}).Count(&count)
	return
}

func ResetPasswordByID(newPassword string, id int) (err error) {
	var m = make(map[string]interface{})
	m["id"] = id
	err = UpdateUserByMaps(m, map[string]interface{}{"password": newPassword})
	return
}
