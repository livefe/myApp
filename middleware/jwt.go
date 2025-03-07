package middleware

import (
	"myApp/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 定义不需要验证的路由
		publicPaths := []string{
			"/api/user/register",
			"/api/user/login",
		}

		// 检查当前路径是否在公共路径列表中
		for _, path := range publicPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		// 将userID转换为uint类型
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID"})
			return
		}
		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}
