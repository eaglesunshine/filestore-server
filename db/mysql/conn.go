package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	//踩坑记：这里不要写成"db, _ := xxx"，否则db将作为函数内局部变量看待，而外面的全局变量db并未得到初始化
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:13306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000) //SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制。
	//SetMaxIdleConns用于设置闲置的连接数。
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connected to mysql, err:%s" + err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}
