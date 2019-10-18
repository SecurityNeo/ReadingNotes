# IO评估及工具 #

[https://mp.weixin.qq.com/s/AZnvrso-HfaqlaRJ_EsVGg](https://mp.weixin.qq.com/s/AZnvrso-HfaqlaRJ_EsVGg)

## IO模型 ##

- IOPS
- 带宽
- IO的尺寸（大小）
- 磁盘IO分别在哪些盘（磁盘）
- 读IO和写IO的比例（磁盘）
- 读IO是顺序的还是随机的（磁盘）
- 写IO是顺序的还是随机的（磁盘）

## 评估工具 ##

**磁盘IO**

- orion

	oracle出品，模拟Oracle数据库IO负载

- iometer

	[http://www.iometer.org/](http://www.iometer.org/)。工作在单系统和集群系统上用来衡量和描述I/O子系统的工具。

- Nmon

	[http://nmon.sourceforge.net/pmwiki.php](http://nmon.sourceforge.net/pmwiki.php)。[nmon analyse可视化展示](https://www.ibm.com/developerworks/community/wikis/home?lang=en#!/wiki/Power+Systems/page/nmon_analyser)
	[https://blog.csdn.net/saraul/article/details/8570781](https://blog.csdn.net/saraul/article/details/8570781)

	Nmon分析文件中各sheet的含义：

		SYS_SUMM    		系统汇总,蓝线为cpu占有率变化情况,粉线为磁盘IO的变化情况
		AAA             	关于操作系统以及nmon本身的一些信息
		CPUnn        		显示执行时间内CPU占用情况，其中包含user%、sys%、wait%和idle%
		CPU_ALL         	所有CPU概述，显示所有CPU平均占用情况
		CPU_SUMM    		每一个CPU在执行时间内的占用情况，其中包含user%、sys%、wait%和idle
		DGBUSY				磁盘组每个hdisk设备平均占用情况
		DGREAD        		每个磁盘组的平均读情况
		DGSIZE				每个磁盘组的平均读写情况（块大小）
		DGWRITE     		每个磁盘组的平均写情况
		DGXFER				每个磁盘组的I/O每秒操作
		MEM					内存相关的主要信息，使用、空闲内存大小等
		NET					本sheet显示系统中每个网络适配器的数据传输速率（千字节/秒）
		MEMUSE				除`%comp`参数外,本sheet包含的所有项都和vmtune命令的报告中一样
		MEMNEW				本sheet显示分配的内存片信息，分三大类：用户进程使用页，文件系统缓存，系统内核使用页
		NETPACKET    		本sheet统计每个适配器网络读写包的数量；这个类似于`netpmon  –O dd`命令
		PAGE        		本sheet统计相关页信息的记录
		BBBB        		系统外挂存储容量以及存储类型
		BBBC        		系统外挂存储位置、状态以及描述信息
		BBBD        		磁盘适配器信息；（包含磁盘适配器名称以及描述）
		BBBE        		包含通过lsdev命令获取的系统设备及其特征，显示vpaths和hdisks之间的映射关系
		BBBG        		显示磁盘组详细的映射关系
		BBBL        		逻辑分区（LPAR）配置细节信息
		BBBN        		网络适配器信息
		BBBP        		vmtune,  schedtune, emstat和lsattr命令的输出信息
		DISKBSIZE    		执行时间内每个hdisk的传输块大小
		DISKBUSY    		每个hdisk设备平均占用情况
		DISKREAD    		每个hdisk的平均读情况
		DISKWRITE    		每个hdisk的平均写情况
		DISKXFER    		每个hdisk的I/O每秒操作
		DISKSERV    		本sheet显示在每个收集间隔中hdisk的评估服务时间（未响应时间）
		DISK_SUMM    		总体disk读、写以及I/O操作
		JFSFILE				本sheet显示对于每一个文件系统中，在每个间隔区间正在被使用的空间百分比
		PROC				本sheet包含nmon内核内部的统计信息。其中RunQueue和Swap-in域是使用的平均时间间隔，其他项的单位是比率/秒
		PROCAIO        		本sheet包含关于可用的和active的异步IO进程数量信息
		TOP					
		ZZZZ        		本sheet自动转换所有nmon的时间戳为现在真实的时间，方便更容易的分析
		EMCBSIZE/FAStBSIZE	执行时间内EMC存储的传输块大小
		EMCBUSY/FAStBUSY	EMC存储设备平均占用情况
		EMCREAD/FAStREAD	EMC存储的平均读情况
		EMCWRITE/FAStWRITE	EMC存储的平均写情况
		ESSBSIZE			本sheet记录在系统中每个vpaths下读写操作的平均数据传输大小  (blocksize) Kbytes
		ESSBUSY				本sheet记录使用ESS系统的每个vpaths下的设备繁忙情况
		ESSREAD				本sheet记录在系统中每个vpaths下读取操作的  data rate (Kbytes/sec)
		ESSWRITE			本sheet记录在系统中每个vpaths下写入操作的  data rate (Kbytes/sec)
		ESSXFER				本sheet记录在系统中每个vpaths下每秒的IO操作
		ESSSERV				本sheet显示在每个收集间隔中vpaths的评估服务时间（未响应时间）
		FILE				本sheet包含nmon内核内部的统计信息的一个子集，跟sar报告的值相同
		IOADAPT				对于BBBCsheet每个IO适配器列表，包含了数据传输速度为读取和写入操作（千字节/秒）和I  / O操作执行的总数量
		JFSFILE				本sheet显示对于每一个文件系统中，在每个间隔区间正在被使用的空间百分比
		JFSINODE			本sheet显示对于每一个文件系统中，在每个间隔区间正在被使用的inode百分比
		LARGEPAGE			本图表显示Usedpages和Freepages随着时间的变化

	指标详解：

		- SYS_SUMM：
			- CPU%		cpu占有率变化情况
			- IO/sec	IO的变化情况
			
		- AAA：
			- AIX		AIX版本号
			- build		build版本号
			- command	执行命令
			- cpus		CPU数量
			- date		执行日期
			- hardware	被测主机处理器技术
			- host		被测主机名
			- interval	监控取样间隔（秒）
			- kernel	被测主机内核信息
			- ML		维护等级
			- progname	执行文件名称
			- runname	运行主机名称
			- snapshots	实际快照次数
			- subversion	nmon版本详情
			- time		执行开始时间戳
			- user		执行命令用户名
			- version	收集数据的nmon版本
			- analyser	nmon analyser版本号
			- environment	所用excel版本
			- parms		excel参数设定
			- settings	excel环境设置
			- elapsed	生成excel消耗时间
			
		- BBBB
			- name		存储磁盘名称
			- size(GB)	磁盘容量
			- disc attach type	磁盘类型
			
		- BBBC
			- hdisknn	各个磁盘信息、状态以及MOUNT位置
			
		- BBBD
			- Adapter_number	磁盘适配器编号
			- Name		磁盘适配器名称
			- Disks		磁盘适配器数量
			- Description	磁盘适配器描述
			
		- BBBN
		
			- NetworkName	网络名称
			- MTU			网络上传送的最大数据包，单位是字节
			- Mbits			带宽
			- Name			名称

		- CPUnn
		
			- CPU nn	执行间隔时间列表
			- User%		显示在用户模式下执行的程序所使用的CPU百分比
			- Sys%		显示在内核模式下执行的程序所使用的CPU百分比
			- Wait%		显示等待IO所花的时间百分比
			- Idle%		显示CPU的空闲时间百分比
			- CPU%		CPU总体占用情况

		- DISKBSIZE
		
			- Disk Block Size Hostname		执行间隔时间列表
			- hdisknn		磁盘传输速度时间间隔采样（读和写的总趋势图）

		- DISKBUSY
		
			- Disk %Busy Hostname		执行间隔时间列表
			- hdisknn		每个磁盘执行采样数据（磁盘设备的占用百分比）

		- DISKREAD
		
			- Disk Read kb/s Hostname		执行间隔时间列表
			- hdisknn		每个磁盘执行采样数据（磁盘设备的读速率）

		- DISKWRITE
		
			- Disk Write kb/s Hostname		执行间隔时间列表
			- hdisknn		每个磁盘执行采样数据（磁盘设备的写速率）

		- DISKXFER
		
			- Disk transfers per second Hostname		执行间隔时间列表
			- hdisknn		每秒钟输出到物理磁盘的传输次数

		- DISK_SUMM
		
			- Disk total kb/s  Hostname		执行间隔时间列表
			- Disk Read kb/s		每个磁盘执行采样数据（磁盘设备的读速率）
			- Disk Write kb/s		每个磁盘执行采样数据（磁盘设备的写速率）
			- IO/sec		每秒钟输出到物理磁盘的传输次数

		- FILE
		
			- iget		在监控期间每秒钟到节点查找例行程序的呼叫数
			- namei		在监控期间每秒钟路径查找例行程序的呼叫数（sar -a）                                                                                                                                              
			- dirblk	在监控期间通过目录搜索例行程序每秒钟扫描到的目录块数  （sar -a）                                                                                                                     
			- readch	在监控期间通过读系统呼叫每秒钟读出的字节数（sar  -c）                                                                                                                                                                                                     
			- writech	在监控期间通过写系统呼叫每秒钟写入的字节数（sar  -c）                                                                                                                                                                                          
			- ttyrawch	在监控期间通过TTYs每秒钟读入的裸字节数（sar  -y ）.                                                                                                                                                                                                                    
			- ttycanch	终端输入队列字符。对于aix Version 4或者更后的版本这个值总是0                                                                                                                                                                                
			- ttyoutch	终端输出队列字符（sar -y ）

		- IOADAPT
		
			- Disk Adapter Hostname(KB/s)		执行间隔时间列表；
			- Disk Adapter_read			磁盘适配器读速率
			- Disk Adapter_write		磁盘适配器写速率
			- Disk  Adapter_xfer-tps	磁盘适配器传输速率（该物理磁盘每秒的IO传输请求数量）

		- JFSFILE
		
			- JFS Filespace %Used Hostname		执行间隔时间列表
			- file  system/LV		文件系统以及mount磁盘设备已使用空间百分比

		- JFSINODE

			- JFS Inode %Used Hostname	执行间隔时间列表
			- file  system/LV		文件系统以及mount磁盘设备的inode已使用空间百分比

		- MEM

			- Memory Hostname		执行间隔时间列表
			- Real Free %		实际剩余内存百分比
			- Virtual free %		虚拟剩余内存百分比
			- Real free(MB)		实际剩余内存大小（MB）
			- Virtual free(MB)		虚拟剩余内存大小（MB）
			- Real total(MB)		实际内存总体大小（MB）
			- Virtual  total(MB)		虚拟内存总体大小（MB）

		- MEMUSE

			- %numperm		分配给文件页的实际内存百分比
			- %minperm		mixperm的缺省值约为20%的物理内存.通常会不断的运行，除非vmtune或rmss命令中使用收集
			- %maxperm		maxperm的缺省值约为80%的物理内存.  通常会不断的运行，除非vmtune或rmss命令中使用收集
			- minfree		空闲页面数的最小值
			- maxfree		空闲页面数的最大值，指定的vmtune命令或系统默认
			- %comp		分配给计算页的内存百分比,NMON分析器计算这个值。计算页是可被page space支持的，包括存储和程序文本段 他们不包括数据，可执行的和共享的库文件

		- MEMNEW

			- Process%		分配给用户进程的内存百分比
			- FSCache%		分配给文件系统缓存的内存百分比
			- System%		系统程序使用的内存百分比
			- Free%		未被分配的内存百分比           
			- User%		非系统程序使用的内存百分比
			
		- NET

			- read/write		显示系统中每个网络适配器的数据传输速率（千字节/秒）

		- NETPACKET

			- reads/s		统计每个适配器网络读包的数量
			- writes/s		统计每个适配器网络写包的数量

		- PAGE

			- aults		每秒的page faults数
			- pgin		每秒钟所读入的页数，包括从文件系统读取的页数
			- pgout		每秒钟所写出的页数，包括写到文件系统的页数
			- pgsin		每秒钟从页面空间所读取的页数
			- pgsout	每秒钟写到页面空间的页数
			- reclaims		从nmon回收这项之前的10个，和vmstat报告的值是一样的，代表了页替换机制释放的pages/sec的数量
			- scans		扫描页替换机制的pages/sec的数量，和vmstat报告的值是一样的，页替换在空闲页数量到达最小值时初始化，在空闲到达最大值时停止
			- cycles		周期times/sec的数值，页替换机制需要扫描整个页表，来补充空闲列表。这和vmstat报告的cy数值一样，只是vmstat报告的这个值是整形值，而nmon报告的是实型值
			- fsin		分析器计算的数据为pgin-pgsin的图形处理所用
			- fsout		分析器计算的数据为pgout-pgsout的图形处理所用
			- sr/fr		分析器计算的数据为scans/reclaims的图形处理所用  

		- PROC

			- RunQueue		运行队列中的内核线程平均数（同`sar -q`中的runq-sz）
			- Swap-in		等待page in的内核线程平均数（同`sar -q`中的swpq-sz）
			- pswitch		上下文开关个数（同`sar -w`中的pswch/s）
			- syscall		系统调用总数（同`sar -c`中的scall/s）
			- read		系统调用中read的数量（同`sar -c`中的sread/s）
			- write		系统调用中write的数量（同`sar -c`中的swrit/s）
			- fork		系统调用中fork的数量（同`sar -c`中的fork/s）
			- exec		系统调用中exec 的数量（同`sar -c`中的exec/s）
			- rcvint		tty接收中断的数量（同`sar -y`中的revin/s）
			- xmtint		tty传输中断的数量（同`sar -y`中的xmtin/s）
			- sem		IPC信号元的数量 创建，使用和消除)（同`sar -m`中的sema/s）
			- msg		IPC消息元的数量 (发送和接收)（同`sar -m`中的sema/s）

		- TOP

			- PID		进程号
			- %CPU		CPU使用的平均数
			- %Usr		显示运行的用户程序所占用的CPU百分比
			- %Sys		显示运行的系统程序所占用的CPU百分比
			- Threads		被使用在这个程序中的线程数
			- Size		对于这个程序一次调用分配给数据段的paging  space平均值                                                                                             
			- ResText		对于这个程序一次调用分配给代码段的内存平均值                                                                                                                                                                                                                                       
			- ResData		对于这个程序一次调用分配给数据段的内存平均值
			- CharIO		通过读写系统调用的每秒字节数                                                                                                                                                              
			- %RAM		此命令所使用的内存百分比（`(ResText + ResData) / Real Mem`）                                                                                                                           
			- Paging		此进程所有page  faults的总数                                                                                                                                                                                                                                                      
			- Command		命令名称                                                                                                                                                                                                                                                                                                                                                                                                        
			- WLMClass		此程序已分配的Workload Manager superclass名称                                                                                                                                                                                                                                                                                                                    
			- IntervalCPU		详细信息中显示在时间间隔中所有调用命令所使用的CPU总数  
			- WSet		详细信息中显示在时间间隔中所有调用命令所使用的内存总数                                                                                                                                                
			- User		运行进程的用户名                                                                                                                                                                                                                                                                                                              
			- Arg		包含完整的参数字符串输入命令


- dd

	仅仅是对文件进行读写，没有模拟应用、业务、场景的效果

- xdd

- iorate
	
	[https://manned.org/iorate/ff1b0b2d](https://manned.org/iorate/ff1b0b2d)

- iostat

	- %user：CPU处在用户模式下的时间百分比
	- %nice：CPU处在带NICE值的用户模式下的时间百分比
	- %system：CPU处在系统模式下的时间百分比
	- %iowait：CPU等待输入输出完成时间的百分比
	- %steal：管理程序维护另一个虚拟处理器时，虚拟CPU的无意识等待时间百分比
	- %idle：CPU空闲时间百分比
	- tps：该设备每秒的传输次数
	- kB_read/s：每秒从设备（drive expressed）读取的数据量
	- kB_wrtn/s：每秒向设备（drive expressed）写入的数据量
	- kB_read：  读取的总数据量
	- kB_wrtn：写入的总数量数据量
	- rrqm/s: 每秒对该设备的读请求被合并次数，文件系统会对读取同块(block)的请求进行合并
	- wrqm/s: 每秒对该设备的写请求被合并次数
	- r/s: 每秒完成的读请求次数
	- w/s: 每秒完成的写请求次数
	- rkB/s: 每秒读数据量(kB为单位)
	- wkB/s: 每秒写数据量(kB为单位)
	- avgrq-sz:平均每次IO操作的数据量(扇区数为单位)
	- avgqu-sz: 平均等待处理的IO请求队列长度
	- await: 平均每次IO请求等待时间(包括等待时间和处理时间，毫秒为单位)
	- svctm: 平均每次IO请求的处理时间(毫秒为单位)
	- %util: 采用周期内用于IO操作的时间比率，即IO队列非空的时间比率

- iozone

	[www.iozone.org](www.iozone.org)
	常用参数：
		-a 全面测试，比如块大小它会自动加
		-i N 用来选择测试项, 比如Read/Write/Random 比较常用的是0 1 2,可以指定成-i 0 -i 1 -i2.这些别的详细内容请查man

			0=write/rewrite
			1=read/re-read
			2=random-read/write
			3=Read-backwards
			4=Re-write-record
			5=stride-read
			6=fwrite/re-fwrite
			7=fread/Re-fread
			8=random mix
			9=pwrite/Re-pwrite
			10=pread/Re-pread
			11=pwritev/Re-pwritev
			12=preadv/Re-preadv

		-r block size 指定一次写入/读出的块大小
		-s file size 指定测试文件的大小
		-f filename 指定测试文件的名字,完成后会自动删除(这个文件必须指定你要测试的那个硬盘中)
		-F file1 file2... 指定多线程下测试的文件名
		-R 代表生成Excel报告文件
		-c 代表每次读写测试完毕都发送关闭连接的命令，主要用于测试NFS系统
		-q 代表最大的记录大小
		-g 代表最大的文件大小
		-n 代表最小的文件大小
		-b 输出的生成的Excel报告文件名字

- postmark

	可以实现文件读写、创建、删除这样的操作。适合小文件应用场景的测试


**网络IO**

- ping：最基本的，可以指定包的大小。
- iperf、ttcp：测试tcp、udp协议最大的带宽、延时、丢包。
- NTttcp（Windows）
- LANBench（Windows）
- pcattcp（Windows）
- LAN Speed Test (Lite)（Windows）
- NETIO（Windows）
- NetStress（Windows）

## 监控指标和工具 ##

### 磁盘IO ###

**IOPS**

- 总IOPS：				Nmon DISK_SUMM Sheet：IO/Sec
- 每个盘对应的读IOPS ：	Nmon DISKRIO Sheet
- 每个盘对应的写IOPS ：	Nmon DISKWIO Sheet
- 总IOPS：				命令行`iostat -Dl：tps`
- 每个盘对应的读IOPS ：	命令行`iostat -Dl：rps`
- 每个盘对应的写IOPS ：	命令行`iostat -Dl：wps`

**带宽**

- 总带宽：	Nmon DISK_SUMM Sheet：Disk Read KB/s，Disk Write KB/s
- 每个盘对应的读带宽：	Nmon DISKREAD Sheet
- 每个盘对应的写带宽：	Nmon DISKWRITE Sheet
- 总带宽：	命令行`iostat -Dl`：bps
- 每个盘对应的读带宽：	命令行`iostat -Dl`：bread
- 每个盘对应的写带宽：	命令行`iostat -Dl`：bwrtn

**响应时间**

- 每个盘对应的读响应时间：	命令行`iostat -Dl`：read avg serv，max serv
- 每个盘对应的写响应时间：	命令行`iostat -Dl`：write avg serv，max serv

### 网络IO ###

**带宽**

- Nmon：	NET Sheet
- 命令行topas：	Network：BPS、B-In、B-Out

