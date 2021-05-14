package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type FileTable struct {
	gorm.Model
	UserID       int    `gorm:"default:0" json:"user_id"`
	FileHash     string `gorm:"default:''" json:"filehash"`
	FileName     string `gorm:"default:''" json:"filename"`
	FileSize     int    `gorm:"default:0" json:"filesize"`
	FileAddr     string `gorm:"default:''" json:"fileaddr"`
	FileType     string `gorm:"default:''" json:"filetype"`
	OrigiName    string `gorm:"default:''" json:"originame"`
	Isimage      int    `gorm:"default:0" json:"isimage"`
	Width        int    `gorm:"default:0" json:"width"`
	Height       int    `gorm:"default:0" json:"height"`
	DownloadsCnt int    `gorm:"default:0" json:"downloads_cnt"`
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string) (*FileTable, error) {
	tfile := FileTable{}
	err := DB.Model(&FileTable{}).Where("file_hash=?", filehash).Find(&tfile).Error
	return &tfile, err
}

// GetFileMetaList : 从mysql批量获取文件元信息
func GetFileMetaList(limit int) ([]FileTable, error) {
	var tfiles []FileTable
	err := DB.Model(&FileTable{}).Limit(limit).Find(&tfiles).Error
	return tfiles, err
}

// OnFileUploadFinished : 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int, fileaddr string) bool {
	file := &FileTable{
		FileHash: filehash,
		FileName: filename,
		FileSize: filesize,
		FileAddr: fileaddr,
	}
	err := DB.Create(file).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// UpdateFileLocation : 更新文件的存储地址(如文件被转移了)
func UpdateFileLocation(filehash string, fileaddr string) bool {
	err := DB.Model(&FileTable{}).Where("file_hash=?", filehash).Update("file_addr", fileaddr).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
