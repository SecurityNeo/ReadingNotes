# Redis监控（INFO） #

`10.125.31.101:9021> info all`

**服务器信息**

- redis_version:3.0.0                              #redis服务器版本
- redis_git_sha1:00000000						   #Git SHA1
- redis_git_dirty:0                                #Git dirty flag
- redis_build_id:6c2c390b97607ff0                  #redis build id
- redis_mode:cluster                               #运行模式，单机或者集群
- os:Linux 2.6.32-358.2.1.el6.x86_64 x86_64        #redis服务器的宿主操作系统
- arch_bits:64                                     #架构(32或64位)
- multiplexing_api:epoll                           #redis所使用的事件处理机制
- gcc_version:4.4.7                                #编译redis时所使用的gcc版本
- process_id:12099                                 #redis服务器进程的pid
- run_id:63bcd0e57adb695ff0bf873cf42d403ddbac1565  #redis服务器的随机标识符(用于sentinel和集群)
- tcp_port:9021                                    #redis服务器监听端口
- uptime_in_seconds:26157730                       #redis服务器启动总时间，单位是秒
- uptime_in_days:302                               #redis服务器启动总时间，单位是天
- hz:10                                            #redis内部调度（进行关闭timeout的客户端，删除过期key等等）频率，程序规定serverCron每秒运行10次。
- lru_clock:14359959                               #自增的时钟，用于LRU管理,该时钟100ms(hz=10,因此每1000ms/10=100ms执行一次定时任务)更新一次。
- config_file:/redis_cluster/etc/9021.conf         #配置文件路径

**已连接客户端信息**

- connected_clients:1081                           #已连接客户端的数量(不包括通过slave连接的客户端)
- client_longest_output_list:0                     #当前连接的客户端当中，最长的输出列表，用client list命令观察omem字段最大值
- client_biggest_input_buf:0                       #当前连接的客户端当中，最大输入缓存，用client list命令观察qbuf和qbuf-free两个字段最大值
- blocked_clients:0                                #正在等待阻塞命令(BLPOP、BRPOP、BRPOPLPUSH)的客户端的数量

**内存信息**

- used_memory:327494024                 		   #由redis分配器分配的内存总量，以字节为单位
- used_memory_human:312.32M                        #以人类可读的格式返回redis分配的内存总量
- used_memory_rss:587247616                        #从操作系统的角度，返回redis已分配的内存总量(俗称常驻集大小)。这个值和top命令的输出一致
- used_memory_peak:1866541112                      #redis的内存消耗峰值(以字节为单位) 
- used_memory_peak_human:1.74G                     #以人类可读的格式返回redis的内存消耗峰值
- used_memory_lua:35840                            #lua引擎所使用的内存大小(以字节为单位)
- mem_fragmentation_ratio:1.79                     #used_memory_rss和used_memory之间的比率，小于1表示使用了swap，大于1表示碎片比较多
- mem_allocator:jemalloc-3.6.0                     #在编译时指定的redis所使用的内存分配器。可以是libc、jemalloc或者tcmalloc

**rdb和aof的持久化相关信息**

- loading:0                                        #服务器是否正在载入持久化文件
- rdb_changes_since_last_save:28900855             #离最近一次成功生成rdb文件，写入命令的个数，即有多少个写入命令没有持久化
- rdb_bgsave_in_progress:0                         #服务器是否正在创建rdb文件
- rdb_last_save_time:1482358115                    #离最近一次成功创建rdb文件的时间戳。当前时间戳rdb_last_save_time=多少秒未成功生成rdb文件
- rdb_last_bgsave_status:ok                        #最近一次rdb持久化是否成功
- rdb_last_bgsave_time_sec:2                       #最近一次成功生成rdb文件耗时秒数
- rdb_current_bgsave_time_sec:-1                   #如果服务器正在创建rdb文件，那么这个域记录的就是当前的创建操作已经耗费的秒数
- aof_enabled:1                                    #是否开启了aof
- aof_rewrite_in_progress:0                        #标识aof的rewrite操作是否在进行中
- aof_rewrite_scheduled:0                          #rewrite任务计划，当客户端发送bgrewriteaof指令，如果当前rewrite子进程正在执行，那么将客户端请求的bgrewriteaof变为计划任务，待aof子进程结束后执行rewrite 
- aof_last_rewrite_time_sec:-1                     #最近一次aof rewrite耗费的时长
- aof_current_rewrite_time_sec:-1                  #如果rewrite操作正在进行，则记录所使用的时间，单位秒
- aof_last_bgrewrite_status:ok                     #上次bgrewriteaof操作的状态
- aof_last_write_status:ok                         #上次aof写入状态
- aof_current_size:4201740                         #aof当前尺寸
- aof_base_size:4201687                            #服务器启动时或者aof重写最近一次执行之后aof文件的大小
- aof_pending_rewrite:0                            #是否有aof重写操作在等待rdb文件创建完毕之后执行
- aof_buffer_length:0                              #aof buffer的大小
- aof_rewrite_buffer_length:0                      #aof rewrite buffer的大小
- aof_pending_bio_fsync:0                          #后台I/O队列里面，等待执行的fsync调用数量
- aof_delayed_fsync:0                              #被延迟的fsync调用数量

