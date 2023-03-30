package models

import (
	"database/sql"
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func init() {

	// 数据库连接信息
	var (
		err                                       error
		dbName, user, password, host, tablePrefix string
	)

	// 读取数据
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	//dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, dbName)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,         // 禁用彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix, // 前缀设定
			SingularTable: true,        // 使用单数表名
		},
		Logger: newLogger, // 日志设置
	}) // 连接数据库

	if err != nil {
		log.Println(err)
	}
	// 直接使用db.DB()方法获取*sql.DB对象
	sqlDB, err1 := db.DB()

	if err1 != nil {
		log.Println(err1)
	}
	sqlDB.SetMaxIdleConns(10)  // 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量

}

func CloseDB() {
	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {

		}
	}(sqlDB)
}
