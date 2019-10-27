# Redis入门指南 #

## Redis特性 ##

Redis支持的键值数据类型：
● 字符串类型
● 散列类型
● 列表类型
● 集合类型
● 有序集合类型

Redis数据库中的所有数据都存储在内存中。同时提供了对持久化的支持，即可以将内存中的数据异步写入到硬盘中，同时不影响继续提供服务。

[Redis的相关命令](https://redis.io/commands)

### 常用命令： ###

**字符串类型：**

- 增加指定的整数（`INCRBY key increment`）
	
	`redis> INCRBY bar 2`

- 减少指定的整数（`DECRBY key decrement`）

	`redis> DECRBY key 5 `

- 增加指定浮点数（`INCRBYFLOAT key increment`）

	`redis> INCRBYFLOAT bar 2.7`

- 向尾部追加值（`APPEND key value`）

	`redis> APPEND key " world!"`

- 获取字符串长度（`STRLEN key`）

	`redis> STRLEN key`

- 同时获得/设置多个键值（`MGET key [key …]` / `MSET key value [key value …]`）

	`redis> MSET key1 v1 key2 v2 key3 v3`
	`redis> MGET key1 key3 `

**散列类型**

- 赋值与取值（`HSET key field value` / `HGET key field` / `HMSET key field value [field value …]` / `HMGET key field [field …]` / `HGETALL key`）

	`redis> HSET car price 500`
	`redis> HSET car name BMW`
	`redis> HGET car name "BMW"`

- 判断字段是否存在（`HEXISTS key field `）,如果存在则返回1，否则返回0（如果键不存 在也会返回0）

	` HEXISTS car model`

- 当字段不存在时赋值（`HSETNX key field value`），如果字段已经存在，HSETNX命令将不 执行任何操作

- 增加数字（`HINCRBY key field increment`）

	`redis> HINCRBY person score 60`

- 删除字段（`HDEL key field [field …]`），可以删除一个或多个字段，返回值是被删除的字段个数

	`redis> HDEL car price`







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
