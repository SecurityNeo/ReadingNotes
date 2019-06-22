# Consul #

[](https://www.consul.io/)

Consul是HashiCorp出品的开源服务发现工具，Consul提供了诸如服务发现，健康检查，KV数据库等功能，方便构建服务集群。

主要特性：

- **服务发现**：Consul提供了通过DNS或者HTTP接口的方式来注册服务和发现服务。一些外部的服务通过Consul很容易的找到它所依赖的服务。
- **健康检查**: Consul客户端可用提供任意数量的健康检查,指定一个服务(比如:webserver是否返回了200 OK 状态码)或者使用本地节点(比如:内存使用是否大于90%). 这个信息可由operator用来监视集群的健康.被服务发现组件用来避免将流量发送到不健康的主机. 
- **Key/Value存储**:应用程序可用根据自己的需要使用Consul的层级的Key/Value存储.比如动态配置,功能标记,协调,领袖选举等等,简单的HTTP API让他更易于使用. 
- **多数据中心**: Consul支持开箱即用的多数据中心.这意味着用户不需要担心需要建立额外的抽象层让业务扩展到多个区域。

相关术语：
[](https://www.cnblogs.com/lsf90/p/6021465.html)

- **Agent**： agent是一直运行在Consul集群中每个成员上的守护进程。通过运行consul agent来启动。agent可以运行在client或者server模式。指定节点作为client或者server是非常简单的，除非有其他agent实例。所有的agent都能运行DNS或者HTTP接口，并负责运行时检查和保持服务同步。
- **Client**： 一个Client是一个转发所有RPC到server的代理。这个client是相对无状态的。client唯一执行的后台活动是加入LAN gossip池。这有一个最低的资源开销并且仅消耗少量的网络带宽。
- **Server**： 一个server是一个有一组扩展功能的代理，这些功能包括参与Raft选举，维护集群状态，响应RPC查询，与其他数据中心交互WAN gossip和转发查询给leader或者远程数据中心。
- **DataCenter**： 虽然数据中心的定义是显而易见的，但是有一些细微的细节必须考虑。例如，在EC2中，多个可用区域被认为组成一个数据中心？我们定义数据中心为一个私有的，低延迟和高带宽的一个网络环境。这不包括访问公共网络，但是对于我们而言，同一个EC2中的多个可用区域可以被认为是一个数据中心的一部分。
- **Consensus**： 在我们的文档中，我们使用Consensus来表明就leader选举和事务的顺序达成一致。由于这些事务都被应用到有限状态机上，Consensus暗示复制状态机的一致性。
- **Gossip**： Consul建立在Serf的基础之上，它提供了一个用于多播目的的完整的gossip协议。Serf提供成员关系，故障检测和事件广播。更多的信息在gossip文档中描述。这足以知道gossip使用基于UDP的随机的点到点通信。
- **LAN Gossip**： 它包含所有位于同一个局域网或者数据中心的所有节点。
- **WAN Gossip**： 它只包含Server。这些server主要分布在不同的数据中心并且通常通过因特网或者广域网通信。
- **RPC**： 远程过程调用。这是一个允许client请求server的请求/响应机制。

架构：

![](img/Consul_Arch.png)

　　向Consul提供服务的每个节点都运行一个Consul代理。 发现其他服务或获取/设置键/值数据不需要运行代理。 代理负责健康检查节点上的服务以及节点本身。代理与一个或多个Consul服务器通信。Consul 服务器是数据存储和复制的地方。 服务器自己选出一个 leader。 虽然Consul可以在一台服务器上运行，但推荐使用3到5台来避免数据丢失的情况。 每个数据中心都建议使用一组Consul服务器。需要发现其他服务或节点的基础架构组件可以查询任何Consul服务器或任何Consul代理。 代理自动将查询转发到服务器。每个数据中心都运行Consul服务器集群。 当跨数据中心服务发现或配置请求时，本地Consul服务器将请求转发到远程数据中心并返回结果。

默认端口：

- 8300(tcp): Server RPC，server 用于接受其他 agent 的请求
- 8301(tcp,udp): Serf LAN，数据中心内gossip交换数据用
- 8302(tcp,udp): Serf WAN，跨数据中心gossip交换数据用
- 8400(tcp): CLI RPC，接受命令行的RPC调用
- 8500(tcp): HTTP API及Web UI
- 8600(tcp udp): DNS服务，可以把它配置到53端口来响应dns请求

