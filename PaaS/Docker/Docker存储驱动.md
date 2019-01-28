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


