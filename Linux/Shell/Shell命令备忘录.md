# Shell命令备忘录 #

## eval ##

语法： `eval command-line`

eval命令将会首先扫描命令行进行所有的替换，类似于C语言中的宏替换，然后再执行命令。该命令使用于那些一次扫描无法实现其功能的变量。该命令对变量进行两次扫描。

示例：

1、`eval echo \$$#`取得函数的最后一个参数值

2、模拟“指针”

```
x=100
ptrx=x
eval echo \$$ptrx
```

## chattr ##

语法： `chattr [ -RVf ] [ -v version ] [ mode ] files…`

```
A：即Atime，告诉系统不要修改对这个文件的最后访问时间。
S：即Sync，一旦应用程序对这个文件执行了写操作，使系统立刻把修改的结果写到磁盘。
a：即Append Only，系统只允许在这个文件之后追加数据，不允许任何进程覆盖或截断这个文件。如果目录具有这个属性，系统将只允许在这个目录下建立和修改文件，而不允许删除任何文件。
b：不更新文件或目录的最后存取时间。
c：将文件或目录压缩后存放。
d：当dump程序执行时，该文件或目录不会被dump备份。
D:检查压缩文件中的错误。
i：即Immutable，系统不允许对这个文件进行任何的修改。如果目录具有这个属性，那么任何的进程只能修改目录之下的文件，不允许建立和删除文件。
s：彻底删除文件，不可恢复，因为是从磁盘上删除，然后用0填充文件所在区域。
u：当一个应用程序请求删除这个文件，系统会保留其数据块以便以后能够恢复删除这个文件，用来防止意外删除文件或目录。
t:文件系统支持尾部合并（tail-merging）。
X：可以直接访问压缩文件的内容。
```

## nsenter ##

nsenter命令可以在指定进程的命令空间下运行指定程序的命令。包含在util-linux包中。

```
nsenter [options] [program [arguments]]

options:
-t, --target pid：指定被进入命名空间的目标进程的pid
-m, --mount[=file]：进入mount命令空间。如果指定了file，则进入file的命令空间
-u, --uts[=file]：进入uts命令空间。如果指定了file，则进入file的命令空间
-i, --ipc[=file]：进入ipc命令空间。如果指定了file，则进入file的命令空间
-n, --net[=file]：进入net命令空间。如果指定了file，则进入file的命令空间
-p, --pid[=file]：进入pid命令空间。如果指定了file，则进入file的命令空间
-U, --user[=file]：进入user命令空间。如果指定了file，则进入file的命令空间
-G, --setgid gid：设置运行程序的gid
-S, --setuid uid：设置运行程序的uid
-r, --root[=directory]：设置根目录
-w, --wd[=directory]：设置工作目录

如果没有给出program，则默认执行$SHELL。
```

示例：在宿主机上抓容器内的包

查询容器PID：

`docker inspect -f {{.State.Pid}} CONTAINER—ID`

使用nsenter命令进入该容器的网络命令空间，然后就可以进行抓容器内的包了：

`nsenter -n -t CONTAINER—PID`


## ip ##

