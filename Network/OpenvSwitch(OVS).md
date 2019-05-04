# OpenvSwitch（OVS） #

摘自[https://www.2cto.com/net/201707/655718.html](https://www.2cto.com/net/201707/655718.html)

## 1、什么是Open vSwitch ##

OpenvSwitch简称OVS，正如其官网(http://openvswitch.org/)所述，OVS是一个高质量、多层的虚拟交换软件。它的目的是通过编程扩展支持大规模网络自动化，同时还支持标准的管理接口和协议。

随着虚拟化应用普及，需要部署更多的虚拟化交换机，而费用昂贵的闭源虚拟交换机让用户不堪重负，多层虚拟化软件交换机Open vSwitch由Nicira Networks开发，主要实现代码为可移植的C代码。它遵循Apache 2.0开源代码版权协议，可用于生产环境，支持跨物理服务器分布式管理、扩展编程、大规模网络自动化和标准化接口，实现了和大多数商业闭源交换机功能类似的软件交换机。

OVS官方的定位是要做一个产品级质量的多层虚拟交换机，通过支持可编程扩展来实现大规模的网络自动化。设计目标是方便管理和配置虚拟机网络，检测多物理主机在动态虚拟环境中的流量情况。针对这一目标，OVS具备很强的灵活性。可以在管理程序中作为软件switch运行，也可以直接部署到硬件设备上作为控制层。

如下图，在某台物理服务器中，运行着4台虚拟机，为了将这4台虚拟机在逻辑上组成我们需要的网络架构，于是就虚拟出了2台交换机，组成图中的网络架构。

![](img/ovs_structure.png)

组件：

![](img/ovs_component.png)

- ovsdb-sever 

	OVS的数据库服务器，用来存储虚拟交换机的配置信息。它与manager和ovs-vswitchd交换信息使用了OVSDB(JSON-RPC)的方式。
 
- ovs-vswitchd

	OVS的核心部件，它和上层controller通信遵从openflow协议，它与ovsdb-server通信使用OVSDB协议，它和内核模块通过netlink通信，它支持多个独立的datapath（网桥），它通过更改flow table实现了绑定，和VLAN等功能。
 
- ovs kernel module

	OVS的内核模块，处理包交换和隧道，缓存flow，如果在内核的缓存中找到转发规则则转发，否则发向用户空间去处理。

有个OVS大概的了解，我们看看OVS到底由哪些模块组成，以及个模块的作用。如图：

![](img/ovs_component_detail.png)

- ovs-vswitchd： 主要模块，实现switch的daemon，包括一个支持流交换的Linux内核模块；
- ovsdb-server： 轻量级数据库服务器，提供ovs-vswitchd获取配置信息，例如vlan、port等信息；
- ovs-brcompatd： 让ovs-vswitch替换linux bridge，包括获取bridge ioctls的Linux内核模块；
- ovs-dpctl：用来配置switch内核模块；
- ovs-vsctl： 查询和更新ovs-vswitchd的配置；
- ovs-appctl： 发送命令消息，运行相关daemon；
- ovs-ofctl： 查询和控制OpenFlow交换机和控制器；
- ovs-openflowd：一个简单的OpenFlow交换机；
- ovs-controller：一个简单的OpenFlow控制器；
- ovs-pki：OpenFlow交换机创建和管理公钥框架；
- ovs-tcpundump：tcpdump的补丁，解析OpenFlow的消息；
- ovs-bugtool：管理openvswitch的bug信息。

![](img/ovs_inside.png)

