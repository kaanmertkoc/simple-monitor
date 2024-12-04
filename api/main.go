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
    // Initialize database client if environment variables are set
    var db *database.Client
    if url := os.Getenv("INFLUXDB_URL"); url != "" {
        db = database.NewClient(
            url,
            os.Getenv("INFLUXDB_TOKEN"),
            os.Getenv("INFLUXDB_ORG"),
            os.Getenv("INFLUXDB_BUCKET"),
        )
        defer db.Close()
    }

    // Initialize metrics collector
    collector := handlers.NewMetricsCollector(db)

    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    r.Use(cors.New(config))
    r.Use(gin.Recovery())

    // Routes
    r.GET("/health", handlers.HealthCheck)
    r.GET("/metrics", func(c *gin.Context) {
        metrics, err := collector.Collect()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, metrics)
    })

    log.Printf("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}