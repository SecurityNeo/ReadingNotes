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


### Configuration File ###


- Specification version

	必选，指定了bundle使用的OCI的版本

- root

	- path：容器的bundle路径，可以是相对路径和绝对路径，该值通常为rootfs
	- readonly：当设置为true时，容器的根文件为只读，默认false

- mount

	- destination：容器中的挂载点，必须是绝对路径
	- source：挂载的设备名称，文件或目录名称(bind mount时)，当option中有bind或rbind时改mount类型为bind mount
	- option：mount的选项，[参考mount选项](http://man7.org/linux/man-pages/man8/mount.8.html)
	
	**示例(Linux)：**
	
	```
	"mounts": [
		{
			"destination": "/tmp",
			"type": "tmpfs",
			"source": "tmpfs",
			"options": ["nosuid","strictatime","mode=755","size=65536k"]
		},
		{
			"destination": "/data",
			"type": "bind",
			"source": "/volumes/testing",
			"options": ["rbind","rw"]
		}
	]

	```


- process：定义了容器的进程信息
	- terminal：默认false，为true时，linux系统会为该进程分配一个pseudoterminal(pts)，并使用标准输入输出流
	- consoleSize：指定terminal的长宽规格
		- height
		- width
	- cwd：执行命令的绝对路径
	- env：环境变量
	- args：命令参数，至少需要指定一个参数，首参数即被execvp执行的文件
	
	其中根据不同的平台还会有不同的参数。

	**POSIX process：**

	- rlimits：设置进程的资源，如cpu，内存，文件大小等，参见getrlimit。docker里面使用--ulimit来设置单个进程的资源
		- type：linux或Solaris
		- soft：内核分配给该进程的资源
		- hard；可配置的资源的最大值，即soft的最大值。unprivileged进程(没有CAP_SYS_RESOURCE capability)可以将soft设置为0-hard之间的值
	
	**Linux process：**

	- apparmorProfile：指定进程的apparmor文件
	- capabilities：指定进程的capabilities
		- effective 
		- bounding 
		- inheritable
		- permitted
		- ambient
	- noNewPrivileges：设置为true后可以防止进程获取额外的权限(如使得suid和文件capabilities失效)，该标记位在内核4.10版本之后可以在/proc/$pid/status中查看NoNewPrivs的设置值。[参考no_new_privs](https://www.kernel.org/doc/Documentation/prctl/no_new_privs.txt)
	- oomScoreAdj ：给进程设置oom_score_adj值
	- selinuxLabel :设置进程的SELinux标签，即MAC值
	
	- user：用于控制运行进程的用户
		- uid：指定容器命名空间的UserID
		- gid：指定容器命名空间的GroupID
		- additionalGids：指定容器命名空间中附加的GroupID
	
	示例（Linux）：

	```
	"process": {
		"terminal": true,
		"consoleSize": {
			"height": 25,
			"width": 80
		},
		"user": {
			"uid": 1,
			"gid": 1,
			"additionalGids": [5, 6]
		},
		"env": [
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm"
		],
		"cwd": "/root",
		"args": [
			"sh"
		],
		"apparmorProfile": "acme_secure_profile",
		"selinuxLabel": "system_u:system_r:svirt_lxc_net_t:s0:c124,c675",
		"noNewPrivileges": true,
		"capabilities": {
			"bounding": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"permitted": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"inheritable": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"effective": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL"
			],
			"ambient": [
				"CAP_NET_BIND_SERVICE"
			]
		},
		"rlimits": [
			{
				"type": "RLIMIT_NOFILE",
				"hard": 1024,
				"soft": 1024
			}
		]
	}
	```


- hostname

	指定容器内的主机名

- Platform-specific configuration

	在linux，Windows，solaris平台上使用namespaces，cgroup等参数项

	- linux

		- Default Filesystems：如下路径需要正确挂载到容器中，以便容器进程的正确执行
		```
		Path        Type
		/proc       proc
		/sys        sysfs
		/dev/pts    devpts
		/dev/shm    tmpfs
		```
		- namespaces
			- type：指定namespace类型，如果没有指定namespace type，则继承父namespace的属性
				- ipc
				- mount
				- user
				- network
				- uts
				- pid
				- cgroup
			- path：namespace的文件，如果没有指定，则生成一个新的namespace
		- User namespace mappings:指定了uid和gid从宿主机到容器的映射关系
			- uidMappings
				- hostID
				- containerID
				- size
			- gidMappings
				- hostID
				- containerID
				- size
		
		- device：列出在容器中的设备
			- type：设备的类型
			- path:容器中的全路径
			- major, minor：设备的主设备号和次设备号，主设备号表示类型，次设备号表示分区，可以使用"ls -al /dev"查看主次设备号。
			- fileMode：文件ADC访问权限
			- uid：容器中设备的uid
			- gid：容器中设备的gid
		
		**示例：**
		```
		"devices": [
			{
				"path": "/dev/fuse",
				"type": "c",
				"major": 10,
				"minor": 229,
				"fileMode": 438,
				"uid": 0,
				"gid": 0
			},
			{
				"path": "/dev/sda",
				"type": "b",
				"major": 8,
				"minor": 0,
				"fileMode": 432,
				"uid": 0,
				"gid": 0
			}
		]
		```

		- Control groups
	
			控制容器的资源以及设备接入等
	
			- cgroupsPath:cgroup的路径，该路径可以是绝对路径，也可以是相对路径。如果没有设置该值，cgroup会使用默认的cgroup路径
			- devices：用于配置设备白名单
				- allow: 设置是否允许接入
				- type：设备类型: a (all), c (char), or b (block)， 默认为all
				- major, minor：设备的主次号。默认all
				- access：设备的cgroup权限,r(read), w(write), 和m(mknod)。
			
			**示例：**
	
			```
			"devices": [
				{
					"allow": false,
					"access": "rwm"
				},
				{
					"allow": true,
					"type": "c",
					"major": 10,
					"minor": 229,
					"access": "rw"
				}
			]
			```
			- memory： 内存限制，参考Linux/CGroups相关笔记
				- limit:设置内存使用limit
				- reservation：设置内存的soft limit
				- swap：设置memory+Swap使用limit
				- kernel：设置内存的hard limit
				- kernelTCP：设置内核TCP buffer的hard limit
				- swapness：设置swap的使用比例
				- disableOOMKiller：是否开启oomkiller
			
			**示例：**
			```
			"memory": {
				"limit": 536870912,
				"reservation": 536870912,
				"swap": 536870912,
				"kernel": -1,
				"kernelTCP": -1,
				"swappiness": 0,
				"disableOOMKiller": false
			}
			```
	
			- cpu: CPU限制，参考Linux/CGroups相关笔记
				- shares:cgroup中task使用的cpu的相对比例
				- quota:一个period中使用的cpu时间
				- period:以毫秒为单位的cpu周期 (CFS scheduler only)
				- realtimeRuntime:以毫秒为单位的cgroup tasks连续使用cpu资源的最长周期
				- realtimePeriod:实时调度的period
				- cpus:CPU列表
				- mems:memory nodes列表
	
			**示例：**
			```
			"cpu": {
				"shares": 1024,
				"quota": 1000000,
				"period": 500000,
				"realtimeRuntime": 950000,
				"realtimePeriod": 1000000,
				"cpus": "2-3",
				"mems": "0-7"
			}
			```
	
			- blockIO
				- weight
				- leafWeight
				- weightDevice
					- major, minor
					- weight
					- leafWeight
				- throttleReadBpsDevice
					- major, minor
					- rate
				- throttleWriteBpsDevice
				- throttleReadIOPSDevice
				- throttleWriteIOPSDevice
	
			**示例：**
			```
			"blockIO": {
				"weight": 10,
				"leafWeight": 10,
				"weightDevice": [
					{
						"major": 8,
						"minor": 0,
						"weight": 500,
						"leafWeight": 300
					},
					{
						"major": 8,
						"minor": 16,
						"weight": 500
					}
				],
				"throttleReadBpsDevice": [
					{
						"major": 8,
						"minor": 0,
						"rate": 600
					}
				],
				"throttleWriteIOPSDevice": [
					{
						"major": 8,
						"minor": 16,
						"rate": 300
					}
				]
			}
			```
	
			- hugepageLimits
			
				- pageSize： 大页大小
				- limit： 大页的使用上限，单位bytes
			
			- Network
	
				- classID：cgroup网络报文的标签
				- priorities
					- name：网卡名称
					- priority：网卡优先级
			
			**示例：**
	
			```
			"network": {
				"classID": 1048577,
				"priorities": [
					{
						"name": "eth0",
						"priority": 500
					},
					{
						"name": "eth1",
						"priority": 1000
			        }
		        ]
	        }
			```
	
			- pids
	
				- limit：cgroup限制的pid的数目
			
		- sysctl: 定义容器运行时的内核参数
			
		**示例：**

		```
		"sysctl": {
			"net.ipv4.ip_forward": "1",
			"net.core.somaxconn": "256"
		}
		```

		- seccomp
		
			在linux内核中为应用提供了一种沙盒机制，[参考seccomp](https://www.kernel.org/doc/Documentation/prctl/seccomp_filter.txt)
		
			- defaultAction：seccomp的默认动作，允许值类型为syscalls[].action

			- architectures：系统调用的平台，如下：

				```
				SCMP_ARCH_X86
				SCMP_ARCH_X86_64
				SCMP_ARCH_X32
				SCMP_ARCH_ARM
				SCMP_ARCH_AARCH64
				SCMP_ARCH_MIPS
				SCMP_ARCH_MIPS64
				SCMP_ARCH_MIPS64N32
				SCMP_ARCH_MIPSEL
				SCMP_ARCH_MIPSEL64
				SCMP_ARCH_MIPSEL64N32
				SCMP_ARCH_PPC
				SCMP_ARCH_PPC64
				SCMP_ARCH_PPC64LE
				SCMP_ARCH_S390
				SCMP_ARCH_S390X
				SCMP_ARCH_PARISC
				SCMP_ARCH_PARISC6
				```

		- syscalls：匹配seccomp的系统调用，可选

			- name：系统调用的名称，至少有一个
			
			- action：seccomp的动作规则。libseccomp v2.3.2中如下：

				```
				SCMP_ACT_KILL
				SCMP_ACT_TRAP
				SCMP_ACT_ERRNO
				SCMP_ACT_TRACE
				SCMP_ACT_ALLOW
				```

			- args
				- index: 系统调用的index
				- value: 系统调用参数的值
				- valueTwo: 系统调用参数的值
				- op: 系统调用参数的动作。 libseccomp v2.3.2如下:

					```
					SCMP_CMP_NE
					SCMP_CMP_LT
					SCMP_CMP_LE
					SCMP_CMP_EQ
					SCMP_CMP_GE
					SCMP_CMP_GT
					SCMP_CMP_MASKED_EQ
					```

		示例：
		
		```
		"seccomp": {
			"defaultAction": "SCMP_ACT_ALLOW",
			"architectures": [
				"SCMP_ARCH_X86",
				"SCMP_ARCH_X32"
			],
			"syscalls": [
				{
					"names": [
						"getcwd",
						"chmod"
					],
					"action": "SCMP_ACT_ERRNO"
				}
			]
		}
		```

		- rootfsPropagation： 设置rootfs的mount Propagation类型，slave，private或shared
			
	- windows
	- solaris


## 容器镜像标准（image spec） ##

摘自[https://segmentfault.com/a/1190000009309347](https://segmentfault.com/a/1190000009309347)

[https://github.com/opencontainers/image-spec](https://github.com/opencontainers/image-spec)


image由manifest、image index(可选)、filesystem layers和configuration四个部分组成，它们之间的关系如下：

![](img/image_spec_relationship.png)

注：
Image Index和Manifest的关系是"1..*"，表示它们是一对多的关系
Image Manifest和Config的关系是"1..1"，表示它们是一对一的关系
Image Manifest和Filesystem Layers是一对多的关系

### Filesystem Layers ###

Filesystem Layer包含了文件系统的信息，即该image包含了哪些文件/目录，以及它们的属性和数据。每个Filesystem Layer都包含了在上一个Layer上的改动情况，主要包含三方面的内容：

- 变化类型：是增加、修改还是删除了文件(Additions,Modifications,Removals)
 
- 文件类型：每个变化发生在哪种文件类型上(regular files,directories,sockets,symbolic links,block devices,character devices,FIFOs)

- 文件属性：文件的修改时间、用户ID、组ID、RWX权限等

最终每个layer都会打包成一个文件，这个文件的格式可以为tar，或者tar+gzip。不同的格式对应不同的Media Types对于Filesystem Layers，有四个Media Types，如下：

`application/vnd.oci.image.layer.v1.tar`

`application/vnd.oci.image.layer.v1.tar+gzip`

`application/vnd.oci.image.layer.nondistributable.v1.tar`

`application/vnd.oci.image.layer.nondistributable.v1.tar+gzip`

名称中有nondistributable的layer，标准要求这种类型的layer不能上传，只能下载。


### Image Config ###

Image Config是一个json文件，是对这个image的整体描述信息，它的Media Types是`application/vnd.oci.image.config.v1+json`。

下面介绍几个比较重要的字段：

- architecture

	CPU架构类型，amd64、arm64等

- os

	操作系统

- config

	当根据这个image启动container时，config里面的配置就是运行container时的默认参数

- rootfs

	指定了image所包含的Filesystem Layers，type的值必须是layers，diff_ids包含了layer的列表（顺序排列），每一个sha256就是每层layer对应tar包的sha256码


官方提供了一个较完整的例子：

```
{
    "created": "2015-10-31T22:22:56.015925234Z",
    "author": "Alyssa P. Hacker <alyspdev@example.com>",
    "architecture": "amd64",
    "os": "linux",
    "config": {
        "User": "alice",
        "ExposedPorts": {
            "8080/tcp": {}
        },
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "FOO=oci_is_a",
            "BAR=well_written_spec"
        ],
        "Entrypoint": [
            "/bin/my-app-binary"
        ],
        "Cmd": [
            "--foreground",
            "--config",
            "/etc/my-app.d/default.cfg"
        ],
        "Volumes": {
            "/var/job-result-data": {},
            "/var/log/my-app-logs": {}
        },
        "WorkingDir": "/home/alice",
        "Labels": {
            "com.example.project.git.url": "https://example.com/project.git",
            "com.example.project.git.commit": "45a939b2999782a3f005621a8d0f29aa387e1d6b"
        }
    },
    "rootfs": {
      "diff_ids": [
        "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
        "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      ],
      "type": "layers"
    },
    "history": [
      {
        "created": "2015-10-31T22:22:54.690851953Z",
        "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
      },
      {
        "created": "2015-10-31T22:22:55.613815829Z",
        "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
        "empty_layer": true
      }
    ]
}
```

### manifest ###

manifest也是一个json文件，Media Types为`application/vnd.oci.image.manifest.v1+json`，这个文件包含了对前面Filesystem Layers和Image Config的描述信息。

- config

	包含了对image config文件的描述，有media type，文件大小，以及sha256码

- layers

	包含了对每一个layer的描述，和对config文件的描述一样，也包含了media type，文件大小，以及sha256码

```
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 7023,
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 32654,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 16724,
      "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b"
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 73109,
      "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736"
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```


### Image Index ###

Image Index也是个json文件，Media Type是`application/vnd.oci.image.index.v1+json`。manifest描述的image只能支持一个平台，也没法支持多个tag，加上index文件的目的就是让这个image能支持多个平台和多tag。

```
{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7682,
      "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
      "platform": {
        "architecture": "amd64",
        "os": "linux",
        "os.features": [
          "sse4"
        ]
      }
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```

