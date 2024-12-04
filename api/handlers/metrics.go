package handlers

import (
    "github.com/gin-gonic/gin"
)

func GetMetrics(c *gin.Context) {
    // Create collector without database for now
    collector := NewMetricsCollector(nil)
    metrics, err := collector.Collect()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, metrics)
}