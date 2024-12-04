package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "monitor/api/handlers"
    "monitor/api/database"
)

func main() {
    // Initialize database client
    dbClient := database.NewClient(
        os.Getenv("INFLUXDB_URL"),
        os.Getenv("INFLUXDB_TOKEN"),
        os.Getenv("INFLUXDB_ORG"),
        os.Getenv("INFLUXDB_BUCKET"),
    )
    defer dbClient.Close()

    // Initialize metrics collector
    collector := handlers.NewMetricsCollector(dbClient)

    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    r.Use(cors.New(config))
    r.Use(gin.Recovery())

    // Basic endpoints
    r.GET("/health", handlers.HealthCheck)
    r.GET("/metrics", func(c *gin.Context) {
        metrics, err := collector.Collect()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, metrics)
    })

    // Historical metrics endpoints
    r.GET("/metrics/history", func(c *gin.Context) {
        hours := 3 // Default to 3 hours
        historical, err := collector.GetHistoricalMetrics(hours)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, historical)
    })

    log.Printf("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}