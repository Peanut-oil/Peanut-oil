package db

import (
	"github.com/gin-gonic/gin/app/def"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MysqlDB *sqlx.DB

func ConnectDB() {
	MysqlDB = sqlx.MustConnect("mysql", def.MysqlAddr)
	MysqlDB.SetMaxOpenConns(100)
}
