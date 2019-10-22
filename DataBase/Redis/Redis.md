# Redis #

## 简单动态字符串 ##

Sds（Simple Dynamic String，简单动态字符串）是Redis底层所使用的字符串表示。对比C字符串， sds有以下特性：

- 可以高效地执行长度计算（strlen）
- 可以高效地执行追加操作（append）
- 二进制安全

**用途**

- 实现字符串对象（StringObject）
- 在Redis程序内部用作`char*`类型的替代品
	在Redis中， 客户端传入服务器的协议内容、 aof缓存、 返回给客户端的回复等等都是由sds类型来保存的。

**sds模块的API**

|函数|作用|算法复杂度|
|sdsnewlen|创建一个指定长度的sds ，接受一个C字符串作为初始化值|O(N)
|sdsempty|创建一个只包含空白字符串 "" 的sds|O(1)
|sdsnew|根据给定C字符串，创建一个相应的sds|O(N)
|sdsdup|复制给定sds|O(N)
|sdsfree|释放给定sds|O(N)
|sdsupdatelen|更新给定sds所对应sdshdr结构的free和len|O(N)
|sdsclear|清除给定sds的内容，将它初始化为 ""|O(1)
|sdsMakeRoomFor|对sds所对应sdshdr结构的buf进行扩展|O(N)
|sdsRemoveFreeSpace|在不改动buf的情况下，将buf内多余的空间释放出去|O(N)
|sdsAllocSize|计算给定sds的buf所占用的内存总数|O(1)
|sdsIncrLen|对sds的buf的右端进行扩展（expand）或修剪（trim）|O(1)
|sdsgrowzero|将给定sds的buf扩展至指定长度，无内容的部分用\0来填充|O(N)
|sdscatlen|按给定长度对sds进行扩展，并将一个C字符串追加到sds的末尾|O(N)
|sdscat|将一个C字符串追加到sds末尾|O(N)
|sdscatsds|将一个sds追加到另一个sds末尾|O(N)
|sdscpylen|将一个C字符串的部分内容复制到另一个sds中，需要时对sds进行扩展|O(N)
|sdscpy|将一个C字符串复制到sds|O(N)
