package models

import (
	"database/sql"
	"fmt"
	"github.com/Asolmn/go-gin-example/pkg/logging"
	"github.com/Asolmn/go-gin-example/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {

	// 数据库连接信息
	var err error

	user := setting.DatabaseSetting.User
	password := setting.DatabaseSetting.Password
	host := setting.DatabaseSetting.Host
	dbName := setting.DatabaseSetting.Name
	tablePrefix := setting.DatabaseSetting.TablePrefix

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, dbName)

	// 创建日志记录器
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
		logging.Warn(err)
	}
	// 直接使用db.DB()方法获取*sql.DB对象
	sqlDB, err1 := db.DB()

	if err1 != nil {
		logging.Warn(err1)
	}
	sqlDB.SetMaxIdleConns(10)  // 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量

	// 执行回调函数
	err2 := db.Callback().Create().Replace("gorm:before_create", updateTimeStampForBeforeCreateCallback)
	if err2 != nil {
		logging.Warn(err2)
		return
	}
	err3 := db.Callback().Update().Replace("gorm:before_update", updateTimeStampForBeforeUpdateCallback)
	if err3 != nil {
		logging.Warn(err3)
		return
	}
	err4 := db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	if err4 != nil {
		logging.Warn(err4)
		return
	}

}

// 关闭连接
func CloseDB() {
	sqlDB, _ := db.DB()
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {

		}
	}(sqlDB)
}

// 更新创建时间
func updateTimeStampForBeforeCreateCallback(db *gorm.DB) {
	db.Statement.SetColumn("CreatedOn", time.Now().Unix())
}

// 更新修改时间
func updateTimeStampForBeforeUpdateCallback(db *gorm.DB) {
	db.Statement.SetColumn("ModifiedOn", time.Now().Unix())
}

// 软删除
func deleteCallback(db *gorm.DB) {
	if db.Error != nil {
		return
	}
	if db.Statement.Schema != nil {

		db.Statement.SQL.Grow(100) // 设置sql语句缓冲区大小

		deletedOnField := db.Statement.Schema.LookUpField("DeletedOn") // 获取DeletedOn字段
		if !db.Statement.Unscoped && deletedOnField != nil {           // 如果字段不为空且软删除为false
			if db.Statement.SQL.String() == "" {
				nowTime := time.Now().Unix()
				db.Statement.AddClause(
					clause.Set{
						{
							Column: clause.Column{Name: deletedOnField.DBName},
							Value:  nowTime,
						}, // 设置set字句
					}) // 添加子句
				db.Statement.AddClauseIfNotExists(clause.Update{})
				db.Statement.Build("UPDATE", "SET", "WHERE") // 构建sql语句
			}
		} else {
			if db.Statement.SQL.String() == "" {
				db.Statement.AddClauseIfNotExists(clause.Delete{})
				db.Statement.AddClauseIfNotExists(clause.From{})
				db.Statement.Build("DELETE", "FROM", "WHERE")
			}
		}
		// 检测有没有where子句
		if _, ok := db.Statement.Clauses["WHERE"]; !db.AllowGlobalUpdate && !ok {
			err := db.AddError(gorm.ErrMissingWhereClause)
			if err != nil {
				return
			}
			return
		}

		// 执行sql
		if !db.DryRun && db.Error == nil {
			result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context,
				db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {

				err1 := db.AddError(err)
				if err1 != nil {
					return
				}
			}
			db.RowsAffected, _ = result.RowsAffected()

		}
	}
}
