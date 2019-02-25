# Docker OCI #

OCI（Open Container Initiative）由 Docker、Coreos以及其他容器相关公司创建于2015年，由Linux基金会进行管理，致力于Container Runtime的标准制定和runc的开发等工作。目前主要有两个标准：容器运行时标准 （runtime spec）和 容器镜像标准（image spec）。

OCI规定了如何下载OCI镜像并解压到`OCI filesystem bundle`，这样OCI runtime就可以运行OCI bundle了。

![](img/image_runtime.jpg)

