package test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// 基础URL常量
const baseURL = "http://localhost:8080"

// 全局认证令牌
var authToken string

// 请求结构体定义
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
}

// 测试辅助函数

// 生成随机用户名
func GenerateRandomUsername() string {
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(10000)
	return "testuser" + strconv.Itoa(randomID)
}

// 注册新用户
func RegisterUser(t *testing.T, username, password, phone string) {
	registerData := RegisterRequest{
		Username: username,
		Password: password,
		Phone:    phone,
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
}

// 登录用户并获取token
func LoginUser(t *testing.T, username, password string) string {
	loginData := LoginRequest{
		Username: username,
		Password: password,
	}
	loginJSON, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/user/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("登录请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			t.Fatalf("登录失败，状态码: %d, 错误信息: %v", resp.StatusCode, errorResp["error"])
		} else {
			t.Fatalf("登录失败，状态码: %d", resp.StatusCode)
		}
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok {
		t.Fatal("登录响应中未找到有效的token")
	}
	return token
}

// 创建认证请求
func CreateAuthRequest(t *testing.T, method, url string, body []byte) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+authToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

// 执行HTTP请求并检查状态码
func ExecuteRequest(t *testing.T, req *http.Request, expectedStatus int) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("执行请求失败: %v", err)
	}
	if resp.StatusCode != expectedStatus {
		t.Errorf("请求失败，期望状态码: %d, 实际状态码: %d", expectedStatus, resp.StatusCode)
	}
	return resp
}

// 初始化测试环境
func SetupTest(t *testing.T) string {
	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 生成随机用户名
	username := GenerateRandomUsername()

	// 注册用户
	RegisterUser(t, username, "password123", "1234567890")

	// 登录用户并获取token
	authToken = LoginUser(t, username, "password123")

	return username
}

// 从响应中解析数据
func ParseResponseData(t *testing.T, resp *http.Response) map[string]interface{} {
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应格式错误，未找到data字段或格式不正确")
	}

	return data
}
