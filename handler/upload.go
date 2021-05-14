package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"filesys/config"
	"filesys/meta"
	"filesys/model"
	"filesys/store/ceph"
	"filesys/store/oss"
	"filesys/util"

	"github.com/gin-gonic/gin"
)

func UploadGet(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

// UploadHandler ： 处理文件上传
func UploadHandler(c *gin.Context) {
	// 接收文件流及存储到本地目录
	file, head, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Failed to get data, err:%s\n", err.Error())
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "./tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("Failed to create file, err:%s\n", err.Error())
		return
	}
	defer newFile.Close()

	data, err := io.Copy(newFile, file)
	fileMeta.FileSize = int(data)
	if err != nil {
		fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
		return
	}

	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	// 游标重新回到文件头部
	newFile.Seek(0, 0)

	if config.CurrentStoreType == config.StoreCeph {
		// 文件写入Ceph存储
		data, _ := ioutil.ReadAll(newFile)
		cephPath := "/ceph/" + fileMeta.FileSha1
		_ = ceph.PutObject("userfile", cephPath, data)
		fileMeta.Location = cephPath
	} else if config.CurrentStoreType == config.StoreOSS {
		// 文件写入OSS存储
		ossPath := "oss/" + fileMeta.FileSha1
		err = oss.Bucket().PutObject(ossPath, newFile)
		if err != nil {
			fmt.Println(err.Error())
			c.Writer.Write([]byte("Upload failed!"))
			return
		}
		fileMeta.Location = ossPath
	}

	// meta.UpdateFileMeta(fileMeta)
	_ = meta.UpdateFileMetaDB(fileMeta)

	// 更新用户文件表记录
	c.Request.ParseForm()
	username := c.Request.Form.Get("username")
	suc := model.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
	if suc {
		c.HTML(http.StatusOK, "files.html", nil)
	} else {
		c.Writer.Write([]byte("Upload failed!"))
		c.Redirect(http.StatusMovedPermanently, "/")
		c.HTML(http.StatusNotFound, "upload.html", "Upload Failed.")
	}
}

// UploadSucHandler : 上传已完成
func UploadSucHandler(c *gin.Context) {
	c.Writer.Write([]byte("Upload Finished!"))
	c.JSON(http.StatusOK, "Upload File Succeed")
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(c *gin.Context) {
	c.Request.ParseForm()

	filehash := c.Request.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if fMeta != nil {
		data, err := json.Marshal(fMeta)
		if err != nil {
			fmt.Println(err)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.Write(data)
	} else {
		c.Writer.Write([]byte(`{"code":-1,"msg":"no such file"}`))
	}
}

// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	c.Request.ParseForm()

	limitCnt, _ := strconv.Atoi(c.Request.Form.Get("limit"))
	username := c.Request.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := model.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(c *gin.Context) {
	c.Request.ParseForm()
	fsha1 := c.Request.Form.Get("filehash")
	fm, _ := meta.GetFileMetaDB(fsha1)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	f, err := os.Open(fm.Location)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	c.Writer.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	c.Writer.Write(data)
}

// FileMetaUpdateHandler ： 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context) {
	c.Request.ParseForm()

	opType := c.Request.Form.Get("op")
	fileSha1 := c.Request.Form.Get("filehash")
	newFileName := c.Request.Form.Get("filename")

	if opType != "0" {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	if c.Request.Method != "POST" {
		c.Writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	// TODO: 更新文件表中的元信息记录

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(c *gin.Context) {
	c.Request.ParseForm()
	fileSha1 := c.Request.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	// 删除文件
	os.Remove(fMeta.Location)
	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)
	// TODO: 删除表文件信息

	c.Writer.WriteHeader(http.StatusOK)
}

// TryFastUploadHandler : 尝试秒传接口
func TryFastUploadHandler(c *gin.Context) {
	c.Request.ParseForm()

	// 1. 解析请求参数
	username := c.Request.Form.Get("username")
	filehash := c.Request.Form.Get("filehash")
	filename := c.Request.Form.Get("filename")
	filesize, _ := strconv.Atoi(c.Request.Form.Get("filesize"))

	// 2. 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. 查不到记录则返回秒传失败
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		c.Writer.Write(resp.JSONBytes())
		return
	}

	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	suc := model.OnUserFileUploadFinished(username, filehash, filename, filesize)
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		c.Writer.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	c.Writer.Write(resp.JSONBytes())
}

// DownloadURLHandler : 生成文件的下载地址
func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.Form.Get("filehash")
	// 从文件表查找记录
	row, _ := model.GetFileMeta(filehash)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

	// TODO: 判断文件存在OSS，还是Ceph，还是在本地
	if strings.HasPrefix(row.FileAddr, "/tmp") {
		username := c.Request.Form.Get("username")
		token := c.Request.Form.Get("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s", c.Request.Host, filehash, username, token)
		c.Writer.Write([]byte(tmpUrl))
	} else if strings.HasPrefix(row.FileAddr, "/ceph") {
		// TODO: ceph下载url
	} else if strings.HasPrefix(row.FileAddr, "oss/") {
		// oss下载url
		signedURL := oss.DownloadURL(row.FileAddr)
		c.Writer.Write([]byte(signedURL))
	}
}
