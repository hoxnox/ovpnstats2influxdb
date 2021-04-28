ovpnstats2influxdb
==================

Read metrics from OpenVPN `openvpn-status.log` and output them as InfluxDB metrics.

## Installation
* Build it yourself with `go get -u ./... && go build .`
* Make the binary executable  `chmod +x ovpnstats`
* Make the OpenVPN status file readable `chmod +r openvpn-status.log`

## Usage
Add it as an exec plugin to your `telegraf.conf`:
```
[[inputs.exec]]
    commands = ["/usr/bin/ovpnstats -path /etc/openvpn/server/openvpn-status.log"]
    data_format = "influx"
    timeout = 1
```

### Read number of connected clients
```flux
from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "openvpn")
  |> filter(fn: (r) => r["_field"] == "clients")
  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
  |> yield(name: "mean")
```

### Read number of active routes
```flux
from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "openvpn")
  |> filter(fn: (r) => r["_field"] == "routes")
  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
  |> yield(name: "mean")
```

### Read Sent and Received
```flux
from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "openvpn")
  |> filter(fn: (r) => r["_field"] == "sent" or r["_field"] == "received")
  |> pivot(
      rowKey:["_time"],
      columnKey: ["_field"],
      valueColumn: "_value"
  )
```

### Read IP Addresses of connected clients
```flux
from(bucket: "telegraf")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "openvpn")
  |> filter(fn: (r) => r["_field"] == "tunnel_address")
  |> distinct(column: "_value")
  |> keep(columns: ["name", "_value"])
`