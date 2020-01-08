# 系统性能分析工具--vmstat #

[https://www.thomas-krenn.com/en/wiki/Linux_Performance_Measurements_using_vmstat](https://www.thomas-krenn.com/en/wiki/Linux_Performance_Measurements_using_vmstat)

**输出参数含义**

```
[root@cloud ~]$ vmstat 1 5
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 3  0      0  44712 110052 623096    0    0    30    28  217  888 13  3 83  1  0
 0  0      0  44408 110052 623096    0    0     0     0   88 1446 31  4 65  0  0
 0  0      0  44524 110052 623096    0    0     0     0   84  872 11  2 87  0  0
 0  0      0  44516 110052 623096    0    0     0     0  149 1429 18  5 77  0  0
 0  0      0  44524 110052 623096    0    0     0     0   60  431 14  1 85  0  0
```

- Procs
	- r: The number of processes waiting for run time.（How many processes are waiting for CPU time.）
	- b: The number of processes in uninterruptible sleep.（Wait Queue - Process which are waiting for I/O (disk, network, user input,etc..) ）
- Memory
	- swpd: the amount of virtual memory used.（shows how many blocks are swapped out to disk (paged). Total Virtual memory usage. ）
	- free: the amount of idle memory.
	- buff: the amount of memory used as buffers,like before/after I/O operations.
	- cache: the amount of memory used as cache by the Operating System.
	- inact: the amount of inactive memory. (-a option)
	- active: the amount of active memory. (-a option)
- Swap
	- si: Amount of memory swapped in from disk (/s).(How many blocks per second the operating system is swapping in. i.e Memory swapped in from the disk (Read from swap area to Memory))
	- so: Amount of memory swapped to disk (/s).
- IO
	- bi: Blocks received from a block device (blocks/s).- Read (like a hard disk) 
	- bo: Blocks sent to a block device (blocks/s).- Write
- System
	- in: The number of interrupts per second, including the clock.
	- cs: The number of context switches per second.
- CPU
	- These are percentages of total CPU time.
	- us: Time spent running non-kernel code. (user time, including nice time)
	- sy: Time spent running kernel code. (system time - network, IO 
     interrupts, etc)
	- id: Time spent idle. Prior to Linux 2.5.41, this includes IO-wait time.
	- wa: Time spent waiting for IO. Prior to Linux 2.5.41, included in idle.
	- st: Time stolen from a virtual machine. Prior to Linux 2.6.11, unknown.
