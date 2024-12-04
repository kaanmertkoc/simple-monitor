package handlers

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "github.com/yourusername/simple-monitor/models"
)

func GetMetrics(c *gin.Context) {
    cpu, err := cpu.Percent(time.Second, false)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    memory, err := mem.VirtualMemory()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    disk, err := disk.Usage("/")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    network, err := net.IOCounters(true)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    netMap := make(map[string]net.IOCountersStat)
    for _, net := range network {
        netMap[net.Name] = net
    }

    metrics := &models.Metrics{
        Timestamp: time.Now().Unix(),
        CPU:      cpu[0],
        Memory:   memory,
        Disk:     disk,
        Network:  netMap,
    }

    c.JSON(200, metrics)
}