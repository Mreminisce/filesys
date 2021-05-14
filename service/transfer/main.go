package main

import (
	"bufio"
	"encoding/json"
	"filesys/config"
	"filesys/model"
	"filesys/rabbitmq"
	"filesys/store/oss"
	"log"
	"os"
)

// ProcessTransfer : 处理文件转移
func ProcessTransfer(msg []byte) bool {
	log.Println(string(msg))
	// 解析 message
	pubData := rabbitmq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 根据临时存储文件路径创建文件句柄
	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 通过文件句柄将文件读取出来上传到OSS
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	// 更新文件的存储路径到文件表
	_ = model.UpdateFileLocation(pubData.FileHash, pubData.DestLocation)
	return true
}

func main() {
	if !config.AsyncTransferEnable {
		log.Println("Aysnc transfer didn't start, need to check config...")
		return
	}
	log.Println("Transfer is running...")
	rabbitmq.StartConsume(config.TransOSSQueueName, "transfer_oss", ProcessTransfer)
}
