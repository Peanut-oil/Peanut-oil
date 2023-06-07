package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.POST("login", Login) // 用户登陆&注册
	}
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	return
}
