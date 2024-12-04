package handlers

import (
    "github.com/gin-gonic/gin"
    "monitor/api/models"
)

const VERSION = "1.0.0"

func HealthCheck(c *gin.Context) {
    response := models.HealthResponse{
        Status:  "ok",
        Version: VERSION,
    }
    c.JSON(200, response)
}