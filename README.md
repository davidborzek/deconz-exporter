[![Go Report Card](https://goreportcard.com/badge/github.com/davidborzek/deconz-exporter)](https://goreportcard.com/report/github.com/davidborzek/deconz-exporter)
# deCONZ Prometheus Exporter

This is a Prometheus exporter for deCONZ / Phoscon.

## Installation

### Using Docker

```bash
docker run \
  -e 'DECONZ_URL=http://127.0.0.1' \
  -e 'DECONZ_API_KEY=mykey' \
  -p 8080:8080 \
  deconz-exporter:latest
```

### Obtaining an deCONZ API Key

You can use the `auth` command to acquire a new key.

The `URL` specifies the url of your deCONZ server.

You can optionally provide `USERNAME` to set a custom key with a minimum length of 10 and a maximum length of 40. 

```bash
docker run \
  -e 'URL=http://127.0.0.1' \
  deconz-exporter:latest auth
```

### Prometheus config

Once you have configured deconz-exporter update your `prometheus.yml` scrape config:

```yaml
scrape_configs:
  - job_name: 'deconz'
    static_configs:
      - targets: ['localhost:8080']
```

### Exported Metrics

Currently the exporter exports all numeric states of a sensor into a single gauge:

```
# HELP deconz_sensor_state Sensor state value
# TYPE deconz_sensor_state gauge
deconz_sensor_state{manufacturername,modelid,name,sensor,state,type}
```

The gauge has multiple labels to identify the sensor and state. 
