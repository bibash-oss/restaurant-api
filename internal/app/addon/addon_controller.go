package addon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddonController struct {
	addonService *AddonService
}

func NewAddonController(addonService *AddonService) *AddonController {
	return &AddonController{addonService: addonService}
}

func (c *AddonController) CreateAddon(ctx *gin.Context) {
	var req CreateAddonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.addonService.CreateAddon(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Addon created successfully",
		"data":    result,
		"success": true,
	})
}

func (c *AddonController) GetAddonsByRestaurant(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID == "" {
		restaurantID = ctx.Param("restaurantId")
	}
	addons, err := c.addonService.GetAddonsByRestaurantID(restaurantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    addons,
		"success": true,
		"message": "Addons fetched successfully",
	})
}

func (c *AddonController) GetAddonByID(ctx *gin.Context) {
	id := ctx.Param("id")
	a, err := c.addonService.GetAddonByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    a,
		"success": true,
		"message": "Addon fetched successfully",
	})
}

func (c *AddonController) UpdateAddon(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateAddonRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.addonService.UpdateAddon(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Addon updated successfully",
		"data":    result,
		"success": true,
	})
}

func (c *AddonController) DeleteAddon(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.addonService.DeleteAddon(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Addon deleted successfully",
		"success": true,
	})
}

func (c *AddonController) AssignAddonsToMenuItem(ctx *gin.Context) {
	menuItemID := ctx.Param("id")
	var req AssignAddonsToMenuItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.addonService.AssignAddonsToMenuItem(menuItemID, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Addons mapped to menu item successfully",
		"success": true,
	})
}

func (c *AddonController) GetAddonsByMenuItem(ctx *gin.Context) {
	menuItemID := ctx.Param("id")
	addons, err := c.addonService.GetAddonsByMenuItemID(menuItemID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    addons,
		"success": true,
		"message": "Menu item addons fetched successfully",
	})
}
