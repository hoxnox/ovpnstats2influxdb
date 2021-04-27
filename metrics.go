package main

import (
	"log"
)

// Metric represents metric events with all required information to write it into an InfluxDB
type Metric struct {
	Fields map[string]interface{}
	Tags   map[string]string
}

func createMetrics(clients []ClientInfo, routes []RoutingInfo) []Metric {
	log.Println("createMetrics ", clients, routes)
	return []Metric{{map[string]interface{}{"clients": len(clients), "routes": len(routes)}, nil}}
}

func createClientMetrics(clients []ClientInfo) []Metric {
	log.Println("createClientMetrics ", clients)
	var metrics []Metric
	for _, client := range clients {
		metrics = append(metrics, Metric{
			map[string]interface{}{"sent": client.BytesSent, "received": client.BytesReceived},
			map[string]string{"name": client.CommonName},
		})
	}
	return metrics
}

func createRoutingMetrics(routes []RoutingInfo) []Metric {
	log.Println("createRoutingMetrics ", routes)
	var metrics []Metric
	for _, route := range routes {
		metrics = append(metrics, Metric{
			map[string]interface{}{"tunnel_address": route.VirtualAddress},
			map[string]string{"name": route.CommonName},
		})
	}
	return metrics
}
