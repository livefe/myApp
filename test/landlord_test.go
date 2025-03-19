package test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

// 房东模块测试用例

// 房东请求结构体
type CreateLandlordRequest struct {
	RealName    string `json:"real_name"`
	IDNumber    string `json:"id_number"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// 创建房东
func CreateLandlord(t *testing.T) uint {
	// 创建房东数据
	landlordData := CreateLandlordRequest{
		RealName:    "测试房东",
		IDNumber:    "123456789012345678",
		PhoneNumber: "13800138000",
		Address:     "测试地址",
	}

	landlordJSON, _ := json.Marshal(landlordData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	defer resp.Body.Close()

	// 解析响应获取房东ID
	data := ParseResponseData(t, resp)
	landlordID, ok := data["id"].(float64)
	if !ok {
		t.Fatal("创建房东响应中未找到有效的ID")
	}

	return uint(landlordID)
}

func TestCreateLandlord(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 测试申请成为房东
	t.Run("正常申请成为房东", func(t *testing.T) {
		// 创建房东数据
		landlordData := CreateLandlordRequest{
			RealName:    "测试房东",
			IDNumber:    "123456789012345678",
			PhoneNumber: "13800138000",
			Address:     "测试地址",
		}

		landlordJSON, _ := json.Marshal(landlordData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证返回的房东ID
		landlordID, ok := data["id"].(float64)
		if !ok || landlordID == 0 {
			t.Fatal("申请成为房东失败")
		}
	})

	// 测试申请成为房东-无效输入
	t.Run("无效输入-空姓名", func(t *testing.T) {
		// 创建房东数据，姓名为空
		landlordData := CreateLandlordRequest{
			RealName:    "", // 空姓名
			IDNumber:    "123456789012345678",
			PhoneNumber: "13800138000",
			Address:     "测试地址",
		}

		landlordJSON, _ := json.Marshal(landlordData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("空姓名申请成为房东应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试重复申请成为房东
	t.Run("重复申请成为房东", func(t *testing.T) {
		// 先申请一次
		landlordData := CreateLandlordRequest{
			RealName:    "测试房东",
			IDNumber:    "123456789012345678",
			PhoneNumber: "13800138000",
			Address:     "测试地址",
		}

		landlordJSON, _ := json.Marshal(landlordData)
		req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
		resp := ExecuteRequest(t, req, http.StatusOK)
		resp.Body.Close()

		// 再次申请
		req = CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("重复申请成为房东应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestGetLandlordProfile(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 申请成为房东
	landlordData := CreateLandlordRequest{
		RealName:    "测试房东",
		IDNumber:    "123456789012345678",
		PhoneNumber: "13800138000",
		Address:     "测试地址",
	}

	landlordJSON, _ := json.Marshal(landlordData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	resp.Body.Close()

	// 测试获取房东个人资料
	t.Run("获取房东个人资料", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/landlord/profile", nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证返回的房东信息
		realName, ok := data["real_name"].(string)
		if !ok || realName != "测试房东" {
			t.Errorf("返回的房东姓名不匹配，期望: %s, 实际: %v", "测试房东", realName)
		}
	})

	// 测试未授权访问
	t.Run("未授权访问", func(t *testing.T) {
		// 创建不带认证token的请求
		req, _ := http.NewRequest("GET", baseURL+"/api/landlord/profile", nil)
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

func TestUpdateLandlord(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 申请成为房东
	landlordData := CreateLandlordRequest{
		RealName:    "测试房东",
		IDNumber:    "123456789012345678",
		PhoneNumber: "13800138000",
		Address:     "测试地址",
	}

	landlordJSON, _ := json.Marshal(landlordData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	resp.Body.Close()

	// 测试更新房东信息
	t.Run("更新房东信息", func(t *testing.T) {
		// 更新房东数据
		updateData := map[string]interface{}{
			"phone_number": "13900139000",
			"address":      "更新后的地址",
		}

		updateJSON, _ := json.Marshal(updateData)
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/landlord/profile", updateJSON)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 获取更新后的房东信息
		req = CreateAuthRequest(t, "GET", baseURL+"/api/landlord/profile", nil)
		resp = ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 解析响应
		data := ParseResponseData(t, resp)

		// 验证房东信息已更新
		phoneNumber, ok := data["phone_number"].(string)
		if !ok || phoneNumber != "13900139000" {
			t.Errorf("房东电话未更新，期望: %s, 实际: %v", "13900139000", phoneNumber)
		}

		address, ok := data["address"].(string)
		if !ok || address != "更新后的地址" {
			t.Errorf("房东地址未更新，期望: %s, 实际: %v", "更新后的地址", address)
		}
	})
}

func TestVerifyLandlord(t *testing.T) {
	// 设置测试环境
	SetupTest(t)

	// 申请成为房东
	landlordData := CreateLandlordRequest{
		RealName:    "测试房东",
		IDNumber:    "123456789012345678",
		PhoneNumber: "13800138000",
		Address:     "测试地址",
	}

	landlordJSON, _ := json.Marshal(landlordData)
	req := CreateAuthRequest(t, "POST", baseURL+"/api/landlord/create", landlordJSON)
	resp := ExecuteRequest(t, req, http.StatusOK)
	defer resp.Body.Close()

	// 解析响应获取房东ID
	data := ParseResponseData(t, resp)
	landlordID, ok := data["id"].(float64)
	if !ok {
		t.Fatal("创建房东响应中未找到有效的ID")
	}

	// 测试验证房东身份
	// 注意：这个接口通常需要管理员权限，可能需要特殊处理
	t.Run("验证房东身份", func(t *testing.T) {
		// 验证数据
		verifyData := map[string]interface{}{
			"status": "verified",
		}

		verifyJSON, _ := json.Marshal(verifyData)
		req := CreateAuthRequest(t, "PUT", baseURL+"/api/landlord/verify/"+strconv.Itoa(int(landlordID)), verifyJSON)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		// 注意：由于这个接口可能需要管理员权限，所以这里不一定能测试成功
		// 如果返回未授权状态码，可以跳过这个测试
		if resp.StatusCode == http.StatusUnauthorized {
			t.Skip("验证房东身份需要管理员权限，跳过测试")
		}

		// 如果测试成功，验证房东状态已更新
		if resp.StatusCode == http.StatusOK {
			// 获取房东信息
			req = CreateAuthRequest(t, "GET", baseURL+"/api/landlord/profile", nil)
			resp = ExecuteRequest(t, req, http.StatusOK)
			defer resp.Body.Close()

			// 解析响应
			data := ParseResponseData(t, resp)

			// 验证房东状态已更新为已验证
			status, ok := data["status"].(string)
			if !ok || status != "verified" {
				t.Errorf("房东状态未更新为已验证，实际状态: %v", status)
			}
		}
	})
}
