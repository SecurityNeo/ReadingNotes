# TCP的SYN queue和accept queue #


socket接收的所有连接都存放在队列类型中，队列有以下两种：

- syns queue（半连接队列，用来保存处于SYN_SENT和SYN_RECV状态的请求）
- accept queue（全连接队列，用来保存处于established状态，但是应用层没有调用accept取走的请求）

![](img/tcp-sync-queue-and-accept-queue-small.jpg)

图片来源于[http://www.cnxct.com/something-about-phpfpm-s-backlog/](http://www.cnxct.com/something-about-phpfpm-s-backlog/)

在TCP握手的第三步，Server收到Client的ACK后，如果全连接队列没满，server就从半连接队列拿出这个连接的信息放入到全连接队列中,如果全连接队列满了并且`net.ipv4.tcp_abort_on_overflow`是0的话，server过一段时间再次发送syn+ack给client，重试的次数由`net.ipv4.tcp_synack_retries`决定，如果全队列满了并且`net.ipv4.tcp_abort_on_overflow`是1的话，Server发送一个reset包给Client，客户端一般会看到connection reset by peer的错误。

备注：关于内核参数，可以到[Linux内核调优部分参数说明](https://github.com/SecurityNeo/ReadingNotes/blob/master/Linux/%E6%80%A7%E8%83%BD/Linux%E5%86%85%E6%A0%B8%E8%B0%83%E4%BC%98%E9%83%A8%E5%88%86%E5%8F%82%E6%95%B0%E8%AF%B4%E6%98%8E.md)复习。

