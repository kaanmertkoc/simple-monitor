package handlers

import (
    "time"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "github.com/kaanmertkoc/simple-monitor/api/models"
)

type MetricsCollector struct{}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{}
}

func (mc *MetricsCollector) Collect() (*models.Metrics, error) {
    cpu, err := cpu.Percent(time.Second, false)
    if err != nil {
        return nil, err
    }

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

    netMap := make(map[string]net.IOCountersStat)
    for _, net := range network {
        netMap[net.Name] = net
    }

    return &models.Metrics{
        Timestamp: time.Now().Unix(),
        CPU:      cpu[0],
        Memory:   memory,
        Disk:     disk,
        Network:  netMap,
    }, nil
}