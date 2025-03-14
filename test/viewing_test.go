package test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// 预约看房模块测试用例

// 预约看房请求结构体
type CreateViewingRequest struct {
	HouseID  uint      `json:"house_id"`
	ViewDate time.Time `json:"view_date"`
	Message  string    `json:"message"`
}

// 创建预约看房
func CreateViewing(t *testing.T, houseID uint) uint {
	// 创建预约看房数据
	viewDate := time.Now().Add(24 * time.Hour) // 预约明天的时间
	viewingData := CreateViewingRequest{
		HouseID:  houseID,
		ViewDate: viewDate,
		Message:  "这是一个测试预约信息",
	}

	viewingJSON, _ := json.Marshal(viewingData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/viewing/create", viewingJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	defer resp.Body.Close()

	// 解析响应获取预约ID
	data := ParseResponseData(t, resp)
	viewingID, ok := data["id"].(float64)
	if !ok {
		t.Fatal("创建预约看房响应中未找到有效的ID")
	}

	return uint(viewingID)
}

func TestCreateViewing(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 测试创建预约看房
	t.Run("正常创建预约看房", func(t *testing.T) {
		viewingID := CreateViewing(t, houseID)
		if viewingID == 0 {
			t.Fatal("创建预约看房失败")
		}
	})

	// 测试创建预约-无效输入
	t.Run("无效输入-过去的日期", func(t *testing.T) {
		// 创建预约看房数据，使用过去的日期
		pastDate := time.Now().Add(-24 * time.Hour) // 昨天的日期
		viewingData := CreateViewingRequest{
			HouseID:  houseID,
			ViewDate: pastDate,
			Message:  "这是一个测试预约信息",
		}

		viewingJSON, _ := json.Marshal(viewingData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/viewing/create", viewingJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("使用过去的日期创建预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试创建不存在的房源预约
	t.Run("创建不存在的房源预约", func(t *testing.T) {
		// 创建预约看房数据，使用不存在的房源ID
		viewDate := time.Now().Add(24 * time.Hour) // 预约明天的时间
		viewingData := CreateViewingRequest{
			HouseID:  99999, // 不存在的房源ID
			ViewDate: viewDate,
			Message:  "这是一个测试预约信息",
		}

		viewingJSON, _ := json.Marshal(viewingData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/viewing/create", viewingJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("创建不存在的房源预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestGetViewing(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	viewingID := CreateViewing(t, houseID)

	// 测试获取预约详情
	t.Run("获取预约详情", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/viewing/"+strconv.Itoa(int(viewingID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证返回的预约ID
		returnedID, ok := data["id"].(float64)
		if !ok || uint(returnedID) != viewingID {
			t.Errorf("返回的预约ID不匹配，期望: %d, 实际: %v", viewingID, returnedID)
		}
	})

	// 测试获取不存在的预约
	t.Run("获取不存在的预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/viewing/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("获取不存在的预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestGetUserViewings(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	CreateViewing(t, houseID)

	// 测试获取用户预约列表
	t.Run("获取用户预约列表", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/viewing/user", nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证响应中包含预约列表
		viewings, ok := data["viewings"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到viewings字段或格式不正确")
		}

		// 验证预约列表不为空
		if len(viewings) == 0 {
			t.Error("用户预约列表为空")
		}
	})
}

func TestGetHouseViewings(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	CreateViewing(t, houseID)

	// 测试获取房源预约列表
	t.Run("获取房源预约列表", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/viewing/house/"+strconv.Itoa(int(houseID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证响应中包含预约列表
		viewings, ok := data["viewings"].([]interface{})
		if !ok {
			t.Fatal("响应格式错误，未找到viewings字段或格式不正确")
		}

		// 验证预约列表不为空
		if len(viewings) == 0 {
			t.Error("房源预约列表为空")
		}
	})

	// 测试获取不存在的房源预约列表
	t.Run("获取不存在的房源预约列表", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/viewing/house/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("获取不存在的房源预约列表应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestConfirmViewing(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	viewingID := CreateViewing(t, houseID)

	// 测试确认预约
	t.Run("确认预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/confirm/"+strconv.Itoa(int(viewingID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 获取预约详情，验证状态已更新
		req = CreateAuthRequest(t, "GET", baseURL+"/api/viewing/"+strconv.Itoa(int(viewingID)), nil)
		resp = ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证预约状态已更新为已确认
		status, ok := data["status"].(string)
		if !ok || status != "confirmed" {
			t.Errorf("预约状态未更新为已确认，实际状态: %v", status)
		}
	})

	// 测试确认不存在的预约
	t.Run("确认不存在的预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/confirm/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("确认不存在的预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestCompleteViewing(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	viewingID := CreateViewing(t, houseID)

	// 先确认预约
	req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/confirm/"+strconv.Itoa(int(viewingID)), nil)
	ExecuteRequest(t, req, http.StatusOK)

	// 测试完成预约
	t.Run("完成预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/complete/"+strconv.Itoa(int(viewingID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 获取预约详情，验证状态已更新
		req = CreateAuthRequest(t, "GET", baseURL+"/api/viewing/"+strconv.Itoa(int(viewingID)), nil)
		resp = ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证预约状态已更新为已完成
		status, ok := data["status"].(string)
		if !ok || status != "completed" {
			t.Errorf("预约状态未更新为已完成，实际状态: %v", status)
		}
	})

	// 测试完成不存在的预约
	t.Run("完成不存在的预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/complete/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("完成不存在的预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestCancelViewing(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 创建一个测试房源
	houseID := CreateHouse(t)

	// 创建一个测试预约
	viewingID := CreateViewing(t, houseID)

	// 测试取消预约
	t.Run("取消预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/cancel/"+strconv.Itoa(int(viewingID)), nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 获取预约详情，验证状态已更新
		req = CreateAuthRequest(t, "GET", baseURL+"/api/viewing/"+strconv.Itoa(int(viewingID)), nil)
		resp = ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证预约状态已更新为已取消
		status, ok := data["status"].(string)
		if !ok || status != "cancelled" {
			t.Errorf("预约状态未更新为已取消，实际状态: %v", status)
		}
	})

	// 测试取消不存在的预约
	t.Run("取消不存在的预约", func(t *testing.T) {
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/viewing/cancel/99999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("取消不存在的预约应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}
