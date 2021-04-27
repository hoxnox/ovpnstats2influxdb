package main

import (
	"fmt"
)

// runTelegraf prints Telegraf-compatible metric-output to stdout
func runTelegraf(path string) error {
	clients, routes, err := ParseStatusFile(path)
	if err != nil {
		return err
	}

	metrics := createMetrics(clients, routes)

	// convert metrics to influxdb line protocol
	points, err := createBatchPoints("openvpn", metrics)
	if err != nil {
		return err
	}

	// output line protocol lines
	for _, point := range points {
		fmt.Println(point.String())
	}

	clientMetrics := createClientMetrics(clients)

	clientPoints, err := createBatchPoints("openvpn", clientMetrics)
	if err != nil {
		return err
	}

	// output line protocol lines
	for _, point := range clientPoints {
		fmt.Println(point.String())
	}

	routingMetrics := createRoutingMetrics(routes)

	routingPoints, err := createBatchPoints("openvpn", routingMetrics)
	if err != nil {
		return err
	}

	// output line protocol lines
	for _, point := range routingPoints {
		fmt.Println(point.String())
	}
	return nil
}
