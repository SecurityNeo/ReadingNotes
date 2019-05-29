# OpenShift Router和Route #

[https://www.cnblogs.com/sammyliu/p/10013461.html](https://www.cnblogs.com/sammyliu/p/10013461.html)

在OpenShift中，Router是路由器，Route是路由器中配置的路由，Router跟Ingress实际上是一类资源。

OpenShift中外部访问Pod服务与集群内访问的流量图大致如下：
![](img/OpenShift_Router01.png)

## Router的部署 ##

使用ansible采用默认配置部署OpenShift集群时，在集群Infra节点上，会以`Host networking`方式运行一个HAProxy的pod，它会在所有网卡的80和443端口上进行监听。

```
[root@infra-node3 cloud-user]# netstat -lntp | grep haproxy
tcp        0      0 127.0.0.1:10443         0.0.0.0:*               LISTEN      583/haproxy         
tcp        0      0 127.0.0.1:10444         0.0.0.0:*               LISTEN      583/haproxy         
tcp        0      0 0.0.0.0:80              0.0.0.0:*               LISTEN      583/haproxy         
tcp        0      0 0.0.0.0:443             0.0.0.0:*               LISTEN      583/haproxy
```

其中，172.0.0.1上的10443和10444端口供HAproxy自己使用。

OpenShift HAProxy Router支持两种部署方式：

- 一种是常见的单Router服务部署，它有一个或多个实例（pod），分布在多个节点上，负责整个集群上部署的服务的对外访问。
- 另一种是分片（sharding）部署。此时，会有多个Router服务，每个Router服务负责指定的若干project，两者之间采用标签（label）进行映射。这是为了解决单个Router的性能不够问题而提出的解决方案。

可以通过`oc adm router`命令创建router服务:

```
[root@master1 cloud-user]# oc adm router router2 --replicas=1 --service-account=router
info: password for stats user admin has been set to J3YyPjlbqf
--> Creating router router2 ...
    warning: serviceaccounts "router" already exists
    clusterrolebinding.authorization.openshift.io "router-router2-role" created
    deploymentconfig.apps.openshift.io "router2" created
    service "router2" created
--> Success
```