**一般统计信息**

- total_connections_received:209561105             #新创建连接个数,如果新创建连接过多，过度地创建和销毁连接对性能有影响，说明短连接严重或连接池使用有问题，需调研代码的连接设置
- total_commands_processed:2220123478              #redis处理的命令数
- instantaneous_ops_per_sec:279                    #redis当前的qps，redis内部较实时的每秒执行的命令数
- total_net_input_bytes:118515678789               #redis网络入口流量字节数
- total_net_output_bytes:236361651271              #redis网络出口流量字节数
- instantaneous_input_kbps:13.56                   #redis网络入口kps
- instantaneous_output_kbps:31.33                  #redis网络出口kps
- rejected_connections:0                           #拒绝的连接个数，redis连接个数达到maxclients限制，拒绝新连接的个数
- sync_full:1                                      #主从完全同步成功次数
- sync_partial_ok:0                                #主从部分同步成功次数
- sync_partial_err:0                               #主从部分同步失败次数
- expired_keys:15598177                            #运行以来过期的key的数量
- evicted_keys:0                                   #运行以来剔除(超过了maxmemory后)的key的数量
- keyspace_hits:1122202228                         #命中次数
- keyspace_misses:577781396                        #没命中次数
- pubsub_channels:0                                #当前使用中的频道数量
- pubsub_patterns:0                                #当前使用的模式的数量
- latest_fork_usec:15679                           #最近一次fork操作阻塞redis进程的耗时数，单位微秒
- migrate_cached_sockets:0                         

**主从信息，master上显示的信息**

- role:master                                      #实例的角色，是master or slave
- connected_slaves:1                               #连接的slave实例个数
- slave0:ip=192.168.64.104,port=9021,state=online,offset=6713173004,lag=0         #lag从库多少秒未向主库发送REPLCONF命令
- master_repl_offset:6713173145                    #主从同步偏移量,此值如果和上面的offset相同说明主从一致没延迟
- repl_backlog_active:1                            #复制积压缓冲区是否开启
- repl_backlog_size:134217728                      #复制积压缓冲大小
- repl_backlog_first_byte_offset:6578955418        #复制缓冲区里偏移量的大小
- repl_backlog_histlen:134217728                   #此值等于master_repl_offset repl_backlog_first_byte_offset,该值不会超过repl_backlog_size的大小

**主从信息，slave上显示的信息**

- role:slave                                       #实例的角色，是master or slave
- master_host:192.168.64.102                       #此节点对应的master的ip
- master_port:9021                                 #此节点对应的master的port
- master_link_status:up                            #slave端可查看它与master之间同步状态,当复制断开后表示down
- master_last_io_seconds_ago:0                     #主库多少秒未发送数据到从库
- master_sync_in_progress:0                        #从服务器是否在与主服务器进行同步
- slave_repl_offset:6713173818                     #slave复制偏移量
- slave_priority:100                               #slave优先级
- slave_read_only:1                                #从库是否设置只读
- connected_slaves:0                               #连接的slave实例个数
- master_repl_offset:0         
- repl_backlog_active:0                            #复制积压缓冲区是否开启
- repl_backlog_size:134217728                      #复制积压缓冲大小
- repl_backlog_first_byte_offset:0                 #复制缓冲区里偏移量的大小
- repl_backlog_histlen:0                           #此值等于master_repl_offset repl_backlog_first_byte_offset,该值不会超过repl_backlog_size的大小

**CPU计算量统计信息**

- used_cpu_sys:96894.66                            #将所有redis主进程在核心态所占用的CPU时求和累计起来
- used_cpu_user:87397.39                           #将所有redis主进程在用户态所占用的CPU时求和累计起来
- used_cpu_sys_children:6.37                       #将后台进程在核心态所占用的CPU时求和累计起来
- used_cpu_user_children:52.83                     #将后台进程在用户态所占用的CPU时求和累计起来

**各种不同类型的命令的执行统计信息**

- cmdstat_get:calls=1664657469                     #call每个命令执行次数
- usec=8266063320                                  #usec总共消耗的CPU时长(单位微秒)
- usec_per_call=4.97                               #平均每次消耗的CPU时长(单位微秒)


**集群相关信息**

- cluster_enabled:1                                #实例是否启用集群模式

**数据库相关的统计信息**

- db0:keys=194690,expires=191702,avg_ttl=3607772262 #db0的key的数量,以及带有生存期的key的数,平均存活时间