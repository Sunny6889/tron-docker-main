# Monitor Java-tron nodes using VictoriaMetrics
This document aims to facilitate the monitoring of metrics for Java tron nodes through VictoriaMetrics.

## Background
At present, Tron's full node and system monitoring are based on a monitoring platform built on Grafana+Prometheus. 
Each full node exposes specified metrics ports, and Prometheus actively pulls and stores data from these ports. 
Then, Prometheus is included in the data source by configuring the IP and port of Prometheus metrics through the Grafana page UI. 
For security reasons, it is not recommended for these full node nodes to expose specific ports to Prometheus for pulling data. 
Therefore, through research, it was found that VictoriaMetrics can store metrics data and provide queries through push mode.

## Evolution of Pull Change Push Architecture:
<img src="../images/metric_pull_to_push.png" alt="Alt Text" width="720" >

## VictoriaMetrics
VictoriaMetrics is a fast, cost-effective, and scalable monitoring solution and time series database.
VictoriaMetrics has the following significant features:
- Can be used as long-term storage for Prometheus.
- It can be integrated into Grafana as a direct alternative to Prometheus, as it is compatible with the Prometheus query API.
- Easy operation and maintenance:
    - Single lightweight executable file with no external dependencies
    - Configure by specifying command-line parameters (including reasonable default values)
    - All data is stored in the directory specified by the - storage Data Path parameter
    - Support vmbackup/vmrestore tools to achieve instant snapshot backup and recovery
- Innovative Query Language Metrics QL: Enhancing Functionality on PromQL
- Global query view: supports unified writing of multiple Prometheus instances or other data sources, and queries through a single interface
- High performance: Data ingestion and querying have excellent vertical/horizontal scalability, with performance up to 20 times that of InfluxDB and TimescaleDB
- Memory efficiency: When processing million level time series (high cardinality), the memory consumption is only 1/10 of InfluxDB and 1/7 of Prometheus/Thanos/Cortex
- High churn rate time series optimization
- Ultimate compression: With the same storage space, the data point capacity is 70 times that of TimescaleDB (benchmark test), and the storage space requirement is reduced by 7 times compared to Prometheus/Thanos/Cortex (benchmark test)
- High latency and low IOPS storage optimization: compatible with HDD and cloud platform network storage (AWS/Google Cloud/Azure, etc., see disk IO testing)
- Single node can replace medium-sized Thanos/M3DB/Cortex/InfluxDB/TimescaleDB clusters (see vertical scaling testing, Thanos comparison testing, and PromCon 2019 presentation)
- Abnormal shutdown protection: prevent data damage caused by OOM/power outage/kill -9 through storage architecture design
- Multi protocol support:
    - Prometheus exporter metric capture
    - Prometheus remote write API
    - Prometheus exposure format
    - InfluxDB line protocol for HTTP/TCP/UDP
    - Label based Graphite plain text protocol
    - OpenTSDB put messages and HTTP/API/put requests
    - JSON line format
    - Any CSV data
    - Native binary format
    - DataDog agent/DogStatsD
    - NewRelic Infrastructure Proxy
    - OpenTelemetry metric format
- Provide open-source cluster versions


The data ingestion and query rate performance between VictoriaMetrics and Prometheus is based on benchmark `node_exporter` testing using metrics. 
The data on memory and disk space usage is applicable to a single Prometheus or VictoriaMetrics server.

|         Feature          |            Prometheus             |               VictoriaMetrics               |
|:------------------------:|:---------------------------------:|:-------------------------------------------:|
|   **Data Collection**    |            Pull-based             |   Supports both pull-based and push-based   |
|    **Data Ingestion**    | Up to 240,000 samples per second  |      Up to 360,000 samples per second       |
|  **Query Performance**   |  Up to 80,000 queries per second  |      Up to 100,000 queries per second       |
|     **Memory Usage**     |          Up to 14GB RAM           |               Up to 4.3GB RAM               |
|   **Data Compression**   |       Uses LZF compression        |           Uses Snappy compression           |
| **Disk Write Frequency** | More frequent data writes to disk |      Less frequent data writes to disk      |
|   **Disk Space Usage**   |     Requires more disk space      |          Requires less disk space           |
|    **Query Language**    |              PromQL               | MetricsQL (backward-compatible with PromQL) |

