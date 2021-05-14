package route

import (
	"filesys/handler"
	"filesys/util"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	sessions "github.com/tommy351/gin-sessions"
)

func InitRouter() *gin.Engine {
	g := gin.New()
	gin.SetMode("debug")
	setMiddleware(g)
	setSession(g)
	setTemplate(g)
	g.NoRoute(handler.Handle404)
	registerApi(g)
	return g
}

func registerApi(g *gin.Engine) {
	// 用户相关接口
	g.GET("/", handler.Home)
	g.GET("/register", handler.RegisterGet)
	g.POST("/register", handler.RegisterPost)
	g.GET("/login", handler.LoginGET)
	g.POST("/login", handler.LoginPOST)
	g.GET("/logout", handler.LogoutGET)
	g.POST("/user/info", handler.UserInfoHandler)

	// 文件存取接口
	g.GET("/file/upload", handler.UploadGet)
	g.POST("/file/upload", handler.UploadHandler)
	g.GET("/file/upload/suc", handler.UploadSucHandler)
	g.POST("/file/meta", handler.GetFileMetaHandler)
	g.GET("/file/query", handler.FileQueryHandler)
	g.POST("/file/download", handler.DownloadHandler)
	g.POST("/file/update", handler.FileMetaUpdateHandler)
	g.POST("/file/delete", handler.FileDeleteHandler)
	// 秒传接口
	g.POST("/file/fastupload", handler.TryFastUploadHandler)
	g.POST("/file/downloadurl", handler.DownloadURLHandler)

	// // 分块上传接口
	// g.POST("/file/mpupload/init", handler.InitialMultipartUploadHandler)
	// g.POST("/file/mpupload/uppart", handler.UploadPartHandler))
	// g.POST("/file/mpupload/complete", handler.CompleteUploadHandler)

}

func setMiddleware(g *gin.Engine) {
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	g.Use(Cors())
	g.Use(limit.MaxAllowed(100))
}

func setSession(g *gin.Engine) {
	store := sessions.NewCookieStore([]byte("FileSysSession"))
	// store.Options(sessions.Options{
	// 	HttpOnly: true,
	// 	Path:     "/",
	// 	MaxAge:   86400 * 30,
	// })
	g.Use(sessions.Middleware("filesys_session", store))
}

func setTemplate(g *gin.Engine) {
	funcMap := template.FuncMap{}
	g.SetFuncMap(funcMap)
	g.StaticFS("/static", http.Dir("./static"))
	g.StaticFS("/upload", http.Dir("./upload"))
	//g.LoadHTMLGlob("views/*/**/***")
	g.LoadHTMLGlob("view/*")
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if util.GetSession(c, "islogin") != "1" {
			c.Redirect(301, "/888")
			return
		}
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
