package menucategory

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuCategoryController struct {
	categoryService *MenuCategoryService
}

func NewMenuCategoryController(categoryService *MenuCategoryService) *MenuCategoryController {
	return &MenuCategoryController{categoryService: categoryService}
}

func (c *MenuCategoryController) CreateCategory(ctx *gin.Context) {
	var req CreateMenuCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.categoryService.CreateCategory(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Menu category created successfully",
		"data":    result,
		"success": true,
	})
}

func (c *MenuCategoryController) GetCategoriesByRestaurant(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID == "" {
		restaurantID = ctx.Param("restaurantId")
	}
	categories, err := c.categoryService.GetCategoriesByRestaurantID(restaurantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    categories,
		"success": true,
		"message": "Menu categories fetched successfully",
	})
}

func (c *MenuCategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	cat, err := c.categoryService.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    cat,
		"success": true,
		"message": "Menu category fetched successfully",
	})
}

func (c *MenuCategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateMenuCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.categoryService.UpdateCategory(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Menu category updated successfully",
		"data":    result,
		"success": true,
	})
}

func (c *MenuCategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.categoryService.DeleteCategory(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Menu category deleted successfully",
		"success": true,
	})
}
