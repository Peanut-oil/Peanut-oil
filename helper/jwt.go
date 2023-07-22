package helper

import (
	"github.com/gin-gonic/gin/app/def"
	"github.com/robbert229/jwt"
	"github.com/sirupsen/logrus"
	"time"
)

func GetJwtToken(uid int, openId string, ip string) (string, error) {
	claims := jwt.NewClaim()
	now := Millisecond() / 1000
	claims.Set("iat", now)
	claims.Set("uid", uid)
	claims.Set("iss", openId)
	claims.Set("ip", ip)
	key := def.JwtEncryptKey
	claims.Set("exp", now+def.UserInfoExpire)
	algorithm := jwt.HmacSha256(key)
	return algorithm.Encode(claims)
}

func GetJwtTokenV2(info map[string]interface{}, encKey string) (string, error) {
	claims := jwt.NewClaim()
	claims.Set("iat", time.Now().Unix())
	for key, val := range info {
		claims.Set(key, val)
	}
	algorithm := jwt.HmacSha256(encKey)
	return algorithm.Encode(claims)
}

func JwtDecode(sid string) string {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warn("panic getUidFromSid")
		}
	}()
	key := def.JwtEncryptKey
	algorithm := jwt.HmacSha256(key)
	claims, err := algorithm.Decode(sid)
	if err != nil {
		logrus.WithField("error", err).Info("JwtDecode err")
		return ""
	}
	if algorithm.Validate(sid) != nil {
		logrus.WithField("error", err).Info("Validate err")
		return ""
	}
	deviceId, err := claims.Get("deviceId")
	if err != nil {
		logrus.WithField("error", err).Warn("no uid")
		return ""
	}
	return deviceId.(string)
}

type JwtDecodeV3Rsp struct {
	Uid    int    `json:"uid"`
	UserId string `json:"user_id"`
	Rid    string `json:"rid"`
}
