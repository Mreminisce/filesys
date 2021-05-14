package handler

import (
	"fmt"
	"net/http"

	"filesys/config"
	"filesys/model"
	"filesys/util"

	"github.com/gin-gonic/gin"
)

func writeJSON(c *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	c.JSON(http.StatusOK, h)
}

func RegisterGet(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func RegisterPost(c *gin.Context) {
	res := gin.H{}
	defer writeJSON(c, res)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) < 3 || len(password) < 5 {
		res["message"] = "Register Message Invalid."
		return
	}
	userob, _ := model.GetUserObjectByName(username)
	if userob.Password != "" {
		fmt.Println("User already exists")
		return
	}
	fmt.Println("Count Users: ", model.CountUsers())
	user := &model.User{
		Username: username,
	}
	user.Password = util.Sha1([]byte(config.UserPwdSalt + password))
	if err := model.CreateUser(user); err != nil {
		res["message"] = "Create User Error"
		return
	}
	res["succeed"] = true
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func LoginGET(c *gin.Context) {
	islogin := util.IsLogin(c)
	sessions := util.GetSessions(c)
	c.HTML(http.StatusOK, "login.html", gin.H{
		"islogin":  islogin,
		"sessions": sessions,
	})
}

func LoginPOST(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")

	// maps := make(map[string]interface{})
	// maps["username"] = username
	// user, err := model.GetUserObjectByMaps(maps)

	user, err := model.GetUserObjectByName(username)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"code":    400,
			"message": "User Not Exist.",
			"data":    nil,
		})
		return
	}
	// 1. 校验用户名及密码
	if username == "" || password == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"code":    400,
			"message": "Login message can not be null.",
			"data":    nil,
		})
		return
	}
	if user.Password != util.Sha1([]byte(config.UserPwdSalt+password)) || err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"code":    400,
			"message": "Username or Password Error.",
			"data":    nil,
		})
		return
	}
	// 2. 生成token
	token := util.CreateToken(username)
	upRes := model.UpdateToken(username, token)
	if !upRes {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"message": "Username or Token Error.",
		})
		return
	}
	// 生成session  使nginx报502错误
	var sok chan int = make(chan int, 1)
	go util.LoginSession(c, user, sok)
	<-sok
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Login Success",
		"data": struct {
			Username string
			Token    string
		}{
			Username: username,
			Token:    token,
		},
	})
	c.Redirect(http.StatusMovedPermanently, "/")
}

func LogoutGET(c *gin.Context) {
	util.LogoutSession(c)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"code":    200,
		"message": "Logout Success",
		"data":    nil,
	})
	// c.JSON(http.StatusOK, gin.H{
	// 	"code":    200,
	// 	"message": nil,
	// 	"data":    nil,
	// })
	c.Redirect(http.StatusMovedPermanently, "/")
	return
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	c.Request.ParseForm()
	username := c.Request.Form.Get("username")
	token := c.Request.Form.Get("token")

	// 2. 验证token是否有效
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		fmt.Println("Token Valid")
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	// 3. 查询用户信息
	user, err := model.GetUserObjectByName(username)
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}
	c.HTML(http.StatusOK, "files.html", gin.H{
		"user":     user,
		"username": user.Username,
		"createat": user.CreatedAt,
	})

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	c.Writer.Write(resp.JSONBytes())
}
