## Kubernetes Operator ##

Operator是CoreOS推出的旨在简化复杂有状态应用管理的框架，它是一个感知应用状态的控制器，通过扩展Kubernetes API来自动创建、管理和配置应用实例。Kubernetes 1.7版本引入了自定义控制器（[CRD](https://github.com/SecurityNeo/ReadingNotes/blob/master/PaaS/Kubernetes/Kubernetes%20CRDs.md)）的概念，该功能可以让开发人员扩展添加新功能，更新现有的功能，并且可以自动执行一些管理任务，这些自定义的控制器就像Kubernetes原生的组件一样，Operator直接使用Kubernetes API进行开发，也就是说他们可以根据这些控制器内部编写的自定义规则来监控集群、更改Pods/Services、对正在运行的应用进行扩缩容。

## Operator Framework ##

Operator Framework是CoreOS的一个开源项目，提供开发人员和Kubernetes运行时工具，帮助我们加速operator的开发。operator Framework包括下面三个部分：

- Operator SDK
使开发人员能够基于他们的专业知识构建operator，而无需了解Kubernetes API的复杂性。

- Operator lifecycle manager
监督Kubernetes集群中运行的所有operator（及其相关服务）的安装，更新和管理整个生命周期。

- Operator Metering
Operator Metering（[2018年加入](http://coreos.com/blog/introducing-operator-framework-metering)）：为提供专业服务的operator启用使用情况报告。

![](img/Operator_SDK.png)


