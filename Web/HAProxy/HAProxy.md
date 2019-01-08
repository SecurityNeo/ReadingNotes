# HAProxy #

[https://www.haproxy.com/documentation/hapee](https://www.haproxy.com/documentation/hapee/)

[https://cbonte.github.io/haproxy-dconv/1.9/intro.html](https://cbonte.github.io/haproxy-dconv/1.9/intro.html)


## 简介 ##

HAProxy是一个使用C语言编写的开源软件，提供高可用性、负载均衡，以及基于TCP和HTTP的应用程序代理，并且支持虚拟主机。对于一些特负载量大的web应用，haproxy非常适用，并且支持会话保持、SSL、ACL等。HAProxy采用一种事件驱动、单一进程的模型，能支持非常大的并发连接数。对于多进程或多线程模型而言，内存限制 、系统调度器限制以及锁限制等会直接影响到其并发处理能力。由此也暴露了此模型的弊端，在多核系统上，程序扩展性较差。

对于HAProxy、Nginx和LVS三种负载均衡软件的优缺点，[三大主流软件负载均衡器对比(LVS VS Nginx VS Haproxy)](https://www.cnblogs.com/ahang/p/5799065.html) ，这篇文章写得很好。

## 负载均衡算法 ##

- roundrobin

基于权重进行轮询.这种算法在设计上限制后端只能由4095台活动服务器

Each server is used in turns, according to their weights.This is the smoothest and fairest algorithm when the server’sprocessing time remains equally distributed. This algorithm is dynamic, which means that server weights may be adjusted on the fly for slow starts for instance. It is limited by design to 4095 active servers per backend. Note that in some large farms, when a server becomes up after having been down for a very short time, it may sometimes take a few hundreds requests for it to be re-integrated into the farm and start receiving traffic. This is normal, though very rare. It is indicated here in case you would have the chance to observe it, so that you don’t worry.


- static-rr 

这个算法与roundrobin类似，也是根据权重进行轮询,但其为一种静态方法,在运行状态下调整服务器权重值不会生效。不过,其在后端服务器数量上没有限制

Each server is used in turns, according to their weights. This algorithm is as similar to roundrobin except that it is static, which means that changing a server’s weight on the fly will have no effect. On the other hand, it has no design limitation on the number of servers, and when a server goes up, it is always immediately reintroduced into the farm, once the full map is recomputed. It also uses slightly less CPU to run (around -1%).



- leastconn 

新的连接请求会被转发到具有最少连接数目的后端服务器上。非常适用于会话时间较长的场景，例如LDAP, SQL, TSE等。

The server with the lowest number of connections receives the connection. Round-robin is performed within groups of servers of the same load to ensure that all servers will be used. Use of this algorithm is recommended where very long sessions are expected, such as LDAP, SQL, TSE, etc… but is not very well suited for protocols using short sessions such as HTTP. This algorithm is dynamic, which means that server weights may be adjusted on the fly for slow starts for instance.


- first 

请求将转发给第一个具有可用连接槽的服务器

The first server with available connection slots receives the connection. The servers are chosen from the lowest numeric identifier to the highest (see server parameter “id”), which defaults to the server’s position in the farm. Once a server reaches its maxconn value, the next server is used. It does not make sense to use this algorithm without setting maxconn. The purpose of this algorithm is to always use the smallest number of servers so that extra servers can be powered off during non-intensive hours. This algorithm ignores the server weight, and brings more benefit to long session such as RDP or IMAP than HTTP, though it can be useful there too. In order to use this algorithm efficiently, it is recommended that a cloud controller regularly checks server usage to turn them off when unused, and regularly checks backend queue to turn new servers on when the queue inflates. Alternatively, using “http-check send-state” may inform servers on the load.


- source 

将请求的源IP地址，作为散列键（Hash Key）进行hash运算，并由后端服务器的权重总数相除后派发至某匹配的服务器。

The source IP address is hashed and divided by the total weight of the running servers to designate which server will receive the request. This ensures that the same client IP address will always reach the same server as long as no server goes down or up. If the hash result changes due to the number of running servers changing, many clients will be directed to a different server. This algorithm is generally used in TCP mode where no cookie may be inserted. It may also be used on the Internet to provide a best-effort stickiness to clients which refuse session cookies. This algorithm is static by default, which means that changing a server’s weight on the fly will have no effect, but this can be changed using “hash-type”.


- uri 

这也是一种动态算法，对URI的左半部分(“?”标记之前的部分)或整个URI进行hash运算,并由服务器的总权重相除后派发至某匹配的服务器。

This algorithm hashes either the left part of the URI (before the question mark) or the whole URI (if the “whole” parameter is present) and divides the hash value by the total weight of the running servers. The result designates which server will receive the request. This ensures that the same URI will always be directed to the same server as long as no server goes up or down. This is used with proxy caches and anti-virus proxies in order to maximize the cache hit rate. Note that this algorithm may only be used in an HTTP backend. This algorithm is static by default, which means that changing a server’s weight on the fly will have no effect, but this can be changed using “hash-type”.

This algorithm supports two optional parameters “len” and “depth”, both followed by a positive integer number. These options may be helpful when it is needed to balance servers based on the beginning of the URI only. The “len” parameter indicates that the algorithm should only consider that many characters at the beginning of the URI to compute the hash. Note that having “len” set to 1 rarely makes sense since most URIs start with a leading “/”.

The “depth” parameter indicates the maximum directory depth to be used to compute the hash. One level is counted for each slash in the request. If both parameters are specified, the evaluation stops when either is reached.



- url_param 

通过<argument>为URL指定的参数在每个HTTP GET请求中将会被检索

The URL parameter specified in argument will be looked up inthe query string of each HTTP GET request.

If the modifier “check_post” is used, then an HTTP POST request entity will be searched for the parameter argument, when it is not found in a query string after a question mark (‘?’) in the URL. The message body will only start to be analyzed once either the advertised amount of data has been received or the request buffer is full. In the unlikely event that chunked encoding is used, only the first chunk is scanned. Parameter values separated by a chunk boundary, may be randomly balanced if at all. This keyword used to support an optional parameter which is now ignored.

If the parameter is found followed by an equal sign (‘=’) and a value, then the value is hashed and divided by the total weight of the running servers. The result designates which server will receive the request.

This is used to track user identifiers in requests and ensure that a same user ID will always be sent to the same server as long as no server goes up or down. If no value is found or if the parameter is not found, then a round robin algorithm is applied. Note that this algorithm may only be used in an HTTP backend. This algorithm is static by default, which means that changing a server’s weight on the fly will have no effect, but this can be changed using “hash-type”.


- hdr(<name>) 

对于每个HTTP请求,通过<name>指定的HTTP首部将会被检索

The HTTP header will be looked up in each HTTP request. Just as with the equivalent ACL ‘hdr()’ function, the header name in parenthesis is not case sensitive. If the header is absent or if it does not contain any value, the roundrobin algorithm is applied instead.

An optional ‘use_domain_only’ parameter is available, for reducing the hash algorithm to the main domain part with some specific headers such as ‘Host’. For instance, in the Host value “haproxy.1wt.eu”, only “1wt” will be considered.

This algorithm is static by default, which means that changing a server’s weight on the fly will have no effect, but this can be changed using “hash-type”.

- rdp-cookie()

The RDP cookie (or “mstshash” if omitted) will be looked up and hashed for each incoming TCP request. Just as with the equivalent ACL ‘req_rdp_cookie()’ function, the name is not case-sensitive. This mechanism is useful as a degraded persistence mode, as it makes it possible to always send the same user (or the same session ID) to the same server. If the cookie is not found, the normal roundrobin algorithm is used instead.

Note that for this to work, the frontend must ensure that an RDP cookie is already present in the request buffer. For this you must use ‘tcp-request content accept’ rule combined with a ‘req_rdp_cookie_cnt’ ACL.

This algorithm is static by default, which means that changing a server’s weight on the fly will have no effect, but this can be changed using “hash-type”. See also the rdp_cookie pattern fetch function.
