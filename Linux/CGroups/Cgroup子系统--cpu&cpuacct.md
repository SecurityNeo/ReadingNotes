# CGroup子系统--cpu&cpuacct #

摘自[Cgroup分析之cpu、cpuacct](https://blog.csdn.net/tanzhe2017/article/details/81001105)

## cpu子系统 ##

Cpu子系统可以通过一些用户态接口文件来实现对cpu资源访问的控制，每个文件都独立存在在cgroup虚拟文件系统的伪文件中，如下：


- cpu.cfs_period_us

	以微秒为单位指定一个period的长度。

- cpu.cfs_quota_us

	以微秒为单位指定一个period中可用运行时间。

- cpu.stat

	当前nr_periods、nr_throttled、throttled_time的状态值，其中：
	
	nr_periods：已执行完的period数。
	
	nr_throttled：当前cgroup中所有进程超过限制值的次数。
	
	throttled_time：当前cgroup中所有进程超过限制值的持续时间。

- cpu.shares

	包含用来指定在cgroup中的任务可用的相对共享cpu时间的整数值。
	例如：在两个cgroup中都将cpu.shares设定为1的任务将有相同的cpu时间，但在cgroup中将cpu.shares设定为2的任务可使用的cpu时间是在cgroup中将cpu.shares设定为1的任务可使用的cpu时间的两倍。

- cpu.rt_period_us

	以微秒（μs，这里以"us"代表）为单位指定在某个时间段中cgroup对cpu资源访问重新分配的频率。如果某个cgroup中的任务应该每5秒钟有4秒时间可访问cpu资源，则将cpu.rt_runtime_us设定为4000000，并将cpu.rt_period_us设定为5000000。

- cpu.rt_runtime_us

	以微秒（μs，这里以"us"代表）为单位指定在某个时间段中cgroup中的任务对cpu资源的最长连续访问时间。建立这个限制是为了防止一个cgroup中的任务独占cpu时间。


## cpuacct子系统 ##

cpuacct主要是根据内核现有的一些接口对cpu使用状况做统计，并自动生成cgroup中进程所使用的CPU资源报告。相关用户态接口如下：

- cpuacct.stat

	报告当前cgroup和子组的所有任务使用用户模式和系统模式消耗的CPU周期数（单位由系统中user_hz定义）。

- cpuacct.usage

	统计这个cgroup中所有任务消耗的总cpu时间（纳秒）。

- cpuacct.usage_percpu

	统计这个cgroup中所有任务在每一个cpu上分别消耗的cpu时间（纳秒）。
