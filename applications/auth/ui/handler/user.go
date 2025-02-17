package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/goda6565/ptf-backends/applications/auth/domain/models"
	"github.com/goda6565/ptf-backends/applications/auth/pkg/logger"
	"github.com/goda6565/ptf-backends/applications/auth/service"
	"github.com/goda6565/ptf-backends/applications/auth/ui/gen"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func userToRegisterResponse(user *models.User) *api.RegisterResponse {
	return &api.RegisterResponse{
		Id:    int(user.ID),
		Email: user.Email,
	}
}

func userToGetUserResponse(user *models.User) *api.UserResponse {
	return &api.UserResponse{
		Id:        int(user.ID),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (h *UserHandler) UserRegister(c *gin.Context) {
	var req api.UserRegisterRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UserRegister(req.Email, req.Password)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusCreated, userToRegisterResponse(user))
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var req api.UserLoginRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.UserLogin(req.Email, req.Password)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, token)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: "token is required", Code: http.StatusBadRequest})
		return
	}

	// Bearer プレフィックスを取り除く
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: "Invalid token format", Code: http.StatusBadRequest})
		return
	}

	user, err := h.userService.GetUser(token)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, userToGetUserResponse(user))
}
