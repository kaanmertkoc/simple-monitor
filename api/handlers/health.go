package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/yourusername/simple-monitor/models"
)

const VERSION = "1.0.0"

func HealthCheck(c *gin.Context) {
    response := models.HealthResponse{
        Status:  "ok",
        Version: VERSION,
    }
    c.JSON(200, response)
}