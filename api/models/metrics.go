package models

import (
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
)

type Metrics struct {
    Timestamp int64                          `json:"timestamp"`
    CPU       float64                        `json:"cpu"`
    Memory    *mem.VirtualMemoryStat         `json:"memory"`
    Disk      *disk.UsageStat               `json:"disk"`
    Network   map[string]net.IOCountersStat `json:"network"`
}

type HealthResponse struct {
    Status  string `json:"status"`
    Version string `json:"version"`
}