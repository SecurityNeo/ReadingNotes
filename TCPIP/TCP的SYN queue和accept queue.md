# TCP的SYN queue和accept queue #

[https://blog.csdn.net/russell_tao/article/details/9950615](https://blog.csdn.net/russell_tao/article/details/9950615)

socket接收的所有连接都存放在队列类型中，队列有以下两种：

- syns queue（半连接队列，用来保存处于SYN_RECV状态的请求）
- accept queue（全连接队列，用来保存处于established状态，但是应用层没有调用accept取走的请求）

队列的长度：

- accept queue的最大值：min(backlog, /proc/sys/net/core/somaxconn)，在linux内核2.2版本以后，backlog参数控制的accept queue的大小,backlog是在socket创建的时候传入的,属于listen函数里的参数；somaxconn是内核的参数，默认是128。`max accept queue size =min(backlog,net.core.somaxconn)`

- syn queue的最大值：4.3版本之前的内核,SYN队列的最大大小是用`net.ipv4.tcp_max_syn_backlog`来配置，但是现在已经不再使用了。现在的取法如下： `max syn queue size = min(backlog,net.core.somaxconn,net.ipv4.tcp_max_syn_backlog)`

![sync-queue-and-accept-queue](img/tcp-queue.jpg)

图片来源于[http://www.cnxct.com/something-about-phpfpm-s-backlog/](http://www.cnxct.com/something-about-phpfpm-s-backlog/)

若SYN queue已满，在收到SYN时

	若设置`net.ipv4.tcp_syncookies = 0`，则直接丢弃当前 SYN 包；
	若设置`net.ipv4.tcp_syncookies = 1`，则令`want_cookie = 1`继续后面的处理；
		若accept queue已满，并且qlen_young的值大于1 ，则直接丢弃当前SYN包；
		若accept queue未满，或者qlen_young的值未大于1 ，则输出 "possible SYN flooding on port %d. Sending cookies.\n"，生成syncookie并在SYN,ACK中带上

在TCP握手的第三步，Server收到Client的ACK后，如果全连接队列没满，server就从半连接队列拿出这个连接的信息放入到全连接队列中,如果全连接队列满了并且`net.ipv4.tcp_abort_on_overflow`是0的话，server过一段时间再次发送syn+ack给client，重试的次数由`net.ipv4.tcp_synack_retries`决定，如果全队列满了并且`net.ipv4.tcp_abort_on_overflow`是1的话，Server发送一个reset包给Client，客户端一般会看到connection reset by peer的错误。

备注：关于内核参数，可以到[Linux内核调优部分参数说明](https://github.com/SecurityNeo/ReadingNotes/blob/master/Linux/%E6%80%A7%E8%83%BD/Linux%E5%86%85%E6%A0%B8%E8%B0%83%E4%BC%98%E9%83%A8%E5%88%86%E5%8F%82%E6%95%B0%E8%AF%B4%E6%98%8E.md)复习。

对于一个SYN包，如果syn queue满了并且没有开启syncookies就丢包，并将`ListenDrops`计数器 +1。如果accept queue满了也会丢包，并将`ListenOverflows`和`ListenDrops`计数器 +1。可以通过命令`netstat -s |grep -E 'overflow|drop'`或者`nstat -az |grep -E 'TcpExtListenOverflows|TcpExtListenDrops'`来查看。另外注意，对于低版本内核，当accept queue满了，并不会完全丢弃SYN包，而是对SYN限速。

