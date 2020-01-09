## netstat与ss ##

[https://www.cnblogs.com/kevingrace/p/6211509.html](https://www.cnblogs.com/kevingrace/p/6211509.html)

ss命令可以用来获取socket统计信息，它可以显示和netstat类似的内容。ss的优势在于它能够显示更多更详细的有关TCP和连接状态的信息，而且比netstat更快速更高效。原因：
- 当服务器的socket连接数量变得非常大时，无论是使用netstat命令还是直接`cat /proc/net/tcp`，执行速度都会很慢。
- ss利用到了TCP协议栈中tcp_diag。tcp_diag是一个用于分析统计的模块，可以获得Linux内核中第一手的信息，这就确保了ss的快捷高效。如果系统中没有tcp_diag，ss也可以正常运行，只是效率会变得稍慢（但仍然比 netstat要快）。

netstat参数和使用:

- -a 显示所有活动的连接以及本机侦听的TCP、UDP端口
- -l 显示监听的server port
- -n 直接使用IP地址，不通过域名服务器
- -p 正在使用Socket的程序PID和程序名称
- -r 显示路由表
- -t 显示TCP传输协议的连线状况
- -u 显示UDP传输协议的连线状况
- -w 显示RAW传输协议的连线状况

ss（socket statistics）参数和使用:

- -a显示所有的sockets
- -l显示正在监听的
- -n显示数字IP和端口，不通过域名服务器
- -p显示使用socket的对应的程序
- -t只显示TCP sockets
- -u只显示UDP sockets
- -4 -6 只显示v4或v6V版本的sockets
- -s打印出统计信息。这个选项不解析从各种源获得的socket。对于解析/proc/net/top大量的sockets计数时很有效
- -0 显示PACKET sockets
- -w 只显示RAW sockets
- -x只显示UNIX域sockets
- -r尝试进行域名解析，地址/端口

ss还可以使用IP地址筛选如ss src xxxxIP:port，以及使用端口筛选ss dport OP PORT，OP支持的运算符有le ge eq ne lt gt。

netstat中的各种状态:

- CLOSED         初始（无连接）状态。
- LISTEN         侦听状态，等待远程机器的连接请求。
- SYN_SEND       在TCP三次握手期间，主动连接端发送了SYN包后，进入SYN_SEND状态，等待对方的ACK包。
- SYN_RECV       在TCP三次握手期间，主动连接端收到SYN包后，进入SYN_RECV状态。
- ESTABLISHED    完成TCP三次握手后，主动连接端进入ESTABLISHED状态。此时，TCP连接已经建立，可以进行通信。
- FIN_WAIT_1     在TCP四次挥手时，主动关闭端发送FIN包后，进入FIN_WAIT_1状态。
- FIN_WAIT_2     在TCP四次挥手时，主动关闭端收到ACK包后，进入FIN_WAIT_2状态。
- TIME_WAIT      在TCP四次挥手时，主动关闭端发送了ACK包之后，进入TIME_WAIT状态，等待最多MSL时间，让被动关闭端收到ACK包。
- CLOSING        在TCP四次挥手期间，主动关闭端发送了FIN包后，没有收到对应的ACK包，却收到对方的FIN包，此时，进入CLOSING状态。
- CLOSE_WAIT     在TCP四次挥手期间，被动关闭端收到FIN包后，进入CLOSE_WAIT状态。
- LAST_ACK       在TCP四次挥手时，被动关闭端发送FIN包后，进入LAST_ACK状态，等待对方的ACK包。
 
主动连接端可能的状态有：    CLOSED        SYN_SEND        ESTABLISHED
主动关闭端可能的状态有：    FIN_WAIT_1    FIN_WAIT_2      TIME_WAIT
被动连接端可能的状态有：    LISTEN        SYN_RECV        ESTABLISHED
被动关闭端可能的状态有：    CLOSE_WAIT    LAST_ACK        CLOSED

查看tomcat的并发数：`netstat -an|grep 10050|awk '{count[$6]++} END{for (i in count) print(i,count[i])}'`