## Deploy VictoriaMetrics
Here, use Docker and deploy VictoriaMetrics in cluster mode.
Deploying a cluster using `docker-compose`
Create `docker-compose.yml` file:
This configuration file contains two vmstorages (vmstorage1 and vmstorage2) configured.

```yaml
services:
  vmstorage1:
    image: victoriametrics/vmstorage:latest
    command:
      - -retentionPeriod=1y
      - -httpListenAddr=:8482
      - -vminsertAddr=:8400
      - -vmselectAddr=:8401
    ports:
      - "8482:8482"
      - "8400:8400"
      - "8401:8401"
    volumes:
      - ./storage1-data:/vmstorage-data
    networks:
      - vm-cluster-net

  vmstorage2:
    image: victoriametrics/vmstorage:latest
    command:
      - -retentionPeriod=1y
      - -httpListenAddr=:8482
      - -vminsertAddr=:8400
      - -vmselectAddr=:8401
    ports:
      - "8483:8482"
      - "8402:8400"
      - "8403:8401"
    volumes:
      - ./storage2-data:/vmstorage-data
    networks:
      - vm-cluster-net

  vminsert:
    image: victoriametrics/vminsert:latest
    command:
      - -httpListenAddr=:8480
      - -storageNode=vmstorage1:8400,vmstorage2:8400
    ports:
      - "8480:8480"
    networks:
      - vm-cluster-net

  vmselect:
    image: victoriametrics/vmselect:latest
    command:
      - -search.latencyOffset=1s
      - -httpListenAddr=:8481
      - -storageNode=vmstorage1:8401,vmstorage2:8401
    ports:
      - "8481:8481"
    networks:
      - vm-cluster-net

networks:
  vm-cluster-net:
    driver: bridge
```
Start cluster:
```shell
docker compose up -d 
```
Deploying scripts on Java-tron nodes to push metrics data to VictoriaMetrics

```shell
#!/bin/bash

# Capture source and destination addresses for indicators
METRICS_PORT=""
METRICS_URL="http://localhost:${METRICS_PORT}/metrics"

# Config VictoriaMetrics
VICTORIA_METRICS_IP=""
VICTORIA_METRICS_PORT=""
VICTORIA_METRICS_URL="http://${VICTORIA_METRICS_IP}:${VICTORIA_METRICS_PORT}/insert/0/prometheus/api/v1/import/prometheus"

# Configurable label variables
GROUP=""
INSTANCE=""
JOB=""
SLEEP_SECONDS=1

# Build additional tag parameters
EXTRA_LABELS="extra_label=group=${GROUP}&extra_label=instance=${INSTANCE}&extra_label=job=${JOB}"

while true; do
  curl -s "$METRICS_URL" | \
  curl -X POST \
      --data-binary @- \
      -H "Content-Type: text/plain" \
      "${VICTORIA_METRICS_URL}?${EXTRA_LABELS}"
  sleep "$SLEEP_SECONDS"
done
```

Example of querying data through the vmselect API interface (where 0 is the tenant ID)
```shell
curl 'http://localhost:8481/select/0/prometheus/api/v1/query?query={metrics}'
```
Key Configuration Description:
- Data directory: Specify the storage path through - storageDataSath, and it is recommended to use persistent volumes
- Port Description:
  - 8480: Cluster version write port (receiving Prometheus data)
  - 8481: Cluster version query port

## Integrate with Grafana
Ensure that the machine where Grafana is located can access the service instance of VictoriaMetrics, and add the data source in the UI of Grafana page.

<img src="../images/grafana_add_datasource.png" alt="Alt Text" width="720" >

Still choose Prometheus type data source:

<img src="../images/grafana_select_datasource.png" alt="Alt Text" width="720" >

Configure the VictoriaMetrics URL, the default port for querying the cluster version URL is 8481.

<img src="../images/grafana_cluster.png" alt="Alt Text" width="720" >

Push data local verification
Open Grafana Panel, click Edit ->Query ->Data Source, switch to VictoriaMetrics, and observe if there is any data in the panel.

<img src="../images/grafana_verify_local.png" alt="Alt Text" width="720" >














