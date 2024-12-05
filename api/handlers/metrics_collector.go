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
    lastCPU       []float64
    lastCPUTime   time.Time
    mu            sync.RWMutex
}

func NewMetricsCollector(db *database.Client) *MetricsCollector {
    mc := &MetricsCollector{
        db:          db,
        lastNetwork: make(map[string]net.IOCountersStat),
    }
    mc.lastCPU, _ = cpu.Percent(0, false)
    mc.lastCPUTime = time.Now()

    // Start background metric collection if DB is available
    if db != nil {
        go mc.startPeriodicCollection()
    }

    return mc
}

func (mc *MetricsCollector) startPeriodicCollection() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        metrics, err := mc.Collect()
        if err != nil {
            continue
        }
        mc.storeMetrics(metrics)
    }
}

func (mc *MetricsCollector) getCPUUsage() float64 {
    mc.mu.Lock()
    defer mc.mu.Unlock()

    currentCPU, err := cpu.Percent(0, false)
    if err != nil {
        return 0
    }

    mc.lastCPU = currentCPU
    mc.lastCPUTime = time.Now()

    if len(currentCPU) > 0 {
        return currentCPU[0]
    }
    return 0
}
func (mc *MetricsCollector) Collect() (*models.Metrics, error) {
    memory, err := mem.VirtualMemory()
    if err != nil {
        return nil, err
    }

    disk, err := disk.Usage("/")
    if err != nil {
        return nil, err
    }

    network, err := net.IOCounters(true)
    if err != nil {
        return nil, err
    }

    networkStats := make(map[string]net.IOCountersStat)
    for _, netStat := range network {
        networkStats[netStat.Name] = netStat
    }

    metrics := &models.Metrics{
        Timestamp: time.Now().Unix(),
        CPU:      mc.getCPUUsage(),
        Memory:   memory,
        Disk:     disk,
        Network:  networkStats,
    }

    return metrics, nil
}

func (mc *MetricsCollector) storeMetrics(metrics *models.Metrics) {
    if mc.db == nil {
        return
    }
    
    // Convert metrics to map for storage
    networkData := make(map[string]interface{})
    for name, stats := range metrics.Network {
        networkData[name] = map[string]interface{}{
            "bytes_sent":      stats.BytesSent,
            "bytes_received":  stats.BytesRecv,
            "packets_sent":    stats.PacketsSent,
            "packets_received": stats.PacketsRecv,
        }
    }

    metricsMap := map[string]interface{}{
        "cpu_percent":     metrics.CPU,
        "memory_percent":  metrics.Memory.UsedPercent,
        "memory_total":    metrics.Memory.Total,
        "memory_used":     metrics.Memory.Used,
        "disk_percent":    metrics.Disk.UsedPercent,
        "disk_total":      metrics.Disk.Total,
        "disk_used":       metrics.Disk.Used,
        "network":         networkData,
    }

    if err := mc.db.WriteMetrics(metricsMap); err != nil {
        println("Error storing metrics:", err.Error())
    }
}