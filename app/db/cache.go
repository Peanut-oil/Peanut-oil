package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"
	"time"
)

var MC cache.Cache

func InitMemoryCache() {
	MC = *cache.New(time.Second, time.Minute)
}

func Set(key, value string, exp time.Duration) {
	MC.Set(key, value, exp)
}

func Get(key string) (string, bool) {
	v, ok := MC.Get(key)
	var str string
	if ok {
		str, _ = v.(string)
	}
	return str, ok
}
