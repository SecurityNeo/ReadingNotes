## Shell命令备忘录 ##

**eval**

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

**chattr**

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

**nsenter**

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


**ip**

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

**变量默认值**

- `${var:-string}和${var:=string}`
	若变量var为空，则用在命令行中用string来替换${var:-string}，变量var不为空时，则用变量var的值来替换${var:-string}

- `${var:=string}`
	${var:=string}的替换规则和${var:-string}是一样的，不同之处是${var:=string}若var为空时，用string替换${var:=string}的同时，把string赋给变量var,${var:=string}很常用的一种用法是，判断某个变量是否赋值，没有的话则给它赋上一个默认值。

- `${var:+string}`
	替换规则和上面的相反，即只有当var不是空的时候才替换成string，若var为空时则不替换或者说是替换成变量 var的值，即空值。(因为变量var此时为空，所以这两种说法是等价的)
 
- `${var:?string}`
	若变量var不为空，则用变量var的值来替换${var:?string}；若变量var为空，则把string输出到标准错误中，并从脚本中退出。我们可利用此特性来检查是否设置了变量的值

**模式匹配替换**

- `${variable%pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最短的匹配模式

- `${variable%%pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最长的匹配模式
 
- `${variable#pattern}`
	shell在variable中查找，看它是否以模式pattern开始，如果是，就从命令行把variable中的内容去掉左边最短的匹配模式
 
- `${variable##pattern}`
	shell在variable中查找，看它是否以模式pattern结尾，如果是，就从命令行把variable中的内容去掉右边最长的匹配模式

这四种模式中都不会改变variable的值，其中，只有在pattern中使用了`*`匹配符号时，%和%%，#和##才有区别。结构中的pattern支持通配符，`*`表示零个或多个任意字符，`?`表示仅与一个任意字符匹配，`[...]`表示匹配中括号里面的字符，`[!...]`表示不匹配中括号里面的字符。

**字符串提取和替换**

- ${var:num}
	shell在var中提取第num个字符到末尾的所有字符。若num为正数，从左边0处开始；若num为负数，从右边开始提取字串，但必须使用在冒号后面加空格或一个数字或整个num加上括号，如${var: -2}、${var:1-3}或${var:(-2)}。

- ${var:num1:num2}
	num1是位置，num2是长度。表示从$var字符串的第$num1个位置开始提取长度为$num2的子串。不能为负数。

- ${var/pattern/pattern}
	表示将var字符串的第一个匹配的pattern替换为另一个pattern。

- ${var//pattern/pattern}
	表示将var字符串中的所有能匹配的pattern替换为另一个pattern。


**查询某个POD的容器ID**

`kubectl describe pod POD_NAME | grep -A10 "^Containers:" | grep -Eo 'docker://.*$' | head -n 1 |sed 's/docker:\/\/\(.*\)$/\1/'`

**Docker相关维护命令**

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

**Harbor空间清理**

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

**Nexus重置密码**

1、进入OrientDB控制台

`java -jar ./lib/support/nexus-orient-console.jar`

2、连接security数据库

`connect plocal:../sonatype-work/nexus3/db/security admin admin`

3、重置admin的密码为admin123

`update user SET password="$shiro1$SHA-512$1024$NE+wqQq/TmjZMvfI7ENh/g==$V4yPw8T64UQ6GfJfxYq2hLsVrBY8D1v+bktfOxGdt4b/9BthpWPNUy/CBk6V9iA0nHpzYzJFWO8v/tZFtES8CA==" UPSERT WHERE id="admin"`