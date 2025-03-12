package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "操作成功",
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string, data ...interface{}) {
	var responseData interface{} = nil
	if len(data) > 0 {
		responseData = data[0]
	}
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string, data ...interface{}) {
	if message == "" {
		message = "无效的请求参数"
	}
	Fail(c, http.StatusBadRequest, message, data...)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string, data ...interface{}) {
	if message == "" {
		message = "未授权访问"
	}
	Fail(c, http.StatusUnauthorized, message, data...)
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string, data ...interface{}) {
	if message == "" {
		message = "资源不存在"
	}
	Fail(c, http.StatusNotFound, message, data...)
}

// ServerError 服务器内部错误
func ServerError(c *gin.Context, message string, data ...interface{}) {
	if message == "" {
		message = "服务器内部错误"
	}
	Fail(c, http.StatusInternalServerError, message, data...)
}

// SuccessWithToken 成功响应并返回token
func SuccessWithToken(c *gin.Context, token string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "操作成功",
		"token":   token,
		"data":    data,
	})
}

// Forbidden 权限不足
func Forbidden(c *gin.Context, message string, data ...interface{}) {
	if message == "" {
		message = "权限不足"
	}
	Fail(c, http.StatusForbidden, message, data...)
}
