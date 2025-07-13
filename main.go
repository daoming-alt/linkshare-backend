// main.go
package main

import (
	"linkshare-backend/db"
	"linkshare-backend/handlers"
	"linkshare-backend/middleware"

	"github.com/gin-gonic/gin"
)

// @title LinkShare Backend API
// @version 1.0
// @description Backend for sharing links between devices.
// @host localhost:8080
// @BasePath /api
func main() {
    // Initialize SQLite database
    db.InitDB()

    // Create Gin router
    r := gin.Default()

    // Define API routes
    api := r.Group("/api")
    {
        // Authentication routes
        api.POST("/register", handlers.Register)
        api.POST("/login", handlers.Login)

        // Device routes (protected by JWT)
        api.POST("/devices", middleware.AuthMiddleware(), handlers.RegisterDevice)
        api.GET("/devices", middleware.AuthMiddleware(), handlers.ListDevices)

        // Link routes (protected by JWT)
        api.POST("/links", middleware.AuthMiddleware(), handlers.SendLink)
    }

    // WebSocket endpoint for real-time link delivery
    r.GET("/ws", handlers.WebSocketHandler)

    // Start server
    r.Run(":8080")
}
