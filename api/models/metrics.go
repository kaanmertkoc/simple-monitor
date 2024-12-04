package models

import (
    "time"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
)

type Metrics struct {
    Timestamp    int64                      `json:"timestamp"`
    CPU          float64                    `json:"cpu"`
    Memory       *mem.VirtualMemoryStat     `json:"memory"`
    Disk         *disk.UsageStat            `json:"disk"`
    NetworkStats map[string]NetworkStats    `json:"network_stats"`
    NetworkSpeed map[string]NetworkSpeed    `json:"network_speed"`
}

type NetworkStats struct {
    BytesSent       uint64  `json:"bytes_sent"`
    BytesReceived   uint64  `json:"bytes_received"`
    PacketsSent     uint64  `json:"packets_sent"`
    PacketsReceived uint64  `json:"packets_received"`
    Errors          uint64  `json:"errors"`
    Dropped         uint64  `json:"dropped"`
}

type NetworkSpeed struct {
    BytesSentPerSec      float64 `json:"bytes_sent_per_sec"`
    BytesReceivedPerSec  float64 `json:"bytes_received_per_sec"`
    PacketsSentPerSec    float64 `json:"packets_sent_per_sec"`
    PacketsReceivedPerSec float64 `json:"packets_received_per_sec"`
    ErrorsInPerSec       float64 `json:"errors_in_per_sec"`
    ErrorsOutPerSec      float64 `json:"errors_out_per_sec"`
}

type TimeseriesPoint struct {
    Timestamp time.Time `json:"timestamp"`
    Value     float64  `json:"value"`
}

type HistoricalMetrics struct {
    CPU    []TimeseriesPoint `json:"cpu"`
    Memory []TimeseriesPoint `json:"memory"`
    Disk   []TimeseriesPoint `json:"disk"`
}

type HealthResponse struct {
    Status  string `json:"status"`
    Version string `json:"version"`
}