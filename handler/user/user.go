package user

import (
	"myApp/config"
	"myApp/model"
	"myApp/pkg/response"
	"myApp/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	createdUser, err := h.service.Register(&user)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, createdUser)
}

func (h *Handler) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		response.BadRequest(c, "无效的请求参数")
		return
	}

	user, err := h.service.Login(credentials.Username, credentials.Password)
	if err != nil {
		response.Unauthorized(c, "无效的凭证")
		return
	}

	// 生成JWT令牌
	claims := jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Duration(config.Conf.JWT.Expire) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Conf.JWT.Secret))
	if err != nil {
		response.ServerError(c, "生成令牌失败")
		return
	}

	response.SuccessWithToken(c, tokenString, user)
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	user, err := h.service.GetUserProfile(userID.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}
