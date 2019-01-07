# HAProxy #

[https://www.haproxy.com/documentation/hapee](https://www.haproxy.com/documentation/hapee/)

[https://cbonte.github.io/haproxy-dconv/1.9/intro.html](https://cbonte.github.io/haproxy-dconv/1.9/intro.html)


## 简介 ##

HAProxy是一个使用C语言编写的开源软件，提供高可用性、负载均衡，以及基于TCP和HTTP的应用程序代理，并且支持虚拟主机。对于一些特负载量大的web应用，haproxy非常适用，并且支持会话保持、SSL、ACL等。HAProxy采用一种事件驱动、单一进程的模型，能支持非常大的并发连接数。对于多进程或多线程模型而言，内存限制 、系统调度器限制以及锁限制等会直接影响到其并发处理能力。由此也暴露了此模型的弊端，在多核系统上，程序扩展性较差。

对于HAProxy、Nginx和LVS三种负载均衡软件的优缺点，[三大主流软件负载均衡器对比(LVS VS Nginx VS Haproxy)](https://www.cnblogs.com/ahang/p/5799065.html) ，这篇文章写得很好。

