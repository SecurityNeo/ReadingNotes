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

再看看扩展链“KUBE-SVC-4CRUJHTV5RT5YMFY”：

```shell
$ iptables -t nat -S KUBE-SVC-R2VK7O5AFVLRAXSH
-N KUBE-SVC-R2VK7O5AFVLRAXSH
-A KUBE-SVC-R2VK7O5AFVLRAXSH -m statistic --mode random --probability 0.50000000000 -j KUBE-SEP-R6MNIMRVL7R2ID27
-A KUBE-SVC-R2VK7O5AFVLRAXSH -j KUBE-SEP-SONECWMN2373EQTO
```

继续往下看：

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


## ipvs ##