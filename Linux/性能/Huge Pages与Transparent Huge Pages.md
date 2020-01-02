# Linux传统Huge Pages与Transparent Huge Pages #

[https://www.cnblogs.com/kerrycode/p/7760026.html](https://www.cnblogs.com/kerrycode/p/7760026.html)

[http://blog.itpub.net/26736162/viewspace-2214374/](http://blog.itpub.net/26736162/viewspace-2214374/)

[https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/performance_tuning_guide/s-memory-transhuge](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/performance_tuning_guide/s-memory-transhuge)

大多数操作系统采用了分段或分页的方式进行管理。分段是粗粒度的管理方式，而分页则是细粒度管理方式，分页方式可以避免内存空间的浪费。相应地，也就存在内存的物理地址与虚拟地址的概念。通过前面这两种方式，CPU必须把虚拟地址转换程物理内存地址才能真正访问内存。为了提高这个转换效率，CPU会缓存最近的虚拟内存地址和物理内存地址的映射关系，并保存在一个由CPU维护的映射表中。为了尽量提高内存的访问速度，需要在映射表中保存尽量多的映射关系。Linux的内存管理采取的是分页存取机制，为了保证物理内存能得到充分的利用，内核会按照LRU算法在适当的时候将物理内存中不经常使用的内存页自动交换到虚拟内存中，而将经常使用的信息保留到物理内存。通常情况下，Linux默认情况下每页是4K，这就意味着如果物理内存很大，则映射表的条目将会非常多，会影响CPU的检索效率。因为内存大小是固定的，为了减少映射表的条目，可采取的办法只有增加页的尺寸。因此Hugepage便因此而来。也就是打破传统的小页面的内存管理方式，使用大页面2M,4M等。如此一来映射条目则明显减少。TLB 缓存命中率将大大提高。
在Linux中大页分为两种：Huge pages (标准大页)和Transparent Huge pages(透明大页) 。内存是以块即页的方式进行管理的，当前大部分系统默认的页大小为4096 bytes即4K 。1MB内存等于256页；1GB内存等于256000页，以此类推。

## Huge pages ##

Huge pages是从Linux Kernel 2.6后被引入的，目的是通过使用大页内存来取代传统的4kb内存页面， 以适应越来越大的系统内存，让操作系统可以支持现代硬件架构的大页面容量功能。Huge pages有两种格式大小：2MB和1GB ，2MB页块大小适合用于GB大小的内存，1GB页块大小适合用于TB级别的内存， 2MB是默认的页大小。

## Transparent Huge Pages ##

Transparent Huge Pages缩写THP，这个是RHEL 6（其它分支版本SUSE Linux Enterprise Server 11, and Oracle Linux 6 with earlier releases of Oracle Linux Unbreakable Enterprise Kernel 2 (UEK2)）开始引入的一个功能，在Linux6上透明大页是默认启用的。由于Huge pages很难手动管理，而且通常需要对代码进行重大的更改才能有效的使用，因此RHEL 6 开始引入了Transparent Huge Pages（ THP ），THP是一个抽象层，能够自动创建、管理和使用传统大页。THP为系统管理员和开发人员减少了很多使用传统大页的复杂性 ,  因为THP的目标是改进性能,因此其它开发人员  ( 来自社区和红帽 )已在各种系统、配置、应用程序和负载中对THP进行了测试和优化。这样可让THP的默认设置改进大多数系统配置性能。但是,不建议对数据库工作负载使用THP 。
Huge pages与Transparent Huge Pages最大的区别是：标准大页管理是预分配的方式，而透明大页管理则是动态分配的方式。

## 相关命令 ##

**查看标准大页（Huage Pages)的页面大小**

```
[root@Neo ~]$ grep Hugepagesize /proc/meminfo
Hugepagesize:     2048 kB
```

**查看是否启用透明大页**

`cat /sys/kernel/mm/redhat_transparent_hugepage/enabled`，该命令适用于Red Hat Enterprise Linux系统

`cat /sys/kernel/mm/transparent_hugepage/enabled`，该命令适用于其它Linux系统

如果输出结果为[always]表示透明大页启用了。[never]表示透明大页禁用、[madvise]表示（只在MADV_HUGEPAGE标志的VMA中使用THP。

`cat /proc/sys/vm/nr_hugepages`，返回0也意味着传统大页禁用了（传统大页和透明大页）。

`grep -i HugePages_Total /proc/meminfo`，如果HugePages_Total返回0，也意味着标准大页禁用了（注意传统/标准大页和透明大页的区别）