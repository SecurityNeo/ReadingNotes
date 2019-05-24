# OpenShift网络 #

[摘自https://www.cnblogs.com/sammyliu/p/10064450.html](https://www.cnblogs.com/sammyliu/p/10064450.html)

## 插件类型 ##

OpenShift使用SDN（软件定义网络）提供集群网络，实现集群POD之间的通信。Openshift SDN使用的是OpenvSwitch（OVS）。有关OVS介绍，可参考[Network/OpenvSwitch(OVS)](https://github.com/SecurityNeo/ReadingNotes/blob/master/Network/OpenvSwitch(OVS).md)

三种插件：

- ovs-subnet：ovs-subnet实现的是一种扁平网络，未实现租户之间的网络隔离，这意味着所有租户之间的pod都可以互访，这使得该实现无法用于绝大多数的生产环境。
- ovs-multitenant：基于OVS和VxLA 等技术实现了项目（project）之间的网络隔离。
- ovs-networkpolicy：介于ovs-subnet和ovs-multitenant之间的一种实现。考虑到ovs-multitenant只是实现了项目级别的网络隔离，这种隔离粒度在一些场景中有些过大，用户没法做更精细的控制，这种需求导致了ovs-networkpolicy的出现。默认地，它和ovs-subnet一样，所有租户之间都没有网络隔离。但是，管理员可以通过定义NetworkPolicy对象来精细地进行网络控制。

当使用ansible部署OpenShift时，默认会启用ovs-subnet。

## Nodes网络 ##

![](img/OpenShift_Network.png)

节点上的主要网络设备：

- br0：OpenShift 创建和管理的 Open vSwitch 网桥, 它会使用 OpenFlow 规则来实现网络隔离和转发。
- vethXXXXX：veth 对，它负责将 pod 的网络命名空间连接到 br0 网桥。
- tun0 ：一OVS 内部端口，它会被分配本机的 pod 子网的网关IP 地址，用于OpenShift pod 以及Docker 容器与集群外部的通信。iptables 的 NAT 规则会作用于tun0。
- docker0：Docker 管理和使用的 linux bridge 网桥，通过 veth 对将不受 OpenShift 管理的Docker 容器的网络地址空间连接到 docker0 上。
- vovsbr/vlinuxbr：将 docker0 和 br0 连接起来的 veth 对，使得Docker 容器能和 OpenShift pod 通信，以及通过 tun0 访问外部网络
- vxlan0：一OVS VXLAN 隧道端点，用于集群内部 pod 之间的网络通信。

