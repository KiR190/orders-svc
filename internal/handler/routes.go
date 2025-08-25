package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupRoutes(router *gin.Engine) {
	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/orders", h.GetAllOrders)
		api.GET("/order/:id", h.GetOrder)
		api.POST("/order", h.CreateOrder)
	}

	// Static files
	router.Static("/static", "./internal/handler/static")
	router.GET("/", h.ServeStatic)
}

func (h *Handler) ServeStatic(c *gin.Context) {
	c.File("./internal/handler/static/index.html")
}

func (h *Handler) Run(port string) error {
	router := gin.Default()
	h.SetupRoutes(router)
	return router.Run(port)
}
