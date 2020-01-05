# Linux panic #

[https://blog.51cto.com/xiexf/1942364](https://blog.51cto.com/xiexf/1942364)

有两种主要类型kernel panic：

1.hard panic(也就是Aieee信息输出)
2.soft panic (也就是Oops信息输出)

只有加载到内核空间的驱动模块才能直接导致kernel panic，可以在系统正常的情况下，使用lsmod查看当前系统加载了哪些模块。除此之外，内建在内核里的组件（比如memory map等）也能导致panic。

常见Linux Kernel Panic报错内容：

- Kernel panic - not syncing fatal exception in interrupt
- kernel panic – not syncing: Attempted to kill the idle task!
- kernel panic – not syncing: killing interrupt handler!
- Kernel Panic – not syncing：Attempted to kill init !

## hard panic ##

一般出现下面的情况，就认为是发生了kernel panic:

- 机器彻底被锁定，不能使用
- 数字键(Num Lock)，大写锁定键(Caps Lock)，滚动锁定键(Scroll Lock)不停闪烁。
- 如果在终端下，应该可以看到内核dump出来的信息（包括一段”Aieee”信息或者”Oops”信息）
- 和Windows蓝屏相似

**原因**

对于hard panic而言，最大的可能性是驱动模块的中断处理(interrupt handler)导致的，一般是因为驱动模块在中断处理程序中访问一个空指针(null pointre)。一旦发生这种情况，驱动模块就无法处理新的中断请求，最终导致系统崩溃。

**信息收集**

根据panic的状态不同，内核将记录所有在系统锁定之前的信息。因为kenrel panic是一种很严重的错误，不能确定系统能记录多少信息，下面是一些需要收集的关键信息，他们非常重要，因此尽可能收集全，当然如果系统启动的时候就kernel panic，那就无法只知道能收集到多少有用的信息了。

- /var/log/messages: 幸运的时候，整个kernel panic栈跟踪信息都能记录在这里。
- 应用程序/库日志: 可能可以从这些日志信息里能看到发生panic之前发生了什么。
- 其他发生panic之前的信息，或者知道如何重现panic那一刻的状态
- 终端屏幕dump信息，一般OS被锁定后，复制，粘贴肯定是没戏了，因此这类信息，你可以需要借助数码相机或者原始的纸笔工具了。
- 如果kernel dump信息既没有在/var/log/message里，也没有在屏幕上，那么尝试下面的方法来获取（当然是在还没有死机的情况下）：
- 如果在图形界面，切换到终端界面，dump信息是不会出现在图形界面的，甚至都不会在图形模式下的虚拟终端里。
- 确保屏幕不黑屏，可以使用下面的几个方法：
	- setterm -blank 0
	- setterm -powerdown 0
	- setvesablank off

完整栈跟踪信息的排查方法：

栈跟踪信息(stack trace)是排查kernel panic最重要的信息，该信息如果在/var/log/messages日志里当然最好，因为可以看到全部的信息，如果仅仅只是在屏幕上，那么最上面的信息可能因为滚屏消失了，只剩下栈跟踪信息的一部分。如果你有一个完整栈跟踪信息的话，那么就可能根据这些充分的信息来定位panic的根本原因。要确认是否有一个足够的栈跟踪信息，你只要查找包含”EIP”的一行，它显示了是什么函数和模块调用时导致panic。

内核调试工具(kenrel debugger ,aka KDB)：

如果跟踪信息只有一部分且不足以用来定位问题的根本原因时，kernel debugger(KDB)就需要请出来了。
KDB编译到内核里，panic发生时，他将内核引导到一个shell环境而不是锁定。这样，我们就可以收集一些与panic相关的信息了，这对我们定位问题的根本原因有很大的帮助。

使用KDB需要注意，内核必须是基本核心版本，比如是2.4.18，而不是2.4.18-5这样子的，因为KDB仅对基本核心有效。

## soft panic ##

特征：

- 没有hard panic严重
- 通常导致段错误(segmentation fault)
- 可以看到一个oops信息，/var/log/messages里可以搜索到’Oops’
- 机器稍微还能用（但是收集信息后，应该重启系统）

**原因**

凡是非中断处理引发的模块崩溃都将导致soft panic。在这种情况下，驱动本身会崩溃，但是还不至于让系统出现致命性失败，因为它没有锁定中断处理例程。导致hard panic的原因同样对soft panic也有用（比如在运行时访问一个空指针)

**信息收集**

当soft panic发生时，内核将产生一个包含内核符号(kernel symbols)信息的dump数据，这个将记录在/var/log/messages里。为了开始排查故障，可以使用ksymoops工具来把内核符号信息转成有意义的数据。

为了生成ksymoops文件,需要：

- 从/var/log/messages里找到的堆栈跟踪文本信息保存为一个新文件。确保删除了时间戳(timestamp)，否则ksymoops会失败。
- 运行ksymoops程序（如果没有，请安装）
- 详细的ksymoops执行用法，可以参考ksymoops(8)手册。