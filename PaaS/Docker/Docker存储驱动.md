# Docker存储驱动 #

## AUFS ##

Aufs最初代表的意思“另一个联合文件系统（another union filesystem）”，试图对当时已经存在的UnionFS实现进行重写。后来叫Alternative UnionFS，再后来就叫成了高大上的Advance UnionFS。

AUFS把若干目录按照顺序和权限挂载为一个目录，默认情况下，只有第一层是可写的，其余都是只读层，新写的文件都会被放在最上面的可写层中，当需要删除只读层中的文件时，AUFS通过在可写层目录下建立对应的whiteout隐藏文件来实现。此外AUFS利用其CoW（copy-on-write）特性来修改只读层中的文件。AUFS工作在文件层面，因此，只要有对只读层中的文件做修改，不管修改数据的量的多少，在第一次修改时，文件都会被拷贝到可写层然后再被修改。AUFS的CoW特性能够允许在多个容器之间共享分层，从而减少物理空间占用。

在Docker中，只读层就是image，可写层就是Container。其结构如下：

![](img/AUFS.png)

AUFS只在Ubuntu或者Debian的内核上才可以启用，因为Aufs从来没有被上游Linux内核社区接受，且原作者已经放弃了让它被内核采纳的努力。不过在Ubtuntu或者Debian上，默认的graphdriver就是aufs，它能满足绝大多数需求。

## Overlay ##

Overlay是一个联合文件系统，相比aufs来说，它的分支模型更为简单。Overlay只有两层：一个是下层目录（lower-dir）”，对应镜像层，另一个是“上层目录（upper-dir）”，对应容器层，同样的，镜像层是只读的，容器层可写。

![](img/Overlay.jpg)

采用Overlay存储驱动时，在路径/var/lib/docker/overlay/下（/var/lib/docker为Docker Root Dir），每个镜像层都有一个对应的目录，包含了本层镜像的内容，而每个镜像层只包含一个root目录，如下：

![](img/Overlay_image.png)

当启动容器后，会在已有的镜像层上创建一层容器层，容器层在路径/var/lib/docker/overlay下也存在对应的目录，在该目录下，文件lower-id记录的即为镜像层，upper包含了容器层的内容。创建容器时将lower-id指向的镜像层目录以及upper目录联合挂载到merged目录。work用来完成如copy-on_write的操作，如下：

![](img/Overlay_container.png)

Overlay从kernel3.18进入主流Linux内核。比AUFS和Device mapper速度快,因为OverlayFS只有两层，不是多层，所以OverlayFS “copy-up”操作快于AUFS,以此可以减少操作延时。另外OverlayFS支持页缓存共享，多个容器访问同一个文件能共享一个页缓存，以此提高内存使用率。不过Overlay有一个非常严重的问题，就是特别消耗inode，随着镜像和容器增加，inode会遇到瓶颈。Overlay2能解决这个问题。在Overlay下，为了解决inode问题，可以考虑将/var/lib/docker挂在单独的文件系统上，或者增加系统inode设置。


## Overlay2 ##

Overlay的硬链接实现方式已经引发了inode耗尽的问题，这阻碍了它的大规模采用，而overlay2可以解决inode耗尽和一些其他的问题。Overlay2也将继续保留overlay的一些特性。

