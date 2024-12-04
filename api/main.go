package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/yourusername/simple-monitor/handlers"
)

func main() {
    // Set Gin to release mode in production
    gin.SetMode(gin.ReleaseMode)

    // Initialize router
    r := gin.Default()

    // Configure CORS
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    r.Use(cors.New(config))

    // Add basic middleware
    r.Use(gin.Recovery())

    // Routes
    r.GET("/health", handlers.HealthCheck)
    r.GET("/metrics", handlers.GetMetrics)

    // Start server
    log.Printf("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}