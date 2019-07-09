# DRBD简介 #

DRBD（Distributed Replicated Block Device）：叫做分布式复制块设备，这是一种基于软件，无共享，复制的解决方案。在服务器之间的块设备（包括硬盘、分区、逻辑卷）进行镜像。DRBD是镜像块设备，是按数据位镜像成一样的数据块。DRBD的位置处于文件系统以下，比文件系统更加靠近操作系统内核及IO栈。

## 工作原理 ##

![](img/drbd.png)

客户端发起一个写操作的系统调用给文件系统，写请求再到达内核缓冲区，然后到达DRBD模块，此时drbd会复制写入磁盘的数据，且进行两步操作，第一步操作是调用磁盘驱动，将数据写入本地的磁盘设备，第二部是通过网卡设备将数据发送给备用节点，备用节点的网卡接受到数据之后，将数据再发送给drbd模块，DRBD模块再调用磁盘驱动将数据写入磁盘。

**工作模式**

- 主从模型master/slave（primary/secondary）

	这种模式下，在某一时刻只允许有一个主节点。主节点可以挂载使用，写入数据等；从节点不可以挂载文件系统，因此，也不可以执行读写操作。这种模式可用在任何的文件系统上。默认这种模式下，一旦主节点发生故障，从节点需要手工将资源进行转移，且主节点变成从节点和从节点变成主节点需要手动进行切换。不能自动进行转移，因此比较麻烦。


- 双主模型 dula primary(primary/primary)

	DRBD8.0之后的新特性。所谓双主模型是2个节点都可以当做主节点来挂载使用。在双主模式下，任何资源在任何特定的时间都存在两个主节点。这种模式需要一个共享的集群文件系统，利用分布式的锁机制进行管理，如GFS和OCFS2。部署双主模式时，DRBD可以是负载均衡的集群，这就需要从两个并发的主节点中选取一个首选的访问数据。这种模式默认是禁用的，如果要是用的话必须在配置文件中进行声明。

**同步协议**

- 协议A

	数据在本地完成写操作且数据已经发送到TCP/IP协议栈的队列中，则认为写操作完成。如果本地节点的写操作完成，此时本地节点发生故障，而数据还处在TCP/IP队列中，则数据不会发送到对端节点上。因此，两个节点的数据将不会保持一致。这种协议虽然高效，但是并不能保证数据的可靠性。

- 协议B
 
	数据在本地完成写操作且数据已到达对端节点则认为写操作完成。如果两个节点同时发生故障，即使数据到达对端节点，这种方式同样也会导致在对端节点和本地节点的数据不一致现象，也不具有可靠性。

- 协议C

	只有当本地节点的磁盘和对端节点的磁盘都完成了写操作，才认为写操作完成。这是集群流行的一种方式，应用也是最多的，这种方式虽然不高效，但是最可靠。