![https://www.cnblogs.com/diantong/p/9511072.html](https://www.cnblogs.com/diantong/p/9511072.html)

`ip [ OPTIONS ] OBJECT { COMMAND | help }`

 - OBJECT={ link | addr | addrlabel | route | rule | neigh | ntable | tunnel | maddr | mroute | mrule | monitor | xfrm | token }

	- link 网络设备
	- address 设备上的协议（IP或IPv6）地址
	- addrlabel 协议地址选择的标签配置
	- neighbour ARP或NDISC缓存条目
	- route 路由表条目
	- rule 路由策略数据库中的规则
	- maddress 组播地址
	- mroute 组播路由缓存条目
	- tunnel IP隧道
	- xfrm IPSec协议框架

 - OPTIONS={ -V[ersion] | -s[tatistics] | -d[etails] | -r[esolve] | -h[uman-readable] | -iec | -f[amily] { inet | inet6 | ipx | dnet | link } | -o[neline] | -t[imestamp] | -b[atch] [filename] | -rc[vbuf] [size] }

	- -V，-Version 显示指令版本信息
	- -s,-stats,statistics 输出详细信息
	- -h,-human,-human-readable 输出人类可读的统计信息和后缀
	- -iec 以IEC标准单位打印人类可读速率（例如1K=1024）
	- -f,-family <FAMILY> 指定要使用的协议族。协议族标识可以是inet、inet6、ipx、dnet或link之一。如果此选项不存在，则从其他参数中推测协议族。如果命令行的其余部分没有提供足够的信息来推测该族，则ip会退回到默认值，通常是inet或any。link是一个特殊的系列标识符，表示不涉及网络协议。
	- -4 –family inet的快捷方式
	- -6 –family inet6的快捷方式
	- -0 –family link的快捷方式
	- -o,-oneline 将每条记录输出到一行，用’\’字符替换换行符。
	- -r,-resolve 使用系统名称解析程序来打印DNS名称而不是主机地址。

示例：

1、启用/禁用网卡

`ip link set enp0s3 up`

2、为网卡分配 IP 地址以及其他网络信息

`ip addr add 192.168.0.50/255.255.255.0 dev enp0s3`

3、为网卡添加别名(为网卡添加多个IP)

`ip addr add 192.168.0.20/24 dev enp0s3 label enp0s3:1`

4、检查路由/默认网关的信息

`ip route show`

`ip route get 192.168.0.1`

5、检查所有的ARP记录

`ip neigh`

6、修改ARP记录

`ip neigh del 192.168.0.106 dev enp0s3`  / 删除对应ARP记录

`ip neigh add 192.168.0.150 lladdr 33:1g:75:37:r3:84 dev enp0s3 nud perm`  / 往ARP缓存中添加新记录  

	nud的意思是 “neghbour state”（网络邻居状态），它的值可以是：
	
	- perm 永久有效并且只能被管理员删除
	- noarp 记录有效，但在生命周期过期后就允许被删除了
	- stale 记录有效，但可能已经过期
	- reachable 记录有效，但超时后就失效了

7、查看网络统计信息

`ip -s link`

## 变量默认值 ##

- `${var:-string}和${var:=string}`
	若变量var为空，则用在命令行中用string来替换${var:-string}，变量var不为空时，则用变量var的值来替换${var:-string}

- `${var:=string}`
	${var:=string}的替换规则和${var:-string}是一样的，不同之处是${var:=string}若var为空时，用string替换${var:=string}的同时，把string赋给变量var,${var:=string}很常用的一种用法是，判断某个变量是否赋值，没有的话则给它赋上一个默认值。

- `${var:+string}`
	替换规则和上面的相反，即只有当var不是空的时候才替换成string，若var为空时则不替换或者说是替换成变量 var的值，即空值。(因为变量var此时为空，所以这两种说法是等价的)
 
- `${var:?string}`
	若变量var不为空，则用变量var的值来替换${var:?string}；若变量var为空，则把string输出到标准错误中，并从脚本中退出。我们可利用此特性来检查是否设置了变量的值

## 模式匹配替换 ##

- `${variable%pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最短的匹配模式

- `${variable%%pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最长的匹配模式
 
- `${variable#pattern}`
	shell在variable中查找，看它是否以模式pattern开始，如果是，就从命令行把variable中的内容去掉左边最短的匹配模式
 
- `${variable##pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最长的匹配模式

这四种模式中都不会改变variable的值，其中，只有在pattern中使用了`*`匹配符号时，%和%%，#和##才有区别。结构中的pattern支持通配符，`*`表示零个或多个任意字符，`?`表示仅与一个任意字符匹配，`[...]`表示匹配中括号里面的字符，`[!...]`表示不匹配中括号里面的字符。

## 字符串提取和替换 ##

- ${var:num}
	shell在var中提取第num个字符到末尾的所有字符。若num为正数，从左边0处开始；若num为负数，从右边开始提取字串，但必须使用在冒号后面加空格或一个数字或整个num加上括号，如${var: -2}、${var:1-3}或${var:(-2)}。

- ${var:num1:num2}
	num1是位置，num2是长度。表示从$var字符串的第$num1个位置开始提取长度为$num2的子串。不能为负数。

- ${var/pattern/pattern}
	表示将var字符串的第一个匹配的pattern替换为另一个pattern。

- ${var//pattern/pattern}
	表示将var字符串中的所有能匹配的pattern替换为另一个pattern。


## 查询某个POD的容器ID ##

`kubectl describe pod POD_NAME | grep -A10 "^Containers:" | grep -Eo 'docker://.*$' | head -n 1 |sed 's/docker:\/\/\(.*\)$/\1/'`

## Docker相关维护命令 ##

- 杀死所有正在运行的容器
- 
`docker kill $(docker ps -a -q)`

- 删除所有已经停止的容器
 
`docker rm $(docker ps -a -q)`

- 删除所有未打 dangling 标签的镜像
 
`docker rmi $(docker images -q -f dangling=true)`

- 删除所有镜像
 
`docker rmi $(docker images -q)`

- 强制删除镜像名称中包含“doss-api”的镜像

`docker rmi --force $(docker images | grep doss-api | awk '{print $3}')`

- 删除所有未使用数据

`docker system prune`

- 只删除未使用的volumes

`docker volume prune`

- 删除所有已退出的容器

`docker rm -v $(docker ps -aq -f status=exited)`

- 删除所有状态为dead的容器

`docker rm -v $(docker ps -aq -f status=dead)`

## Harbor空间清理 ##

1、查看仓库镜像信息

`skopeo inspect docker://docker.io/fedora`

2、拷贝镜像

`skopeo copy docker://busybox:1-glibc atomic:myns/unsigned:streaming`

3、删除镜像

`skopeo delete docker://localhost:5000/imagename:latest`

4、进入harbor的registry容器执行垃圾回收命令

`registry garbage-collect /etc/docker/registry/config.yml`

5、如果执行垃圾回收遇到一些文件校验的错时尝试先将大小为0的文件删除再执行

`find . -name "*" -type f -size 0c`

## Nexus重置密码 ##

1、进入OrientDB控制台

`java -jar ./lib/support/nexus-orient-console.jar`

2、连接security数据库

`connect plocal:../sonatype-work/nexus3/db/security admin admin`

3、重置admin的密码为admin123

`update user SET password="$shiro1$SHA-512$1024$NE+wqQq/TmjZMvfI7ENh/g==$V4yPw8T64UQ6GfJfxYq2hLsVrBY8D1v+bktfOxGdt4b/9BthpWPNUy/CBk6V9iA0nHpzYzJFWO8v/tZFtES8CA==" UPSERT WHERE id="admin"`

## strace ##

strace会追踪程序运行时的整个生命周期，输出每一个系统调用的名字、参数、返回值和执行所消耗的时间等

常见参数

```
-p 跟踪指定的进程
-f 跟踪由fork子进程系统调用
-F 尝试跟踪vfork子进程系统调吸入，与-f同时出现时, vfork不被跟踪
-o filename 默认strace将结果输出到stdout。通过-o可以将输出写入到filename文件中
-ff 常与-o选项一起使用，不同进程(子进程)产生的系统调用输出到filename.PID文件
-r 打印每一个系统调用的相对时间
-t 在输出中的每一行前加上时间信息。 -tt 时间确定到微秒级。还可以使用-ttt打印相对时间
-v 输出所有系统调用。默认情况下，一些频繁调用的系统调用不会输出
-s 指定每一行输出字符串的长度,默认是32。文件名一直全部输出
-c 统计每种系统调用所执行的时间，调用次数，出错次数。
-e expr 输出过滤器，通过表达式，可以过滤出掉你不想要输出

        常见选项：
        -e trace=[set]    只跟踪指定的系统调用
        -e trace=file     只跟踪与文件操作有关的系统调用
        -e trace=process  只跟踪与进程控制有关的系统调用
        -e trace-network  只跟踪与网络有关的系统调用
        -e trace=signal   只跟踪与系统信号有关的系统调用
        -e trace=desc     只跟踪与文件描述符有关的系统调用
        -e trace=ipc      只跟踪与进程通信有关的系统调用
        -e abbrev=[set]   设定strace输出的系统调用的结果集
        -e raw=[set]      将指定的系统调用的参数以十六进制显示
        -e signal=[set]   指定跟踪的系统信号
        -e read=[set]     输出从指定文件中读出的数据
        -e write=[set]    输出写入到指定文件中的数据
```

## ltrace ##

ltrace 能够跟踪进程的库函数调用，它会显现出调用了哪个库函数，而 strace则是跟踪进程的每个系统调用。 

常见参数：

```
-c    统计库函数每次调用的时间，最后程序退出时打印摘要
-C    解码低级别名称（内核级）为用户级名称
-d    打印调试信息
-e expr    输出过滤器，通过表达式，可以过滤掉你不想要的输出
          -e printf  表示只查看printf函数调用
          -e !printf 表示查看除printf函数以外的所有函数调用
-f         跟踪子进程
-o filename   将ltrace的输出写入文件filename
-p pid     指定要跟踪的进程pid
-r         输出每一个调用的相对时间
-S         显示系统调用
-t         在输出中的每一行前加上时间信息。例如16：45：28
-tt        在输出中的每一行前加上时间信息，精确到微秒。例如11：18：59.759546
-ttt       在输出中的每一行前加上时间信息，精确到微秒，而且时间表示为UNIX时间截。例如1486111461.650434
-T          显示每次调用所花费的时间
-u username 以username的UID和GID执行所跟踪的命令
```

## IPSET ##
### 创建(create) ###

命令：

`ipset create SETNAME TYPENAME`

注解：

	SETNAME： ipset的名称
	TYPENAME： 类型，格式为： method:datatype[,datatype[,datatype]]
		method: 指定ipset中的entry存放的方式，随后的datatype约定了每个entry的格式。bitmap, hash, list。bitmap和list使用固定大小的存储，hash使用hash表来存储元素。但为了避免Hash表键冲突,在ipset会在hash表key用完后，若又有新增条目，则ipset将自动对hash表扩大,假如当前哈希表大小为100条,则它将扩展为200条。当在iptables/ip6tables中使用了ipset hash类型的集合，则该集合将不能再新增条目。
		datatype: 可以为ip, net, mac, port, iface。[官网](http://ipset.netfilter.org/ipset.man.html)

示例：

```shell
[root@VM-0-4-centos ~]# ipset create blacklist hash:ip
[root@VM-0-4-centos ~]# ipset create webserver hash:ip,port
[root@VM-0-4-centos ~]# ipset create database hash:net
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16784
References: 0
Members:
```

### 添加条目(add) ###

命令：

`ipset add SETNAME ENTRY`

注解：
	
	ENTRY： 形式为ip/port/ip-ip等，注意：创建的集合属于哪种类型，在添加时的数据就要符合对应的类型

示例：

```
[root@VM-0-4-centos ~]# ipset add blacklist 192.168.1.2
[root@VM-0-4-centos ~]# ipset add blacklist 192.168.1.3,10.10.10.10
ipset v6.29: Syntax error: Elem separator in 192.168.1.3,10.10.10.10, but settype hash:ip supports none.
[root@VM-0-4-centos ~]# ipset add webserver 10.10.10.10,80
[root@VM-0-4-centos ~]# ipset add database 172.25.0.0/16
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16544
References: 0
Members:
192.168.1.2

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16560
References: 0
Members:
10.10.10.10,tcp:80

Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16816
References: 0
Members:
172.25.0.0/16

```

### 删除条例(del/flush/destroy) ###

命令&示例：

`ipset del SETNAME ENTRY`: 删除某个IP条目

```shell
[root@VM-0-4-centos ~]# ipset list blacklist
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16576
References: 0
Members:
192.168.1.2
192.168.1.3
192.168.1.4
[root@VM-0-4-centos ~]# ipset del blacklist 192.168.1.4
[root@VM-0-4-centos ~]# ipset list blacklist
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16576
References: 0
Members:
192.168.1.2
192.168.1.3
```

`ipset flush SETNAME`：删除某个集合的所有IP条目

```shell
[root@VM-0-4-centos ~]# ipset list database
Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16848
References: 0
Members:
10.125.0.0/16
172.25.0.0/16
[root@VM-0-4-centos ~]# ipset flush database
[root@VM-0-4-centos ~]# ipset list database
Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16784
References: 0
Members:
```

`ipset flush`：清空ipset中所有集合的ip条目（删条目，不删集合）

```shell
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16576
References: 0
Members:
192.168.1.2
192.168.1.3

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16560
References: 0
Members:
10.10.10.10,tcp:80

Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16784
References: 0
Members:
[root@VM-0-4-centos ~]# ipset flush
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16784
References: 0
Members:
```

`ipset destroy SETNAME`： 删除ipset中的某个集合或者所有集合

```shell
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: database
Type: hash:net
Revision: 3
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16784
References: 0
Members:
[root@VM-0-4-centos ~]# ipset destroy database
[root@VM-0-4-centos ~]# ipset list
Name: blacklist
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:

Name: webserver
Type: hash:ip,port
Revision: 2
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:
```

### IPSET选项 ###

#### timeout选项 ####

timeout设置超时时间，如果设置为0，表示永久生效，条目的超时时间也可以通过`-exist`来进行修改。可以为整个集合设置超时时间，即为加入此集合的条目默认超时时间，也可为具体条目设置超时时间。超时到期后自动清除条目。

示例：

```shell
[root@VM-0-4-centos ~]# ipset list webserver
Name: webserver
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536 timeout 50
Size in memory: 16592
References: 0
Members:
[root@VM-0-4-centos ~]# ipset add webserver 10.125.31.80 timeout 100
[root@VM-0-4-centos ~]# ipset add webserver 10.125.31.81
[root@VM-0-4-centos ~]# ipset list webserver
Name: webserver
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536 timeout 50
Size in memory: 16720
References: 0
Members:
10.125.31.81 timeout 30
10.125.31.80 timeout 75
[root@VM-0-4-centos ~]# ipset -exist add webserver 10.125.31.81 timeout 200
[root@VM-0-4-centos ~]# ipset list webserver
Name: webserver
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536 timeout 50
Size in memory: 16720
References: 0
Members:
10.125.31.81 timeout 197
10.125.31.80 timeout 20
```

#### counters, packets, bytes选项 ####

示例：

```shell
[root@VM-0-4-centos ~]# ipset create test hash:ip counters
[root@VM-0-4-centos ~]# ipset add test 100.1.1.2 packets 100 bytes 200
[root@VM-0-4-centos ~]# ipset add test 100.1.1.3
[root@VM-0-4-centos ~]# ipset list test
Name: test
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536 counters
Size in memory: 16720
References: 0
Members:
100.1.1.2 packets 100 bytes 200
100.1.1.3 packets 0 bytes 0
```

#### hashsize ####

定义集合的初始哈希大小，默认值为1024。哈希大小必须是2的幂，内核会自动舍入两个哈希大小的非幂到第一个正确的值。

示例：

```shell
[root@VM-0-4-centos ~]# ipset creat2 test hash:ip
[root@VM-0-4-centos ~]# ipset list test
Name: test
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:
[root@VM-0-4-centos ~]# ipset create test2 hash:ip hashsize 2048
[root@VM-0-4-centos ~]# ipset list test2
Name: test2
Type: hash:ip
Revision: 1
Header: family inet hashsize 2048 maxelem 65536
Size in memory: 32912
References: 0
Members:
```

#### maxelem ####

定义可以存储在集合中的元素的最大数量，默认值为65536.

```shell
[root@VM-0-4-centos ~]# ipset create test hash:ip maxelem 100000
[root@VM-0-4-centos ~]# ipset list test
Name: test
Type: hash:ip
Revision: 1
Header: family inet hashsize 1024 maxelem 100000
Size in memory: 16528
References: 0
Members:
```

#### family {inet|inet6} ####

定义要存储在集合中的IP地址的协议族，不指定时默认为IPV4。这个参数对于除hash:mac之外的所有hash类型集的create命令都是有效的。

```shell
[root@VM-0-4-centos ~]# ipset create test hash:ip family inet6
[root@VM-0-4-centos ~]# ipset list test
Name: test
Type: hash:ip
Revision: 1
Header: family inet6 hashsize 1024 maxelem 65536
Size in memory: 16528
References: 0
Members:
```

## veth-pair ##

veth-pair 是成对出现的一种虚拟网络设备，一端连接着协议栈，一端彼此相连，数据从一端出，从另一端进。

1. 创建一对名为 veth0 和 veth1 的 veth 接口：

	`ip link add veth0 type veth peer name veth1`

2. 创建 ns1 网络命名空间：

	`ip netns add ns1`

3. 将 veth0 接口加到 ns1 网络命名空间里：

	`ip link set veth0 netns ns1`
 
4. 为 veth0 接口配置 IP 地址：
 
	`ip -n ns1 addr add 10.1.1.1/24 dev veth0`
 
5. 将 veth0 接口和 lo 口 up 起来：

	`ip -n ns1 link set veth0 up`
	`ip -n ns1 link set lo up`

```shell
[root@VM-0-4-centos ~]# ip -n ns1 addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
5: veth0@if4: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state LOWERLAYERDOWN group default qlen 1000
    link/ether 52:a3:29:82:70:75 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.1.1.1/24 scope global veth0
       valid_lft forever preferred_lft forever
[root@VM-0-4-centos ~]# ip netns exec ns1 ping -c2 10.1.1.1
PING 10.1.1.1 (10.1.1.1) 56(84) bytes of data.
64 bytes from 10.1.1.1: icmp_seq=1 ttl=64 time=0.041 ms
64 bytes from 10.1.1.1: icmp_seq=2 ttl=64 time=0.059 ms

```

现在配置另一个接口veth1

1. 创建 ns2 网络命名空间：

	`ip netns add ns2`

2. 将 veth1 接口加到 ns2 网络命名空间里：

	`ip link set veth1 netns ns2`

3. 为 veth1 接口配置 IP 地址：

	`ip -n ns2 addr add 10.2.1.1/24 dev veth1`

4. 将 veth1 接口和 lo 口 up 起来：

	`ip -n ns2 link set veth1 up`
	`ip -n ns2 link set lo up`

```shell
[root@VM-0-4-centos ~]# ip -n ns2 addr show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
4: veth1@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default qlen 1000
    link/ether d6:c4:03:51:f9:98 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 10.2.1.1/24 scope global veth1
       valid_lft forever preferred_lft forever
    inet6 fe80::d4c4:3ff:fe51:f998/64 scope link
       valid_lft forever preferred_lft forever
[root@VM-0-4-centos ~]# ip netns exec ns2 ping -c2 10.2.1.1
PING 10.2.1.1 (10.2.1.1) 56(84) bytes of data.
64 bytes from 10.2.1.1: icmp_seq=1 ttl=64 time=0.040 ms
64 bytes from 10.2.1.1: icmp_seq=2 ttl=64 time=0.060 ms
```

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 ethtool -S veth0
NIC statistics:
     peer_ifindex: 4
[root@VM-0-4-centos ~]# ip netns exec ns2 ethtool -S veth1
NIC statistics:
     peer_ifindex: 5
```

此时两个命名空间的路由如下：

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
[root@VM-0-4-centos ~]# ip netns exec ns2 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1

```

此时两个接口之间没法通信，我们需要分别为其添加路由

```shell
[root@VM-0-4-centos ~]# ip -n ns1 route add 10.2.1.0/24 dev veth0
[root@VM-0-4-centos ~]# ip -n ns2 route add 10.1.1.0/24 dev veth1
[root@VM-0-4-centos ~]# ip netns exec ns1 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth0
[root@VM-0-4-centos ~]# ip netns exec ns2 route
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
10.1.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1
10.2.1.0        0.0.0.0         255.255.255.0   U     0      0        0 veth1
```

```shell
[root@VM-0-4-centos ~]# ip netns exec ns1 ping -c2 10.2.1.1
PING 10.2.1.1 (10.2.1.1) 56(84) bytes of data.
64 bytes from 10.2.1.1: icmp_seq=1 ttl=64 time=0.068 ms
64 bytes from 10.2.1.1: icmp_seq=2 ttl=64 time=0.074 ms

--- 10.2.1.1 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 999ms
rtt min/avg/max/mdev = 0.068/0.071/0.074/0.003 ms
[root@VM-0-4-centos ~]# ip netns exec ns2 ping -c2 10.1.1.1
PING 10.1.1.1 (10.1.1.1) 56(84) bytes of data.
64 bytes from 10.1.1.1: icmp_seq=1 ttl=64 time=0.052 ms
64 bytes from 10.1.1.1: icmp_seq=2 ttl=64 time=0.071 ms

--- 10.1.1.1 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 999ms
rtt min/avg/max/mdev = 0.052/0.061/0.071/0.012 ms
```
