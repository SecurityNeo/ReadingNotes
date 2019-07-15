# Prometheus简介 #

[https://prometheus.io/](https://prometheus.io/)
[https://github.com/prometheus](https://github.com/prometheus)

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

## 部署配置 ##

[https://www.cnblogs.com/chenqionghe/p/10494868.html](https://www.cnblogs.com/chenqionghe/p/10494868.html)

[https://songjiayang.gitbooks.io/prometheus/content/configuration/global.html](https://songjiayang.gitbooks.io/prometheus/content/configuration/global.html)

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

