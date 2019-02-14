# Kubernetes服务质量（QoS）#

摘自[http://dockone.io/article/2592](http://dockone.io/article/2592)

QoS的英文全称为"Quality of Service",中文名为"服务质量"。在Kubernetes中，每个Pod都会被标注一个Qos级别，Kubernetes针对不同服务质量的预期，通过QoS来对Pod进行服务质量管理。其中有两个指标，一个是CPU，一个是内存，体现在配置上就是`requests`和`limits`两种类型。


## `requests`和`limits` ##

`requests`申请范围从0到Node的最大配置，而`limits`申请范围是从`requests`到无限，即

`0 <= requests <=Node Allocatable`

`requests <= limits <= Infinity`

- 对于CPU，如果容器中服务使用的CPU超过设置的limits，容器不会被kill掉但会被限制。如果没有设置limits，容器可以使用全部空闲的CPU资源。

- 对于内存，当一个容器使用内存超过了设置的limits，Pod中container的进程会被kernel因OOM kill掉。当container因为OOM被kill掉时，系统倾向于在其原所在的机器上重启该container或在本机上重新创建一个Pod。





