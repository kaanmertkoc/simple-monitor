package database

import (
    "time"
    "context"
    "github.com/influxdata/influxdb-client-go/v2"
    "github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
    client   influxdb2.Client
    writeAPI api.WriteAPIBlocking
    queryAPI api.QueryAPI
}

func NewClient(url, token, org, bucket string) *Client {
    client := influxdb2.NewClient(url, token)
    writeAPI := client.WriteAPIBlocking(org, bucket)
    queryAPI := client.QueryAPI(org)

    return &Client{
        client:   client,
        writeAPI: writeAPI,
        queryAPI: queryAPI,
    }
}

func (c *Client) WriteMetrics(metrics interface{}, measurement string) error {
    p := influxdb2.NewPoint(
        measurement,
        map[string]string{"host": "server1"},
        metrics.(map[string]interface{}),
        time.Now(),
    )
    
    return c.writeAPI.WritePoint(context.Background(), p)
}

func (c *Client) QueryLastHours(measurement string, hours int) ([]map[string]interface{}, error) {
    query := `
    from(bucket:"monitoring")
        |> range(start: -` + string(hours) + `h)
        |> filter(fn: (r) => r["_measurement"] == "` + measurement + `")
        |> yield(name: "last")
    `
    
    result, err := c.queryAPI.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }

    var metrics []map[string]interface{}
    for result.Next() {
        metrics = append(metrics, result.Record().Values())
    }

    return metrics, nil
}

func (c *Client) Close() {
    c.client.Close()
}