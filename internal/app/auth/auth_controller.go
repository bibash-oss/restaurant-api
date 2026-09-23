package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (controller *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	res, err := controller.authService.Login(&req)
	if err != nil {
		c.JSON(401, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    res,
	})
}

func (controller *AuthController) Me(c *gin.Context) {
	val, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"error":   "Unauthorized",
		})
		return
	}

	userID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{
			"success": false,
			"error":   "Invalid user identity in context",
		})
		return
	}

	u, err := controller.authService.GetProfile(userID)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "User profile not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Profile fetched successfully",
		"data":    u,
	})
}
