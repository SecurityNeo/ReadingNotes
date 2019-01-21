# Docker简介 #

## Docker概念 ##

**image**
The basis of a Docker container.The content at rest.

**Container**

The image when it is running.The standard unit for app service.

**Engine**

The software that executes commands for  contaners.Networking and volumes are part of Engine.Can be clustered together.

**Registry**

Stores,distributes and manages Docker images

**Control Plane**

Management plane for container and cluster orchestration.

## 相关历史 ##

- 2010年 Solomon Hykes成立dotcloud公司
- 2013年 dotcloud公司更为Docker，同年发布Docker-compose
- 2014年 Docker发布1.0版本
- 2015年 提供Docker Machine，支持windows，mac等平台
- 2015年 OCI开源社区成立，容器管理工具runc由社区维护（前身是Libcontainer）
- 2017年 Docker发布版分为Docker CE、EE、Moby
- 2018年 Docker支持Kubernetes

## Docker底层实现 ##

- OCI Open Container Initiative，也就是常说的OCI，是由多家公司共同成立的项目，并由linux基金会进行管理，致力于container runtime的标准的制定和runc的开发等工作。
- [runCrunc](https://github.com/opencontainers/runc)， 前身是libcontainer，是对于OCI标准的一个参考实现，是一个可以用于创建和运行容器的CLI(command-line interface)工具。runc直接与容器所依赖的cgroup/linux kernel等进行交互，负责为容器配置cgroup/namespace等启动容器所需的环境，创建启动容器的相关进程。为了兼容oci标准，docker也做了架构调整。将容器运行时相关的程序从docker daemon剥离出来，形成了containerd。Containerd向docker提供运行容器的API，二者通过grpc进行交互。containerd最后会通过runc来实际运行容器。