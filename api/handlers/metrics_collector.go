package handlers

import (
    "time"
    "sync"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "monitor/api/models"
    "monitor/api/database"
)

type MetricsCollector struct {
    db            *database.Client
    lastNetwork   map[string]net.IOCountersStat
    lastCollected time.Time
    mu            sync.Mutex
}

func NewMetricsCollector(db *database.Client) *MetricsCollector {
    return &MetricsCollector{
        db:          db,
        lastNetwork: make(map[string]net.IOCountersStat),
    }
}

func (mc *MetricsCollector) Collect() (*models.Metrics, error) {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    // Collect CPU
    cpuPercent, err := cpu.Percent(time.Second, false)
    if err != nil {
        return nil, err
    }

    // Collect Memory
    memory, err := mem.VirtualMemory()
    if err != nil {
        return nil, err
    }

    // Collect Disk
    disk, err := disk.Usage("/")
    if err != nil {
        return nil, err
    }

    // Collect Network
    network, err := net.IOCounters(true)
    if err != nil {
        return nil, err
    }

    // Process network stats
    networkStats := make(map[string]net.IOCountersStat)
    for _, netStat := range network {
        networkStats[netStat.Name] = netStat
    }

    metrics := &models.Metrics{
        Timestamp: time.Now().Unix(),
        CPU:      cpuPercent[0],
        Memory:   memory,
        Disk:     disk,
        Network:  networkStats,
    }

    // Store metrics in database if available
    if mc.db != nil {
        go mc.storeMetrics(metrics)
    }

    return metrics, nil
}

func (mc *MetricsCollector) storeMetrics(metrics *models.Metrics) {
    if mc.db == nil {
        return
    }
    
    metricsMap := map[string]interface{}{
        "cpu_percent":     metrics.CPU,
        "memory_percent":  metrics.Memory.UsedPercent,
        "disk_percent":    metrics.Disk.UsedPercent,
    }

    if err := mc.db.WriteMetrics(metricsMap, "system_metrics"); err != nil {
        println("Error storing metrics:", err.Error())
    }
}