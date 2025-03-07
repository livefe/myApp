package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
	"math/rand"
	"strconv"
)

const baseURL = "http://localhost:8080"
var authToken string

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateCommunityRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
}

type CreateOrderRequest struct {
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

func TestAPIEndpoints(t *testing.T) {
	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 生成随机用户名，避免重复
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(10000)
	username := "testuser" + strconv.Itoa(randomID)

	// 1. 测试用户注册
	registerData := RegisterRequest{
		Username: username,
		Password: "password123",
		Email:    "test" + strconv.Itoa(randomID) + "@example.com",
		Phone:    "1234567890",
	}
	registerJSON, _ := json.Marshal(registerData)
	resp, err := http.Post(baseURL+"/api/user/register", "application/json", bytes.NewBuffer(registerJSON))
	if err != nil {
		t.Fatalf("注册请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("注册失败，状态码: %d", resp.StatusCode)
	}

	// 2. 测试用户登录
	loginData := LoginRequest{
		Username: username,
		Password: "password123",
	}
	loginJSON, _ := json.Marshal(loginData)
	resp, err = http.Post(baseURL+"/api/user/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			t.Errorf("登录失败，状态码: %d, 错误信息: %v", resp.StatusCode, errorResp["error"])
		} else {
			t.Errorf("登录失败，状态码: %d", resp.StatusCode)
		}
		return
	}

	// 从响应中获取token
	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok {
		t.Fatal("登录响应中未找到有效的token")
	}
	authToken = token

	// 3. 测试获取用户信息
	req, _ := http.NewRequest("GET", baseURL+"/api/user/info", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("获取用户信息请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("获取用户信息失败，状态码: %d", resp.StatusCode)
	}

	// 4. 测试创建社区
	communityData := CreateCommunityRequest{
		Name:        "测试社区",
		Description: "这是一个测试社区",
		LogoURL:     "https://example.com/logo.png",
	}
	communityJSON, _ := json.Marshal(communityData)
	req, _ = http.NewRequest("POST", baseURL+"/api/community/create", bytes.NewBuffer(communityJSON))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("创建社区请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("创建社区失败，状态码: %d", resp.StatusCode)
	}

	// 5. 测试获取社区列表
	req, _ = http.NewRequest("GET", baseURL+"/api/community/list", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("获取社区列表请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("获取社区列表失败，状态码: %d", resp.StatusCode)
	}

	// 6. 测试创建订单
	orderData := CreateOrderRequest{
		ProductID:  1,
		Quantity:   2,
		TotalPrice: 199.99,
	}
	orderJSON, _ := json.Marshal(orderData)
	req, _ = http.NewRequest("POST", baseURL+"/api/order/create", bytes.NewBuffer(orderJSON))
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("创建订单请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("创建订单失败，状态码: %d", resp.StatusCode)
	}

	// 7. 测试获取订单信息
	req, _ = http.NewRequest("GET", baseURL+"/api/order/1", nil)
	req.Header.Set("Authorization", "Bearer "+authToken)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("获取订单信息请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("获取订单信息失败，状态码: %d", resp.StatusCode)
	}

	fmt.Println("所有API接口测试完成！")
}