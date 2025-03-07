package model

import (
	"fmt"
	"myApp/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// BaseModel 定义了所有模型共享的基础字段
// 这个结构体可以被其他模型嵌入，以提供统一的ID、时间戳和软删除功能
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Database.User,
		config.Conf.Database.Password,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.DBName)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	return db
}

// GetDB 返回已初始化的数据库连接实例
func GetDB() *gorm.DB {
	if db == nil {
		db = InitDB()
	}
	return db
}
