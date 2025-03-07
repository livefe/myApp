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
	UserCount         = 50
	CommunityCount    = 10
	MaxMembersPerComm = 20
	OrdersPerUser     = 5
	MaxOrderAmount    = 1000.0
	MinOrderAmount    = 10.0
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

	// 生成社区数据
	communities := generateCommunities(db)
	fmt.Printf("✓ 已生成 %d 个社区\n", len(communities))

	// 生成社区成员关系
	generateCommunityMembers(db, users, communities)
	fmt.Println("✓ 已生成社区成员关系")

	// 生成订单数据
	generateOrders(db, users)
	fmt.Println("✓ 已生成订单数据")

	fmt.Println("\n🎉 测试数据生成完成！")
}

func generateUsers(db *gorm.DB) []model.User {
	users := make([]model.User, UserCount)
	for i := 0; i < UserCount; i++ {
		now := time.Now()
		user := model.User{
			Username:  faker.Username(),
			Email:     faker.Email(),
			Phone:     faker.Phonenumber(),
			Password:  "$2a$10$IZkFRxQr2oceXrF.Zl.pBeu4FnuWPFjxog7eaQoXeP5/WlEFONYe2", // 默认密码：123456
			Avatar:    fmt.Sprintf("https://i.pravatar.cc/150?img=%d", i+1),
			LastLogin: &now,
		}
		db.Create(&user)
		users[i] = user
	}
	return users
}

func generateCommunities(db *gorm.DB) []model.Community {
	communities := make([]model.Community, CommunityCount)
	for i := 0; i < CommunityCount; i++ {
		community := model.Community{
			Name:         faker.Word() + " Community",
			Description:  faker.Sentence(),
			LogoURL:      fmt.Sprintf("https://picsum.photos/200/200?random=%d", i+1),
			Status:       1,
			MembersCount: 0,
			CreatorID:    uint(rand.Intn(UserCount) + 1),
		}
		db.Create(&community)
		communities[i] = community
	}
	return communities
}

func generateCommunityMembers(db *gorm.DB, users []model.User, communities []model.Community) {
	for _, community := range communities {
		// 随机选择成员数量
		memberCount := rand.Intn(MaxMembersPerComm) + 1
		addedMembers := make(map[uint]bool)

		// 使用事务确保数据一致性
		err := db.Transaction(func(tx *gorm.DB) error {
			actualCount := 0
			// 随机选择用户作为成员
			for i := 0; i < memberCount*2 && actualCount < memberCount; i++ {
				user := users[rand.Intn(len(users))]
				// 检查是否已经是成员
				if addedMembers[user.ID] {
					continue
				}

				member := model.CommunityMember{
					UserID:      user.ID,
					CommunityID: community.ID,
					Role:        0, // 普通成员
				}

				// 尝试创建成员关系
				if err := tx.Create(&member).Error; err == nil {
					addedMembers[user.ID] = true
					actualCount++
				}
			}

			// 更新社区成员数量
			return tx.Model(&community).Update("members_count", actualCount).Error
		})

		if err != nil {
			fmt.Printf("生成社区 %d 的成员时发生错误: %v\n", community.ID, err)
		}
	}
}

func generateOrders(db *gorm.DB, users []model.User) {
	orderStatus := []string{model.OrderPending, model.OrderPaid, model.OrderCompleted, model.OrderCancelled}

	// 检查产品表是否存在并获取有效的产品ID
	var productIDs []uint
	type Product struct {
		ID uint
	}

	// 尝试查询产品表
	var products []Product
	result := db.Table("products").Select("id").Find(&products)

	// 如果产品表存在且有数据，使用实际的产品ID
	if result.Error == nil && len(products) > 0 {
		for _, product := range products {
			productIDs = append(productIDs, product.ID)
		}
		fmt.Printf("✓ 已找到 %d 个产品\n", len(productIDs))
	} else {
		// 如果产品表不存在或没有数据，使用默认的产品ID
		productIDs = []uint{1, 2, 3, 4, 5}
		fmt.Println("⚠️ 未找到产品表或产品数据，使用默认产品ID")
	}

	for _, user := range users {
		// 为每个用户生成随机数量的订单
		for i := 0; i < OrdersPerUser; i++ {
			// 生成随机金额和数量
			quantity := rand.Intn(5) + 1
			pricePerItem := MinOrderAmount + rand.Float64()*(MaxOrderAmount-MinOrderAmount)
			totalPrice := pricePerItem * float64(quantity)

			order := model.Order{
				UserID:     user.ID,
				ProductID:  productIDs[rand.Intn(len(productIDs))],
				Quantity:   quantity,
				TotalPrice: totalPrice,
				Status:     orderStatus[rand.Intn(len(orderStatus))],
			}
			db.Create(&order)
		}
	}
}
