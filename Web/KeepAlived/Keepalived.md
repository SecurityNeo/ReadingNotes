# Keepalived #

[http://www.keepalived.org/manpage.html](http://www.keepalived.org/manpage.html)


Keepalived起初是为LVS设计的，用于监控集群系统中各个服务节点的状态，后来Keepalived又加入了VRRP的功能，VRRP（Vritrual Router Redundancy Protocol,虚拟路由冗余协议)可以解决静态路由出现的单点故障问题，通过VRRP可以实现网络不间断稳定运行

注：VRRP相关内容移步到Network-->VRRP中复习

## Keepalived体系架构 ##

![Keepalived](img/keepalived.png)


