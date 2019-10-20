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

	

- dd

	仅仅是对文件进行读写，没有模拟应用、业务、场景的效果

- xdd

- iorate
	
	[https://manned.org/iorate/ff1b0b2d](https://manned.org/iorate/ff1b0b2d)

- iostat

	

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

