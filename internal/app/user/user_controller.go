package user

import "github.com/gin-gonic/gin"

type UserController struct {
	userServce *UserService
}

func NewUserController(userServce *UserService) *UserController {
	return &UserController{userServce: userServce}
}

func (controller *UserController) CreateUser(c *gin.Context) {
	var userReq CreateUserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := controller.userServce.CreateUser(&userReq)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{
		"message": "User created successfully",
		"data":    createdUser,
		"success": true,
	})
}

func (controller *UserController) GetAllUsers(c *gin.Context) {
	users, err := controller.userServce.GetAllUsers()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": users, "success": true, "message": "Users Fetched Successfully"})
}
