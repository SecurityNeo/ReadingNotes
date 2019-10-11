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





- dd
	仅仅是对文件进行读写，没有模拟应用、业务、场景的效果
- xdd
- iorate
- iozone
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
- 总IOPS：				命令行iostat -Dl：tps
- 每个盘对应的读IOPS ：	命令行iostat -Dl：rps
- 每个盘对应的写IOPS ：	命令行iostat -Dl：wps
