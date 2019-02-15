# Kubernetes服务质量（QoS）

摘自[http://dockone.io/article/2592](http://dockone.io/article/2592)

QoS的英文全称为"Quality of Service",中文名为"服务质量"。在Kubernetes中，每个Pod都会被标注一个Qos级别，Kubernetes针对不同服务质量的预期，通过QoS来对Pod进行服务质量管理。其中有两个指标，一个是CPU，一个是内存，体现在配置上就是`requests`和`limits`两种类型。


## `requests`和`limits` ##

`requests`申请范围从0到Node的最大配置，而`limits`申请范围是从`requests`到无限，即

`0 <= requests <=Node Allocatable`

`requests <= limits <= Infinity`

**超限额处理**

- 对于CPU，如果容器中服务使用的CPU超过设置的limits，容器不会被kill掉但会被限制。如果没有设置limits，容器可以使用全部空闲的CPU资源。

- 对于内存，当一个容器使用内存超过了设置的limits，Pod中container的进程会被kernel因OOM kill掉。当container因为OOM被kill掉时，系统倾向于在其原所在的机器上重启该container或在本机/其它机器上重新创建一个Pod。

## QoS分类 ##

**Guaranteed**

Pod中所有容器都必须统一设置limits，并且设置参数都一致，如果有一个容器要设置requests，那么所有容器都要设置，并设置参数同limits一致，那么这个Pod的QoS就是Guaranteed级别。

注：如果一个容器只指明limit而未设定request，则request的值等于limit值。

**Burstable**

Pod中只要有一个容器的requests和limits的设置不相同，该Pod的QoS即为Burstable。

**Best-Effort**

全部resources的requests与limits均未设置，该Pod的QoS即为Best-Effort。

### QoS优先级 ###

`Best-Effort  -> Burstable  -> Guaranteed `


## 资源回收策略 ##

Kubernetes根据资源能否伸缩进行分类，划分为可压缩资源和不可以压缩资源两种。CPU资源是目前支持的一种可压缩资源，而内存资源和磁盘资源为目前所支持的不可压缩资源。

Kubernetes根据QoS类型来进行资源回收（以内存为例）：

- Best-Effort：系统用完了全部内存时，该类型Pods会最先被kill掉。

- Burstable：系统用完了全部内存，且没有Best-Effort container可以被kill时，该类型Pods会被kill掉。

- Guaranteed：系统用完了全部内存、且没有Burstable与Best-Effort container可以被kill时，该类型的Pods会被kill掉。


### OOM打分 ###

有关Linux OOM参考[郭建：Linux内存管理系统参数配置之OOM（内存耗尽） ](https://www.sohu.com/a/238012686_467784)

各类容器的`OOM_SCORE_ADJ`参数定义如下：

```
pkg/kubelet/qos/policy.go:21
 
const (
    PodInfraOOMAdj        int = -998
    KubeletOOMScoreAdj    int = -999
    DockerOOMScoreAdj     int = -999
    KubeProxyOOMScoreAdj  int = -999
    guaranteedOOMScoreAdj int = -998
    besteffortOOMScoreAdj int = 1000
)

```

容器`OOM_SCORE_ADJ`计算规则如下：

```
pkg/kubelet/qos/policy.go:40

func GetContainerOOMScoreAdjust(pod *v1.Pod, container *v1.Container, memoryCapacity int64) int {
        switch GetPodQOS(pod) {
        case Guaranteed:
                // Guaranteed containers should be the last to get killed.
                return guaranteedOOMScoreAdj
        case BestEffort:
                return besteffortOOMScoreAdj
        }

        // Burstable containers are a middle tier, between Guaranteed and Best-Effort. Ideally,
        // we want to protect Burstable containers that consume less memory than requested.
        // The formula below is a heuristic. A container requesting for 10% of a system's
        // memory will have an OOM score adjust of 900. If a process in container Y
        // uses over 10% of memory, its OOM score will be 1000. The idea is that containers
        // which use more than their request will have an OOM score of 1000 and will be prime
        // targets for OOM kills.
        // Note that this is a heuristic, it won't work if a container has many small processes.
        memoryRequest := container.Resources.Requests.Memory().Value()
        oomScoreAdjust := 1000 - (1000*memoryRequest)/memoryCapacity
        // A guaranteed pod using 100% of memory can have an OOM score of 10. Ensure
        // that burstable pods have a higher OOM score adjustment.
        if int(oomScoreAdjust) < (1000 + guaranteedOOMScoreAdj) {
                return (1000 + guaranteedOOMScoreAdj)
        }
        // Give burstable pods a higher chance of survival over besteffort pods.
        if int(oomScoreAdjust) == besteffortOOMScoreAdj {
                return int(oomScoreAdjust - 1)
        }
        return int(oomScoreAdjust)
}
```

由上可知几种类型的`OOM_SCORE_ADJ`值：


1. **Best-effort**

	OOM_SCORE_ADJ: 1000

	best-effort容器的OOM_SCORE 值为1000

2. **Guaranteed**

	OOM_SCORE_ADJ: -998

	guaranteed容器的OOM_SCORE 值为0 或 1

3. **Burstable**

	- 如果总的memory request大于99.9%的可用内存，OOM_SCORE_ADJ设置为2。否则，OOM_SCORE_ADJ = 1000-10 * (% of memory requested)，这确保了burstable的POD OOM_SCORE > 1
	- 如果memory request设置为0，OOM_SCORE_ADJ默认设置为999。所以如果burstable pods和guaranteed pods冲突时，前者会被kill。
	- 如果burstable pod使用的内存少于request值，那它的OOM_SCORE < 1000。如果best-effort pod和这些 burstable pod冲突时，best-effort pod会先被kill掉。
	- 如果burstable pod容器中进程使用比request值的内存更多，OOM_SCORE设置为1000。反之，OOM_SCORES少于1000。
	- 在一堆burstable pod中，使用内存超过request值的pod，优先于内存使用少于request值的pod被kill。
	- 如果burstable pod有多个进程冲突，则OOM_SCORE会被随机设置，不受“request & limit”限制。

4. **Pod infra containers or Special Pod init process**

	OOM_SCORE_ADJ: -998

5. **Kubelet, Docker**

	OOM_SCORE_ADJ: -999 (won’t be OOM killed)
	系统上的关键进程，如果和guranteed 进程冲突，则会优先被kill 。将来会被放到一个单独的cgroup中，并且限制内存。


## 节点可用资源 ##

理想情况下，我们总是希望节点所有的资源都能提供给Pod使用，但实际上，节点上还会运行很多其它进程，包括一些系统进程、Kubelet、Docker进程等，而这些进程是保障Kubernetes集群能正常运行，甚至是这个集群能正常运行的，所以，我们需要为这些进程预留一部分资源，保障其正常运行，剩余的资源再分配给集群Pod使用。

节点资源分配大致如下：
摘自[从一次集群雪崩看Kubelet资源预留的正确姿势](https://my.oschina.net/jxcdwangtao/blog/1629059)

![](img/Node_Resources.png)

- Node Capacity: Node的所有硬件资源
- Kube-Reserved: 预留给kube组件的资源
- System-Reserved： 预留给System进程的资源
- Hard-Eviction-Threshold： kubelet eviction设定的阈值
- Allocatable： scheduler调度Pod时的参考值

计算方式： ` [Allocatable] = [Node Capacity] - [Kube-Reserved] - [System-Reserved] - [Hard-Eviction-Threshold]`

[参考官方说明](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node/node-allocatable.md#recommended-cgroups-setup)

Kubernetes相关配置：

- --enforce-node-allocatable

	默认为pods，要为kube组件和System进程预留资源，则需要设置为`pods,kube-reserved,system-reserve`。

- --cgroups-per-qos

	默认开启。开启后，kubelet会将管理所有workload Pods的cgroups。

- --cgroup-driver

	默认为cgroupfs，另一可选项为systemd。取决于容器运行时使用的cgroup driver，kubelet与其保持一致。比如配置docker使用systemd cgroup driver，那么kubelet也需要配置`--cgroup-driver=systemd`。

- --kube-reserved
 
	用于配置为kube组件（kubelet,kube-proxy,dockerd等）预留的资源量，比如`—kube-reserved=cpu=1000m,memory=8Gi，ephemeral-storage=16Gi`。

- --kube-reserved-cgroup
 
	如果你设置了--kube-reserved，那么一定要设置对应的cgroup，并且该cgroup目录要事先创建好，否则kubelet将不会自动创建导致kubelet启动失败。比如设置为`kube-reserved-cgroup=/kubelet.service`。

- --system-reserved

	用于配置为System进程预留的资源量，比如`—system-reserved=cpu=500m,memory=4Gi,ephemeral-storage=4Gi`。

- --system-reserved-cgroup

	如果你设置了--system-reserved，那么一定要设置对应的cgroup，并且该cgroup目录要事先创建好，否则kubelet将不会自动创建导致kubelet启动失败。比如设置为`system-reserved-cgroup=/system.slice`。

- --eviction-hard
 
	用来配置kubelet的hard eviction条件，只支持memory和ephemeral-storage两种不可压缩资源。当出现MemoryPressure时，Scheduler不会调度新的Best-Effort QoS Pods到此节点。当出现DiskPressure时，Scheduler不会调度任何新Pods到此节点。





