# Use Prometheus Remote Write to Monitor java-tron Node

In this README, we will introduce how to use Prometheus remote-write to monitor java-tron node more securely.

## Background
The previous [README](../README.md) explains how to monitor a java-tron node using Grafana and Prometheus. It can be illustrated by the image below:
![image](../../images/metric_pull_simple.png)
Basically, the Prometheus service pulls metrics from java-tron node through an exposed port. Subsequently, Grafana retrieves these metrics from Prometheus to provide visualized insights and alerts.

There are some limitations to this approach. From a security perspective, it is essential to separate java-tron services and monitoring services into different network zones. Specifically, we need to isolate java-tron nodes, especially SR nodes, from external exposure to reduce risks such as Denial of Service (DoS) attacks. However, monitoring metrics and similar indicators of TRON blockchain status can be made more accessible to a broader range of users.
To address these concerns, we need to change the pull mode either from java-tron or Prometheus service to push mode. Refer to Prometheus official documentation of ["Why do you pull rather than push"](https://prometheus.io/docs/introduction/faq/#why-do-you-pull-rather-than-push) and ["When to use the Pushgateway"](https://prometheus.io/docs/practices/pushing/#when-to-use-the-pushgateway), the best practise for long-live observation target is to use Prometheus pull mode, and put java-tron and Prometheus service in the same failure domain.

### New Architecture
Given these considerations, we will implement a push mode for the data flow from Prometheus to Grafana. Prometheus offers a **remote-write** feature that supports push mode, facilitating this transition. We have selected [Thanos](https://github.com/thanos-io/thanos) as an intermediate component. Thanos not only supports remote write but also provides additional features such as long-term storage, high availability, and global querying, thereby improving the overall architecture and functionality of our monitoring system.

Below is the new architecture of the monitoring system. We will introduce how to set up the Prometheus remote-write feature and Thanos in the following sections.
![image](../../images/metric_push_with_thanos.png)

## Use Prometheus remote write with Thanos guidance
This section introduces the steps of setting up Prometheus remote write with Thanos.

### Prerequisites

- Docker and Docker Compose: Installation refer [prerequisites](../../README.md#prerequisites).
- Clone the tron-docker repository, then navigate to the `push_mode` directory.
```sh
git clone https://github.com/tronprotocol/tron-docker.git
cd tron-docker/metric_monitor/push_mode
```
### Main Components
Before we start, let's list the main components of the monitoring system:
- **TRON FullNode**: TRON FullNode service with metrics enabled
- **Prometheus**: Monitoring service that collects metrics from java-tron node
- **Thanos Receive**: A component of Thanos that receives data from Prometheus’s remote write write-ahead log, exposes it, and/or uploads it to cloud storage.
- **Thanos Query**: A component of Thanos that implements Prometheus’s v1 API to aggregate data from the underlying components.
- **Grafana**: Visualization service that retrieves metrics from **Thanos Query** to provide visualized insights and alerts.

### Step 1: Set up Thanos Receive
As we can see from the above architecture, Thanos Receive is the intermediate component we need to set up first. The [Thanos Receive](https://thanos.io/tip/components/receive.md/#receiver) service implements the Prometheus Remote Write API. It builds on top of existing Prometheus TSDB and retains its usefulness while extending its functionality with long-term-storage, horizontal scalability, and downsampling.

Run below command to start Thanos Receive service and a minio service for long-term metric storage:
```sh
docker-compose -f docker-compose-receive.yml up -d
```

Core configuration in [docker-compose-receive.yml](docker-compose-receive.yml):
```
  thanos-receive:
    ...
    container_name: thanos-receive
    volumes:
      - ./receive-data:/receive/data
      - ./conf:/receive
    ports:
      - "10907:10907"
      - "10908:10908"
      - "10909:10909"
    command:
      - "receive"
      - "--tsdb.path=/receive/data"
      - "--tsdb.retention=15d" # How long to retain raw samples on local storage.
      - "--grpc-address=0.0.0.0:10907"
      - "--http-address=0.0.0.0:10909"
      - "--remote-write.address=0.0.0.0:10908"
      - "--label=receive_replica=\"0\""
      - "--label=receive_cluster=\"java-tron-mainnet\""
      - "--objstore.config-file=/receive/bucket_storage_minio.yml"
```
#### Key Configuration Elements:
##### 1. Storage Configuration
- Local Storage:
`./receive-data:/receive/data` maps host directory for metric TSDB storage.
  - Retention Policy: `--tsdb.retention=15d` auto-purges data older than 15 days. As observed, it takes about 0.5GB of Receive disk space per month for one java-tron(v4.7.6) FullNode connecting Mainnet.

- External Storage:
`./conf:/receive` mounts configuration files. The `--objstore.config-file` flag enables long-term storage in MinIO/S3-compatible buckets. In this case it is [bucket_storage_minio.yml](conf/bucket_storage_minio.yml).
  - Thanos Receive uploads TSDB blocks to an object storage bucket every 2 hours by default.
  - Fallback Behavior: Omitting this flag keeps data local-only.

##### 2. Network Configuration
- Remote Write `--remote-write.address=0.0.0.0:10908`: Receives Prometheus metrics. Prometheus instances are configured to continuously write metrics to it.
- Thanos Receive exposes the StoreAPI so that Thanos Query can query received metrics in **real-time**.
  - The `ports` combined with flags `--grpc-address, --http-address` expose the ports for Thanos Query service.
- Security Note: `0.0.0.0` means it accepts all incoming connections from any IP address. For production, consider restricting access to specific IP addresses.

##### 3. Operational Parameters

- `--label=receive_replica=.` and `--label=receive_cluster=.`: Cluster labels ensure unique identification in Thanos ecosystem. You could find these labels in Grafana dashboards. You could add any key value pairs as labels.

<img src="../../images/metric_pull_receive_label.png" alt="Alt Text" width="880" >

For more flags explanation and default value can be found in official [Thanos documentation](https://thanos.io/tip/components/receive.md/#flags).

### Step 2: Set up TRON and Prometheus services
Run below command to start java-tron and Prometheus services:
```sh
docker-compose -f docker-compose-tron-prometheus.yml up -d
```
Review the [docker-compose-tron-prometheus.yml](docker-compose-tron-prometheus.yml) file, the command explanation of java-tron service can be found in the [README](../single_node/README.md#run-the-container).

Below are the core configurations for Prometheus service:
```yaml
  ports:
    - "9090:9090"  # Used for local Prometheus status check
  volumes:
    - ./conf/prometheus-remote-write.yml:/etc/prometheus/prometheus.yml  # Main config
    - ./prometheus_data:/prometheus  # Persistent metric local storage
  command:
    - --config.file=/etc/prometheus/prometheus.yml
    - --storage.tsdb.path=/prometheus
    - --storage.tsdb.retention.time=30d  # ~1GB/month storage requirement at Prometheus side
    - --storage.tsdb.max-block-duration=30m  # TSDB block generation frequency
    - --storage.tsdb.min-block-duration=30m
    - --web.enable-lifecycle  # Enable config reload endpoints "/-/reload"
    - --web.enable-admin-api  # Expose administration endpoints
```
#### Key Configuration Elements:
##### 1. Storage configurations
- The volumes command `- ./prometheus_data:/prometheus` mounts a local directory used by Prometheus to store metrics data.
  - Although in this case, we use Prometheus with remote-write, it also stores metrics data locally. Through http://localhost:9090/, you can check the running status of the Prometheus service and observe targets.
- The `--storage.tsdb.retention.time=30d` flag specifies the retention period for the metrics data. Prometheus will automatically delete data older than 30 days.
- Other storage flags can be found in the [official documentation](https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects). For a quick start, you could use the default values.

##### 2. Prometheus remote-write configuration

Prometheus configuration file is set
to use the [prometheus-remote-write.yml](conf/prometheus-remote-write.yml) by volume mapping `./conf/prometheus-remote-write.yml:...` and flag `--config.file=...`.
It contains configuration of `scrape_configs` and `remote_write`.
You need to fill the `url` with the IP address of the Thanos Receive service started at the first step.
Check the official documentation [remote_write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) for all configurations'
explanation.

```
global:
  external_labels:
    monitor: 'java-tron-node1-remote-write'  # Unique identifier in Thanos

scrape_configs:
  - job_name: java-tron
    scrape_interval: 3s  # High-frequency monitoring
    ...
    static_configs:
      - targets: [tron-node1:9527]
        labels:
          group: group-tron
          instance: fullnode-01

remote_write:
  - url: http://[THANOS_RECEIVE_IP]:10908/api/v1/receive
    metadata_config:
      send: true  # Enable metric metadata transmission
      send_interval: 3s  # How frequently metric metadata is sent to remote Receive.
      max_samples_per_send: 500  # Batch size optimization
```
The `external_labels` defined in Prometheus configuration file are propagated with all metric data to Thanos Receive.
It uniquely identifies the Prometheus instance,
acting as critical tracing metadata that ultimately correlates metrics to their originating java-tron node.
You can use it in Grafana dashboards using label-based filtering (e.g., {monitor="java-tron-node1-remote-write"}).

<img src="../../images/metric_push_external_label.png" alt="Alt Text" width="680" >

### step 3: Set up Thanos Query, Grafana
So far, we have Thanos Receive、Prometheus and java-tron services running.
java-tron as the origin produces metrics, then pulled by Prometheus service, which then keeping push metrics to Thanos Receive.
As Grafana cannot directly query Thanos Receive, we need Thanos Query that implements Prometheus’s v1 API to aggregate data from the underlying components.

```sh
docker-compose -f docker-compose-querier-grafana.yml up -d
```
