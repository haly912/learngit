package common

import (
	"blog_server/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	user := "root"
	password := "123456"
	host := "localhost"
	port := "3306"
	database := "blog"
	charset := "utf8"
	loc := "Asia/Shanghai"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		user,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	// 连接数据库
	db, err := gorm.Open(driverName, args)
	//这个open函数是自带的 我们把我们要用的驱动名称 和数据库的详细东西的东西传进去了 然后就会连接数据库了
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	// 迁移数据表
	db.AutoMigrate(&model.User{})
	//自动建表
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
