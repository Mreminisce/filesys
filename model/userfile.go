package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// UserFile : 用户文件表结构体
type UserFile struct {
	gorm.Model
	UserID   uint   `gorm:"not null" json:"userid"`
	FileID   uint   `gorm:"not null" json:"fileid"`
	Username string `gorm:"not null" json:"username"`
	FileHash string `gorm:"default:''" json:"filehash"`
	FileName string `gorm:"default:''" json:"filename"`
	FileSize int    `gorm:"default:0" json:"filesize"`
	UploadAt string `gorm:"default:''" json:"uploadat"`
}

// OnUserFileUploadFinished : 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int) bool {
	ufile := &UserFile{
		Username: username,
		FileHash: filehash,
		FileName: filename,
		FileSize: filesize,
	}
	err := DB.Create(ufile).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// QueryUserFileMetas : 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	var userFiles []UserFile
	err := DB.Model(&FileTable{}).Where("username=?", username).Limit(limit).Find(&userFiles).Error
	return userFiles, err
}
