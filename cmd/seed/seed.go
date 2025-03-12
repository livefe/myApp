package main

import (
	"fmt"
	"math/rand"
	"myApp/config"
	"myApp/model"
	"time"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"
)

const (
	UserCount        = 50
	LandlordCount    = 10 // 房东数量
	HousePerLandlord = 3  // 每个房东的房源数量
)

func main() {
	// 初始化配置和数据库连接
	config.InitConfig()
	db := model.InitDB()

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成测试数据
	fmt.Println("开始生成测试数据...")

	// 生成用户数据
	users := generateUsers(db)
	fmt.Printf("✓ 已生成 %d 个用户\n", len(users))

	// 生成房东数据
	landlords := generateLandlords(db, users)
	fmt.Printf("✓ 已生成 %d 个房东\n", len(landlords))

	fmt.Println("\n🎉 测试数据生成完成！")
}

func generateUsers(db *gorm.DB) []model.User {
	users := make([]model.User, UserCount)
	for i := 0; i < UserCount; i++ {
		now := time.Now()
		user := model.User{
			Username:  faker.Username(),
			Phone:     faker.Phonenumber(),
			Password:  "$2a$10$IZkFRxQr2oceXrF.Zl.pBeu4FnuWPFjxog7eaQoXeP5/WlEFONYe2", // 默认密码：123456
			Avatar:    fmt.Sprintf("https://i.pravatar.cc/150?img=%d", i+1),
			LastLogin: &now,
			RealName:  faker.Name(),
			IdCard:    fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d%d%d%d%d%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10)),
			Email:     faker.Email(),
		}
		db.Create(&user)
		users[i] = user
	}
	return users
}

// 生成房东数据
func generateLandlords(db *gorm.DB, users []model.User) []model.Landlord {
	landlords := make([]model.Landlord, LandlordCount)

	// 从用户中随机选择一些作为房东
	selectedUserIndices := make(map[int]bool)
	for i := 0; i < LandlordCount; i++ {
		// 随机选择一个未被选中的用户
		userIndex := rand.Intn(len(users))
		for selectedUserIndices[userIndex] {
			userIndex = rand.Intn(len(users))
		}
		selectedUserIndices[userIndex] = true

		// 创建房东信息
		landlord := model.Landlord{
			UserID:       users[userIndex].ID,
			Verified:     rand.Intn(2) == 1, // 随机设置认证状态
			IdCardFront:  fmt.Sprintf("https://example.com/id_card_front_%d.jpg", i+1),
			IdCardBack:   fmt.Sprintf("https://example.com/id_card_back_%d.jpg", i+1),
			BankAccount:  faker.CCNumber(),
			BankName:     "中国银行",
			AccountName:  users[userIndex].RealName,
			Introduction: faker.Paragraph(),
			Rating:       4.0 + rand.Float64(), // 4.0-5.0之间的随机评分
		}

		// 如果评分大于5，则设为5
		if landlord.Rating > 5.0 {
			landlord.Rating = 5.0
		}

		// 保存到数据库
		db.Create(&landlord)
		landlords[i] = landlord

		// 更新用户类型为房东
		db.Model(&users[userIndex]).Update("user_type", 1)
	}

	return landlords
}

// 这里可以添加生成房源数据的函数
