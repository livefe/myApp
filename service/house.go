package service

import (
	"encoding/json"
	"fmt"
	"myApp/model"
	"myApp/pkg/redis"
	"myApp/repository"
	"time"
)

type HouseService interface {
	CreateHouse(house *model.House) error
	GetHouseByID(id uint) (*model.House, error)
	GetAllHouses(params map[string]interface{}) ([]model.House, error)
	UpdateHouse(house *model.House) error
	DeleteHouse(id uint) error
	GetHousesByLandlordID(landlordID uint) ([]model.House, error)
	IncrementViewCount(id uint) error
}

type houseService struct {
	repo repository.HouseRepository
}

func NewHouseService(repo repository.HouseRepository) HouseService {
	return &houseService{repo: repo}
}

func (s *houseService) CreateHouse(house *model.House) error {
	return s.repo.Create(house)
}

func (s *houseService) GetHouseByID(id uint) (*model.House, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("house:%d", id)

	// 尝试从缓存获取
	cacheData, err := redis.Get(cacheKey)
	if err == nil {
		// 缓存命中，反序列化数据
		var house model.House
		if err := json.Unmarshal([]byte(cacheData), &house); err == nil {
			return &house, nil
		}
		// 如果是空JSON对象，表示空结果
		if cacheData == "{}" {
			return nil, nil
		}
	} else if err != redis.Nil {
		// 如果是其他错误，记录但不影响主流程
		fmt.Printf("Redis获取缓存错误: %v\n", err)
	}

	// 缓存未命中或反序列化失败，从数据库获取
	house, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 将数据存入缓存
	if house != nil {
		houseData, err := json.Marshal(house)
		if err == nil {
			// 设置缓存，过期时间30分钟
			_ = redis.Set(cacheKey, string(houseData), 30*time.Minute)
		}
	} else {
		// 缓存空结果，设置较短的过期时间（5分钟）
		_ = redis.Set(cacheKey, "{}", 5*time.Minute)
	}

	return house, nil
}

func (s *houseService) GetAllHouses(params map[string]interface{}) ([]model.House, error) {
	// 构造缓存键，基于查询参数
	var cacheKey string
	if params == nil || len(params) == 0 {
		cacheKey = "houses:list:all"
	} else {
		// 对于有参数的查询，生成唯一的缓存键
		paramsData, err := json.Marshal(params)
		if err == nil {
			cacheKey = fmt.Sprintf("houses:list:%x", paramsData)
		} else {
			// 如果无法序列化参数，使用默认键
			cacheKey = "houses:list:default"
		}
	}

	// 尝试从缓存获取
	cacheData, err := redis.Get(cacheKey)
	if err == nil {
		// 缓存命中，反序列化数据
		var houses []model.House
		if err := json.Unmarshal([]byte(cacheData), &houses); err == nil {
			return houses, nil
		}
		// 检查是否是空结果的标记
		if cacheData == "empty_list" {
			return []model.House{}, nil
		}
	} else if err != redis.Nil {
		// 如果是其他错误，记录但不影响主流程
		fmt.Printf("Redis获取缓存错误: %v\n", err)
	}

	// 缓存未命中或反序列化失败，从数据库获取
	houses, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// 将数据存入缓存
	if len(houses) > 0 {
		housesData, err := json.Marshal(houses)
		if err == nil {
			// 设置缓存，过期时间15分钟
			_ = redis.Set(cacheKey, string(housesData), 15*time.Minute)
		}
	} else {
		// 缓存空结果，设置较短的过期时间（5分钟）
		_ = redis.Set(cacheKey, "[]", 5*time.Minute)
	}

	return houses, nil
}

func (s *houseService) UpdateHouse(house *model.House) error {
	// 更新数据库
	err := s.repo.Update(house)
	if err != nil {
		return err
	}

	// 删除缓存，强制下次请求重新从数据库加载
	cacheKey := fmt.Sprintf("house:%d", house.ID)
	_ = redis.Delete(cacheKey)

	// 删除相关列表缓存
	_ = redis.DeleteByPattern("houses:list:*")
	_ = redis.DeleteByPattern(fmt.Sprintf("houses:landlord:%d", house.LandlordID))

	return nil
}

func (s *houseService) DeleteHouse(id uint) error {
	// 先获取房源信息，用于后续清除相关缓存
	house, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除数据库记录
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("house:%d", id)
	_ = redis.Delete(cacheKey)

	// 删除相关列表缓存
	_ = redis.DeleteByPattern("houses:list:*")
	if house != nil {
		_ = redis.DeleteByPattern(fmt.Sprintf("houses:landlord:%d", house.LandlordID))
	}

	return nil
}

func (s *houseService) GetHousesByLandlordID(landlordID uint) ([]model.House, error) {
	// 构造缓存键
	cacheKey := fmt.Sprintf("houses:landlord:%d", landlordID)

	// 尝试从缓存获取
	cacheData, err := redis.Get(cacheKey)
	if err == nil {
		// 缓存命中，反序列化数据
		var houses []model.House
		if err := json.Unmarshal([]byte(cacheData), &houses); err == nil {
			return houses, nil
		}
		// 检查是否是空结果的标记
		if cacheData == "empty_list" {
			return []model.House{}, nil
		}
	} else if err != redis.Nil {
		// 如果是其他错误，记录但不影响主流程
		fmt.Printf("Redis获取缓存错误: %v\n", err)
	}

	// 缓存未命中或反序列化失败，从数据库获取
	houses, err := s.repo.GetHousesByLandlordID(landlordID)
	if err != nil {
		return nil, err
	}

	// 将数据存入缓存
	if len(houses) > 0 {
		housesData, err := json.Marshal(houses)
		if err == nil {
			// 设置缓存，过期时间20分钟
			_ = redis.Set(cacheKey, string(housesData), 20*time.Minute)
		}
	} else {
		// 缓存空结果，设置较短的过期时间（5分钟）
		_ = redis.Set(cacheKey, "[]", 5*time.Minute)
	}

	return houses, nil
}

func (s *houseService) IncrementViewCount(id uint) error {
	// 更新数据库
	err := s.repo.IncrementViewCount(id)
	if err != nil {
		return err
	}

	// 删除缓存，强制下次请求重新从数据库加载
	cacheKey := fmt.Sprintf("house:%d", id)
	_ = redis.Delete(cacheKey)

	return nil
}
