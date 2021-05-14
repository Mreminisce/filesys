package handler

import (
	"filesys/model"
	"filesys/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	islogin := util.IsLogin(c)
	sessions := util.GetSessions(c)
	user, _ := model.GetUserObjectByID(sessions.Userid)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"islogin":  islogin,
		"sessions": sessions,
		"user":     user,
	})
}
func Handle404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error.html", gin.H{
		"message": "404 not found...",
	})
}
