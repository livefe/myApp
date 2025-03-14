package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// 用户模块测试用例
func TestUserRegister(t *testing.T) {
	// 生成随机用户名
	username := GenerateRandomUsername()

	// 测试用户注册
	t.Run("正常注册", func(t *testing.T) {
		RegisterUser(t, username, "password123", "1234567890")
	})

	// 测试重复注册
	t.Run("重复注册", func(t *testing.T) {
		// 使用相同用户名再次注册
		registerData := RegisterRequest{
			Username: username,
			Password: "password123",
			Phone:    "1234567890",
		}
		registerJSON, _ := json.Marshal(registerData)
		resp, err := http.Post(baseURL+"/api/user/register", "application/json", bytes.NewBuffer(registerJSON))
		if err != nil {
			t.Fatalf("注册请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码，因为用户名已存在
		if resp.StatusCode == http.StatusOK {
			t.Errorf("重复注册应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试无效输入
	t.Run("无效输入-空用户名", func(t *testing.T) {
		registerData := RegisterRequest{
			Username: "", // 空用户名
			Password: "password123",
			Phone:    "1234567890",
		}
		registerJSON, _ := json.Marshal(registerData)
		resp, err := http.Post(baseURL+"/api/user/register", "application/json", bytes.NewBuffer(registerJSON))
		if err != nil {
			t.Fatalf("注册请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("空用户名注册应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestUserLogin(t *testing.T) {
	// 生成随机用户名并注册
	username := GenerateRandomUsername()
	RegisterUser(t, username, "password123", "1234567890")

	// 测试正常登录
	t.Run("正常登录", func(t *testing.T) {
		token := LoginUser(t, username, "password123")
		if token == "" {
			t.Fatal("登录成功但未获取到token")
		}
		// 保存token供其他测试使用
		authToken = token
	})

	// 测试错误密码
	t.Run("错误密码", func(t *testing.T) {
		loginData := LoginRequest{
			Username: username,
			Password: "wrong_password",
		}
		loginJSON, _ := json.Marshal(loginData)
		resp, err := http.Post(baseURL+"/api/user/login", "application/json", bytes.NewBuffer(loginJSON))
		if err != nil {
			t.Fatalf("登录请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("错误密码登录应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})

	// 测试不存在的用户
	t.Run("不存在的用户", func(t *testing.T) {
		loginData := LoginRequest{
			Username: "nonexistent_user",
			Password: "password123",
		}
		loginJSON, _ := json.Marshal(loginData)
		resp, err := http.Post(baseURL+"/api/user/login", "application/json", bytes.NewBuffer(loginJSON))
		if err != nil {
			t.Fatalf("登录请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回错误状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("不存在用户登录应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}

func TestUserInfo(t *testing.T) {
	// 设置测试环境
	username := SetupTest(t)

	// 测试获取用户信息
	t.Run("获取用户信息", func(t *testing.T) {
		req := CreateAuthRequest(t, "GET", baseURL+"/api/user/info", nil)
		resp := ExecuteRequest(t, req, http.StatusOK)
		defer resp.Body.Close()

		// 验证返回的用户信息
		data := ParseResponseData(t, resp)
		returnedUsername, ok := data["username"].(string)
		if !ok || returnedUsername != username {
			t.Errorf("返回的用户名不匹配，期望: %s, 实际: %v", username, returnedUsername)
		}
	})

	// 测试未授权访问
	t.Run("未授权访问", func(t *testing.T) {
		// 创建不带认证token的请求
		req, _ := http.NewRequest("GET", baseURL+"/api/user/info", nil)
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

	// 测试无效token
	t.Run("无效token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/user/info", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()
		// 期望返回未授权状态码
		if resp.StatusCode == http.StatusOK {
			t.Errorf("无效token请求应该失败，但返回状态码: %d", resp.StatusCode)
		}
	})
}
