# Keepalived #

[http://www.keepalived.org/manpage.html](http://www.keepalived.org/manpage.html)

[http://www.keepalived.org/documentation.html](http://www.keepalived.org/documentation.html)


Keepalived起初是为LVS设计的，用于监控集群系统中各个服务节点的状态，后来Keepalived又加入了VRRP的功能，VRRP（Vritrual Router Redundancy Protocol,虚拟路由冗余协议)可以解决静态路由出现的单点故障问题，通过VRRP可以实现网络不间断稳定运行

注：VRRP相关内容移步到Network-->VRRP中复习

## Keepalived体系架构 ##


Keepalived 启动只有会有三个进程：

111	Keepalived	<-- Parent process monitoring children

112	\_ Keepalived	<-- VRRP child

113	\_ Keepalived	<-- Healthchecking child

父进程fork出子进程并进行监控，父进程也称为WatchDog。两个子进程会开启本地套接字Unix Domain Socket。keepalived服务启动后，父进程通过unxi domain socket每隔5秒发送一个"Hello"消息给子进程，如果父进程无法发送消息给子进程，将认为子进程出现问题，并会重启子进程。

Keepalived软件设计架构如下：

![Keepalived](img/keepalived.png)


- Control Plane：keepalived配置文件是keepalived.conf，使用专门设计的编译器进行解析。解析器使用关键字树层次结构来使用特定的处理程序映射每个配置关键字。在解析过程中，配置文件被加载进内存中。

 
- Scheduler-I/O Multiplexer：Keepalived中所有事件都被发送到同一个进程中。keepalived是一个独立的进程，负责调度所有内部任务。


- Memory Management：这个框架提供一些通用的内存管理功能，比如内存分配、重新分配、释放等等。此框架有两种功能模式：正常模式和调试模式，当处于debug模式时，可以跟踪内存泄漏问题。


- WatchDog：此框架提供了子进程监控功能（VRRP & HealthChecking），父进程通过unxi domain socket每隔5秒发送一个"Hello"消息给子进程来检测子进程的健康状态。


- Checkers：这是keepalived的主要功能之一，其负责RealServer的健康检查，通过检测结构在LVS的拓扑中移除、添加RealServer，支持layer4/5/7层的协议检查。


- VRRP Stack：这是keepalived另一个主要功能，通过VRRP协议（RFC2338）实现Director的高可用。


- System call：提供读取自定义脚本的功能，主要用于MISC检查，其将临时产生一个子进程来执行对应任务，不影响全局调度计时器。


- SMTP：SMTP协议用于管理通知，为HealthChecker活动和VRRP协议状态转换发送通知给管理员。


- IPVS wrapper：这个框架负责将用户定于的配置文件中IPVS相关规则发送到内核的ipvs模块。


- Netlink Reflector：此模块用于VIP的设置、监控。