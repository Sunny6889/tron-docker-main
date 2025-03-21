# Use Prometheus Remote Write with Thanos to Monitor java-tron Node

In this document, we will introduce how to use Prometheus remote-write to monitor a java-tron node more securely.

## Background
The previous [README](README.md) explains how to monitor a java-tron node using Grafana and Prometheus. It can be illustrated by the image below:

<img src="../images/metric_pull_simple.png" alt="Alt Text" width="600" >

Basically, the Prometheus service pulls metrics from the java-tron node through an exposed port. Subsequently, Grafana retrieves these metrics from Prometheus to provide visualized insights and alerts.

This quick-start setup is only recommed on test enviornment, as there are some limitations to this approach. From a security perspective, it is essential to separate java-tron services and monitoring services into different network zones. Specifically, we need to isolate java-tron nodes, especially SR nodes, from external exposure to reduce risks such as Denial of Service (DoS) attacks. However, monitoring metrics and similar indicators of TRON blockchain status can be made more accessible to a broader range of users.
To address these concerns, we need to change the pull mode either from the java-tron or Prometheus service to push mode. Refer to Prometheus official documentation of ["Why do you pull rather than push"](https://prometheus.io/docs/introduction/faq/#why-do-you-pull-rather-than-push) and ["When to use the Pushgateway"](https://prometheus.io/docs/practices/pushing/#when-to-use-the-pushgateway), the best practice for long-live observation target is to use Prometheus pull mode, and put java-tron and Prometheus service in the same failure domain.

### New Architecture
Given these considerations, we will implement a push mode for the data flow from Prometheus to Grafana. Prometheus offers a **remote-write** feature that supports push mode, facilitating this transition. We have selected [Thanos](https://github.com/thanos-io/thanos) as an intermediate component. Thanos not only supports remote write but also provides additional features such as long-term storage, high availability, and global querying, thereby improving the overall architecture and functionality of our monitoring system.

Below is the new more reliable architecture of the monitoring system used on production. We will introduce how to set up the Prometheus remote-write feature and Thanos cluster in the following sections.
<img src="../images/metric_thanos_arch.png" alt="Alt Text" width="950" >

## Implementation Guide
This section introduces the steps of setting up Prometheus remote write with Thanos.

### Prerequisites
Before starting, ensure you have:
- Docker and Docker Compose installed (refer to [prerequisites](../README.md#prerequisites))
  - Make sure docker resource memory is at least 16GB, as java-tron node needs at least 12GB memory.
- The tron-docker repository cloned locally
```sh
git clone https://github.com/tronprotocol/tron-docker.git
cd tron-docker/metric_monitor/push_mode
```
### Main components
The monitoring system consists of:
- **TRON FullNode**: TRON FullNode service with metrics enabled
- **Prometheus**: Monitoring service that collects metrics from the java-tron node
- **Thanos Receive**: A component of Thanos that receives data from Prometheus’s remote write write-ahead log, exposes it, and/or uploads it to cloud storage
- **Thanos Query**: A component of Thanos that implements Prometheus’s v1 API to aggregate data from the underlying components
- **Grafana**: Visualization service that retrieves metrics from **Thanos Query** to provide visualized insights and alerts

### Step 1: Set up TRON and Prometheus services
Run the below command to start java-tron and Prometheus services:
```sh
docker-compose up -d tron-node1 prometheus
```

Review the [docker-compose-thanos.yml](docker-compose-thanos.yml) file, the command explanation of the java-tron service can be found in [Run Single Node](../single_node/README.md#run-the-container).

Below are the core configurations for the Prometheus service:
```yaml
  ports:
    - "9090:9090"  # Used for local Prometheus status check
  volumes:
    - ./conf:/etc/prometheus
    - ./prometheus_data:/prometheus
  command:
    - "--config.file=/etc/prometheus/prometheus-remote-write.yml" # Default path to the configuration file
    - "--storage.tsdb.path=/prometheus" # The path where Prometheus stores its metric database
    - "--storage.tsdb.retention.time=30d"
    - "--storage.tsdb.max-block-duration=30m" # The maximum duration for a block of time series data that can be stored in the time series database (TSDB)
    - "--storage.tsdb.min-block-duration=30m"
    - "--web.enable-lifecycle" # Makes Prometheus expose the /-/reload HTTP endpoints
    - "--web.enable-admin-api"
```
#### Key configuration elements:
##### 1. Prometheus configuration file

The Prometheus configuration file is set to use the [prometheus-remote-write.yml](conf/prometheus-remote-write.yml) by volume mapping `./conf/prometheus-remote-write.yml:...` and flag `--config.file=...`. It contains the configuration of `scrape_configs` and `remote_write`.
```yml
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
    headers:
      Content-Encoding: snappy    # Explicit compression, as Prometheus sends raw uncompressed samples by default (2.7× reduction)
```
- For `global.external_labels`:
  - The `external_labels` defined in the Prometheus configuration file are propagated with all metric data to Thanos Receive. You could add multiple Prometheus services with remote-write to the same Receive service, just make sure the `external_labels` are unique. It uniquely identifies the Prometheus instance, acting as critical tracing metadata that ultimately correlates metrics to their originating java-tron node. You can use it in Grafana dashboards using label-based filtering (e.g., `{monitor="java-tron-node1-remote-write"}`).

    <img src="../images/metric_push_external_label.png" alt="Alt Text" width="680" >

- For `scrape_configs`:
  - The `scrape_interval` is set to 3s for high-frequency monitoring. You can adjust it based on your requirements. As the metric is collected every time the promethus service call java-tron service on `service:9527/metrics`, it is recommend not to set the `send_interval` too small, as it will increase the load on the java-tron service. As the recommendated scrape interval for high-resolution requirements is 1-5s, here we set the `send_interval` to 3s.
  - The `targets` are the java-tron node's IP address and port. The Prometheus service scrapes metrics from this target.
  - The `labels` are key-value pairs that help identify the target in the Prometheus service. You can use it in Grafana dashboards using label-based filtering (e.g., `{group="group-tron"}`).

- For `remote_write`:
  - Fill the `url` with the IP address of the Thanos Receive service started in the first step.
  - The `headers` section specifies the compression method for the data sent to Thanos Receive. The `Content-Encoding: snappy` flag compresses the data before sending it to the Receive service. This reduces the network traffic and storage space required by the Receive service.
  - Check the official documentation [remote_write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write) for all configurations' explanation.

##### 2. Storage configurations
- The volumes command `- ./prometheus_data:/prometheus` mounts a local directory used by Prometheus to store metrics data.
  - Although in this case, we use Prometheus with remote-write, it also stores metrics data locally. Through http://host_IP:9090/, you can check the running status of the Prometheus service and observe targets.
- The `--storage.tsdb.retention.time=30d` flag specifies the retention period for the metrics data. Prometheus will automatically delete data older than 30 days. For each metric request of a java-tron(v4.7.6+) FullNode connecting Mainnet, it returns ~9KB raw metric data. With the `scrape_interval: 3s`, it takes about 7GB per month raw data. TSDB will do compression, **thus for 30d retention, one java-tron takes about promethus less than 1GB**. The TSDB is useful for long-term storage and querying.
- The `--storage.tsdb.max-block-duration=30m` and `--storage.tsdb.min-block-duration=30m` flags specify the block generation duration for the TSDB. In this case, it is set to 30 minutes. This means that the TSDB will generate a new block every 30 minutes with about 900KB raw metric data.
- Other storage flags can be found in the [official documentation](https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects). For a quick start, you could use the default values.

### Step 2: Set up Thanos Receive
 The [Thanos Receive](https://thanos.io/tip/components/receive.md/#receiver) service implements the Prometheus Remote Write API. It builds on top of the existing Prometheus TSDB and retains its usefulness while extending its functionality with long-term-storage, horizontal scalability, and downsampling. Prometheus instances are configured to continuously write metrics to it, and then Thanos Receive uploads TSDB blocks to an object storage bucket every 2 hours by default. Thanos Receive exposes the StoreAPI so that Thanos Queriers can query received metrics in real-time.

Run the below command to start the Thanos Receive and [Minio](https://github.com/minio/minio) service for long-term metric storage. Minio is S3 compatible object storage service.
```sh
docker-compose up -d thanos-receive minio
```

Core configuration for Thanos Receive in [docker-compose-thanos.yml](docker-compose-thanos.yml):
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
#### Key configuration elements:
##### 1. Storage configuration
- Local Storage:
  `./receive-data:/receive/data` maps the host directory for metric TSDB storage.
  - Retention Policy: `--tsdb.retention=15d` auto-purges data older than 15 days. **It takes about 0.5GB of Receive disk space per month for one java-tron(v4.7.6) FullNode connecting Mainnet**.

- External Storage:
  `./conf:/receive` mounts configuration files. The `--objstore.config-file` flag enables long-term storage in MinIO/S3-compatible buckets. In this case, it is [bucket_storage_minio.yml](conf/bucket_storage_minio.yml).
  - Thanos Receive uploads TSDB blocks to an object storage bucket every 2 hours by default.
  - Fallback Behavior: Omitting this flag keeps data local-only.

##### 2. Network configuration
- Remote Write `--remote-write.address=0.0.0.0:10908`: Receives Prometheus metrics. Prometheus instances are configured to continuously write metrics to it.
- Thanos Receive exposes the StoreAPI so that Thanos Query can query received metrics in **real-time**.
  - The `ports` combined with flags `--grpc-address, --http-address` expose the ports for the Thanos Query service.
- Security Note: `0.0.0.0` means it accepts all incoming connections from any IP address. For production, consider restricting access to specific IP addresses.

For more flags explanation and default value can be found in the official [Thanos Receive Flags](https://thanos.io/tip/components/receive.md/#flags) documentation.

#### Thanos Receive reliable deployment
For systems monitoring multiple services with increasing scale, there are two approaches to prevent Thanos Receive from becoming a single point of failure or performance bottleneck:

1. Deploy multiple independent Thanos Receive instances, each dedicated to different monitoring targets. The configurations outlined in this document provide clear guidance for setting up this distributed architecture.
2. Implement Thanos Receive in cluster mode, which provides a unified entry point with automatic scaling capabilities as illustrated in the architecture diagram. We are working to provide guidance.

### Step 3: Set up Thanos Query
As Grafana cannot directly query Thanos Receive, we need Thanos Query that implements Prometheus’s v1 API to aggregate data from the Receive services. **Querier is fully stateless and horizontally scalable**.

Run the below command to start the Thanos Query service:
```sh
docker-compose -f docker-compose-thanos.yml up -d querier
```

Below are the core configurations for the Thanos Query service:
``` yaml
  querier:
    ...
    container_name: querier
    ports:
      - "9091:9091"
    command:
      - query
      - --endpoint.info-timeout=30s
      - --http-address=0.0.0.0:9091
      - --query.replica-label=replica # Deduplication turned on for identical series except the replica label.
      - --store=thanos-receive:10907 # --store: The grpc-address of the Thanos Receive service，if Receive run remotely replace container name "thanos-receive" with the real ip
     # - --store=thanos-receive2:10907 # Add more thanos-receive services
```
It will set up the Thanos Query service
that listens to port 9091 and queries metrics from the Thanos Receive service from `--store=[Thanos Receive IP]:10907`.
Make sure the IP address is correct.

You could add multiple Thanos Receive services to the Querier service. It will do duplication based on the
For more complex usage, please refer to the [official Query document](https://thanos.io/tip/components/query.md/).

### Step 4: Monitor through Grafana
To start the Grafana service on the host machine, run the following command:
```sh
docker-compose -f docker-compose-thanos.yml up -d grafana
```
Then log in to the Grafana web UI through http://localhost:3000/ or your host machine's IP address. The initial username and password are both `admin`.
Click the **Connections** on the left side of the main page and select Prometheus as datasource. Enter the IP and port of the Query service in URL with `http://[Query service IP]:9091`.

<img src="../images/metric_grafana_datasource_query.png" alt="Alt Text" width="680" >

Follow the same instruction as [Import Dashboard](https://github.com/tronprotocol/tron-docker/blob/main/metric_monitor/README.md#import-dashboard) to import the dashboard.
Then you can play with it with different Thanos Receive/Query, Prometheus configurations.

### Step 5: Clean up
To stop and remove all or part of the services, you could run the below commands:
```sh
docker-compose -f docker-compose-thanos.yml down # Stop and remove all services
docker-compose -f docker-compose-thanos.yml down thanos-receive # Thanos Receive service only
docker-compose -f docker-compose-thanos.yml down prometheus, tron-node1, querier, grafana # Multiple Services at once
```
### Multiple targets metric monitoring
As stated in [Thanos Scaling](https://thanos.io/tip/thanos/design.md/#scaling):
"Overall, first-class horizontal sharding is possible but will not be considered for the time being since there’s no evidence that it is required in practical setups".
Thus for multiple services/targets monitoring, it is recommended to use Thanos Recieve + Querier per target.

## Troubleshooting

### Common Issues
1. **Container Config Error (Linux)**
- If you encounter a `KeyError: 'ContainerConfig'`, check for conflicting container names and remove them:
   ```bash
  # List all containers
  docker ps -a

  # Remove conflicting containers
  docker rm [container-name]
  ```
2. **Network Connectivity**

- Verify all services can communicate by checking logs:
  ```bash
  docker-compose logs [service-name]
  ```

- Ensure all IP addresses are correctly configured in the docker-compose file
- Ensure all exposed ports are accessible from external services

3. **Storage Issues**

- Check available disk space: `df -h`
- Monitor storage usage in Prometheus, Thanos Receive, and Minio (if used) directories

### Getting Help
For additional support:
- Raise an issue on [GitHub](https://github.com/tronprotocol/tron-docker/issues)
- Consult the official Thanos documentation
- Review Docker logs for specific service issues

## Conclusion
This guide provides a secure and scalable solution for monitoring java-tron nodes. For custom configurations beyond this setup, refer to the [official Thanos documentation](https://thanos.io/tip/thanos/quick-tutorial.md/) or engage with the community on [GitHub](https://github.com/tronprotocol/tron-docker/issues) Issue.
