[Envoy中文文档](http://www.servicemesher.com/envoy/)

[网易云刘超的《深入解读 Service Mesh 背后的技术细节》](https://www.cnblogs.com/163yun/p/8962278.html)

[理解Istio Service Mesh中Envoy代理Sidecar注入及流量劫持](https://jimmysong.io/posts/envoy-sidecar-injection-in-istio-service-mesh-deep-dive/)


## Envoy threading model ##

[https://my.oschina.net/u/1262062/blog/3011063](https://my.oschina.net/u/1262062/blog/3011063)

[https://blog.envoyproxy.io/envoy-threading-model-a8d44b922310](https://blog.envoyproxy.io/envoy-threading-model-a8d44b922310)

![](img/Envoy_ThreadingMode.jpg)

**Envoy三种线程**

- Main Thread：此线程拥有服务器启动和关闭、所有xDS API处理（包括DNS，运行状况检查和常规集群管理）、运行时、统计刷新、管理和一般进程管理（信号，热启动等）的功能。 在此线程上发生的所有事情都是异步的并且是“非阻塞的”。通常，主线程协调所有不需要大量CPU来完成的关键过程功能。 这允许将大多数管理代码编写为单线程编写。

- Worker Thread：默认情况下，Envoy为系统中的每个硬件线程生成一个工作线程。 （这可以通过--concurrency 选项控制）。 每个工作线程运行一个“非阻塞”事件循环，负责监听每个侦听器（当前没有侦听器分片），接受新连接，为连接实例化过滤器堆栈，以及处理所有IO的生命周期。 连接。 同样，这允许将大多数连接处理代码写成好像是单线程的。

- File Flush：Envoy写入的每个文件（主要是访问日志）当前都有一个独立的阻塞刷新线程。 这是因为即使使用O_NONBLOCK写入文件系统缓存文件有时也会阻塞（叹息）。 当工作线程需要写入文件时，数据实际上被移入内存缓冲区，最终通过文件刷新线程刷新。 这是代码的一个区域，技术上所有工作人员都可以阻止同一个锁尝试填充内存缓冲区。 还有一些其他的将在下面进一步讨论。



