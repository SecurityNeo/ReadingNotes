# K8S的kube-proxy模式 #

Kube-proxy支持三种工作模式：userspace、iptables和ipvs。userspace模式有较大的性能损耗，几乎被淘汰。现在用得比较多的是iptables模式，从K8S的1.8版本开始，kube-proxy引入了IPVS模式，并于1.11版本GA。IPVS主要是想解决iptables模式在超大规模K8S集群体现出来的性能问题。

## iptables ##

kube-proxy对iptables的链进行了扩充，主要定义了以下几条链： [源码](https://github.com/kubernetes/kubernetes/blob/master/pkg/proxy/iptables/proxier.go)

注意：K8S不是只创建这几条链，会有很多扩展链。

```go
const (
	// the services chain
	kubeServicesChain utiliptables.Chain = "KUBE-SERVICES"

	// the external services chain
	kubeExternalServicesChain utiliptables.Chain = "KUBE-EXTERNAL-SERVICES"

	// the nodeports chain
	kubeNodePortsChain utiliptables.Chain = "KUBE-NODEPORTS"

	// the kubernetes postrouting chain
	kubePostroutingChain utiliptables.Chain = "KUBE-POSTROUTING"

	// KubeMarkMasqChain is the mark-for-masquerade chain
	KubeMarkMasqChain utiliptables.Chain = "KUBE-MARK-MASQ"

	// KubeMarkDropChain is the mark-for-drop chain
	KubeMarkDropChain utiliptables.Chain = "KUBE-MARK-DROP"

	// the kubernetes forward chain
	kubeForwardChain utiliptables.Chain = "KUBE-FORWARD"
)
```

### KUBE-MARK-MASQ与KUBE-MARK-DROP链 ###

这两条链主要是为报文打上标签，便于后续做NAT或者DROP。集群内的相关规则如下：

```shell
-A KUBE-MARK-DROP -j MARK --set-xmark 0x8000/0x8000
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
-A KUBE-POSTROUTING -m comment --comment "kubernetes service traffic requiring SNAT" -m mark --mark 0x4000/0x4000 -j MASQUERADE
-A KUBE-FIREWALL -m comment --comment "kubernetes firewall for dropping marked packets" -m mark --mark 0x8000/0x8000 -j DROP
```

注：将相关规则放到了这一块，实际上不长这样

**KUBE-MARK-MASQ**

主要工作就是给报文打上标签“0x4000/0x4000”，便于后续报文出去时自动做SNAT。

注：MASQUERADE实际上是一种特殊的SNAT，当报文在出网卡时会自动获取网卡的IP，并修改报文的源IP为网卡IP。后续在iptables相关笔记中补充。

**KUBE-MARK-DROP**

主要工作是给报文打上标签“0x8000/0x8000”，后续报文会在KUBE-FIREWALL处被DROP掉。比如创建了SVC，但是没有对应的endpoint，访问这种SVC的报文就会走此链被无情DROP掉。


### KUBE-SERVICES链 ####

K8S会在`KUBE-SERVICES`链中为每个service创建规则，并为其创建新的链，名称以`KUBE-SVC-`开头。另外会为挂载在此service下边的endpoint创建新的链，名称以`KUBE-SEP-`开头， 看例子理解：

```shell
$ kubectl get pod -o wide
NAME                     READY   STATUS    RESTARTS   AGE    IP           NODE       NOMINATED NODE   READINESS GATESnginx-7bdbbfb5cf-88knd   1/1     Running   0          118s   172.18.0.5   minikube   <none>           <none>nginx-7bdbbfb5cf-p2h2s   1/1     Running   0          118s   172.18.0.4   minikube   <none>           <none>
$ kubectl get svc | grep nginx
nginx-svc    ClusterIP   10.102.40.108   <none>        80/TCP    9m20s
```

KUBE-SERVICES中相关规则如下：

```shell
$ iptables -t nat -S KUBE-SERVICES | grep nginx
-A KUBE-SERVICES -d 10.102.40.108/32 -p tcp -m comment --comment "default/nginx-svc: cluster IP" -m tcp --dport 80 -j KUBE-SVC-R2VK7O5AFVLRAXSH
```

再看看扩展链“KUBE-SVC-R2VK7O5AFVLRAXSH”：

```shell
$ iptables -t nat -S KUBE-SVC-R2VK7O5AFVLRAXSH
-N KUBE-SVC-R2VK7O5AFVLRAXSH
-A KUBE-SVC-R2VK7O5AFVLRAXSH -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-R6MNIMRVL7R2ID27
-A KUBE-SVC-R2VK7O5AFVLRAXSH -j KUBE-SEP-SONECWMN2373EQTO
```

可以看到，service -> pod的负载就在KUBE-SVC-XXXX链规则内完成的。继续往下看：

```
$ iptables -t nat -S KUBE-SEP-R6MNIMRVL7R2ID27
-N KUBE-SEP-R6MNIMRVL7R2ID27
-A KUBE-SEP-R6MNIMRVL7R2ID27 -s 172.18.0.4/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-R6MNIMRVL7R2ID27 -p tcp -m tcp -j DNAT --to-destination 172.18.0.4:80
$ iptables -t nat -S KUBE-SEP-SONECWMN2373EQTO
-N KUBE-SEP-SONECWMN2373EQTO
-A KUBE-SEP-SONECWMN2373EQTO -s 172.18.0.5/32 -j KUBE-MARK-MASQ
-A KUBE-SEP-SONECWMN2373EQTO -p tcp -m tcp -j DNAT --to-destination 172.18.0.5:80
```

当service类型为nodePort时，KUBE-NODEPORTS链会将数据包导入“KUBE-SVC-XXX”链。其余跟ClusterIP类型的一致。看个示例：

```shell
$ kubectl get pod -o wide
NAME                     READY   STATUS    RESTARTS   AGE    IP           NODE       NOMINATED NODE   READINESS GATES
nginx-57c6bff7f6-7tg9s   1/1     Running   0          9m3s   172.18.0.5   minikube   <none>           <none>
nginx-57c6bff7f6-kxk6v   1/1     Running   0          9m3s   172.18.0.4   minikube   <none>           <none>
$ kubectl get svc nginx-svc
NAME        TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
nginx-svc   NodePort   10.110.25.59   <none>        80:30168/TCP   4m31s
```

我们可以看到在KUBE-NODEPORTS链中会增加两条规则：

```shell
$ iptables -t nat -S KUBE-NODEPORTS-N KUBE-NODEPORTS
-A KUBE-NODEPORTS -p tcp -m comment --comment "default/nginx-svc:" -m tcp --dport 30168 -j KUBE-MARK-MASQ
-A KUBE-NODEPORTS -p tcp -m comment --comment "default/nginx-svc:" -m tcp --dport 30168 -j KUBE-SVC-R2VK7O5AFVLRAXSH
```

注意：实际上系统是优先处理的ClusterIP类型流量，我们可以看到在KUBE-SERVICES链的最后一个才将流量转发至KUBE-NODEPORTS链中。

```shell
$ iptables -t nat -S KUBE-SERVICES
-N KUBE-SERVICES
......
-A KUBE-SERVICES -d 10.110.25.59/32 -p tcp -m comment --comment "default/nginx-svc: cluster IP" -m tcp --dport 80 -j KUBE-SVC-R2VK7O5AFVLRAXSH
......
-A KUBE-SERVICES -m comment --comment "kubernetes service nodeports; NOTE: this must be the last rule in this chain" -m addrtype --dst-type LOCAL -j KUBE-NODEPORTS
```

## ipvs ##

ipvs模式下并不是完全不使用iptables，ipvs只有DNAT功能，所以跟多功能仍然依靠iptables来完成。在ipvs模式下，kube-proxy会在节点上新添加一个网卡，一般名称为`kube-ipvs0`，集群内所有服务的ClusterIP都会设置在此网卡上，如下：

```shell
[root@VM-0-16-centos ~]# ip addr show kube-ipvs0
4: kube-ipvs0: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default
    link/ether 62:87:11:34:94:61 brd ff:ff:ff:ff:ff:ff
    inet 172.16.254.132/32 brd 172.16.254.132 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
    inet 172.16.252.1/32 brd 172.16.252.1 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
    inet 172.16.255.254/32 brd 172.16.255.254 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
    inet 172.16.252.40/32 brd 172.16.252.40 scope global kube-ipvs0
       valid_lft forever preferred_lft forever
```

对于service -> pod的流量转发由ipvs完成，相关规则可通过命令`ipvsadm -Ln`查看，如下：
```shell
[root@VM-0-16-centos ~]# ipvsadm -Ln
IP Virtual Server version 1.2.1 (size=4096)
Prot LocalAddress:Port Scheduler Flags
  -> RemoteAddress:Port           Forward Weight ActiveConn InActConn
TCP  172.16.252.40:80 rr
  -> 172.16.0.196:80              Masq    1      0          0
  -> 172.16.0.197:80              Masq    1      0          0
```

梳理一下ClusterIP和NodePort两种service类型的流量转发：

**ClusterIP类型**


```shell
[root@10-206-0-3 ~]# kubectl get pod -o wide
NAME                    READY   STATUS    RESTARTS   AGE   IP             NODE          NOMINATED NODE   READINESS GATES
nginx-9d776f4cf-5psws   1/1     Running   0          78s   172.16.0.197   10.206.0.16   <none>           <none>
nginx-9d776f4cf-vcq86   1/1     Running   0          89s   172.16.0.196   10.206.0.16   <none>           <none>
[root@10-206-0-3 ~]# kubectl get svc nginx
NAME    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
nginx   ClusterIP   172.16.252.40   <none>        80/TCP    2m32s
[root@10-206-0-3 ~]# kubectl get ep nginx
NAME    ENDPOINTS                         AGE
nginx   172.16.0.196:80,172.16.0.197:80   2m34s
```

梳理一下数据包流向：

1、数据包首先进入NAT表的PREROUTING链，然后将所有请求都交由KUBE-SERVICES链。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S PREROUTING
-P PREROUTING ACCEPT
-A PREROUTING -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
```

2、在KUBE-SERVICES链上，如果请求的目的ip和port在KUBE-CLUSTER-IP对应的ipset里面，则会命中规则继续跳往KUBE-MARK-MASQ链。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-SERVICES
-N KUBE-SERVICES
-A KUBE-SERVICES -m comment --comment "Kubernetes service cluster ip + port for masquerade purpose" -m set --match-set KUBE-CLUSTER-IP src,dst -j KUBE-MARK-MASQ
-A KUBE-SERVICES -m addrtype --dst-type LOCAL -j KUBE-NODE-PORT
-A KUBE-SERVICES -m set --match-set KUBE-CLUSTER-IP dst,dst -j ACCEPT
```

注：`-m set --match-set`是iptables的一种扩展模式，依靠ipset大大降低iptables规则数目。

```shell
[root@VM-0-16-centos ~]# ipset list KUBE-CLUSTER-IP
Name: KUBE-CLUSTER-IP
Type: hash:ip,port
Revision: 5
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 440
References: 2
Number of entries: 5
Members:
172.16.255.254,udp:53
172.16.252.1,tcp:443
172.16.254.132,tcp:443
172.16.255.254,tcp:53
172.16.252.40,tcp:80
```

3、数据会到KUBE-MARK-MASQ链上过一遍，主要就是为其打上标签，方便后续做自动SNAT。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-MARK-MASQ
-N KUBE-MARK-MASQ
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
```

注：标签不一定是`0x4000/0x4000`

4、接下来数据会进入filter表的INPUT链，在此链上会将所有数据转向KUBE-FIREWALL链。

```
[root@VM-0-16-centos ~]# iptables -S INPUT
-P INPUT ACCEPT
-A INPUT -j KUBE-FIREWALL
```

5、在KUBE-FIREWALL链上，如果发现数据包打了0x8000/0x8000就将包DROP掉。

```shell
[root@VM-0-16-centos ~]# iptables -S KUBE-FIREWALL
-N KUBE-FIREWALL
-A KUBE-FIREWALL -m comment --comment "kubernetes firewall for dropping marked packets" -m mark --mark 0x8000/0x8000 -j DROP
-A KUBE-FIREWALL ! -s 127.0.0.0/8 -d 127.0.0.0/8 -m comment --comment "block incoming localnet connections" -m conntrack ! --ctstate RELATED,ESTABLISHED,DNAT -j DROP
```

注：与iptables模式一样，异常数据包会在此处丢掉。比如创建了SVC，但是没有对应的endpoint，访问这种SVC的报文就会走此链被无情DROP掉。

6、ipvs工作在INPUT链上，在其做完DNAT之后数据库会直接转发到KUBE-POSTROUTING链上。进入KUBE-POSTROUTING链，对打了标记`0x4000/0x4000`的数据包做SNAT转换。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-POSTROUTING
-N KUBE-POSTROUTING
-A KUBE-POSTROUTING -m comment --comment "kubernetes service traffic requiring SNAT" -m mark --mark 0x4000/0x4000 -j MASQUERADE
-A KUBE-POSTROUTING -m comment --comment "Kubernetes endpoints dst ip:port, source ip for solving hairpin purpose" -m set --match-set KUBE-LOOP-BACK dst,dst,src -j MASQUERADE
```

**NodePort**

```shell
[root@10-206-0-3 ~]# kubectl get pod -o wide
NAME                    READY   STATUS    RESTARTS   AGE   IP             NODE          NOMINATED NODE   READINESS GATES
nginx-9d776f4cf-5psws   1/1     Running   0          23m   172.16.0.197   10.206.0.16   <none>           <none>
nginx-9d776f4cf-vcq86   1/1     Running   0          24m   172.16.0.196   10.206.0.16   <none>           <none>
[root@10-206-0-3 ~]# kubectl get svc nginx
NAME    TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
nginx   NodePort   172.16.252.40   <none>        80:32684/TCP   25m
[root@10-206-0-3 ~]# kubectl get ep nginx
NAME    ENDPOINTS                         AGE
nginx   172.16.0.196:80,172.16.0.197:80   25m
```

梳理一下数据包流向：

1、与ClusterIP一样，数据包首先进入NAT表的PREROUTING链，然后将所有请求都交由KUBE-SERVICES链。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S PREROUTING
-P PREROUTING ACCEPT
-A PREROUTING -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
```

2、在KUBE-SERVICES链上，如果进来的流量是通过NodePort访问的，则会命中`-A KUBE-SERVICES -m addrtype --dst-type LOCAL -j KUBE-NODE-PORT`规则。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-SERVICES
-N KUBE-SERVICES
-A KUBE-SERVICES -m comment --comment "Kubernetes service cluster ip + port for masquerade purpose" -m set --match-set KUBE-CLUSTER-IP src,dst -j KUBE-MARK-MASQ
-A KUBE-SERVICES -m addrtype --dst-type LOCAL -j KUBE-NODE-PORT
-A KUBE-SERVICES -m set --match-set KUBE-CLUSTER-IP dst,dst -j ACCEPT
```

3、在`KUBE-NODE-PORT`链上就将流量转入`KUBE-MARK-MASQ`链，也是为了为其打上标签`0x4000/0x4000`，方便后续做自动SNAT。

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-NODE-PORT
-N KUBE-NODE-PORT
-A KUBE-NODE-PORT -p tcp -m comment --comment "Kubernetes nodeport TCP port for masquerade purpose" -m set --match-set KUBE-NODE-PORT-TCP dst -j KUBE-MARK-MASQ
```

```shell
[root@VM-0-16-centos ~]# iptables -t nat -S KUBE-MARK-MASQ
-N KUBE-MARK-MASQ
-A KUBE-MARK-MASQ -j MARK --set-xmark 0x4000/0x4000
```

注：此处也使用了iptables扩展。


```shell
[root@VM-0-16-centos ~]# ipset list KUBE-NODE-PORT-TCP
Name: KUBE-NODE-PORT-TCP
Type: bitmap:port
Revision: 3
Header: range 0-65535
Size in memory: 8300
References: 1
Number of entries: 1
Members:
32684
```

后续的流程跟ClusterIP一致。