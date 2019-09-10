# Prometheus Operator  #

[https://github.com/coreos/prometheus-operator](https://github.com/coreos/prometheus-operator)

[https://blog.csdn.net/ygqygq2/article/details/83655552](https://blog.csdn.net/ygqygq2/article/details/83655552)

Kubernetes的Prometheus Operator为Kubernetes服务和Prometheus实例的部署和管理提供了简单的监控定义。Prometheus Operator的功能：

- 在Kubernetes namespace中更容易启动一个Prometheus实例，一个特定的应用程序或团队更容易使用Operator。
- 配置Prometheus的基础东西，比如在Kubernetes的本地资源versions, persistence, retention policies, 和replicas。
- 基于常见的Kubernetes label查询，自动生成监控target 配置；不需要学习普罗米修斯特定的配置语言。

架构：

![](img/Prometheus_Operator_Arch.png)

- Operator： Operator资源会根据自定义资源（Custom Resource Definition/CRDs）来部署和管理Prometheus Server，同时监控这些自定义资源事件的变化来做相应的处理，是整个系统的控制中心。
- Prometheus： Prometheus资源是声明性地描述Prometheus部署的期望状态。
- Prometheus Server： Operator根据自定义资源Prometheus类型中定义的内容而部署的Prometheus Server集群，这些自定义资源可以看作是用来管理Prometheus Server集群的StatefulSets资源。
- ServiceMonitor： ServiceMonitor也是一个自定义资源，它描述了一组被Prometheus监控的targets列表。该资源通过Labels来选取对应的Service Endpoint，让Prometheus Server通过选取的Service来获取Metrics信息。
- Service： Service资源主要用来对应Kubernetes集群中的Metrics Server Pod，来提供给ServiceMonitor选取让Prometheus Server来获取信息。简单的说就是Prometheus监控的对象，例如Node Exporter Service、Mysql Exporter Service等等。
- Alertmanager： Alertmanager也是一个自定义资源类型，由Operator根据资源描述内容来部署Alertmanager集群。

[Prometheus Operator API文档](https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md)

