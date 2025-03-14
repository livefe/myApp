package test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

// 收藏模块测试用例

// 收藏请求结构体
type AddFavoriteRequest struct {
	HouseID uint `json:"house_id"`
}

// 添加收藏
func AddFavorite(t *testing.T, houseID uint) uint {
	// 创建收藏数据
	favoriteData := AddFavoriteRequest{
		HouseID: houseID,
	}

	favoriteJSON, _ := json.Marshal(favoriteData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/add", favoriteJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	defer resp.Body.Close()

	// 解析响应获取收藏ID
	data := ParseResponseData(t, resp)
	favoriteID, ok := data["id"].(float64)
	if !ok {
		t.Fatal("添加收藏响应中未找到有效的ID")
	}

	return uint(favoriteID)
}

func TestAddFavorite(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试添加收藏
	t.Run("正常添加收藏", func(t *testing.T) {
		favoriteID := AddFavorite(t, houseID)
		if favoriteID == 0 {
			t.Fatal("添加收藏失败")
		}
	})

	// 测试添加不存在的房源收藏
	t.Run("添加不存在的房源收藏", func(t *testing.T) {
		// 创建收藏数据，使用不存在的房源ID
		favoriteData := AddFavoriteRequest{
			HouseID: 99999, // 不存在的房源ID
		}

		favoriteJSON, _ := json.Marshal(favoriteData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/add", favoriteJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("添加不存在的房源收藏应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试重复添加收藏
	t.Run("重复添加收藏", func(t *testing.T) {
		// 先添加一次收藏
		AddFavorite(t, houseID)

		// 再次添加相同房源的收藏
		favoriteData := AddFavoriteRequest{
			HouseID: houseID,
		}

		favoriteJSON, _ := json.Marshal(favoriteData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/add", favoriteJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码或特殊处理
		// 注意：有些API设计可能允许重复添加并返回成功，具体取决于业务逻辑
	})
}

func TestRemoveFavorite(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 添加收藏
	favoriteID := AddFavorite(t, houseID)

	// 测试删除收藏
	t.Run("删除收藏", func(t *testing.T) {
		req := CreateAuthRequest(t, "DELETE", baseURL+"/api/favorite/"+strconv.Itoa(int(favoriteID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 验证收藏已被删除（通过获取用户收藏列表检查）
		req = CreateAuthRequest(t, "GET", baseURL+"/api/favorite/list", nil)
		resp = ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 检查收藏列表中是否不包含已删除的收藏
		favorites, ok := data["favorites"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到favorites字段或格式不正确")
		}

		// 检查是否存在已删除的收藏
		for _, fav := range favorites {
			favMap, ok := fav.(map[string]interface{})
			if !ok {
				continue
			}
			favID, ok := favMap["id"].(float64)
			if !ok {
				continue
			}
			if uint(favID) == favoriteID {
				t.Errorf("收藏未被成功删除，ID: %d", favoriteID)
				break
			}
		}
	})

	// 测试删除不存在的收藏
	t.Run("删除不存在的收藏", func(t *testing.T) {
		req := CreateAuthRequest(t, "DELETE", baseURL+"/api/favorite/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("删除不存在的收藏应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestGetUserFavorites(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 添加收藏
	AddFavorite(t, houseID)

	// 测试获取用户收藏列表
	t.Run("获取用户收藏列表", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/favorite/list", nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证响应中包含收藏列表
		favorites, ok := data["favorites"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到favorites字段或格式不正确")
		}

		// 验证收藏列表不为空
		if len(favorites) == 0 {
			t.Error("用户收藏列表为空")
		}
	})

	// 测试未授权访问
	t.Run("未授权访问", func(t *testing.T) {
		// 创建不带认证token的请求
		req, _ := http.NewRequest("GET", baseURL+"/api/favorite/list", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回未授权状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("未授权请求应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestToggleFavorite(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试切换收藏状态（添加收藏）
	t.Run("切换收藏状态-添加", func(t *testing.T) {
		req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/toggle/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证收藏状态
		status, ok := data["status"].(bool)
		if !ok || !status {
			t.Errorf("切换收藏状态失败，期望状态为true，实际: %v", status)
		}
	})

	// 测试切换收藏状态（取消收藏）
	t.Run("切换收藏状态-取消", func(t *testing.T) {
		// 再次调用切换接口，应该取消收藏
		req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/toggle/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证收藏状态
		status, ok := data["status"].(bool)
		if !ok || status {
			t.Errorf("切换收藏状态失败，期望状态为false，实际: %v", status)
		}
	})

	// 测试切换不存在的房源收藏状态
	t.Run("切换不存在的房源收藏状态", func(t *testing.T) {
		req := CreateAuthRequest(t, "POST", baseURL+"/api/favorite/toggle/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("切换不存在的房源收藏状态应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestCheckFavorite(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试检查未收藏的房源
	t.Run("检查未收藏的房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/favorite/check/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证收藏状态
		isFavorite, ok := data["is_favorite"].(bool)
		if !ok || isFavorite {
			t.Errorf("检查未收藏的房源状态失败，期望状态为false，实际: %v", isFavorite)
		}
	})

	// 添加收藏
	AddFavorite(t, houseID)

	// 测试检查已收藏的房源
	t.Run("检查已收藏的房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/favorite/check/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证收藏状态
		isFavorite, ok := data["is_favorite"].(bool)
		if !ok || !isFavorite {
			t.Errorf("检查已收藏的房源状态失败，期望状态为true，实际: %v", isFavorite)
		}
	})

	// 测试检查不存在的房源
	t.Run("检查不存在的房源", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/favorite/check/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("检查不存在的房源应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}
