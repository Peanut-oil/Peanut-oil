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
		user.POST("/get_user_info", GetUserInfo)   // 获取用户信息
		user.POST("/get_rank_list", GetRankList)   // 获取排行榜单
		user.POST("/add_rank_score", AddRankScore) // 增加排行分数
	}
}

func GetRankList(c *gin.Context) {
	var ps struct {
		RankTypeOneClass int `json:"rank_type_one_class" form:"rank_type_one_class" binding:"required,min=1,max=3"`
		RankTypeTwoClass int `json:"rank_type_two_class" form:"rank_type_one_class" binding:"required,min=1,max=3"`
	}
	err := c.ShouldBind(&ps)
	if err != nil {
		c.JSON(http.StatusOK, helper.Response(def.CodeError, def.MsgParamErr, nil))
		return
	}
	rankList, err := service.GetRankList(ps.RankTypeOneClass, ps.RankTypeTwoClass)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", rankList))
	return
}

func AddRankScore(c *gin.Context) {
	deviceId := c.GetString("did")
	var ps struct {
		RankTypeOneClass int    `json:"rank_type_one_class" form:"rank_type_one_class" binding:"required,min=1,max=3"`
		RankTypeTwoClass int    `json:"rank_type_two_class" form:"rank_type_two_class" binding:"required,min=1,max=3"`
		Score            int    `json:"score" form:"score" binding:"required"`
		NickName         string `json:"nick_name" form:"nick_name"`
		Avatar           string `json:"avatar" form:"avatar"`
	}
	err := c.ShouldBind(&ps)
	if err != nil {
		c.JSON(http.StatusOK, helper.Response(def.CodeError, def.MsgParamErr, nil))
		return
	}
	err = store.AddRankScore(ps.Score, ps.RankTypeOneClass, ps.RankTypeTwoClass, ps.NickName, ps.Avatar, deviceId)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, def.MsgSystemErr, nil))
		return
	}

	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", nil))
	return
}

// 统一通过设备id登陆
func GetUserInfo(c *gin.Context) {
	deviceId := c.GetString("did")
	uInfo, err := service.UserLoginByDeviceId(deviceId)
	if err != nil {
		c.JSON(http.StatusOK, util.Response(def.CodeError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, util.Response(def.CodeSucc, "ok", uInfo))
	return
}
