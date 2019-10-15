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

- `fs.file-max=1000000`
	max-file表示系统级别的能够打开的文件句柄的数量，一般如果遇到文件句柄达到上限时，会碰到"Too many open files"或者Socket/File: Can’t open so many files等错误。
- `net.ipv4.neigh.default.gc_thresh1=1024`
	存在于ARP高速缓存中的最少层数，如果少于这个数，垃圾收集器将不会运行。缺省值是128。
- `net.ipv4.neigh.default.gc_thresh2=4096`
	保存在ARP高速缓存中的最多的记录软限制。垃圾收集器在开始收集前，允许记录数超过这个数字5秒。缺省值是512。
- `net.ipv4.neigh.default.gc_thresh3=8192`
	保存在ARP高速缓存中的最多记录的硬限制，一旦高速缓存中的数目高于此，垃圾收集器将马上运行。缺省值是1024。
- `net.netfilter.nf_conntrack_max=10485760`
	允许的最大跟踪连接条目，是在内核内存中netfilter可以同时处理的“任务”（连接跟踪条目），默认 `nf_conntrack_buckets * 4`。
	[http://www.mamicode.com/info-detail-2422830.html](http://www.mamicode.com/info-detail-2422830.html)
- `net.netfilter.nf_conntrack_tcp_timeout_established=300`
	established的超時時間
- `net.netfilter.nf_conntrack_buckets=655360`
	哈希表大小（只读）（64位系统、8G内存默认65536，16G翻倍，以此类推）
- `net.core.netdev_max_backlog=10000`
	每个网络接口接收数据包的速率比内核处理这些包的速率快时，允许送到队列的数据包的最大数目。
- `fs.inotify.max_user_instances=524288`
	每一个real user ID可创建的inotify instatnces的数量上限，默认为128
- `fs.inotify.max_user_watches=524288`
	每个inotify instance相关联的watches的上限，默认值8192
	

## 二、镜像拉取相关配置 ##

**Docker配置**
- `max-concurrent-downloads=10`
	配置每个pull操作的最大并行下载数，提高镜像拉取效率，默认值是3。
- 使用SSD存储。
- 预加载pause镜像

**Kubelet配置**
- `--serialize-image-pulls=false`
	该选项配置串行拉取镜像，默认值时true，配置为false可以增加并发度。 但是如果docker daemon版本小于1.9，且使用aufs存储则不能改动该选项。
- `--image-pull-progress-deadline=30`
	配置镜像拉取超时。 默认值时1分，对于大镜像拉取需要适量增大超时时间。
- `--max-pods=110`
	Kubelet单节点允许运行的最大Pod数,默认是110，可以根据实际需要设置。
**镜像registry p2p分发**

## Kubernetes配置 ##

**Kube APIServer配置**

node节点数量 >=3000， 推荐设置如下配置：

- `--max-requests-inflight=3000`
- `--max-mutating-requests-inflight=1000`

node节点数量在1000--3000， 推荐设置如下配置：

- `--max-requests-inflight=1500`
- `--max-mutating-requests-inflight=500`

内存配置选项和node数量的关系，单位是MB：

- `--target-ram-mb=node_nums * 60`

**Kube-scheduler配置**

- `--kube-api-qps=100`
	默认值是50

**Kube-controller-manager配置**

- `--kube-api-qps=100`
	默认值是20
- `--kube-api-burst=100`
	默认值是30


	
