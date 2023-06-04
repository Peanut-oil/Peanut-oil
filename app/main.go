package main

import (
	"github.com/gin-gonic/gin/app/pkg"
	"math/rand"
	"time"
)

func main() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 设置日志输出
	pkg.InitLogger("/root/serverlog/api")

}
