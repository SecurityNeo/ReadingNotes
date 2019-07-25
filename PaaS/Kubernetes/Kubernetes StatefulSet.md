## StatefulSet ##

[https://blog.51cto.com/newfly/2140004](https://blog.51cto.com/newfly/2140004)

StatefulSet本质上是Deployment的一种变体，在v1.9版本中已成为GA版本，它为了解决有状态服务的问题，它所管理的Pod拥有固定的Pod名称，启停顺序，在StatefulSet中，Pod名字称为网络标识(hostname)，还必须要用到共享存储。在Deployment中，与之对应的服务是service，而在StatefulSet中与之对应的headless service，与service的区别就是它没有Cluster IP，解析它的名称时将返回该Headless Service对应的全部Pod的Endpoint列表。除此之外，StatefulSet在Headless Service的基础上又为StatefulSet控制的每个Pod副本创建了一个DNS域名，这个域名的格式为：

```
$(podname).(headless server name)   
FQDN： $(podname).(headless server name).namespace.svc.cluster.local
```

**稳定的网络标识**

StatefulSet中反复强调的“稳定的网络标识”，主要指Pods的hostname以及对应的DNS Records。


- **HostName：** StatefulSet的Pods的hostname按照这种格式生成：$(statefulset name)-$(ordinal)， ordinal从0 ~ N-1(N为期望副本数)。

	StatefulSet Controller在创建pods时，会给pod加上一个pod name label：statefulset.kubernetes.io/pod-name, 然后设置到Pod的pod name和hostname中。我们可以创建独立的Service匹配到这个指定的pod，然后方便我们单独对这个pod进行debug等处理。

- **DNS Records：**

	- Headless Service的DNS解析：$(service name).$(namespace).svc.cluster.local 通过DNS RR解析到后端其中一个Pod。SRV Records只包含对应的Running and Ready的Pods，不Ready的Pods不会在对应的SRV Records中。
	- Pod的DNS解析：$(hostname).$(service name).$(namespace).svc.cluster.local解析到对应hostname的Pod。

**Statefulset的启停顺序**

- 有序部署：部署StatefulSet时，如果有多个Pod副本，它们会被顺序地创建（从0到N-1）并且，在下一个Pod运行之前所有之前的Pod必须都是Running和Ready状态。
- 有序删除：当Pod被删除时，它们被终止的顺序是从N-1到0。
- 有序扩展：当对Pod执行扩展操作时，与部署一样，它前面的Pod必须都处于Running和Ready状态。

**Statefulset Pod管理策略**

	在v1.7以后，通过允许修改Pod排序策略，同时通过`.spec.podManagementPolicy`字段确保其身份的唯一性。

	- **OrderedReady**：默认设置，参照**Statefulset的启停顺序**所述。
	- **Parallel**：告诉StatefulSet控制器并行启动或终止所有Pod，并且在启动或终止另一个Pod之前不等待前一个Pod变为Running and Ready或完全终止。

**更新策略**

	在Kubernetes 1.7及更高版本中，通过`.spec.updateStrategy`字段允许配置或禁用Pod、labels、source request/limits、annotations自动滚动更新功能。

	- **OnDelete**：通过`.spec.updateStrategy.type`字段设置为OnDelete，StatefulSet控制器不会自动更新StatefulSet中的Pod。用户必须手动删除Pod，以使控制器创建新的Pod。
	- **RollingUpdate**：通过`.spec.updateStrategy.type`字段设置为RollingUpdate，实现了Pod的自动滚动更新，如果`.spec.updateStrategy`未指定，则此为默认策略。 StatefulSet控制器将删除并重新创建StatefulSet中的每个Pod。它将以Pod终止（从最大序数到最小序数）的顺序进行，一次更新每个Pod。在更新下一个Pod之前，必须等待这个Pod Running and Ready。
	- **Partitions**：通过指定`.spec.updateStrategy.rollingUpdate.partition`来对 RollingUpdate 更新策略进行分区，如果指定了分区，则当StatefulSet的`.spec.template`更新时，具有大于或等于分区序数的所有Pod将被更新。具有小于分区的序数的所有Pod将不会被更新，即使删除它们也将被重新创建。如果StatefulSet的`.spec.updateStrategy.rollingUpdate.partition`大于其`.spec.replicas`，则其`.spec.template`的更新将不会传播到Pod。在大多数情况下，不需要使用分区。

## 实现原理 ##

[https://draveness.me/kubernetes-statefulset](https://draveness.me/kubernetes-statefulset)

StatefulSet控制器主要由`StatefulSetController`、`StatefulSetControl`和`StatefulPodControl`三个组件协作来完成StatefulSet的管理，StatefulSetController会同时从`PodInformer`和`ReplicaSetInformer`中接受增删改事件并将事件推送到队列中：

![](img/kubernetes_statefulset.png)

控制器StatefulSetController会在Run方法中启动多个Goroutine协程，这些协程会从队列中获取待处理的StatefulSet资源进行同步

	
## 代码分析 ##

摘自[https://segmentfault.com/a/1190000019488735](https://segmentfault.com/a/1190000019488735)和[https://my.oschina.net/jxcdwangtao/blog/1784739?fromerr=XjbBhQrv](https://my.oschina.net/jxcdwangtao/blog/1784739?fromerr=XjbBhQrv)

