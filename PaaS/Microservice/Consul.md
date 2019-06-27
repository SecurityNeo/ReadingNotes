# Consul #

## 简介 ##

[官网](https://www.consul.io/)

[API接口文档](https://www.consul.io/api/index.html)

[Consul的Helm包](https://github.com/hashicorp/consul-helm)

Consul是HashiCorp出品的开源服务发现工具，Consul提供了诸如服务发现，健康检查，KV数据库等功能，方便构建服务集群。

主要特性：

- **服务发现**：Consul提供了通过DNS或者HTTP接口的方式来注册服务和发现服务。一些外部的服务通过Consul很容易的找到它所依赖的服务。
- **健康检查**: Consul客户端可用提供任意数量的健康检查,指定一个服务(比如:webserver是否返回了200 OK 状态码)或者使用本地节点(比如:内存使用是否大于90%). 这个信息可由operator用来监视集群的健康.被服务发现组件用来避免将流量发送到不健康的主机. 
- **Key/Value存储**:应用程序可用根据自己的需要使用Consul的层级的Key/Value存储.比如动态配置,功能标记,协调,领袖选举等等,简单的HTTP API让他更易于使用.。
- **多数据中心**: Consul支持开箱即用的多数据中心。这意味着用户不需要担心需要建立额外的抽象层让业务扩展到多个区域。

相关术语：
[https://www.cnblogs.com/lsf90/p/6021465.html](https://www.cnblogs.com/lsf90/p/6021465.html)

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

[https://blog.csdn.net/liuzhuchen/article/details/81913562](https://blog.csdn.net/liuzhuchen/article/details/81913562)

![](img/Consul_Arch.png)

向Consul提供服务的每个节点都运行一个Consul代理。 发现其他服务或获取/设置键/值数据不需要运行代理。 代理负责健康检查节点上的服务以及节点本身。代理与一个或多个Consul服务器通信。Consul服务器是数据存储和复制的地方。 服务器自己选出一个leader。 虽然Consul可以在一台服务器上运行，但推荐使用3到5台来避免数据丢失的情况。 每个数据中心都建议使用一组Consul服务器。需要发现其他服务或节点的基础架构组件可以查询任何Consul服务器或任何Consul代理。 代理自动将查询转发到服务器。每个数据中心都运行Consul服务器集群。 当跨数据中心服务发现或配置请求时，本地Consul服务器将请求转发到远程数据中心并返回结果。

默认端口：

- 8300(tcp): Server RPC，server 用于接受其他 agent 的请求
- 8301(tcp,udp): Serf LAN，数据中心内gossip交换数据用
- 8302(tcp,udp): Serf WAN，跨数据中心gossip交换数据用
- 8400(tcp): CLI RPC，接受命令行的RPC调用
- 8500(tcp): HTTP API及Web UI
- 8600(tcp udp): DNS服务，可以把它配置到53端口来响应dns请求


## 配置 ##

**运行Agent**

完成Consul的安装后,必须运行agent,agent可以运行为server或client模式.每个数据中心至少必须拥有一台server.建议在一个集群中有3或者5个server,部署单一的server,在出现失败时会不可避免的造成数据丢失。一个client是一个非常轻量级的进程，用于注册服务、运行健康检查和转发对server的查询。agent必须在集群中的每个主机上运行。

**启动Consul Server**

以server模式运行cosnul agent：

`consul agent -server -bootstrap-expect 3 -data-dir /tmp/consul -node=s1 -bind=10.201.102.198 -ui-dir ./consul_ui/ -rejoin -config-dir=/etc/consul.d/ -client 0.0.0.0`

- server ： 定义agent运行在server模式
- bootstrap-expect ：在一个datacenter中期望提供的server节点数目，当该值提供的时候，consul一直等到达到指定sever数目的时候才会引导整个集群，该标记不能和bootstrap共用
- bind：该地址用来在集群内部的通讯，集群内的所有节点到地址都必须是可达的，默认是0.0.0.0
- node：节点在集群中的名称，在一个集群中必须是唯一的，默认是该节点的主机名
- ui-dir： 提供存放web ui资源的路径，该目录必须是可读的
- rejoin：使consul忽略先前的离开，在再次启动后仍旧尝试加入集群中。
- config-dir：配置文件目录，里面所有以.json结尾的文件都会被加载
- client：consul服务侦听地址，这个地址提供HTTP、DNS、RPC等服务，默认是127.0.0.1所以不对外提供服务，如果你要对外提供服务改成0.0.0.0

**启动Consul Client**

`consul agent -data-dir /tmp/consul -node=c1 -bind=10.201.102.248 -config-dir=/etc/consul.d/ -join 10.201.102.198`

**WEB管理界面**

`consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul -node=s1 -bind=10.201.102.198 -ui-dir ./consul_ui/ -rejoin -config-dir=/etc/consul.d/ -client 0.0.0.0`

- ui-dir： 提供存放web ui资源的路径，指向该目录必须是可读的
- client： consul服务侦听地址，这个地址提供HTTP、DNS、RPC等服务，默认是127.0.0.1所以不对外提供服务，如果要对外提供服务改成0.0.0.0 

**配置文件参数**

- acl_datacenter：只用于server，指定的datacenter的权威ACL信息，所有的servers和datacenter必须同意ACL datacenter
- acl_default_policy：默认是allow
- acl_down_policy：
- acl_master_token：
- acl_token：agent会使用这个token和consul server进行请求
- acl_ttl：控制TTL的cache，默认是30s
- addresses：一个嵌套对象，可以设置以下key：dns、http、rpc
- advertise_addr：等同于-advertise
- bootstrap：等同于-bootstrap
- bootstrap_expect：等同于-bootstrap-expect
- bind_addr：等同于-bind
- ca_file：提供CA文件路径，用来检查客户端或者服务端的链接
- cert_file：必须和key_file一起
- check_update_interval：
- client_addr：等同于-client
- datacenter：等同于-dc
- data_dir：等同于-data-dir
- disable_anonymous_signature：在进行更新检查时禁止匿名签名
- disable_remote_exec：禁止支持远程执行，设置为true，agent会忽视所有进入的远程执行请求
- disable_update_check：禁止自动检查安全公告和新版本信息
- dns_config：是一个嵌套对象，可以设置以下参数：allow_stale、max_stale、node_ttl 、service_ttl、enable_truncate
- domain：默认情况下consul在进行DNS查询时，查询的是consul域，可以通过该参数进行修改
- enable_debug：开启debug模式
- enable_syslog：等同于-syslog
- encrypt：等同于-encrypt
- key_file：提供私钥的路径
- leave_on_terminate：默认是false，如果为true，当agent收到一个TERM信号的时候，它会发送leave信息到集群中的其他节点上。
- log_level：等同于-log-level
- node_name:等同于-node
- ports：这是一个嵌套对象，可以设置以下key：dns(dns地址：8600)、http(http api地址：8500)、rpc(rpc:8400)、serf_lan(lan port:8301)、serf_wan(wan port:8302)、server(server rpc:8300)
- protocol：等同于-protocol
- recursor：
- rejoin_after_leave：等同于-rejoin
- retry_join：等同于-retry-join
- retry_interval：等同于-retry-interval
- server：等同于-server
- server_name：会覆盖TLS CA的node_name，可以用来确认CA name和hostname相匹配
- skip_leave_on_interrupt：和leave_on_terminate比较类似，不过只影响当前句柄
- start_join：一个字符数组提供的节点地址会在启动时被加入
- statsd_addr：
- statsite_addr：
- syslog_facility：当enable_syslog被提供后，该参数控制哪个级别的信息被发送，默认Local0
- ui_dir：等同于-ui-dir
- verify_incoming：默认false，如果为true，则所有进入链接都需要使用TLS，需要客户端使用ca_file提供ca文件，只用于consul server端，因为client从来没有进入的链接
- verify_outgoing：默认false，如果为true，则所有出去链接都需要使用TLS，需要服务端使用ca_file提供ca文件，consul server和client都需要使用，因为两者都有出去的链接
- watches：watch一个详细名单


## 命令行 ##

```
[root@dhcp-10-201-102-198 ~]# consul
usage: consul [--version] [--help] <command> [<args>]
Available commands are:
    agent          agent指令是consul的核心，它运行agent来维护成员的重要信息、运行检查、服务宣布、查询处理等等。
    configtest     Validate config file
    event          Fire a new event
    exec           Executes a command on Consul nodes  在consul节点上执行一个命令
    force-leave    Forces a member of the cluster to enter the "left" state   强制节点成员在集群中的状态转换到left状态
    info           Provides debugging information for operators  提供操作的debug级别的信息
    join           Tell Consul agent to join cluster   加入consul节点到集群中
    keygen         Generates a new encryption key  生成一个新的加密key
    keyring        Manages gossip layer encryption keys
    kv             Interact with the key-value store
    leave          Gracefully leaves the Consul cluster and shuts down
    lock           Execute a command holding a lock
    maint          Controls node or service maintenance mode
    members        Lists the members of a Consul cluster    列出集群中成员
    monitor        Stream logs from a Consul agent  打印consul节点的日志信息
    operator       Provides cluster-level tools for Consul operators
    reload         Triggers the agent to reload configuration files   触发节点重新加载配置文件
    rtt            Estimates network round trip time between nodes
    snapshot       Saves, restores and inspects snapshots of Consul server state
    version        Prints the Consul version    打印consul的版本信息
    watch          Watch for changes in Consul   监控consul的改变
```

**event**

event命令提供了一种机制，用来fire自定义的用户事件，这些事件对consul来说是不透明的，但它们可以用来构建自动部署、重启服务或者其他行动的脚本。

```
- http-addr：http服务的地址，agent可以链接上来发送命令，如果没有设置，则默认是127.0.0.1:8500。
- datacenter：数据中心。
- name：事件的名称
- node：一个正则表达式，用来过滤节点
- service：一个正则表达式，用来过滤节点上匹配的服务
- tag：一个正则表达式，用来过滤节点上符合tag的服务，必须和-service一起使用。
```

**exec**

exec指令提供了一种远程执行机制，比如你要在所有的机器上执行uptime命令，远程执行的工作通过job来指定，存储在KV中，agent使用event系统可以快速的知道有新的job产生，消息是通过gossip协议来传递的，因此消息传递是最佳的，但是并不保证命令的执行。事件通过gossip来驱动，远程执行依赖KV存储系统(就像消息代理一样)。

```
- http-addr：http服务的地址，agent可以链接上来发送命令，如果没有设置，则默认是127.0.0.1:8500。
- datacenter：数据中心。
- prefix：key在KV系统中的前缀，用来存储请求数据，默认是_rexec
- node：一个正则表达式，用来过滤节点，评估事件
- service：一个正则表达式，用来过滤节点上匹配的服务
- tag：一个正则表达式，用来过滤节点上符合tag的服务，必须和-service一起使用。
- wait：在节点多长时间没有响应后，认为job已经完成。
- wait-repl：
- verbose：输出更多信息
```

**force-leave**

force-leave治疗可以强制consul集群中的成员进入left状态(空闲状态)，记住，即使一个成员处于活跃状态，它仍旧可以再次加入集群中，这个方法的真实目的是强制移除failed的节点。如果failed的节点还是网络的一部分，则consul会周期性的重新链接failed的节点，如果经过一段时间后(默认是72小时)，consul则会宣布停止尝试链接failed的节点。force-leave指令可以快速的把failed节点转换到left状态。

```
- rpc-addr:一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400。
```

**info**

info指令提供了各种操作时可以用到的debug信息，对于client和server，info有返回不同的子系统信息，目前有以下几个KV信息：agent(提供agent信息)，consul(提供consul库的信息)，raft(提供raft库的信息)，serf_lan(提供LAN gossip pool),serf_wan(提供WAN gossip pool)

```
- rpc-addr：一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400
```

**join**

join指令告诉consul agent加入一个已经存在的集群中，一个新的consul agent必须加入一个已经有至少一个成员的集群中，这样它才能加入已经存在的集群中，如果你不加入一个已经存在的集群，则agent是它自身集群的一部分，其他agent则可以加入进来。agents可以加入其他agent多次。consul join [options] address。如果你想加入多个集群，则可以写多个地址，consul会加入所有的地址。

```
- wan：agent运行在server模式，xxxxxxx
- rpc-addr：一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400。
```

**keygen**

keygen指令生成加密的密钥，可以用在consul agent通讯加密

**leave**

leave指令触发一个优雅的离开动作并关闭agent，节点离开后不会尝试重新加入集群中。运行在server状态的节点，节点会被优雅的删除，这是很严重的，在某些情况下一个不优雅的离开会影响到集群的可用性。

```
- rpc-addr:一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400。
```

**members**

members指令输出consul agent目前所知道的所有的成员以及它们的状态，节点的状态只有alive、left、failed三种状态。

```
detailed：输出每个节点更详细的信息。
rpc-addr：一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400。
status：过滤出符合正则规则的节点
wan：xxxxxx
```

**monitor**

monitor指令用来链接运行的agent，并显示日志。monitor会显示最近的日志，并持续的显示日志流，不会自动退出，除非你手动或者远程agent自己退出。

```
- log-level：显示哪个级别的日志，默认是info
- rpc-addr：一个rpc地址，agent可以链接上来发送命令，如果没有指定，默认是127.0.0.1:8400
```

**reload**

reload指令可以重新加载agent的配置文件。SIGHUP指令在重新加载配置文件时使用，任何重新加载的错误都会写在agent的log文件中，并不会打印到屏幕。

**watch**

watch指令提供了一个机制，用来监视实际数据视图的改变(节点列表、成员服务、KV)，如果没有指定进程，当前值会被dump出来

```
http-addr：http服务的地址，agent可以链接上来发送命令，如果没有设置，则默认是127.0.0.1:8500。
datacenter：数据中心查询。
token：ACL token
key：监视key，只针对key类型
name：监视event，只针对event类型
prefix：监视key prefix，只针对keyprefix类型
service：监控service，只针对service类型
state：过略check state
tag：过滤service tag
type：监控类型，一般有key、keyprefix、service、nodes、checks、event
```

