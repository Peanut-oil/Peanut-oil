package auth

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/helper"
	"github.com/gin-gonic/gin/app/pkg/crypto"
	"github.com/gin-gonic/gin/app/pkg/serialize"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
)

func CheckSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			return
		}
		signKey := def.SignKey
		sign := c.Request.FormValue("sign")
		method := c.Request.Method
		path := c.Request.URL.Path
		for _, m := range def.SignWhiteMethod {
			if strings.Contains(path, m) {
				return
			}
		}
		params := c.Request.Form
		var paramStr string
		keys := make([]string, 0, len(params))
		for k := range params {
			if k == "sign" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, v := range keys {
			paramStr += "&" + v + "=" + params[v][0]
		}
		sourceStr := method + "&" + path + paramStr
		hashSign, _ := hex.DecodeString(crypto.Hmac(signKey, sourceStr))
		sign = strings.Replace(sign, " ", "+", -1)
		sign2 := base64.StdEncoding.EncodeToString(hashSign)
		if sign != sign2 {
			logrus.WithFields(logrus.Fields{
				"sign":      sign,
				"sign2":     sign2,
				"sourceStr": sourceStr,
			}).Info("签名不正确")
			c.AbortWithStatusJSON(http.StatusOK, serialize.Response(500, "签名不正确", nil))
			return
		}
	}
}

func CheckLogin(c *gin.Context) {
	path := c.Request.URL.Path
	for _, m := range def.LogWhiteMethod {
		if strings.Contains(path, m) {
			return
		}
	}
	token := c.Request.FormValue("sid")
	// 统一使用设备号登陆
	deviceId := helper.JwtDecode(token)
	if deviceId == "" {
		logrus.Info(c.Request.RequestURI, c.Request.Form)
		c.JSON(http.StatusOK, serialize.Response(def.CodeUnAuth, "请先登录", nil))
		c.Abort()
		return
	}
	c.Set("deviceId", deviceId)
	// 加锁防并发
	key := def.StringUserLock + deviceId + "-" + path
	res, err := db.MainRedis.Do("set", key, 1, "ex", 5, "nx")
	if err != nil || res != "OK" {
		c.JSON(http.StatusOK, serialize.Response(def.CodeRequestTooFast, "请求太快了，请稍后再试~", nil))
		c.Abort()
		return
	}
	defer func() {
		db.MainRedis.Do("del", key)
	}()
	c.Next()
}