Derek McGowan在[PR22126](https://github.com/moby/moby/pull/22126)中添加了overlay2的graphdriver，于2016年6月被合并进Docker 1.12版本。Linux在内核4.0上添加的[Multiple lower layers in overlayfs特性](https://kernelnewbies.org/Linux_4.0)，也即是说如果需要安装Docker使用Overlay2，需要先将Linux内核升级至4.0版本。

采用Overlay2存储驱动时，在路径/var/lib/docker/overlay2/下（/var/lib/docker为Docker Root Dir）即存储了镜像文件。在此目录下，有一个l目录，其中包含了很多软连接，使用短名称指向了其他层。采用短名称是用来避免mount参数时达到页面大小的限制。

![](img/Overlay2_img01.png)

在容器镜像目录内，有一个link文件，其中包含了上边提到的短名称，另外还有个diff目录，这其中包含了当前镜像的内容。

![](img/Overlay2_img02.png)

当启动容器之后，也是会在/var/lib/docker/overlay2目录下生成一层容器层，其中包括diff，merged和work目录，link和lower文件。diff目录中记录了每一层自己的数据，link文件中记录了该层链接目录，在lower文件中，使用:符号来分割不同的底层，并且顺序是从高层到底层。

![](img/Overlay_img03.png)

## Device Mapper ##

在Linux Kernel 2.6.9之后支持Device Mapper，Device Mapper提供一种从逻辑设备到物理设备的映射框架机制,为实现用于存储资源管理的块设备驱动提供了一个高度模块化的内核架构。

Device Mapper包含三个比较重要的对象概念，Mapped Device、Mapping Table和Target Device。Mapped Device是一个抽象出来的逻辑设备，通过Mapping Table映射关系与Target Device建立映射，Target Device即为Mapped Device所映射的物理空间段。

![](img/DeviceMapper_structure.png)

Device Mapper在内核中实现了诸多Target Driver插件，包括软Raid、加密、多路径、镜像、快照等，上图中linear、mirror、snapshot、multipath表示的就是这些Target Driver。

尹洋老师写了一篇[Linux 内核中的 Device Mapper 机制](https://www.ibm.com/developerworks/cn/linux/l-devmapper/index.html)，非常详细，复习时需仔细品味。


Docker的Device mapper利用Thin-provisioning snapshot管理镜像和容器。Snapshot是Lvm的一种特性，它可以在线为the origin（original device）创建一个虚拟快照(Snapshot)。Thin-Provisioning就是精简置备，逻辑上为其分配足够的空间，但实际上是真正占用多少空间就为其分配多少空间。Thin-provisioning Snapshot将Thin-Provisioning和Snapshoting两种技术结合起来，将多个虚拟设备同时挂载到一个数据卷从而实现数据共享。

上边提到的AUFS和OverlayFS是文件级存储，Device mapper是块级存储，所有的操作都是直接对块进行操作，而不是文件。

Docker Daemon第一次启动时会创建thin-pool，thin-pool的命名规则为"docker-dev_num-inode_num-pool"（dev是/var/lib/docker/devicemapper目录所在块设备的设备号，形式为主设备号:二级设备号；inode是这个目录的inode号），如下：

![](img/NeedToAddImg.png)

thin-pool基于块设备或者loop设备创建，这取决于使用loop-lvm还是direct-lvm，默认情况下是使用loop-lvm，但这仅仅适用于测试环境，若是生产环境强烈建议使用direct-lvm。块设备有两个，一个为data，存放数据，另一个为metadata，存放元数据(通过--storage-opt dm.datadev和--storage-optdm.metadatadev指定块设备)。

在 Docker 17.06或更高版本，Docker可以自动管理块设备，简化direct-lvm模式的配置。而这仅适用于Docker的首次设置，并且只能使用一个块设备，如果需要使用多个块设备，需手动配置direct-lvm模式。参考以下配置选项：

![](img/direct-lvm.png)

**Device Mapper的工作流程**
Docker采用devicemapper存储驱动时，所有和镜像及容器层相关的对象都保存在/var/lib/docker/devicemapper/里。base device是最底层的对象，就是上边说到的thin-pool（可使用`docker info` 命令查看），它包含一个文件系统，每个容器镜像层和容器层都从它开始。base device的元数据和所有的镜像及容器层都以JSON格式存储在/var/lib/docker/devicemapper/metadata/中，这些层是CoW的快照。每个容器的可写层都会挂载到/var/lib/docker/devicemapper/mnt/中的一个挂载点。对每个只读镜像层和每个停止状态的容器，都对应一个空目录。

对于一个镜像来说，每个镜像层都是其下一层的snapshot，而每个镜像的最底层是Pool中的一个已存在base device的snapshot。当一个容器运行起来时，容器就是它所依赖镜像的snapshot。下图是官方提供的两个容器的层级结构图：

![](img/Container_Layer.png)

**容器的读写操作（devicemapper）**

**读取文件**：

devicemapper的读取操作也发生在块级别，以官方给的例子为例：

1、容器内的一个APP要读`0x44f`块的内容，容器的dm设备是基于镜像的snapshot，容器里没有这个数据，但是它有一个指针指向这个块数据在镜像dm设备的位置，`0x44f`这个块的内容在`a005e`设备的`0xf33`块上。

2、从镜像`a005e`中读取`0xf33`块的内容到内存中，把数据返回给APP。

3、一个镜像有很多层，应用所需要的数据不一定在镜像的最上层就能找到，如果找不到就会依次往下层去寻找。

![](img/DeviceMapper_ReadingFile.png)

**写入文件**：


- 写入新文件

使用devicemapper驱动程序，通过allocate-on-demand操作实现将新数据写到容器中，新文件的每个块都分配到容器的可写层。

- 更新文件

将待更新文件的相关块从最近的镜像层中读取出来，当容器写入文件时，只有修改后的块被写入容器的可写层。

- 删除文件或目录

当删除容器可写层中的文件或目录，或者镜像层删除其父层中存在的文件时，devicemapper存储驱动程序会截获对该文件或目录的进一步读取尝试，并回应文件或目录不存在，不会真的删除，并且在读取相关文件或目录时告诉程序其已不存在。

- 写然后删除文件

如果在容器中写入文件并稍后删除该文件，所有这些操作都发生在容器的可写层中。在这种情况下，如果使用`direct-lvm`，块将被释放。如果使用`loop-lvm`，块可能不会被释放。这是不在生产环境中使用`loop-lvm`的另一个原因。

