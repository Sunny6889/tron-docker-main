# Use Prometheus Remote Write to Monitor java-tron Node

In this README, we will introduce how to use Prometheus remote-write to monitor java-tron node more securely.

## Background
The previous [README](../README.md) explains how to monitor a java-tron node using Grafana and Prometheus. It can be illustrated by the image below:
![image](../../images/metric_pull_simple.png)
Basically, the Prometheus service pulls metrics from java-tron node through an exposed port. Subsequently, Grafana retrieves these metrics from Prometheus to provide visualized insights and alerts.

There are some limitations to this approach. From a security perspective, it is essential to separate java-tron services and monitoring services into different network zones. Specifically, we need isolate java-tron nodes, especially SR nodes, from external exposure to reduce risks such as Denial of Service (DoS) attacks. However, monitoring metrics and similar indicators of TRON blockchain status can be made more accessible to a broader range of users.
To address these concerns, we need to change the pull mode either from java-tron or Prometheus service to push mode. Refer to Prometheus official documentation of ["Why do you pull rather than push"](https://prometheus.io/docs/introduction/faq/#why-do-you-pull-rather-than-push) and ["When to use the Pushgateway"](https://prometheus.io/docs/practices/pushing/#when-to-use-the-pushgateway), the best practise for long-live observation target is to use Prometheus pull mode, and put java-tron and Prometheus service in the same failure domain.

### New Architecture
Given these considerations, we will implement a push mode for the data flow from Prometheus to Grafana. Prometheus offers a remote-write feature that supports push mode, facilitating this transition. We have selected [Thanos](https://github.com/thanos-io/thanos) as an intermediate component. Thanos not only supports remote write but also provides additional features such as long-term storage, high availability, and global querying, thereby improving the overall architecture and functionality of our monitoring system.

Below is the new architecture of the monitoring system. We will introduce how to set up the Prometheus remote-write feature and Thanos in the following sections.
![image](../../images/metric_push_with_thanos.png)

## Use Prometheus remote write with Thanos guidance
This section introduce the steps of setting up Prometheus remote write with Thanos.

### Prerequisites

- Docker and Docker Compose: Installation refer [prerequisites](../../README.md#prerequisites).
- Clone the tron-docker repository, then navigate to the `push_mode` directory.
```sh
git clone https://github.com/tronprotocol/tron-docker.git
cd tron-docker/metric_monitor/push_mode
```

### Step 1: Set up TRON and Prometheus services
Run below command to start java-tron and Prometheus services:
```sh
docker-compose -f docker-compose-tron-prometheus.yml up -d
```
Review the [docker-compose-tron-prometheus.yml](docker-compose-tron-prometheus.yml) file, the command explanation of java-tron service can be found in the [README](../single_node/README.md#run-the-container).
Below is the Prometheus service command:
```
  prometheus:
    image: quay.io/prometheus/prometheus:v2.38.0
    container_name: prometheus
    networks:
      - tron_network
    deploy:
      resources:
        limits:
          memory: 1g
    volumes:
      - ./conf/prometheus-remote-write.yml:/etc/prometheus/prometheus.yml
      - ./prometheus_data:/prometheus
    command:
      - --config.file=/etc/prometheus/prometheus.yml
      - --storage.tsdb.path=/prometheus # the path where Prometheus stores its metric database
      - --storage.tsdb.retention.time=30d # takes about 1GB of disk space per month
      - --storage.tsdb.max-block-duration=30m #the maximum duration for a block of time series data that can be stored in the time series database (TSDB)
      - --storage.tsdb.min-block-duration=30m
      - --web.enable-lifecycle # tell Prometheus to expose the /-/reload HTTP endpoints
      - --web.enable-admin-api
    restart: unless-stopped
```

### Step 2: Set up Thanos Receive
Receiver: receives data from Prometheus’s remote write write-ahead log, exposes it, and/or uploads it to cloud storage.

### step 3: Set up Thanos Query, Grafana
Querier/Query: implements Prometheus’s v1 API to aggregate data from the underlying components.
