package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/util"
	"github.com/gin-gonic/gin/auth"
	"net/http"
)

func AddUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user", auth.CheckLogin)
	{
		user.POST("/login", Login)              // 用户登陆&注册
		user.POST("/create_order", CreateOrder) // 创建订单
	}
}

// 统一通过设备id登陆
func Login(c *gin.Context) {
	deviceId := c.GetString("deviceId")

	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", nil))
	return
}

func CreateOrder(c *gin.Context) {

}
