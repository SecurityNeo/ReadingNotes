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


**列表类型**

- 向列表两端增加元素（`LPUSH key value [value …]` / `RPUSH key value [value …]`）,LPUSH命令用来向列表左边增加元素，返回值表示增加元素后列表的长度。 向列表右边增加元素的话则使用RPUSH命令。

	`redis> LPUSH numbers 1`
	`redis> LPUSH numbers 2 3`
	`redis> RPUSH numbers 0 −1`

- 从列表两端弹出元素（`LPOP key` / `RPOP key `），LPOP命令可以从列表左边弹出一个元素。LPOP命令执行两步操作：第一步 是将列表左边的元素从列表中移除，第二步是返回被移除的元素值。

	`redis> LPOP numbers`

- 获取列表中元素的个数（`LLEN key`），当键不存在时LLEN会返回0

	`redis> LLEN numbers`

- 获得列表片段（`LRANGE key start stop`），LRANGE 命令将返回索引从 start到 stop之间的所有元素（包含两端的元素）。

	`redis> LRANGE numbers 0 2`

- 删除列表中指定的值（`LREM key count value `），LREM命令会删除列表中前count个值为value的元素，返回值是实际删除的元素个数。根据count值的不同，LREM命令的执行方式会略有差异。

	- 当count > 0时，LREM命令会从列表左边开始删除前count个值为value的元素。 
	- 当count < 0时，LREM命令会从列表右边开始删除前|count|个值为value的元素。 
	- 当count = 0是，LREM命令会删除所有值为value的元素。

- 获得/设置指定索引的元素值（`LINDEX key index` / `LSET key index value`），。LINDEX命令用来返回 指定索引的元素，索引从0开始，如果index是负数则表示从右边开始计算的索引，最右边元素的索引是−1。LSET是另一个通过索引操作列表的命令，它会将索引为index的元素赋值为value。

	`redis> LINDEX numbers 0`
	`redis> LSET numbers 1 7` 

- 只保留列表指定片段（`LTRIM key start end`），LTRIM 命令可以删除指定索引范围之外的所有元素，其指定列表范围的方法和LRANGE 命令相同。

	`redis> LRANGE numbers 0 1`

- 向列表中插入元素（`LINSERT key BEFORE|AFTER pivot value`）。LINSERT 命令首先会在列表中从左到右查找值为pivot的元素，然后根据第二个参数是BEFORE还是AFTER来决定将value插入到该元素的前面还是后面。LINSERT命令的返回值是插入后列表的元素个数。

	`redis> LINSERT numbers AFTER 7 3`
	`(integer) 4` 

- 将元素从一个列表转到另一个列表（`RPOPLPUSH source destination`）。先执行RPOP命令再 执行LPUSH命令。RPOPLPUSH命令会先从source列表类型键的右边弹出一个元素，然后将其加入到destination列表类型键的左边，并返回这个元素的值，整个过程是原子的。


**集合类型**

- 增加/删除元素（`SADD key member [member …]` / `SREM key member [member …]`），SADD命令返回值是成功加入的元素数量，SREM命令用来从集合中删除一个或多个元素，并返回删除成功的个数。

	`redis> SADD letters a b c`
	`redis> SREM letters c d`

- 获得集合中的所有元素（`SMEMBERS key`）

	`redis> SMEMBERS letters`

- 判断元素是否在集合中（`SISMEMBER key member`）,当值存在时SISMEMBER命令返回1，当值不存在或键不存在时返回0

	`redis> SISMEMBER letters a`

- 集合间运算

	`SDIFF key [key „]`，对多个集合执行差集运算
	`SINTER key [key „]`，对多个集合执行交集运算
	`SUNION key [key „]`，对多个集合执行并集运算

- 进行集合运算并将结果存储

	`SDIFFSTORE destination key [key …]`
	`SINTERSTORE destination key [key …]`
	`SUNIONSTORE destination key [key …]`

- 获得集合中元素个数（`SCARD key`）

	`redis> SCARD letters`

- 随机获得集合中的元素（`SRANDMEMBER key [count]`），可以传递count参数来一次随机获得多个元素，根据count的正负不同，具体表现也不同：
	- 当count为正数时，SRANDMEMBER会随机从集合里获得count个不重复的元素。如果count的值大于集合中的元素个数，则SRANDMEMBER会返回集合中的全部元素。
	- 当count为负数时，SRANDMEMBER会随机从集合里获得|count|个的元素，这些元素有可能相同。

- 从集合中弹出一个元素（`SPOP key`）






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
