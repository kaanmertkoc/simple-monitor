package database

import (
    "database/sql"
    "time"
    _ "github.com/mattn/go-sqlite3"
)

type Client struct {
    db *sql.DB
}

type MetricRow struct {
    Timestamp     int64
    CPUPercent    float64
    MemoryPercent float64
    MemoryTotal   uint64
    MemoryUsed    uint64
    DiskPercent   float64
    DiskTotal     uint64
    DiskUsed      uint64
}

func NewClient(dbPath string) (*Client, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    // Create tables if they don't exist
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS metrics (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp INTEGER NOT NULL,
            cpu_percent REAL NOT NULL,
            memory_percent REAL NOT NULL,
            memory_total INTEGER NOT NULL,
            memory_used INTEGER NOT NULL,
            disk_percent REAL NOT NULL,
            disk_total INTEGER NOT NULL,
            disk_used INTEGER NOT NULL
        );

        CREATE TABLE IF NOT EXISTS network_metrics (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            metric_id INTEGER NOT NULL,
            interface_name TEXT NOT NULL,
            bytes_sent INTEGER NOT NULL,
            bytes_received INTEGER NOT NULL,
            packets_sent INTEGER NOT NULL,
            packets_received INTEGER NOT NULL,
            FOREIGN KEY(metric_id) REFERENCES metrics(id) ON DELETE CASCADE
        );

        CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp);
    `)

    if err != nil {
        return nil, err
    }

    return &Client{db: db}, nil
}

func (c *Client) WriteMetrics(metrics map[string]interface{}) error {
    tx, err := c.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert main metrics
    result, err := tx.Exec(`
        INSERT INTO metrics (
            timestamp, cpu_percent, memory_percent, memory_total, 
            memory_used, disk_percent, disk_total, disk_used
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
        time.Now().Unix(),
        metrics["cpu_percent"],
        metrics["memory_percent"],
        metrics["memory_total"],
        metrics["memory_used"],
        metrics["disk_percent"],
        metrics["disk_total"],
        metrics["disk_used"],
    )
    if err != nil {
        return err
    }

    // Get the ID of inserted metrics
    metricId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Insert network metrics if present
    if networkData, ok := metrics["network"].(map[string]interface{}); ok {
        stmt, err := tx.Prepare(`
            INSERT INTO network_metrics (
                metric_id, interface_name, 
                bytes_sent, bytes_received,
                packets_sent, packets_received
            ) VALUES (?, ?, ?, ?, ?, ?)`)
        if err != nil {
            return err
        }
        defer stmt.Close()

        for iface, data := range networkData {
            netData := data.(map[string]interface{})
            _, err = stmt.Exec(
                metricId,
                iface,
                netData["bytes_sent"],
                netData["bytes_received"],
                netData["packets_sent"],
                netData["packets_received"],
            )
            if err != nil {
                return err
            }
        }
    }

    return tx.Commit()
}

func (c *Client) GetMetrics(duration time.Duration) ([]map[string]interface{}, error) {
    startTime := time.Now().Add(-duration).Unix()

    rows, err := c.db.Query(`
        SELECT 
            timestamp, cpu_percent, memory_percent,
            memory_total, memory_used,
            disk_percent, disk_total, disk_used
        FROM metrics
        WHERE timestamp > ?
        ORDER BY timestamp DESC
        LIMIT 1000`, startTime)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var metrics []map[string]interface{}
    for rows.Next() {
        var m MetricRow
        err := rows.Scan(
            &m.Timestamp,
            &m.CPUPercent,
            &m.MemoryPercent,
            &m.MemoryTotal,
            &m.MemoryUsed,
            &m.DiskPercent,
            &m.DiskTotal,
            &m.DiskUsed,
        )
        if err != nil {
            return nil, err
        }

        metrics = append(metrics, map[string]interface{}{
            "timestamp":      m.Timestamp,
            "cpu_percent":    m.CPUPercent,
            "memory_percent": m.MemoryPercent,
            "memory_total":   m.MemoryTotal,
            "memory_used":    m.MemoryUsed,
            "disk_percent":   m.DiskPercent,
            "disk_total":     m.DiskTotal,
            "disk_used":      m.DiskUsed,
        })
    }

    return metrics, nil
}

func (c *Client) Cleanup() error {
    _, err := c.db.Exec(`
        DELETE FROM metrics 
        WHERE timestamp < unixepoch() - (60 * 24 * 60 * 60)
    `)
    return err
}

func (c *Client) Close() error {
    return c.db.Close()
}