# Skywalking #

[https://github.com/apache/skywalking](https://github.com/apache/skywalking)
[Skywalking Helm](https://github.com/apache/skywalking-kubernetes)

SkyWalking是一款优秀的国产APM工具，包括了分布式追踪、性能指标分析、应用和服务依赖分析等。

**架构**

![](img/skywalking_arch.png)

[https://www.jianshu.com/p/2fd56627a3cf](https://www.jianshu.com/p/2fd56627a3cf)
[https://www.liangzl.com/get-article-detail-37412.html](https://www.liangzl.com/get-article-detail-37412.html)

SkyWalking的核心是数据分析和度量结果的存储平台，通过HTTP或gRPC方式向SkyWalking Collecter提交分析和度量数据，SkyWalking Collecter对数据进行分析和聚合，存储到Elasticsearch、H2、MySQL、TiDB等其一即可，最后可以通过SkyWalking UI的可视化界面对最终的结果进行查看。Skywalking支持从多个来源和多种格式收集数据：多种语言的Skywalking Agent 、Zipkin v1/v2 、Istio勘测、Envoy度量等数据格式。


- **Skywalking Agent**：使用Javaagent做字节码植入，无侵入式的收集，并通过HTTP或者gRPC方式发送数据到Skywalking Collector。
- **Skywalking Collector**：链路数据收集器，对agent传过来的数据进行整合分析处理并落入相关的数据存储中。
- **Storage**：Skywalking的存储，在6.x版本中支持以ElasticSearch、Mysql、TiDB、H2、作为存储介质进行数据存储。
- **UI**：Web可视化平台，用来展示落地的数据。