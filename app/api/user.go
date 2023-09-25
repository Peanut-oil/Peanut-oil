package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/app/auth"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gin-gonic/gin/app/service"
	"github.com/gin-gonic/gin/app/store"
	"github.com/gin-gonic/gin/app/util"
	"net/http"
)

func AddUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user", auth.CheckLogin)
	{
		user.POST("/login", Login)                 // 用户登陆&注册
		user.POST("/create_order", CreateOrder)    // 创建订单
		user.POST("/get_rank_list", GetRankList)   // 获取排行榜单
		user.POST("/add_rank_score", AddRankScore) // 增加排行分数
	}
}

func GetRankList(c *gin.Context) {
	var ps struct {
		RankType int `json:"rank_type" form:"rank_type" binding:"required,min=1,max=3"`
	}
	err := c.ShouldBind(&ps)
	if err != nil {
		c.JSON(http.StatusOK, helper.Response(def.CodeError, "param error", nil))
		return
	}
	rankList, err := service.GetRankList(ps.RankType)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", rankList))
	return
}

func AddRankScore(c *gin.Context) {
	var ps struct {
		Uid      int `json:"uid" form:"uid" binding:"required"`
		RankType int `json:"rank_type" form:"rank_type" binding:"required"`
		Score    int `json:"score" form:"score" binding:"required"`
	}
	err := c.ShouldBind(&ps)
	if err != nil {
		c.JSON(http.StatusOK, helper.Response(def.CodeError, "param error", nil))
		return
	}
	err = store.AddRankScore(ps.Score, ps.Uid, ps.RankType)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, "系统异常", nil))
		return
	}

	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", nil))
	return
}

// 统一通过设备id登陆
func Login(c *gin.Context) {
	deviceId := c.GetString("deviceId")
	_, err := service.UserLoginByDeviceId(deviceId)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", nil))
	return
}

func CreateOrder(c *gin.Context) {

}
