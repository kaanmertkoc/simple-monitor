package main

import (
    "log"
    "os"
    "path/filepath"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/shirou/gopsutil/cpu"
    "monitor/api/handlers"
    "monitor/api/database"
)

func main() {
    // Set up data directory
    dataDir := "./simple-monitor"
    if dir := os.Getenv("DATA_DIR"); dir != "" {
        dataDir = dir
    }

    // Create data directory if it doesn't exist
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        log.Fatalf("Failed to create data directory: %v", err)
    }

    // Initialize SQLite database
    dbPath := filepath.Join(dataDir, "metrics.db")
    db, err := database.NewClient(dbPath)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

    // Initialize metrics collector
    collector := handlers.NewMetricsCollector(db)

    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()

    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    r.Use(cors.New(config))
    r.Use(gin.Recovery())

    // Basic routes
    r.GET("/health", handlers.HealthCheck)
    r.GET("/metrics", func(c *gin.Context) {
        metrics, err := collector.Collect()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, metrics)
    })

    // Summary endpoint
    r.GET("/summary", func(c *gin.Context) {
        // Get the latest metrics
        metrics, err := collector.Collect()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        // Get CPU count
        cpuCount, err := cpu.Counts(true)
        if err != nil {
            cpuCount = 1 // Default to 1 if we can't get the count
        }

        // Get CPU info
        cpuInfo, err := cpu.Info()
        var cpuVendor string
        if err != nil || len(cpuInfo) == 0 {
            cpuVendor = "unknown"
        } else {
            cpuVendor = cpuInfo[0].VendorID
        }

        // Create summary response
        summary := gin.H{
            "total_vcpu":   cpuCount,
            "cpu_vendor":   cpuVendor,    // Will show "GenuineIntel" for Intel or "AuthenticAMD" for AMD
            "total_disk":   metrics.Disk.Total,    // in bytes
            "total_memory": metrics.Memory.Total,  // in bytes
        }

        c.JSON(200, summary)
    })

    // Historical data route
    r.GET("/metrics/history", func(c *gin.Context) {
        // Default to 24 hours if duration not specified
        duration := "24h"
        if d := c.Query("duration"); d != "" {
            duration = d
        }

        // Parse duration string (e.g., "24h", "7d", "1h")
        d, err := time.ParseDuration(duration)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid duration format"})
            return
        }

        history, err := db.GetMetrics(d)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

        c.JSON(200, history)
    })

    /* Cleanup old data periodically (runs every hour)
    go func() {
        ticker := time.NewTicker(1 * time.Hour)
        for range ticker.C {
                
            if err := db.Cleanup(); err != nil {
                log.Printf("Error cleaning up old data: %v", err)
            }
        }
    }()
    */
    
    log.Printf("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}