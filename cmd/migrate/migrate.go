package main

import (
	"fmt"
	"myApp/config"
	"myApp/model"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 获取数据库连接
	db := model.InitDB()

	// 执行数据库迁移
	fmt.Println("开始执行数据库迁移...")

	// 自动迁移所有模型
	err := db.AutoMigrate(
		&model.User{},
		&model.Community{},
		&model.CommunityMember{},
		&model.Order{},
		&model.Product{},
		&model.ProductCategory{},
	)

	if err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	fmt.Println("数据库迁移完成！")
}
