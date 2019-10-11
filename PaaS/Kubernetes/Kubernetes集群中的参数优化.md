# Kubernetes集群中的参数优化 #

[摘自 @zhuhaiqiang  5个维度对Kubernetes集群优化，写得非常棒](http://www.talkwithtrend.com/Article/246185)

## 一、节点配额和内核参数调整 ##

大规模集群相关配额：

- 虚拟机个数
- vCPU个数
- 内网IP地址个数
- 公网IP地址个数
- 安全组条数
- 路由表条数
- 持久化存储大小

GCE随着node节点的增加master节点的配置：

- 1-5 nodes: n1-standard-1
- 6-10 nodes: n1-standard-2
- 11-100 nodes: n1-standard-4
- 101-250 nodes: n1-standard-8
- 251-500 nodes: n1-standard-16
- more than 500 nodes: n1-standard-32

阿里云随着node节点的增加master节点的配置：：

- 1-5个节点: 4C8G(不建议2C4G)
- 6-20个节点: 4C16G
- 21-100个节点: 8C32G
- 100-200个节点: 16C64G

系统参数配置：

- fs.file-max=1000000
	max-file表示系统级别的能够打开的文件句柄的数量，一般如果遇到文件句柄达到上限时，会碰到"Too many open files"或者Socket/File: Can’t open so many files等错误。
- net.ipv4.neigh.default.gc_thresh1=1024
	存在于ARP高速缓存中的最少层数，如果少于这个数，垃圾收集器将不会运行。缺省值是128。
- net.ipv4.neigh.default.gc_thresh2=4096
	保存在ARP高速缓存中的最多的记录软限制。垃圾收集器在开始收集前，允许记录数超过这个数字5秒。缺省值是512。
- net.ipv4.neigh.default.gc_thresh3=8192
	保存在ARP高速缓存中的最多记录的硬限制，一旦高速缓存中的数目高于此，垃圾收集器将马上运行。缺省值是1024。
- net.netfilter.nf_conntrack_max=10485760
	允许的最大跟踪连接条目，是在内核内存中netfilter可以同时处理的“任务”（连接跟踪条目），默认 `nf_conntrack_buckets * 4`。
	[http://www.mamicode.com/info-detail-2422830.html](http://www.mamicode.com/info-detail-2422830.html)
- net.netfilter.nf_conntrack_tcp_timeout_established=300
	established的超時時間
- net.netfilter.nf_conntrack_buckets=655360
	哈希表大小（只读）（64位系统、8G内存默认65536，16G翻倍，以此类推）
- net.core.netdev_max_backlog=10000
	每个网络接口接收数据包的速率比内核处理这些包的速率快时，允许送到队列的数据包的最大数目。
- fs.inotify.max_user_instances=524288
	每一个real user ID可创建的inotify instatnces的数量上限，默认为128
- fs.inotify.max_user_watches=524288
	每个inotify instance相关联的watches的上限，默认值8192
	


	