## 常用命令 ##
[https://blog.csdn.net/jailman/article/details/81479891](https://blog.csdn.net/jailman/article/details/81479891)


- 使用`drbd-overview`命令观察状态:

	```# drbd-overview
	  0:r0  Connected Primary/Secondary UpToDate/UpToDate C r—– /nfs ext4 20G 45M 19G 1%
	```

- 通过伪文件系统/proc/drbd 文件来运行状态

	```
	# cat /proc/drbd 
	version: 8.3.16 (api:88/proto:86-97)
	GIT-hash: a798fa7e274428a357657fb52f0ecf40192c1985 build by phil@Build64R6, 2014-11-24 14:51:37
	 0: cs:Connected ro:Primary/Secondary ds:UpToDate/UpToDate C r—–
	    ns:0 nr:0 dw:664 dr:2017 al:1 bm:0 lo:0 pe:0 ua:0 ap:0 ep:1 wo:f oos:0
	```

	- cs(Connect State)：表示网络连接的状态
	- ro(Role)：表示运行节点的角色，Primary/Secondary　表示本机为主
	- ds(Disk State)：表示当前的硬盘状态
	- Replication protocol：表示当前复制所使用的协议，可以是ABC
	- I/O Flags：6个I/O输入输出标志，从各个方面反映了本地资源的状态
	- Performance indicators：性能指标，这是一组统计数据和计数器，反映出资源的利用情况和性能
	- IO状态标记表示当前资源的IO操作状态。共有6种状态：
		- IO挂起：r或s都可能表示IO挂起，一般是r。r=running，s=suspended。
		- 串行重新同步：资源正在等待进行重新同步，但被resync-after选项延迟了同步进度。该状态标记为"a"，通常该状态栏应该处于"-"。
		- 对端初始化同步挂起：资源正在等待进行重新同步，但对端节点因为某些原因而IO挂起。该状态标记为"p"，通常该状态栏应该处于"-"。
		- 本地初始化同步挂起：资源正在等待进行重新同步，但本节点因为某些原因而IO挂起。该状态标记为"u"，通常该状态栏应该处于"-"。
		- 本地IO阻塞：通常该状态栏应该处于"-"。可能有以下几种标记：
			- d：因为DRBD内部原因导致的IO阻塞。
			- b：后端设备正处于IO阻塞。
			- n：网络套接字阻塞。
			- a：网络套接字和后端块设备同时处于阻塞状态。
		- Activity Log更新挂起：当al更新被挂起时，处于该状态，标记为"s"，通常该状态栏应该处于"-"。(如果不知道什么是Active Log，请无视本标记)
	- ns(network send)：通过网络连接发送到对端的数据量，单位KB.
	- nr(network receive)：通过网络连接从对点接收的数据量，单位KB.
	- dw(disk write)：向本地硬盘写入网络数据，单位KB.
	- dr(disk read)：网络从本地硬盘读取的数据量，单位KB.
	- al(activity log)：元数据活动日志的更新次数。
	- bm(bit map)：元数据区域更新的资源。
	- lo(local count)：由DRBD产生的本地I/O请求数据。
	- pe(pending)：就是等待响应，已经发送到圣战，但是还没有得到对端回应的数量。
	- ua(unacknow wledged)：就是未确认，通过网络连接收到对方的请求，但是还没有做出处理的数量.
	- ap(application pending)：转发到DRBD的I/O请求，仍然没有被DRBD所响应。
	- ep(epochs)：epoch对象的数，通常为1。当使用barrier或者none写顺序方法时，可能会增加底层I/O负荷。
	- wo(write order)：当前使用的写顺序的方法：b(barrier)/f(flush)/d(drain)/n(none)。
	- oos(out of sync)：当前没有同步的数据总数量，单位为KB.

	drbd9中添加了以下几个指标：

	- resync-suspended：重新同步操作当前是否被挂起。可能的值为no/user/peer/dependency。
	- blocked：本地IO的拥挤情况。
	- no：本地IO不拥挤。
	- upper：DRBD层之上的IO被阻塞。例如到文件系统上的IO阻塞。可能有以下几种原因：
	- 管理员使用drbdadm suspend-io命令挂起了IO操作。
	- 短暂的IO阻塞，例如attach/detach导致的。
	- 删除了缓冲区。
	- bitmap的IO等待。
	- lower：底层设备处于拥挤状态。

- 查看磁盘状态

	```drbdadm dstate git```
	
	本地和对等节点的硬盘状态（首先输出的是本地硬盘状态，后面的是远程硬盘状态）：
	
	- Diskless 无盘：本地没有块设备分配给DRBD使用，这表示没有可用的设备，或者使用drbdadm命令手工分离或是底层的I/O错误导致自动分离  
	- Attaching：读取无数据时候的瞬间状态 
	- Failed 失败：本地块设备报告I/O错误的下一个状态，其下一个状态为Diskless无盘  
	- Negotiating：在已经连接的DRBD设置进行Attach读取无数据前的瞬间状态 
	- Inconsistent：数据是不一致的，在两个节点上（初始的完全同步前）这种状态出现后立即创建一个新的资源。此外，在同步期间（同步目标）在一个节点上出现这种状态 
	- Outdated：数据资源是一致的，但是已经过时 
	- DUnknown：当对等节点网络连接不可用时出现这种状态 
	- Consistent：一个没有连接的节点数据一致，当建立连接时，它决定数据是UpToDate或是Outdated 
	- UpToDate：一致的最新的数据状态，这个状态为正常状态 

- 查看资源连接状态

	```drbdadm cstate git```
	
	- StandAlone 独立的：网络配置不可用；资源还没有被连接或是被管理断开（使用 drbdadm disconnect 命令），或是由于出现认证失败或是脑裂的情况
	- Disconnecting 断开：断开只是临时状态，下一个状态是StandAlone独立的
	- Unconnected 悬空：是尝试连接前的临时状态，可能下一个状态为WFconnection和WFReportParams
	- Timeout 超时：与对等节点连接超时，也是临时状态，下一个状态为Unconected悬空
	- BrokerPipe：与对等节点连接丢失，也是临时状态，下一个状态为Unconected悬空
	- NetworkFailure：与对等节点推动连接后的临时状态，下一个状态为Unconected悬空
	- ProtocolError：与对等节点推动连接后的临时状态，下一个状态为Unconected悬空
	- TearDown 拆解：临时状态，对等节点关闭，下一个状态为Unconected悬空
	- WFConnection：等待和对等节点建立网络连接
	- WFReportParams：已经建立TCP连接，本节点等待从对等节点传来的第一个网络包
	- Connected 连接：DRBD已经建立连接，数据镜像现在可用，节点处于正常状态
	- StartingSyncS：完全同步，有管理员发起的刚刚开始同步，未来可能的状态为SyncSource或PausedSyncS
	- StartingSyncT：完全同步，有管理员发起的刚刚开始同步，下一状态为WFSyncUUID
	- WFBitMapS：部分同步刚刚开始，下一步可能的状态为SyncSource或PausedSyncS
	- WFBitMapT：部分同步刚刚开始，下一步可能的状态为WFSyncUUID
	- WFSyncUUID：同步即将开始，下一步可能的状态为SyncTarget或PausedSyncT
	- SyncSource：以本节点为同步源的同步正在进行
	- SyncTarget：以本节点为同步目标的同步正在进行
	- PausedSyncS：以本地节点是一个持续同步的源，但是目前同步已经暂停，可能是因为另外一个同步正在进行或是使用命令(drbdadm pause-sync)暂停了同步
	- PausedSyncT：以本地节点为持续同步的目标，但是目前同步已经暂停，这可以是因为另外一个同步正在进行或是使用命令(drbdadm pause-sync)暂停了同步
	- VerifyS：以本地节点为验证源的线上设备验证正在执行
	- VerifyT：以本地节点为验证目标的线上设备验证正在执行

- 启用/禁用资源

	```
	//启用资源r0
	# drbdadm up r0
	//禁用资源r0
	# drbdadm down r0
	提示:也可以将r0更改为all
	```
- 升级和降级资源
	```
	升级资源
	# drbdadm primary <resource>
	降级资源
	# drbdadm secondary <resource>
	注意：在单主模式下的DRBD，两个节点同时处于连接状态，任何一个节点都可以在特定的时间内变成主；但两个节点中只能一为主，如果已经有一个主，需先降级才可能升级；在双主模式下没有这个限制


- 重新配置资源

	```
	DRBD在运行时，允许用户重新配置资源，为了实现这个目的，需要进行以下操作：
	1、在DRBD的配置文件/etc/drbd.conf(包括所有资源)中进行有必要的改变
	2、在两个节点之间同步DRBD的配置文件
	3、在两个节点上执行drbdadm adjust <source>命令 (在执行此命令时，建议添加-d参数)
	```

- 导出当前资源配置信息

	```
	drbdadm dump all
	```

- 连接与断开

	```
	# drbdadm connect r0
	# drbdadm disconnect r0
	```

## 配置 ##

[http://blog.chinaunix.net/uid-20346344-id-3491536.html](http://blog.chinaunix.net/uid-20346344-id-3491536.html)

**global_common.conf**：

- global：
	在配置文件中global段仅出现一次，且若所有配置信息都保存至同一个配置文件中而不分开为多个文件的话则global段必须位于配置文件最开始处，目前global段中可以定义的参数仅有minor-count, dialog-refresh, disable-ip-verification和usage-count。
	- minor-count：从（设备）个数，取值范围1~255，默认值为32。该选项设定了允许定义的resource个数，当要定义的resource超过了此选项的设定时，需要重新载入drbd内核模块。
	- dialog-refresh time：time取值0，或任一正数。默认值为1。
	- disable-ip-verification：是否禁用ip检查
	- usage-count：是否参加用户统计，合法参数为yes、no或ask。
	
- common：
    用于定义被每一个资源默认继承的参数，可在资源定义中使用的参数都可在common定义。实际应用中common段并非必须但建议将多个资源共享的参数定义在common段以降低配置复杂度，common配置段中可以包含如下配置段：disk、net、startup、syncer和handlers。（参考后面的**common配置项**）
	- startup： 用来更加精细地调节drbd属性，它作用于配置节点在启动或重启时。
	- disk： 用来精细地调节drbd底层存储的属性。
	- syncer： 用来更加精细地调节服务的同步进程。
	- net： 用来精细地调节drbd的网络相关的属性。

- resource：
    用于定义drbd资源，每个资源通常定义在一个单独的位于/etc/drbd.d目录中的以.res结尾的文件中。资源在定义时必须为其命名，每个资源段的定义中至少要包含两个host子段，以定义此资源关联至的节点，其它参数均可以从common段或drbd的默认中进行继承而无须定义。

**xxx.res**：

```
resource Mysqls {     　                   #resource关键字指定资源名称为Mysql    注：resource段一般写入*.res结尾的文件
　　　protocol C;          　               #使用的协议类型
　　　  meta-disk internal;                 #meta-data和数据存放在同一个底层
 　　   disk {  on-io-error detach; }       #当磁盘出现IO错误时如何处理
 　　   startup { degr-wfc-timeout 60;  }   #启动时连接资源的超时时间
 　　   on Mysql1 {         　              #集群中的其中一个节点：Mysql1
　　        device    /dev/drbd1;           #物理设备的逻辑路径（参数最后必须有数字，用于global的minor-count）
　　        disk     /dev/sda1;             #物理设备  
　　        address 10.0.0.7:7788;          #监听地址和端口，用于与另一台主机通信（主从的角色由drbdadm命令指定）
            meta-disk  internal;            #定义metadata的存储方式，有2种（参考本节**metadata存储方式**）

　　     }
　　     on Mysql2 {       　             　#集群中的其中一个节点：Mysql2
 　　       device    /dev/drbd1;     
 　　       disk     /dev/sda1;         
  　　      address 10.0.0.8:7788;         #监听地址和端口，用于与另一台主机通信（主从的角色由drbdadm命令指定）双方配置文件中均需写入各自与对方的监听地址

  　　      } 
　　　　}

```

**metadata存储方式**：

- Internal metadata：
	一个resource被配置成使用internal metadata，意味着DRBD把它的metadata，和实际生产数据存储于相同的底层物理设备中。该存储方式是在设备的最后位置留出一个区域来存储metadata。
	优点：因为metadata是和实际生产数据紧密联系在一起的，如果发生了硬盘损坏，不需要管理员做额外的工作，因为metadata会随实际生产数据的丢失而丢失，同样会随着生产数据的恢复而恢复。
	缺点：如果底层设备只有一块物理硬盘（和RAID相反），这种存储方式对写操作的吞吐量有负面影响，因为应用程序的写操作请求会触发DRBD的metadata的更新。如果metadata存储于硬盘的同一块盘片上，那么，写操作会导致额外的两次磁头读写移动。
	要注意的是：如果你打算在已有数据的底层设备中使用internal metadata，需要计算并留出DRBD的metadata所占的空间大小，并采取一些特殊的操作，否则很有可能会破坏掉原有的数据！至于需要什么样的 特殊操作，可以参考DRBD的官方文档。我要说的是，最好不要这样做！

- external metadata：
	该存储方式比较简单，就是把metadata存储于一个和生产数据分开的专门的设备块中。
	优点：对某些写操作，提供某些潜在的改进。
	缺点：因为metadata和生产数据是分开的，如果发生了硬盘损坏，在更换硬盘后，需要管理员进行人工干预，从其它存活的节点向刚替换的硬盘进行完全的数据同步。
	什么时候应该使用exteranl的存储方式：设备中已经存有数据，而该设备不支持扩展（如LVM），也不支持收缩（shrinking）。

**common配置项**

- **startup**

	- **wfc-timeout**：该选项设定一个时间值，单位是秒。在启用DRBD块时，初始化脚本drbd会阻塞启动进程的运行，直到对等节点的出现。该选项就是用来限制这个等待时间的，默认为0，即不限制，永远等待。
	- **degr-wfc-timeout**：该选项也设定一个时间值，单位为秒。也是用于限制等待时间，只是作用的情形不同：它作用于一个降级集群（即那些只剩下一个节点的集群）在重启时的等待时间。
	- **outdated-wfc-timeout**：同上，也是用来设定等待时间，单位为秒。它用于设定等待过期节点的时间。

- **disk（重要）**

	- **on-io-error选项**：此选项设定了一个策略，如果底层设备向上层设备报告发生I/O错误，将按照该策略进行处理。有效的策略包括：
		- **pass_on**： 把I/O错误报告给上层设备。如果错误发生在primary节点，把它报告给文件系统，由上层设备处理这些错误（例如，它会导致文件系统以只读方式重新挂载），它可能会导致drbd停止提供服务；如果发生在secondary节点，则忽略该错误（因为secondary节点没有上层设备可以报告）。该策略曾经是默认策略，但现在已被detach所取代。
		- **call-local-io-error**：调用预定义的本地local-io-error脚本进行处理。该策略需要在resource配置段的handlers部分，预定义一个相应的`local-io-error`命令调用。该策略完全由管理员通过local-io-error命令（或脚本）调用来控制如何处理I/O错误。
		- **detach**：发生I/O错误的节点将放弃底层设备，以diskless mode继续工作。在diskless mode下，只要还有网络连接，drbd将从secondary node读写数据，而不需要failover。该策略会导致一定的损失，但好处也很明显，drbd服务不会中断。官方推荐和默认策略。
	- **fencing选项**：该选项设定一个策略来避免split brain的状况。有效的策略包括：
		- **dont-care**：默认策略。不采取任何隔离措施。
		- **resourceonly**：在此策略下，如果一个节点处于split brain状态，它将尝试隔离对等端的磁盘。这个操作通过调用fence-peer处理器来实现。fence-peer处理器将通过其它通信路径到达对等节点，并在这个对等节点上调用`drbdadm outdate res`命令。
		- **resource-and-stonith**：在此策略下，如果一 个节点处于split brain状态，它将停止I/O操作，并调用fence-peer处理器。处理器通过其它通信路径到达对等节点，并在这个对等节点上调用`drbdadm outdate res`命令。如果无法到达对等节点，它将向对等端发送关机命令。一旦问题解决，I/O操作将重新进行。如果处理器失败，你可以使用`resume-io`命令来重新开始I/O操作。

- **syncer**

	- **rate**：设置同步时的速率，默认为250KB。默认的单位是KB/sec，也允许使用K、M和G，如40M。注意：syncer中的速率是以bytes，而不是bits来设定的。配置文件中的这个选项设置的速率是永久性的，但可使用命令`drbdsetup /dev/drbdnum syncer -r 100M`临时地改变rate的值(只能在所有节点的其中一个节点上运行此命令。)如果想重新恢复成drbd.conf配置文件中设定的速率，执行命令`drbdadm adjust resource`。注意：速率的设定，最好设为有效可用带宽的30%。所谓有效可用带宽，是指网络带宽和磁盘读写速度中的最小者
	- **al-extents**：该选项用来设定hot area（即active set）的块数，取值范围是7~3843，默认为127，每个块标志4M大小的底层存储（即底层设备）。DRBD会自动检测host area，如果主节点意外地从集群中断开，当该节点重新加入集群时，hot area所覆盖的区域必须被重新同步。hot area的每次变化，实际上都是对metadata区域的写操作，因此，该选项的值越大，重新同步的时间越长，但需要更新的meta-data也越少。
	- **verify-alg**：该选项指定一个用于在线校验的算法，内核一般都会支持md5、sha1和crc32c校验算法。在线校验默认关闭，必须在此选项设定参数，以明确启用在线设备校验。 DRBD支持在线设备校验，它以一种高效的方式对不同节点的数据进行一致性校验。在线校验会影响CPU负载和使用，但影响比较轻微。drbd 8.2.5及以后版本支持此功能。一旦启用了该功能，就可以使用下列命令进行一个在线校验： `drbdadm verify resource`,该命令对指定的resource进行检验，如果检测到有数据块没有同步，它会标记这些块，并往内核日志中写入一条信息。这个过程不会影响正在使用该设备的程序。
	- **csums-alg**：该选项指定一个校验算法，用来标志数据块。如果不启用该选项，resync会从source发送所有的数据块到destination；而如果启用了此选项，那么resync将只交换那些校验值不同的数据块，当网络带宽有限时，此选项非常有用。而且，在重启一个崩溃的primary节点时，该选项 会降低CPU带宽的占用。

- **net**

	- **sndbuf-size**：该选项用来调节TCP send buffer的大小，drbd 8.2.7以前的版本，默认值为0，意味着自动调节大小；新版本的drbd的默认值为128KiB。高吞吐量的网络（例如专用的千兆网卡，或负载均衡中绑 定的连接）中，增加到512K比较合适，或者可以更高，但是最好不要超过2M。
	- **timeout**：该选项设定一个时间值，单位为0.1秒。如果搭档节点没有在此时间内发来应答包，那么就认为搭档节点已经死亡，因此将断开这次TCP/IP连接。默认值为60，即6秒。该选项的值必须小于connect-int和ping-int的值。
	- **connect-int**：如果无法立即连接上远程DRBD设备，系统将断续尝试连接。该选项设定的就是两次尝试间隔时间。单位为秒，默认值为10秒。
	- **ping-int**：该选项设定一个时间值，单位为秒。如果连接到远程DRBD设备的TCP/IP的空闲时间超过此值，系统将生成一个keep-alive包来检测对等节点是否还存活。默认值为10秒。
	- **ping-timeout**：该选项设定一个时间值，单位是0.1秒。如果对等端没有在此时间内应答keep-alive包，它将被认为已经死亡。默认值是500ms。
	- **max-buffers**：该选项设定一个由drbd分配的最大请求数，单位是页面大小（PAGE_SIZE），大多数系统中，页面大小为4KB。这些buffer用来存储那些即将写入磁盘的数据。最小值为32（即128KB）。这个值大一点好。
	- **max-epoch-size**：该选项设定了两次write barriers之间最大的数据块数。如果选项的值小于10，将影响系统性能。大一点好。