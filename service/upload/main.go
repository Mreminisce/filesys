package main

import (
	"filesys/config"
	"filesys/model"
	"filesys/route"
	"fmt"
)

func main() {
	db := model.InitDB()
	db.AutoMigrate(
		&model.User{},
		&model.FileTable{},
		&model.UserFile{},
	)
	defer db.Close()
	g := route.InitRouter()
	g.Run(":8088")
	fmt.Printf("Upload is running: [%s]...\n", config.UploadServiceHost)
}
