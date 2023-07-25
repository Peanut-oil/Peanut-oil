// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/pkg/crypto"
	"github.com/robbert229/jwt"
	"sort"
	"testing"
)

func TestLoadSign(t *testing.T) {
	sid := generateSid()
	signKey := def.SignKey
	method := "POST"
	path := "/genus/user/login"
	params := map[string][]string{
		"sid": {sid},
	}
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
	sign2 := base64.StdEncoding.EncodeToString(hashSign)
	fmt.Println(sign2)
}

func generateSid() string {
	key := def.JwtEncryptKey
	algorithm := jwt.HmacSha256(key)
	claims := jwt.NewClaim()
	claims.Set("deviceId", "qzgjjjjjj")
	encode, _ := algorithm.Encode(claims)
	fmt.Println("sid:", encode)
	return encode
}
