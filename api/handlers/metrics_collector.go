package handlers

import (
    "time"
    "sync"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/host"
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

func (mc *MetricsCollector) calculateNetworkSpeed(current, last net.IOCountersStat, duration time.Duration) models.NetworkSpeed {
    seconds := duration.Seconds()
    return models.NetworkSpeed{
        BytesSentPerSec:     float64(current.BytesSent-last.BytesSent) / seconds,
        BytesReceivedPerSec: float64(current.BytesRecv-last.BytesRecv) / seconds,
        PacketsSentPerSec:   float64(current.PacketsSent-last.PacketsSent) / seconds,
        PacketsReceivedPerSec: float64(current.PacketsRecv-last.PacketsRecv) / seconds,
        ErrorsInPerSec:      float64(current.Errin-last.Errin) / seconds,
        ErrorsOutPerSec:     float64(current.Errout-last.Errout) / seconds,
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

    // Calculate network speeds
    now := time.Now()
    networkSpeeds := make(map[string]models.NetworkSpeed)
    networkStats := make(map[string]models.NetworkStats)

    for _, netStat := range network {
        if last, exists := mc.lastNetwork[netStat.Name]; exists {
            duration := now.Sub(mc.lastCollected)
            networkSpeeds[netStat.Name] = mc.calculateNetworkSpeed(netStat, last, duration)
        }
        
        networkStats[netStat.Name] = models.NetworkStats{
            BytesSent:    netStat.BytesSent,
            BytesReceived: netStat.BytesRecv,
            PacketsSent:   netStat.PacketsSent,
            PacketsReceived: netStat.PacketsRecv,
            Errors:        netStat.Errin + netStat.Errout,
            Dropped:       netStat.Dropin + netStat.Dropout,
        }

        mc.lastNetwork[netStat.Name] = netStat
    }
    mc.lastCollected = now

    metrics := &models.Metrics{
        Timestamp:     now.Unix(),
        CPU:          cpuPercent[0],
        Memory:       memory,
        Disk:         disk,
        NetworkStats: networkStats,
        NetworkSpeed: networkSpeeds,
    }

    // Store metrics in database
    if mc.db != nil {
        go mc.storeMetrics(metrics)
    }

    return metrics, nil
}

func (mc *MetricsCollector) storeMetrics(metrics *models.Metrics) {
    metricsMap := map[string]interface{}{
        "cpu_percent":     metrics.CPU,
        "memory_percent":  metrics.Memory.UsedPercent,
        "disk_percent":    metrics.Disk.UsedPercent,
    }

    if err := mc.db.WriteMetrics(metricsMap, "system_metrics"); err != nil {
        // Log error but don't fail the request
        println("Error storing metrics:", err.Error())
    }
}

func (mc *MetricsCollector) GetHistoricalMetrics(hours int) (*models.HistoricalMetrics, error) {
    results, err := mc.db.QueryLastHours("system_metrics", hours)
    if err != nil {
        return nil, err
    }

    historical := &models.HistoricalMetrics{
        CPU:    make([]models.TimeseriesPoint, 0),
        Memory: make([]models.TimeseriesPoint, 0),
        Disk:   make([]models.TimeseriesPoint, 0),
    }

    for _, result := range results {
        timestamp := result["_time"].(time.Time)
        historical.CPU = append(historical.CPU, models.TimeseriesPoint{
            Timestamp: timestamp,
            Value:     result["cpu_percent"].(float64),
        })
        // Add memory and disk points similarly
    }
	for _, result := range results {
		timestamp := result["_time"].(time.Time)
		historical.Memory = append(historical.Memory, models.TimeseriesPoint{
			Timestamp: timestamp,
			Value:     result["memory_percent"].(float64),
		})
		historical.Disk = append(historical.Disk, models.TimeseriesPoint{
			Timestamp: timestamp,
			Value:     result["disk_percent"].(float64),
		})
	}

    return historical, nil
}
