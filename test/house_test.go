package test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// 房源模块测试用例

// 房源请求结构体
type CreateHouseRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Price       float64  `json:"price"`
	Area        float64  `json:"area"`
	Rooms       int      `json:"rooms"`
	Bathrooms   int      `json:"bathrooms"`
	ImageURLs   []string `json:"image_urls"`
}

// 生成随机房源标题
func GenerateRandomHouseTitle() string {
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(10000)
	return "测试房源" + strconv.Itoa(randomID)
}

// 创建房源
func CreateHouse(t *testing.T) uint {
	// 创建房源数据
	houseData := CreateHouseRequest{
		Title:       GenerateRandomHouseTitle(),
		Description: "这是一个测试房源描述",
		Address:     "测试地址",
		Price:       5000.0,
		Area:        100.0,
		Rooms:       3,
		Bathrooms:   2,
		ImageURLs:   []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
	}

	houseJSON, _ := json.Marshal(houseData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/house/create", houseJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	defer resp.Body.Close()

	// 解析响应获取房源ID
	data := ParseResponseData(t, resp)
	houseID, ok := data["id"].(float64)
	if !ok {
		t.Fatal("创建房源响应中未找到有效的ID")
	}

	return uint(houseID)
}

func TestHouseCreate(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 测试创建房源
	t.Run("正常创建房源", func(t *testing.T) {
		houseID := CreateHouse(t)
		if houseID == 0 {
			t.Fatal("创建房源失败")
		}
	})

	// 测试创建房源-无效输入
	t.Run("无效输入-空标题", func(t *testing.T) {
		// 创建房源数据，标题为空
		houseData := CreateHouseRequest{
			Title:       "", // 空标题
			Description: "这是一个测试房源描述",
			Address:     "测试地址",
			Price:       5000.0,
			Area:        100.0,
			Rooms:       3,
			Bathrooms:   2,
			ImageURLs:   []string{"http://example.com/image1.jpg"},
		}

		houseJSON, _ := json.Marshal(houseData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/house/create", houseJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("空标题创建房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestHouseList(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	CreateHouse(t)

	// 测试获取房源列表
	t.Run("获取房源列表", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/house/list", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 验证状态码
		if resp.StatusCode != http.StatusOK {
			t.Errorf("获取房源列表失败，状态码: %d", resp.StatusCode)
		}

		// 解析响应
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("解析响应失败: %v", err)
		}

		// 验证响应中包含房源列表
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到data字段或格式不正确")
		}

		houses, ok := data["houses"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到houses字段或格式不正确")
		}

		// 验证房源列表不为空
		if len(houses) == 0 {
			t.Error("房源列表为空")
		}
	})
}

func TestHouseDetail(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试获取房源详情
	t.Run("获取房源详情", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/house/"+strconv.Itoa(int(houseID)), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 验证状态码
		if resp.StatusCode != http.StatusOK {
			t.Errorf("获取房源详情失败，状态码: %d", resp.StatusCode)
		}

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证返回的房源ID
		returnedID, ok := data["id"].(float64)
		if !ok || uint(returnedID) != houseID {
			t.Errorf("返回的房源ID不匹配，期望: %d, 实际: %v", houseID, returnedID)
		}
	})

	// 测试获取不存在的房源
	t.Run("获取不存在的房源", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/house/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("获取不存在的房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestHouseUpdate(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试更新房源
	t.Run("更新房源", func(t *testing.T) {
		// 更新房源数据
		updateData := map[string]interface{}{
			"title":       "更新后的房源标题",
			"description": "更新后的房源描述",
			"price":       6000.0,
		}

		updateJSON, _ := json.Marshal(updateData)
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/house/"+strconv.Itoa(int(houseID)), updateJSON)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 获取更新后的房源详情
		req, _ = http.NewRequest("GET", baseURL+"/api/house/"+strconv.Itoa(int(houseID)), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 验证状态码
		if resp.StatusCode != http.StatusOK {
			t.Errorf("获取房源详情失败，状态码: %d", resp.StatusCode)
		}

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证房源信息已更新
		title, ok := data["title"].(string)
		if !ok || title != "更新后的房源标题" {
			t.Errorf("房源标题未更新，期望: %s, 实际: %v", "更新后的房源标题", title)
		}
	})

	// 测试更新不存在的房源
	t.Run("更新不存在的房源", func(t *testing.T) {
		updateData := map[string]interface{}{
			"title": "更新后的房源标题",
		}

		updateJSON, _ := json.Marshal(updateData)
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/house/99999", updateJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("更新不存在的房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestHouseDelete(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试删除房源
	t.Run("删除房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "DELETE", baseURL+"/api/house/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 尝试获取已删除的房源
		req, _ = http.NewRequest("GET", baseURL+"/api/house/"+strconv.Itoa(int(houseID)), nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("获取已删除的房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试删除不存在的房源
	t.Run("删除不存在的房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "DELETE", baseURL+"/api/house/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("删除不存在的房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestLandlordHouses(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	CreateHouse(t)

	// 测试获取房东的所有房源
	t.Run("获取房东的所有房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/house/landlord", nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证响应中包含房源列表
		houses, ok := data["houses"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到houses字段或格式不正确")
		}

		// 验证房源列表不为空
		if len(houses) == 0 {
			t.Error("房东的房源列表为空")
		}
	})
}
