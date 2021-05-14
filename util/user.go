package util

import (
	"filesys/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserSession struct {
	Username   string `json:"username"`
	Userid     int    `json:"userid"`
	Useravatar string `json:"useravatar"`
}

func GetSessions(c *gin.Context) (sessions *UserSession) {
	username := GetSession(c, "username")
	userid, _ := strconv.Atoi(GetSession(c, "userid"))
	useravatar := GetSession(c, "useravatar")
	sessions = &UserSession{
		Username:   username,
		Userid:     userid,
		Useravatar: useravatar,
	}
	return
}

func LoginSession(c *gin.Context, user model.User, sok chan int) {
	SetSession(c, "username", user.Username)
	SetSession(c, "userid", strconv.Itoa(int(user.ID)))
	SetSession(c, "useravatar", user.Avatar)
	sok <- 1
}

func LogoutSession(c *gin.Context) {
	DeleteSession(c, "username")
	DeleteSession(c, "userid")
	DeleteSession(c, "useravatar")
}

func IsLogin(c *gin.Context) bool {
	username := GetSession(c, "username")
	return len(username) > 0
}
