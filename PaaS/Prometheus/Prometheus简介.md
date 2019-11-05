# Prometheus简介 #

[https://prometheus.io/](https://prometheus.io/)

[https://github.com/prometheus](https://github.com/prometheus)

[https://yunlzheng.gitbook.io/prometheus-book/](https://yunlzheng.gitbook.io/prometheus-book/)

[https://www.bookstack.cn/read/prometheus-book/quickstart-why-monitor.md](https://www.bookstack.cn/read/prometheus-book/quickstart-why-monitor.md)

摘自[https://www.ibm.com/developerworks/cn/cloud/library/cl-lo-prometheus-getting-started-and-practice/index.html](https://www.ibm.com/developerworks/cn/cloud/library/cl-lo-prometheus-getting-started-and-practice/index.html)

**Prometheus特点**
- 强大的多维度数据模型：
	- 时间序列数据通过metric名和键值对来区分。
	- 所有的metrics都可以设置任意的多维标签。
	- 数据模型更随意，不需要刻意设置为以点分隔的字符串。
	- 可以对数据模型进行聚合，切割和切片操作。
	- 支持双精度浮点类型，标签可以设为全unicode。
- 灵活而强大的查询语句（PromQL）：在同一个查询语句，可以对多个metrics进行乘法、加法、连接、取分数位等操作。
- 易于管理： Prometheus server 是一个单独的二进制文件，可直接在本地工作，不依赖于分布式存储。
- 高效：平均每个采样点仅占3.5 bytes，且一个Prometheus server可以处理数百万的metrics。
- 使用pull模式采集时间序列数据，这样不仅有利于本机测试而且可以避免有问题的服务器推送坏的metrics。
- 可以采用push gateway的方式把时间序列数据推送至Prometheus server端。
- 可以通过服务发现或者静态配置去获取监控的targets。
- 有多种可视化图形界面。
- 易于伸缩。

**架构**

![](img/Prometheus.png)
![](img/Prometheus_Arch.png)

主要组件：

- Prometheus Server: 用于收集和存储时间序列数据，提供PromQL查询语言的支持。
- Client Library: 客户端库，为需要监控的服务生成相应的metrics并暴露给Prometheus server。当Prometheus server来pull时，直接返回实时状态的metrics。
- Push Gateway: 主要用于短期的Jobs。由于这类jobs存在时间较短，可能在Prometheus来pull之前就消失了。为此，这些jobs可以直接向Prometheus server端推送它们的metrics。这种方式主要用于服务层面的metrics，对于机器层面的metrices，需要使用node exporter。
- Exporters: 用于暴露已有的第三方服务的metrics给Prometheus。
- Alertmanager: 从Prometheus server端接收到alerts后，会进行去除重复数据，分组，并路由到对收的接受方式，发出报警。常见的接收方式有：电子邮件，pagerduty，OpsGenie, webhook 等。

Prometheus的工作流

- Prometheus Daemon负责定时去目标上抓取metrics(指标)数据，每个抓取目标需要暴露一个http服务的接口给它定时抓取。Prometheus支持通过配置文件、文本文件、Zookeeper、Consul、DNS SRV Lookup等方式指定抓取目标。Prometheus采用PULL的方式进行监控，即服务器可以直接通过目标PULL数据或者间接地通过中间网关来Push数据。
- Prometheus server定期从配置好的jobs或者exporters中拉取metrics，或者接收来自Pushgateway发送过来的metrics，或者从其它的Prometheus server中拉metrics。
- Prometheus server在本地存储收集到的metrics，并运行定义好的alerts.rules，记录新的时间序列或者向Alert manager推送警报。
- Alertmanager根据配置文件，对接收到的警报进行处理，发出告警。
- 在图形界面中，可视化采集数据。

## 数据模型 ##

**时序索引**

时序(time series) 是由指标名字(Metric)，以及一组key/value标签定义的，具有相同的名字以及标签属于相同时序。时序的名字由ASCII字符，数字，下划线，以及冒号组成，它必须满足正则表达式`[a-zA-Z_:][a-zA-Z0-9_:]*`, 其名字应该具有语义化，一般表示一个可以度量的指标，例如: `http_requests_total`, 可以表示http请求的总数。时序的标签可以使Prometheus的数据更加丰富，能够区分具体不同的实例，例如：`http_requests_total{method="POST"}`可以表示所有http中的POST请求。标签名称由ASCII字符，数字，以及下划线组成，其中__开头属于Prometheus保留，标签的值可以是任何Unicode字符，支持中文。

**时序样本**

按照某个时序以时间维度采集的数据，称之为样本，其值包含：
- 一个float64值
- 一个毫秒级的unix时间戳

**格式**

Prometheus时序格式与OpenTSDB相似：

`<metric name>{<label name>=<label value>, ...}`

其中包含时序名字以及时序的标签。

**四种时序类型**

Prometheus时序数据分为Counter（变化的增减量）,Gauge（瞬时值）,Histogram（采样并统计）,Summary（采样结果）四种类型。

- Counter：表示收集的数据是按照某个趋势（增加／减少）一直变化的，我们往往用它记录服务请求总量、错误总数等。 例如Prometheus server中`http_requests_total`, 表示Prometheus处理的http请求总数，我们可以使用delta, 很容易得到任意区间数据的增量
- Gauge:表示搜集的数据是一个瞬时的值，与时间没有关系，可以任意变高变低，往往可以用来记录内存使用率、磁盘使用率等。 例如Prometheus server中`go_goroutines`, 表示Prometheus当前goroutines的数量。
- Histogram:主要用于表示一段时间范围内对数据进行采样（通常是请求持续时间或响应大小），并能够对其指定区间以及总数进行统计，通常它采集的数据展示为直方图。Histogram由`<basename>_bucket{le="<upper inclusive bound>"}，<basename>_bucket{le="+Inf"}, <basename>_sum，<basename>_count`组成，例如Prometheus server中`prometheus_local_storage_series_chunks_persisted`, 表示Prometheus中每个时序需要存储的chunks数量，我们可以用它计算待持久化的数据的分位数。
- Summary:主要用于表示一段时间内数据采样结果（通常是请求持续时间或响应大小），它直接存储了quantile数据，而不是根据统计区间计算出来的。Summary和Histogram类似，由`<basename>{quantile="<φ>"}，<basename>_sum，<basename>_count`组成，例如Prometheus server中`prometheus_target_interval_length_seconds`。


## 文本数据格式 ##

**注释**

文本内容，如果以#开头通常表示注释。
- 以`# HELP`开头表示metric帮助说明。
- 以`# TYPE`开头表示定义metric类型，包含counter, gauge, histogram, summary, 和untyped类型。
- 其他表示一般注释，供阅读使用，将被Prometheus忽略。

**采样数据**

内容如果不以#开头，表示采样数据。它通常紧挨着类型定义行，满足以下格式：

```
metric_name [
  "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } [ "," ] "}"
] value [ timestamp ]
```

注意：

假设采样数据metric叫做x, 如果x是histogram或summary类型必需满足以下条件：

- 采样数据的总和应表示为`x_sum`。
- 采样数据的总量应表示为`x_count`。
- summary类型的采样数据的quantile应表示为`x{quantile="y"}`。
- histogram类型的采样分区统计数据将表示为`x_bucket{le="y"}`。
- histogram类型的采样必须包含`x_bucket{le="+Inf"}`, 它的值等于x_count的值。
- summary和historam中quantile和le必需按从小到大顺序排列。

示例：

```golang
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363000
http_requests_total{method="post",code="400"}    3 1395066363000

# Escaping in label values:
msdos_file_access_time_seconds{path="C:\\DIR\\FILE.TXT",error="Cannot find file:\n\"FILE.TXT\""} 1.458255915e9

# Minimalistic line:
metric_without_timestamp_and_labels 12.47

# A weird metric from before the epoch:
something_weird{problem="division by zero"} +Inf -3982045

# A histogram, which has a pretty complex representation in the text format:
# HELP http_request_duration_seconds A histogram of the request duration.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.05"} 24054
http_request_duration_seconds_bucket{le="0.1"} 33444
http_request_duration_seconds_bucket{le="0.2"} 100392
http_request_duration_seconds_bucket{le="0.5"} 129389
http_request_duration_seconds_bucket{le="1"} 133988
http_request_duration_seconds_bucket{le="+Inf"} 144320
http_request_duration_seconds_sum 53423
http_request_duration_seconds_count 144320

# Finally a summary, which has a complex representation, too:
# HELP rpc_duration_seconds A summary of the RPC duration in seconds.
# TYPE rpc_duration_seconds summary
rpc_duration_seconds{quantile="0.01"} 3102
rpc_duration_seconds{quantile="0.05"} 3272
rpc_duration_seconds{quantile="0.5"} 4773
rpc_duration_seconds{quantile="0.9"} 9001
rpc_duration_seconds{quantile="0.99"} 76656
rpc_duration_seconds_sum 1.7560473e+07
rpc_duration_seconds_count 2693
```

## Pushgateway ##

[https://github.com/prometheus/pushgateway](https://github.com/prometheus/pushgateway)

![](img/Promethues_pushgateway.png)

Pushgateway是Prometheus生态中一个重要工具，Prometheus采用pull模式，可能由于不在一个子网或者防火墙原因，导致Prometheus无法直接拉取各个target数据。在监控业务数据的时候，需要将不同数据汇总, 由Prometheus统一收集。pushgateway就是为了解决这些问题，但在使用之前，有必要了解一下它的一些弊端：

- 将多个节点数据汇总到pushgateway, 如果pushgateway挂了，受影响比多个target大。
- Prometheus拉取状态up只针对pushgateway, 无法做到对每个节点有效。
- Pushgateway可以持久化推送给它的所有监控数据。

因此，即使你的监控已经下线，prometheus还会拉取到旧的监控数据，需要手动清理pushgateway不要的数据。

## 部署配置 ##

[https://www.cnblogs.com/chenqionghe/p/10494868.html](https://www.cnblogs.com/chenqionghe/p/10494868.html)

[https://studygolang.com/articles/13522?fr=sidebar](https://studygolang.com/articles/13522?fr=sidebar)

[https://songjiayang.gitbooks.io/prometheus/content/configuration/global.html](https://songjiayang.gitbooks.io/prometheus/content/configuration/global.html)

[https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Ctls_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Ctls_config)

**全局配置**：

global属于全局的默认配置，它主要包含4个属性：

- scrape_interval: 抓取间隔，默认为1m
- scrape_timeout: 抓取超时时间，默认为10s
- evaluation_interval: 规则评估间隔，默认为1m
- external_labels: 额外的属性，会添加到拉取的数据并存到数据库中。

```
global:
  scrape_interval:     15s 
  evaluation_interval: 15s 
  scrape_timeout: 10s
  external_labels:
    monitor: 'codelab-monitor'
```

**抓取配置**:

scrape_configs可以有多个，一般来说每个任务（Job）对应一个配置。单个抓取配置的格式如下：

- job_name：任务名称
- honor_labels： 用于解决拉取数据标签有冲突，当设置为true, 以拉取数据为准，否则以服务配置为准。即当抓取回来的采样值的标签值跟服务端配置的不一致时，如果该配置为true，则以抓取回来的为准。否则以服务端的为准，抓取回来的值会保存到一个新标签下，该新标签名在原来的前面加上了“exported_”，比如 exported_job。
- params：数据拉取访问时带的请求参数
- scrape_interval： 拉取时间间隔,默认为对应全局配置
- scrape_timeout: 拉取超时时间,默认为对应全局配置
- metrics_path： 拉取节点的metric 路径,默认为`/metrics`
- scheme： 拉取数据访问协议
- sample_limit： 存储的数据标签个数限制，如果超过限制，该数据将被忽略，不入存储；默认值为0，表示没有限制
- relabel_configs： 拉取数据重置标签配置
- metric_relabel_configs：metric重置标签配置
- static_configs:静态目标配置，配置了该任务要抓取的所有实例，按组配置，包含相同标签的实例可以分为一组，以简化配置。

**relabel_configs**

发生在采集样本数据之前，对Target实例的标签进行重写的机制在Prometheus被称为Relabeling。

一些重要标签：

- __address__：当前Target实例的访问地址<host>:<port>
- __scheme__：采集目标服务访问地址的HTTP Scheme，HTTP或者HTTPS
- __metrics_path__：采集目标服务访问地址的访问路径
- __param_<name>：采集任务目标服务的中包含的请求参数

配置：

```golang
# The source labels select values from existing labels. Their content is concatenated
# using the configured separator and matched against the configured regular expression
# for the replace, keep, and drop actions.
[ source_labels: '[' <labelname> [, ...] ']' ]
 
# Separator placed between concatenated source label values.
[ separator: <string> | default = ; ]
 
# Label to which the resulting value is written in a replace action.
# It is mandatory for replace actions. Regex capture groups are available.
[ target_label: <labelname> ]
 
# Regular expression against which the extracted value is matched.
[ regex: <regex> | default = (.*) ]
 
# Modulus to take of the hash of the source label values.
[ modulus: <uint64> ]
 
# Replacement value against which a regex replace is performed if the
# regular expression matches. Regex capture groups are available.
[ replacement: <string> | default = $1 ]
 
# Action to perform based on regex matching.
[ action: <relabel_action> | default = replace ]

```


**规则配置**

记录规则允许我们把一些经常需要使用并且查询时计算量很大的查询表达式，预先计算并保存到一个新的时序。查询这个新的时序比从原始一个或多个时序实时计算快得多，并且还能够避免不必要的计算。在一些特殊场景下这甚至是必须的，比如仪表盘里展示的各类定时刷新的数据，数据种类多且需要计算非常快。

```
groups:
  [ - <rule_group> ]
```
单个组的配置如下：

```
name: <string> 
[ interval: <duration> | default = global.evaluation_interval ]
rules:
  [ - <rule> ... ]
```

每个组下包含多个规则：

```
record: <string>
expr: <string>
labels:
  [ <labelname>: <labelvalue> ]
```

**告警配置**

alerting主要包含两个个参数：

- alert_relabel_configs: 动态修改alert属性的规则配置。
- alertmanagers: 用于动态发现Alertmanager的配置。


**远程存储**

在metrics的存储这块，prometheus提供了本地存储，即tsdb时序数据库。本地存储的优势就是运维简单，启动prometheus只需一个命令，下面两个启动参数指定了数据路径和保存时间。

- storage.tsdb.path: tsdb数据库路径，默认 data/
- storage.tsdb.retention: 数据保留时间，默认15天

本地存储的一个缺点就是无法将大量的metrics持久化。当然prometheus2.0以后压缩数据能力得到了很大的提升。为了解决单节点存储的限制，prometheus没有自己实现集群存储，而是提供了远程读写的接口，让用户自己选择合适的时序数据库来实现prometheus的扩展性。prometheus通过下面两张方式来实现与其他的远端存储系统对接

- Prometheus按照标准的格式将metrics写到远端存储
- prometheus按照标准格式从远端的url来读取metrics

现在社区已经实现了以下的远程存储方案(以下只支持远程读):

- AppOptics: write
- Chronix: write
- Cortex: read and write
- CrateDB: read and write
- Elasticsearch: write
- Gnocchi: write
- Graphite: write
- InfluxDB: read and write
- OpenTSDB: write
- PostgreSQL/TimescaleDB: read and write
- SignalFx: write

下面几种同时支持远程读写：

- Cortex来源于weave公司，整个架构对prometheus做了上层的封装，用到了很多组件。稍微复杂。
- InfluxDB 开源版不支持集群。对于metrics量比较大的,写入压力大，然后influxdb-relay方案并不是真正的高可用。当然饿了么开源了influxdb-proxy，有兴趣的可以尝试一下。
- CrateDB 基于es。具体了解不多
- TimescaleDB 个人比较中意该方案。传统运维对pgsql熟悉度高，运维靠谱。目前支持 streaming replication方案支持高可用。

配置文件：

-  远程写:

```
url: <string>  //访问地址
[ remote_timeout: <duration> | default = 30s ]  //超时时间，默认30S
write_relabel_configs:                          //标签重置配置, 拉取到的数据，经过重置处理后，发送给远程存储
  [ - <relabel_config> ... ]
basic_auth:
  [ username: <string> ]
  [ password: <string> ]
  [ password_file: <string> ]

[ bearer_token: <string> ]


[ bearer_token_file: /path/to/bearer/token/file ]

tls_config:
  [ <tls_config> ]

[ proxy_url: <string> ]

queue_config:

  [ capacity: <int> | default = 100000 ]

  [ max_shards: <int> | default = 1000 ]

  [ max_samples_per_send: <int> | default = 100]

  [ batch_send_deadline: <duration> | default = 5s ]

  [ max_retries: <int> | default = 10 ]

  [ min_backoff: <duration> | default = 30ms ]

  [ max_backoff: <duration> | default = 100ms ]
```

配置中的`write_relabel_configs`配置项，充分利用了prometheus强大的relabel的功能。可以过滤需要写到远端存储的metrics。

例如：选择指定的metrics。

```
remote_write:
      - url: "http://prometheus-remote-storage-adapter-svc:9201/write"
        write_relabel_configs:
        - action: keep
          source_labels: [__name__]
          regex: container_network_receive_bytes_total|container_network_receive_packets_dropped_total
```

global配置中external_labels，在prometheus的联邦和远程读写的可以考虑设置该配置项，从而区分各个集群。

```
global:
      scrape_interval: 20s
      # The labels to add to any time series or alerts when communicating with
      # external systems (federation, remote storage, Alertmanager).
      external_labels:
        cid: '9'
```

- 远程读

```
url: <string>  // 访问地址

required_matchers:
  [ <labelname>: <labelvalue> ... ]

[ remote_timeout: <duration> | default = 1m ]  //超时时间，默认1m

[ read_recent: <boolean> | default = false ]

basic_auth:
  [ username: <string> ]
  [ password: <string> ]
  [ password_file: <string> ]

[ bearer_token: <string> ]

[ bearer_token_file: /path/to/bearer/token/file ]

tls_config:
  [ <tls_config> ]

[ proxy_url: <string> ]
```

**服务发现**

在Prometheus的配置中，一个最重要的概念就是数据源target，而数据源的配置主要分为静态配置和动态发现, 大致为以下几类：
- static_configs: 静态服务发现
- dns_sd_configs: DNS 服务发现
- file_sd_configs: 文件服务发现
- consul_sd_configs: Consul服务发现
- serverset_sd_configs: Serverset服务发现
- nerve_sd_configs: Nerve服务发现
- marathon_sd_configs: Marathon服务发现
- kubernetes_sd_configs: Kubernetes服务发现
- gce_sd_configs: GCE服务发现
- ec2_sd_configs: EC2服务发现
- openstack_sd_configs: OpenStack服务发现
- azure_sd_configs: Azure服务发现
- triton_sd_configs: Triton服务发现

## Exporter ##

在Prometheus中负责数据汇报的程序统一叫做Exporter, 而不同的Exporter负责不同的业务。 它们具有统一命名格式，即xx_exporter, 例如负责主机信息收集的node_exporter。
[社区上支持的Exporter参考]（https://prometheus.io/docs/instrumenting/exporters/#exporters-and-integrations）

**Node Exporter**

- 默认开启的功能

| 名称 | 说明 | 系统 |
| ------| ------ | ------ |
| arp | 从 `/proc/net/arp` 中收集 ARP 统计信息 | Linux |
| conntrack | 从 `/proc/sys/net/netfilter/` 中收集 conntrack 统计信息 | Linux |
| cpu | 收集 cpu 统计信息 | Darwin, Dragonfly, FreeBSD, Linux |
| diskstats | 从 `/proc/diskstats` 中收集磁盘 I/O 统计信息  | Linux |
| edac | 错误检测与纠正统计信息 | Linux |
| entropy | 可用内核熵信息 | Linux |
| exec | execution 统计信息 | Dragonfly, FreeBSD |
| filefd | 从 `/proc/sys/fs/file-nr` 中收集文件描述符统计信息 | Linux |
| filesystem | 文件系统统计信息，例如磁盘已使用空间 | Darwin, Dragonfly, FreeBSD, Linux, OpenBSD |
| hwmon | 从 `/sys/class/hwmon/` 中收集监控器或传感器数据信息 | Linux |
| infiniband | 从 InfiniBand 配置中收集网络统计信息 | Linux |
| loadavg | 收集系统负载信息 | 	Darwin, Dragonfly, FreeBSD, Linux, NetBSD, OpenBSD, Solaris |
| mdadm | 从 `/proc/mdstat` 中获取设备统计信息 | Linux |
| meminfo | 内存统计信息 | Darwin, Dragonfly, FreeBSD, Linux |
| netdev | 网口流量统计信息，单位 bytes | Darwin, Dragonfly, FreeBSD, Linux, OpenBSD |
| netstat | 从 `/proc/net/netstat` 收集网络统计数据，等同于 `netstat -s` | Linux |
| sockstat | 从 `/proc/net/sockstat` 中收集 socket 统计信息 | Linux |
| stat | 从 `/proc/stat` 中收集各种统计信息，包含系统启动时间，forks, 中断等 | Linux |
| textfile | 通过 `--collector.textfile.directory` 参数指定本地文本收集路径，收集文本信息 | any |
| time | 系统当前时间 | any |
| uname | 通过 `uname` 系统调用, 获取系统信息  | any |
| vmstat | 从 `/proc/vmstat` 中收集统计信息  | Linux |
| wifi | 收集 wifi 设备相关统计数据  | Linux |
| xfs | 收集 xfs 运行时统计信息  | Linux (kernel 4.4+) |
| zfs | 收集 zfs 性能统计信息 | Linux |

- 默认关闭的功能

| 名称 | 说明 | 系统 |
| ------| ------ | ------ |
| bonding | 收集系统配置以及激活的绑定网卡数量 | Linux |
| buddyinfo | 从 `/proc/buddyinfo` 中收集内存碎片统计信息 | Linux |
| devstat | 收集设备统计信息 | Dragonfly, FreeBSD |
| drbd |  收集远程镜像块设备（DRBD）统计信息  | Linux |
| interrupts | 收集更具体的中断统计信息 | Linux，OpenBSD |
| ipvs | 从 `/proc/net/ip_vs` 中收集 IPVS 状态信息，从 `/proc/net/ip_vs_stats` 获取统计信息 | Linux |
| ksmd | 从 `/sys/kernel/mm/ksm` 中获取内核和系统统计信息 | Linux |
| logind | 从 `logind` 中收集会话统计信息 | Linux |
| meminfo_numa | 从 `/proc/meminfo_numa` 中收集内存统计信息 | Linux |
| mountstats | 从 `/proc/self/mountstat` 中收集文件系统统计信息，包括 NFS 客户端统计信息 | Linux |
| nfs | 从 `/proc/net/rpc/nfs` 中收集 NFS 统计信息，等同于 `nfsstat -c` | Linux |
| qdisc | 收集队列推定统计信息 | Linux |
| runit | 收集 runit 状态信息 | any |
| supervisord | 收集 supervisord 状态信息 | any |
| systemd | 从 `systemd` 中收集设备系统状态信息 | Linux |
| tcpstat | 从 `/proc/net/tcp` 和 `/proc/net/tcp6` 收集 TCP 连接状态信息 | Linux |

- 将被废弃功能：

| 名称 | 说明 | 系统 |
| ------| ------ | ------ |
| gmond | 收集 Ganglia 统计信息 | any |
| megacli | 从 MegaCLI 中收集 RAID 统计信息 | Linux |
| ntp | 从 NTP 服务器中获取时钟 | any |

注意：我们可以使用 `--collectors.enabled` 运行参数指定node_exporter收集的功能模块, 如果不指定，将使用默认模块。

- 数据存储

可以利用Prometheus的static_configs来拉取node_exporter的数据。
打开 prometheus.yml 文件, 在 scrape_configs 中添加如下配置：
```
- job_name: "node"
    static_configs:
      - targets: ["127.0.0.1:9100"]
```
重启加载配置，然后到Prometheus Console查询，你会看到node_exporter的数据。

- 常用查询语句

CPU使用率:

```
100 - (avg by (instance) (irate(node_cpu{instance="xxx", mode="idle"}[5m])) * 100)
```

CPU各mode占比率:

```
avg by (instance, mode) (irate(node_cpu{instance="xxx"}[5m])) * 100
```

机器平均负载:

```
node_load1{instance="xxx"} // 1分钟负载
node_load5{instance="xxx"} // 5分钟负载
node_load15{instance="xxx"} // 15分钟负载
```

内存使用率:

```
100 - ((node_memory_MemFree{instance="xxx"}+node_memory_Cached{instance="xxx"}+node_memory_Buffers{instance="xxx"})/node_memory_MemTotal) * 100
```

磁盘使用率:

```
100 - node_filesystem_free{instance="xxx",fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*"} / node_filesystem_size{instance="xxx",fstype!~"rootfs|selinuxfs|autofs|rpc_pipefs|tmpfs|udev|none|devpts|sysfs|debugfs|fuse.*"} * 100
```

或者你也可以直接使用 {fstype="xxx"} 来指定想查看的磁盘信息

网络IO:

```
// 上行带宽
sum by (instance) (irate(node_network_receive_bytes{instance="xxx",device!~"bond.*?|lo"}[5m])/128)

// 下行带宽
sum by (instance) (irate(node_network_transmit_bytes{instance="xxx",device!~"bond.*?|lo"}[5m])/128)
```

网卡出/入包:

```
// 入包量
sum by (instance) (rate(node_network_receive_bytes{instance="xxx",device!="lo"}[5m]))

// 出包量
sum by (instance) (rate(node_network_transmit_bytes{instance="xxx",device!="lo"}[5m]))
```



