package server

import (
	"net/http"

	"kitchen-api/internal/app/addon"
	"kitchen-api/internal/app/auth"
	menucategory "kitchen-api/internal/app/menu_category"
	menuitem "kitchen-api/internal/app/menu_item"
	"kitchen-api/internal/app/order"
	orderitem "kitchen-api/internal/app/order_item"
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/app/table"
	"kitchen-api/internal/app/user"

	"github.com/gin-gonic/gin"
	"github.com/ktmbeestech/yanshi/config"
	"github.com/ktmbeestech/yanshi/database"
)

func NewRouter(db *database.OrmDb) *gin.Engine {
	router := gin.Default()

	origin := config.Default().GetString("cors.origin")
	router.Use(func(c *gin.Context) {
		reqOrigin := c.Request.Header.Get("Origin")
		if reqOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", reqOrigin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		headers := c.Request.Header.Get("Access-Control-Request-Headers")
		if headers == "" {
			headers = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	auth.RegisterAuthRoutes(router, db)
	user.RegisterUserRoutes(router, db)
	restaurant.RegisterRestaurantRoutes(router, db)
	table.RegisterTableRoutes(router, db)
	menucategory.RegisterMenuCategoryRoutes(router, db)
	menuitem.RegisterMenuItemRoutes(router, db)
	order.RegisterOrderRoutes(router, db)
	orderitem.RegisterOrderItemRoutes(router, db)
	addon.RegisterAddonRoutes(router, db)

	return router
}
