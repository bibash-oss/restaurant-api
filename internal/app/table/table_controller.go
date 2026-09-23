package table

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TableController struct {
	tableService *TableService
}

func NewTableController(tableService *TableService) *TableController {
	return &TableController{tableService: tableService}
}

func (c *TableController) CreateTable(ctx *gin.Context) {
	var req CreateTableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.tableService.CreateTable(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Table created successfully",
		"data":    result,
		"success": true,
	})
}

func (c *TableController) GetTablesByRestaurant(ctx *gin.Context) {
	restaurantID := ctx.Param("id")
	if restaurantID == "" {
		restaurantID = ctx.Param("restaurantId")
	}
	tables, err := c.tableService.GetTablesByRestaurantID(restaurantID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    tables,
		"success": true,
		"message": "Tables fetched successfully",
	})
}

func (c *TableController) GetTableByID(ctx *gin.Context) {
	id := ctx.Param("id")
	t, err := c.tableService.GetTableByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    t,
		"success": true,
		"message": "Table fetched successfully",
	})
}

func (c *TableController) UpdateTable(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateTableRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.tableService.UpdateTable(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Table updated successfully",
		"data":    result,
		"success": true,
	})
}

func (c *TableController) DeleteTable(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.tableService.DeleteTable(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Table deleted successfully",
		"success": true,
	})
}
