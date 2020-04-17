# TCP/IP的握手与挥手 #

[https://mp.weixin.qq.com/s/pqUFksNEwT9UWDpcKdGpQg](https://mp.weixin.qq.com/s/pqUFksNEwT9UWDpcKdGpQg)

## TCP头部 ##

![](img/TCP_Header.PNG)

- TCP序列号（序列码SN,Sequence Number）

	在建立连接时由计算机生成的随机数作为其初始值，通过SYN包传给接收端主机，每发送一次数据，就累加一次该数据字节数的大小。用来解决网络包乱序问题。比如发送端发送的一个TCP包净荷(不包含TCP头)为12byte，SN为5，则发送端接着发送的下一个数据包的时候，SN应该设置为5+12=17。通过系列号，TCP接收端可以识别出重复接收到的TCP包，从而丢弃重复包，同时对于乱序数据包也可以依靠系列号进行重排序，进而对高层提供有序的数据流。

- TCP应答号(Acknowledgment Number简称ACK Number或简称为ACK Field)

	ACK Number标识报文发送端期望接收的字节序列。比如当前接收端接收到一个净荷为12byte的数据包，SN为5，则发送端可能会回复一个确认收到的数据包，如果这个数据包之前的数据也都已经收到了，这个数据包中的ACK Number则设置为12+5=17，表示17byte之前的数据都已经收到了。

- 校验位(Checksum)
 
	发送端基于数据内容计算一个数值，接收端checksum校验失败的时候会直接丢掉这个数据包。CheckSum是根据伪头+TCP头+TCP数据三部分进行计算。另外对于大的数据包，checksum并不能可靠的反应比特错误，应用层应该再添加自己的校验方式。

- CWR(Congestion Window Reduce)

	拥塞窗口减少标志被发送主机设置，用来表明它接收到了设置ECE标志的TCP包。拥塞窗口是被TCP维护的一个内部变量，用来管理发送窗口大小。

- ECE(ECN Echo)

	用来在TCP三次握手时表明一个TCP端是具备ECN功能的。在数据传输过程中，它也用来表明接收到的TCP包的IP头部的ECN被设置为11，即网络线路拥堵。

- URG(Urgent)

	该标志位置位表示紧急(The urgent pointer) 标志有效。

- ACK(Acknowledgment)

	取值1代表Acknowledgment Number字段有效

- PSH(Push)

	该标志置位时，一般是表示发送端缓存中已经没有待发送的数据，接收端不将该数据进行队列处理，而是尽可能快将数据转由应用处理。在处理telnet或rlogin等交互模式的连接时，该标志总是置位的。

- RST(Reset)

	用于复位相应的TCP连接。通常在发生异常或者错误的时候会触发复位TCP连接。

- SYN(Synchronize)

	同步序列编号(Synchronize Sequence Numbers)有效。该标志仅在三次握手建立TCP连接时有效。

- FIN(Finish)
 
	带有该标志置位的数据包用来结束一个TCP会话，但对应端口仍处于开放状态，准备接收后续数据。