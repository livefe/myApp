package middleware

import (
	"myApp/config"
	"myApp/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "无效的访问令牌")
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		// 将userID转换为uint类型
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			response.ServerError(c, "无效的用户ID")
			c.Abort()
			return
		}
		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}
