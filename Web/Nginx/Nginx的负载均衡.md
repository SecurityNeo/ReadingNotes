# Nginx的负载均衡 #

## round robin ##
round robin（轮询）也是Nginx默认使用的负载均衡算法，Nginx会将请求依次分配给后端的实际服务器，当后端服务器宕掉时，Nginx会将其剔除其集群。可以为每一台服务器定义一个weight值（权重），这样Nginx就会根据backend的权重大小进行请求分发。

## Least Connections ##
Nginx会记录和backend当前的活跃连接数，根据活跃连接数大小进行请求的分配，当然还要综合考虑每个upstream分配的weight权重信息。

## IP Hash ##
Nginx记录请求的源IP，对源IP进行hash计算（对于IPv4地址，取前三个Octes进行hash计算，对于IPv6地址，取整个源IP进行hash计算），并根据计算出的hash值进行请求分发。这种方式能将同一个客户端的请求分发到同一台服务器上，直到对应服务器不可用。

## Generic Hash ##

用户可以自定义资源，Nginx根据用户自定义的资源进行hash，资源可以是字符串、变量或者数组，例如可以指定资源为源IP:Port、URI等等。

## Least Time ##

Nginx会将请求分配给平均响应延迟最小和活跃连接数最少的backend，可以在least_time里配置以下三个参数：

- header -- 计算从server收到第一个byte的时间
- last_byte –- 计算从server收到一个完整响应的时间
- last_byte inflight –- 计算从server接收到完整响应的时间（考虑不完整的请求）

## Random ##
Nginx将每个请求随机分配给后端的服务器，它提供了一个参数two，当指定了这个参数时，Nginx会先随机地选择两个server(考虑weight)，然后用以下几种方法选择其中的一个服务器：

- least_conn -- 最少连接数
- least_time=header(NGINX PLUS only) -- 接收到响应header的最短平均时间
- least_time=last_byte(NGINX PLUS only) -- 接收到完整请求的最短平均时间



