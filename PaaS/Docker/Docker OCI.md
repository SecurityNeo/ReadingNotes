# Docker OCI #

OCI（Open Container Initiative）由 Docker、Coreos以及其他容器相关公司创建于2015年，由Linux基金会进行管理，致力于Container Runtime的标准制定和runc的开发等工作。目前主要有两个标准：容器运行时标准 （runtime spec）和 容器镜像标准（image spec）。

OCI规定了如何下载OCI镜像并解压到`OCI filesystem bundle`，这样OCI runtime就可以运行OCI bundle了。

![](img/image_runtime.jpg)

## 容器运行时标准 （runtime spec） ##

[https://github.com/opencontainers/runtime-spec](https://github.com/opencontainers/runtime-spec)

### Filesystem Bundle ###

Filesystem Bundle是个目录，用于给runtime提供启动容器必备的配置文件和文件系统。标准的Container Bundle包含以下两个内容：

- config.json

	此文件必须存在于Bundle的根目录中，且名字必须为config.json，其中包含了容器运行的配置信息

- 容器的Root Filesystem

	可以由config.json文件中的root.path指定

### Runtime and Lifecycle ###

**state：**

- ociVersion：创建容器时的OCI版本

- id：容器唯一的ID

- status：容器的runtime状态，可取值如下：

	- creating：容器正在被创建（lifecycle的第2步）
	- created：容器完成创建，没有返回错误且没有执行用户程序（lifecycle的第2步之后）
	- running：容器正在执行用户程序且没有返回错误（lifecycle的第5步之后）
	- stoped：容器进程退出（lifecycle的第7步）
	
	![](img/Container_Status.png)

- pid：容器在宿主机上的进程PID(容器状态为`created`或者`running`时)

- bundle：宿主机上容器bundle目录的绝对路径

- annotation：容器相关的标注（可选）

示例：

```json
{
    "ociVersion": "0.2.0",
    "id": "oci-container1",
    "status": "running",
    "pid": 4422,
    "bundle": "/containers/redis",
    "annotations": {
        "myKey": "myValue"
    }
}

```

**Lifecycle**

- 1、OCI runtime的create调用与bundle的路径和ID相关

- 2、OCI runtime依据config.json中的设置来创建环境，如果无法创建config.json中指定的环境，则返回错误。此阶段主要创建config.json中的资源，并没有执行用户程序。该步骤之后对config.json文件的修改都不会影响容器

- 3、runtime使用容器的唯一ID来执行start容器的命令

- 4、runtine执行prestart hooks，如果prestart hooks执行失败，则返回错误，并停止容器，执行第9条操作

- 5、runtime执行用户程序

- 6、runtime执行poststart hooks，如果poststart hooks执行失败，则记录warning日志，而poststart hooks和lifecycle继续运行

- 7、容器进程退出，可能是由错误退出，人为退出，程序崩溃或runtime 执行kill命令引起

- 8、runtime使用容器的唯一id来执行delete容器操作

- 9、如果在容器创建阶段（第2步）没有完成某些步骤，则容器必须被销毁

- 10、runtime执行poststop hooks，如果poststop hooks执行失败，则记录warning日志，而poststop hooks和lifecycle继续运行

**Operations**

- Query State

	`state <container-id>`，返回上述state内容

- Create

	`create <container-id> <path-to-bundle>`，该操作中会用到config.json除process之外的配置属性(因为process是在start阶段用到的)。实现中可能会与本规范不一致，如在create操作之前实现了pre-create。

- Start

	`start <container-id>`，执行config.json的process中定义的程序，如果process没有设定，则返回错误

- Kill

	`kill <container-id> <signal>`，向一个非running状态的容器发送的信号会被忽略。此操作用于向容器进程发送Kill信号

- Delete

	`delete <container-id>`，尝试删除一个非stopped的容器会返回错误。容器删除后其ID可能会被后续的容器使用

**Hooks**

定义每个Operations前后执行的操作，共有`prestart`、`poststart`和`poststop`三个操作，每个操作中都可包含`path`、`args`、`env`和`timeout`四个参数。[参考posix-platform-hooks部分](https://github.com/opencontainers/runtime-spec/blob/master/config.md#posix-platform-hooks)

示例：

```
"hooks": {
        "prestart": [
            {
                "path": "/usr/bin/fix-mounts",
                "args": ["fix-mounts", "arg1", "arg2"],
                "env":  [ "key1=value1"]
            },
            {
                "path": "/usr/bin/setup-network"
            }
        ],
        "poststart": [
            {
                "path": "/usr/bin/notify-start",
                "timeout": 5
            }
        ],
        "poststop": [
            {
                "path": "/usr/sbin/cleanup.sh",
                "args": ["cleanup.sh", "-f"]
            }
        ]
    }
```

