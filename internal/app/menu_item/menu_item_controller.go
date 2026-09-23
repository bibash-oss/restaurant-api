package menuitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuItemController struct {
	itemService *MenuItemService
}

func NewMenuItemController(itemService *MenuItemService) *MenuItemController {
	return &MenuItemController{itemService: itemService}
}

func (c *MenuItemController) CreateMenuItem(ctx *gin.Context) {
	var req CreateMenuItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.itemService.CreateMenuItem(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Menu item created successfully",
		"data":    result,
		"success": true,
	})
}

func (c *MenuItemController) GetMenuItemsByRestaurant(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID == "" {
		restaurantID = ctx.Param("restaurantId")
	}
	items, err := c.itemService.GetMenuItemsByRestaurantID(restaurantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    items,
		"success": true,
		"message": "Menu items fetched successfully",
	})
}

func (c *MenuItemController) GetMenuItemsByCategory(ctx *gin.Context) {
	categoryID := ctx.Param("id")
	if categoryID == "" {
		categoryID = ctx.Param("categoryId")
	}
	items, err := c.itemService.GetMenuItemsByCategoryID(categoryID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    items,
		"success": true,
		"message": "Menu items fetched successfully",
	})
}

func (c *MenuItemController) GetMenuItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := c.itemService.GetMenuItemByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    item,
		"success": true,
		"message": "Menu item fetched successfully",
	})
}

func (c *MenuItemController) UpdateMenuItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateMenuItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.itemService.UpdateMenuItem(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Menu item updated successfully",
		"data":    result,
		"success": true,
	})
}

func (c *MenuItemController) DeleteMenuItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.itemService.DeleteMenuItem(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Menu item deleted successfully",
		"success": true,
	})
}
