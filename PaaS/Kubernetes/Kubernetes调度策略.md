# Kubernetes调度 #

## 调度策略 ##

kube-scheduler是Kubernetes中的调度组件，被称为三大核心组件之一，它的作用是根据特定的调度算法将pod调度到最优的工作节点上，这一过程叫做绑定（bind）。Scheduler启动之后会监听APIServer，获取PodSpec.NodeName为空的Pod，经过一系列调度算法计算后会对每个Pod创建一个绑定（binding）。

Kubernetes的调度器是以插件形式存在的，这样做的好处是可以方便用户定制和进行二次开发。用户可以自己编写调度器并以插件形式与Kubernetes进行集成，或着去集成其他调度器。

Kubernetes的调度策略分为Predicates（预选策略）和Priorites（优选策略），整个调度过程也分为两步进行：

- 预选策略 : 遍历所有Node，根据策略过滤掉所有不满足条件的Node，这一阶段的输出作为优选策略的输入。如没有Node符合Predicates策略规则，那么这个Pod就会被挂起，直到有Node能够满足

- 优选策略 : 按照预选策略将第一阶段输入的Node进行打分排序，得分最高的Node就是最合适的Node，Pod将Bind到此Node。

## 预选策略 ##

**MatchInterPodAffinity**

根据亲和性规则筛选Node。

**CheckVolumeBinding**

检查存储卷是否能绑定到Node上。

**CheckNodeCondition**

检查节点状态是否正常。
 
**GeneralPredicates**

包含一些基本的筛选规则，主要考虑Kubernetes资源是否充足，比如CPU和内存是否足够，端口是否冲突、selector是否匹配等。（PodFitsResources、PodFitsHostPorts、HostName、MatchNodeSelector）

- **HostName**

  检查Node是否满足PodSpec的NodeName字段中指定节点主机名，不满足的Node会被过滤掉。

- **PodFitsHostPorts**
 
  检查Pod定义的HostPort是否已被该Node上其它容器或服务占用。如果存在已被占用的情况，那么Pod将不会调度到这个Node上。（1.0版本被称之为PodFitsPorts，1.0之后版本变更为PodFitsHostPorts，为了向前兼容，PodFitsPorts名称仍然保留。）
 
- **MatchNodeSelector**
 
  检查Node节点的label定义是否满足Pod的NodeSelector属性需求。
 
- **PodFitsResources**
 
  检查Node上的空闲资源(CPU、Memory、GPU资源)是否满足Pod的需求，注意其是根据实际已经分配的资源量做调度，而不是使用已实际使用的资源量做调度。
 
**NoDiskConflict**

检查Node上是否存在卷冲突。如果存在卷冲突，其它使用这个卷的Pod不能调度到这个Node上。
ISCSI、GCE、AWS EBS和Ceph RBD的规则如下:

- ISCSI：在卷都是只读挂载的情况下，才允许挂载多个IQN相同的卷，否则都会产生卷冲突。

- GCE：与ISCSI一样允许多个Pod挂载同一个卷，但这些卷都必须是只读挂载的。
 
- AWS EBS：不允许不同的Pod挂载同一个卷。
 
- Ceph RBD：不允许任何两个Pod分享相同的monitor，match pool和image。

**PodToleratesNodeTaints**

根据taints和toleration的关系来判断Pod是否可以调度到该Node上。
 
**CheckNodeUnschedulable**

检查节点是否处于不可调度状态。

**PodToleratesNodeNoExecuteTaints**

检查Pod是否容忍节点上有NoExecute污点。如果一个Pod上运行在一个没有污点的Node上后，这个Node又给加上污点了，那么NoExecute表示这个新加污点的Node会祛除其上正在运行的Pod；不加NoExecute不会祛除节点上运行的Pod，表示接受。

**CheckNodeLabelPresence**

检查所有指定的Label是否存在于Node上（此处不考虑Label的值）。

**CheckServiceAffinity**
 
检查服务亲和性。多个Pod可以绑定到一个Service上，如果这些Pod都集中在集群中某一部分Node上，那新加入的Pod也会调度到这些Node上。

**MaxEBSVolumeCount**

 确保已挂载的EBS存储卷不超过设置的最大值。（默认值是39，Amazon推荐最大卷数量为40，其中一个卷为root卷，[参考官网说明](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/volume_limits.html#linux-specific-volume-limits)）。调度器会检查直接使用和间接使用这种类型存储的PVC，计算不同卷的总和，如果待调度的Pod部署上去后卷的数目会超过设置的最大值，那么该Pod就不能调度到这个Node上。我们可以通过环境变量`KUBE_MAX_PD_VOLS`设置最大卷的数量。
 
**MaxGCEPDVolumeCount**

确保已挂载的GCE存储卷不超过预设的最大值。（默认值是16,[参考官网说明](https://cloud.google.com/compute/docs/disks/add-persistent-disk#limits_for_predefined_machine_types))。与MaxEBSVolumeCount类似，最大卷的数量同样可通过环境变量`KUBE_MAX_PD_VOLS`设置。

**MaxAzureDiskVolumeCount**

确保已挂载的Azure存储卷不超过设置的最大值,默认为16。与MaxEBSVolumeCount类似，最大卷的数量同样可通过环境变量`KUBE_MAX_PD_VOLS`设置。
 
**MaxCinderVolumeCount**

确保已挂载的Cinder存储卷不超过设置的最大值。

**MaxCSIVolumeCountPred**

确保已挂载的CSI存储卷不超过设置的最大值。

**NoVolumeZoneConflict**

在给定的zone限制前提下，检查在此Node上部署Pod是否存在卷冲突。假定一些volumes可能有zone调度约束，VolumeZonePredicate根据volumes自身需求来评估pod是否满足条件。必要条件就是任何volumes的zone-labels必须与节点上的zone-labels完全匹配。节点上可以有多个zone-labels的约束（比如一个假设的复制卷可能会允许进行区域范围内的访问）。目前，这个只对PersistentVolumeClaims支持，而且只在PersistentVolume的范围内查找标签。

**CheckNodeMemoryPressure**

判断Node是否已经进入到内存压力状态，如果是则只允许调度内存为0的Pod到该Node上。

**CheckNodeDiskPressure**

判断Node是否已经进入到磁盘压力状态，如果是，则不调度新的Pod到该Node上。

**CheckNodePIDPressure**

检查Node上PID数量压力是否过大，但是一般PID时可以重复使用的。