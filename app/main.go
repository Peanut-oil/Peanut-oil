package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/app/api"
	"github.com/gin-gonic/gin/app/db"
	"github.com/gin-gonic/gin/app/def"
	"github.com/gin-gonic/gin/app/pkg"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())
	// 设置日志输出
	pkg.InitLogger(def.ServerLog)
	// 初始化redis
	db.ConnectRedis()
	// 初始化mysql
	db.ConnectDB()
	// 初始化内存缓存
	db.InitMemoryCache()
	// 初始化路由
	port := def.ServerPort
	serverRun(port)
	log.Println("Genus服务启动，port:", port, ",PID:", strconv.Itoa(os.Getpid()))
}

func serverRun(port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	v1 := router.Group("/genus")
	api.AddUserRoutes(v1)
	router.Run(":" + port)
}
