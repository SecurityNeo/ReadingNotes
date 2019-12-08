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


