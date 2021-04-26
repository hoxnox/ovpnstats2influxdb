package main

import (
	"github.com/thor77/ovpnstats"
)

// Metric represents metric events with all required information to write it into an InfluxDB
type Metric struct {
	Fields map[string]interface{}
	Tags   map[string]string
}

func createMetrics(clients []ovpnstats.ClientInfo, routes []ovpnstats.RoutingInfo) []Metric {
	return []Metric{{map[string]interface{}{"clients": len(clients), "routes": len(routes)}, nil}}
}

func createClientMetrics(clients []ovpnstats.ClientInfo) []Metric {
	var metrics []Metric
	for _, client := range clients {
		metrics = append(metrics, Metric{
			map[string]interface{}{"sent": client.BytesSent, "received": client.BytesReceived},
			map[string]string{"name": client.Name, "username": client.Username, "virtual_ip4": client.VirtualAddress},
		})
	}
	return metrics
}
