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

- 增加/删除元素（`

-  key member [member …]` / `SREM key member [member …]`），SADD命令返回值是成功加入的元素数量，SREM命令用来从集合中删除一个或多个元素，并返回删除成功的个数。

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


**有序集合类型**

有序集合类型和列表类型的异同：

	- 二者都是有序的。
	- 二者都可以获得某一范围的元素。但是二者有着很大的区别，这使得它们的应用场景也是不同的。
	- 列表类型是通过链表实现的，获取靠近两端的数据速度极快，而当元素增多后，访问中间数据的速度会较慢，所以它更加适合实现如“新鲜事”或“日志”这样很少访问中间元素的应用。
	- 有序集合类型是使用散列表和跳跃表（Skip list）实现的，所以即使读取位于中间部分的数据速度也很快（时间复杂度是O(log(N))）。
	- 列表中不能简单地调整某个元素的位置，但是有序集合可以（通过更改这个元素的分数）。
	- 有序集合要比列表类型更耗费内存。

- 增加元素（`ZADD key score member [score member …]`），ZADD 命令用来向有序集合中加入一个元素和该元素的分数，如果该元素已经存在则会用新的分数替换原有的分数。ZADD命令的返回值是新加入到集合中的元素个数。

	`redis> ZADD scoreboard 89 Tom 67 Peter 100 David`

- 获得元素的分数（`ZSCORE key member`）

	`redis> ZSCORE scoreboard Tom`

- 获得排名在某个范围的元素列表（`ZRANGE key start stop [WITHSCORES]` / `ZREVRANGE key start stop [WITHSCORES]`），ZRANGE命令会按照元素分数从小到大的顺序返回索引从 start到stop之间的所有元素（包含两端的元素）。ZREVRANGE命令和ZRANGE的唯一不同在于ZREVRANGE命令是按照元素分数从大到
小的顺序给出结果的。

	`ZRANGE scoreboard 0 2`

- 获得指定分数范围的元素（`ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]`），如果希望分数范围不包含端点值，可以在分数前加上“(”符号。min和max还支持无穷大，同ZADD命令一样，-inf和+inf分别表示负无穷和正无穷。

	`redis> ZRANGEBYSCORE scoreboard (80 +inf`

- 增加某个元素的分数（`ZINCRBY key increment member`）

	`redis> ZINCRBY scoreboard 4 Jerry`

- 获得集合中元素的数量（`ZCARD key`）

	`redis> ZCARD scoreboard`

- 获得指定分数范围内的元素个数（`ZCOUNT key min max`）

	`redis> ZCOUNT scoreboard 90 100`

- 删除一个或多个元素（`ZREM key member [member …]`），ZREM命令的返回值是成功删除的元素数量（不包含本来就不存在的元素）。

	`redis> ZREM scoreboard Wendy`

- 按照排名范围删除元素（`ZREMRANGEBYRANK key start stop`），ZREMRANGEBYRANK命令按照元素分数从小到大的顺序（即索引0表示最小的值）删除处在指定排名范围内的所有元素，并返回删除的元素数量。

	`redis> ZREMRANGEBYRANK testRem 0 2`

- 按照分数范围删除元素（`ZREMRANGEBYSCORE key min max`），ZREMRANGEBYSCORE命令会删除指定分数范围内的所有元素

	`redis> ZREMRANGEBYSCORE testRem (4 5`

- 获得元素的排名（`ZRANK key member` / `ZREVRANK key member`），ZRANK命令会按照元素分数从小到大的顺序获得指定的元素的排名（从0开始，即分数最小的元素排名为0）。

	`redis> ZRANK scoreboard Peter`


## 进阶 ##

### 事务 ###

**错误处理**

- 语法错误。

	语法错误指命令不存在或者命令参数的个数不对。如果事务中有一个命令执行错误，执行EXEC命令后Redis就会直接返回错误，连语法正确的命令也不会执行。注意：Redis2.6.5之前的版本会忽略有语法错误的命令，然后执行事务中其他语法正确的命令。

- 运行错误。

	运行错误指在命令执行时出现的错误。这种错误在实际执行之前 Redis 是无法发现的，所以在事务里这样的命令是会被Redis接受并执行的。如果事务里的一条命令出现了运行错误，事务里其他的命令依然会继续执行（包括出错命令之后的命令）。

注意： Redis的事务没有关系数据库事务提供的回滚（rollback）功能。


**WATCH命令**

WATCH 命令可以监控一个或多个键，一旦其中有一个键被修改（或删除），之后的事务就不会执行。监控一直持续到EXEC命令（事务中的命令是在EXEC之后才执行）。执行EXEC命令后会取消对所有键的监控，如果不想执行事务中的命令也可以使用
UNWATCH命令来取消监控。

```
redis> SET key 1
OK
redis> WATCH key
OK
redis> SET key 2
OK
redis> MULTI
OK
redis> SET key 3
QUEUED
redis> EXEC
(nil)
redis> GET key
"2"
```

### 过期时间 ###

**EXPIRE命令**

使用方法为`EXPIRE key seconds`，其中seconds参数表示键的过期时间，单位是秒。可以使用PERSIST命令取消键的过期时间设置。EXPIRE命令的seconds参数必须是整数，所以最小单位是1秒。如果想要更精确的控制键的过期时间应该使用 PEXPIRE命令，PEXPIRE命令与EXPIRE的唯一区别是前者的时间单位是毫秒。

```
redis> SET session:29e3d uid1314
OK
redis> EXPIRE session:29e3d 900
(integer) 1
```

使用TTL命令可以知道一个键还有多久的时间会被删除。当键不存在时TTL命令会返回−2。

如果使用WATCH命令监测了一个拥有过期时间的键，该键时间到期自动删除并不会被WATCH命令认为该键被改变。

**缓存过期**

LRU(Least Recently Used) 算法是众多置换算法中的一种。Redis用到的LRU算法，是一种近似的LRU算法。。 Redis配置文件的maxmemory参数，限制Redis最大可用内存大小（单位是字节），当超出了这个限制时Redis会依据maxmemory policy参数指定的策略来删除不需要的键直到Redis占用的内存小于指定内存。

置换策略：

- noeviction: 不进行置换，表示即使内存达到上限也不进行置换，所有能引起内存增加的命令都会返回error
- allkeys-lru: 优先删除掉最近最不经常使用的key，用以保存新数据
- volatile-lru: 只从设置失效（expire set）的key中选择最近最不经常使用的key进行删除，用以保存新数据
- allkeys-random: 随机从all-keys中选择一些key进行删除，用以保存新数据
- volatile-random: 只从设置失效（expire set）的key中，选择一些key进行删除，用以保存新数据
- volatile-ttl: 只从设置失效（expire set）的key中，选出存活时间（TTL）最短的key进行删除，用以保存新数据


### 排序 ###

SORT命令可以对列表类型、集合类型和有序集合类型键进行排序，并且可以完成与关系数据库中的连接查询相类似的任务。

BY参数的语法为BY参考键。其中参考键可以是字符串类型键或者是散列类型键的某个 字段（表示为键名->字段名）。如果提供了BY参数，SORT命令将不再依据元素自身的值 进行排序，而是对每个元素使用元素的值替换参考键中的第一个“*”并获取其值，然后依据该值对元素排序。

```
redis> SORT tag:ruby:posts BY post:*->time DESC 
1) "12" 
2) "26" 
3) "6" 
4) "2" 
``` 

GET参数不影响排序，它的作用是使SORT命令的返回结果不再是元素自身的值，而是GET参数中指定的键值。GET参数的规则和BY参数一样，GET参数也支持字符串类型和散列 类型的键，并使用“*”作为占位符。

`redis> SORT tag:ruby:posts BY post:*->time DESC GET post:*->title GET post:*->time`

默认情况下SORT会直接返回排序结果，如果希望保存排序结果，可以使用STORE参数。保存后的键的类型为列表类型，如果键已经存在则会覆盖它。加上STORE参数后SORT命令的返回值为结果的个数。 STORE参数常用来结合EXPIRE命令缓存排序结果。


### 消息通知 ###

**任务队列**

让生产者将任务使用LPUSH命令加入到某个键中，另一边让消费者不断地使用RPOP命令从该键中取出任务即可。使用DROP命令有一个问题，当任务队列中没有任务时消费者会频繁使用RPOP命令查看是否有新任务。BRPOP命令和RPOP命令相似，唯一的区别是当列表中没有元素时BRPOP命令会一直阻 塞住连接，直到有新元素加入。BRPOP命令接收两个参数，第一个是键名，第二个是超时时间，单位是秒。当超过了此时间仍然没有获得新元素的话就会返回nil。如果超时时间为"0"，表示不限制等待的时间，即如果没有新元素加入列表就会永远阻塞下去。 当获得一个元素后 BRPOP命令返回两个值，分别是键名和元素值。

**优先级队列**

BRPOP命令可以同时接收多个键，其完整的命令格式为`BLPOP key [key …] timeout`， 如`BLPOP queue:1 queue:2 0`。意义是同时检测多个键，如果所有键都没有元素则阻塞，如果其中有一个键有元素则会从该键中弹出元素。如果多个键都有元素则按照从左到右的顺序取第一个键中的一个元素。

**“发布/订阅”模式**

“发布/订阅”模式中包含两种角色，分别是发布者和订阅者。订阅者可以订阅一个或若干个频道（channel），而发布者可以向指定的频道发送消息，所有订阅此频道的订阅者都会收到此消息。 

发布者发布消息的命令是PUBLISH，用法是`PUBLISH channel message`。。PUBLISH命令的返回值表示接收到这条消息的订阅者数量。发出去的消息不会被持久化，也就是说客户端只能收到后续发布到该频道的消息，无法收到之前发送的消息。 

订阅频道的命令是SUBSCRIBE，可以同时订阅多个频道，用法是`SUBSCRIBE channel [channel …]`。执行SUBSCRIBE命令后客户端会进入订阅状态，处于此状态下客户端不能使用除`SUBSCRIBE、UNSUBSCRIBE、PSUBSCRIBE和PUNSUBSCRIBE`这4个属于“发布/订阅”模式的命令之外的命令，否则会报错。 进入订阅状态后客户端可能收到3种类型的回复。每种类型的回复都包含3个值，第一个值是消息的类型，根据消息类型的不同，第二、三个值的含义也不同。消息类型可能的取值有以下3个。 

- （1）subscribe。表示订阅成功的反馈信息。第二个值是订阅成功的频道名称，第三个值是当前客户端订阅的频道数量。 
- （2）message。这个类型的回复是我们最关心的，它表示接收到的消息。第二个值表示产生消息的频道名称，第三个值是消息的内容。 
- （3）unsubscribe。表示成功取消订阅某个频道。第二个值是对应的频道名称，第三个值是当前客户端订阅的频道数量，当此值为0时客户端会退出订阅状态，之后就可以执行其他非“发布/订阅”模式的命令了。 

使用UNSUBSCRIBE命令可以取消订阅指定的频道，用法为`UNSUBSCRIBE [channel [channel …]]`，如果不指定频道则会取消订阅所有频道。

除了可以使用SUBSCRIBE命令订阅指定名称的频道外，还可以使用PSUBSCRIBE命令订阅指定的规则。规则支持glob风格通配符格式。

	**注意**：使用PUNSUBSCRIBE命令只能退订通过PSUBSCRIBE命令订阅的规则，不会影响直接通过SUBSCRIBE命令订阅的频道；同样UNSUBSCRIBE命令也不会影响通过PSUBSCRIBE命令订阅的规则。另外容易出错的一点是使用PUNSUBSCRIBE命令退订某个规 则时不会将其中的通配符展开，而是进行严格的字符串匹配，所以`PUNSUBSCRIBE *`无法退订`channel.*`规则，而是必须使用`PUNSUBSCRIBE channel.*`才能退订。


## 持久化 ##

### RBD ###

RDB方式的持久化是通过快照（snapshotting）完成的，当符合一定条件时Redis会自动将内存中的所有数据生成一份副本并存储在硬盘上，这个过程即为“快照”。Redis会在以下几种情况下对数据进行快照：
- 根据配置规则进行自动快照；
- 用户执行SAVE或BGSAVE命令；
- 执行FLUSHALL命令；
- 执行复制（replication）时。

**根据配置规则进行自动快照**

进行快照的条件可以由用户在配置文件中自定义，由两个参数构成：时间窗口M和改动的键的个数N。每当时间M内被更改的键的个数大于N时，即符合自动快照条件。每条快照条件占一行，并且以save参数开头。同时可以存在多个条件，条件之间是“或”的关系。

```
save 900 1
save 300 10
save 60 10000
```

**用户执行SAVE或BGSAVE命令**

- SAVE命令
	当执行SAVE命令时，Redis同步地进行快照操作，在快照执行的过程中会阻塞所有来自客户端的请求。当数据库中的数据比较多时，这一过程会导致Redis较长时间不响应，所以要尽量避免在生产环境中使用这一命令。

- BGSAVE命令
	需要手动执行快照时推荐使用BGSAVE命令。BGSAVE命令可以在后台异步地进行快照操作，快照的同时服务器还可以继续响应来自客户端的请求。执行BGSAVE后Redis会立即返回OK表示开始执行快照操作，如果想知道快照是否完成,可以通过LASTSAVE命令获取最近一次成功执行快照的时间，返回结果是一个Unix时间戳。

**执行FLUSHALL命令**

当执行FLUSHALL命令时，Redis会清除数据库中的所有数据。需要注意的是，不论清空数据库的过程是否触发了自动快照条件，只要自动快照条件不为空，Redis就会执行一次快照操作。当没有定义自动快照条件时，执行FLUSHALL则不会进行快照。

**快照原理**

Redis默认会将快照文件存储在Redis当前进程的工作目录中的dump.rdb文件中，可以通过配置dir和dbfilename两个参数分别指定快照文件的存储路径和文件名。快照的过程：
- Redis使用fork函数复制一份当前进程（父进程）的副本（子进程）；
- 父进程继续接收并处理客户端发来的命令，而子进程开始将内存中的数据写入硬盘中的临时文件；
- 当子进程写入完所有数据后会用该临时文件替换旧的RDB文件，至此一次快照操作完成。

几个注意事项：

- 执行fork的时候操作系统（类Unix操作系统）会使用写时复制（copy-onwrite）策略，即fork函数发生的一刻父子进程共享同一内存数据，当父进程要更改其中某片数据时（如执行一个写命令），操作系统会将该片数据复制一份以保证子进程的数据不受影响，所以新的RDB文件存储的是执行fork一刻的内存数据。
- 写时复制策略也保证了在fork的时刻虽然看上去生成了两份内存副本，但实际上内存的占用量并不会增加一倍。
- 当进行快照的过程中，如果写入操作较多，造成fork前后数据差异较大，是会使得内存使用量显著超过实际数据大小的，因为内存中不仅保存了当前的数据库数据，而且还保存着fork时刻的内存数据。
- 通过RDB方式实现持久化，一旦Redis异常退出，就会丢失最后一次快照以后更改的所有数据。


### AOF ###

AOF可以将Redis执行的每一条写命令追加到硬盘文件中，默认情况下Redis没有开启AOF（append only file）方式的持久化，可以通过appendonly参数启用：`appendonly yes`。AOF文件的保存位置和RDB文件的位置相同，都是通过dir参数设置，默认的文件名是appendonly.aof，可以通过appendfilename参数修改：`appendfilename appendonly.aof`。

**AOF的实现**

AOF文件的内容是Redis客户端向Redis发送的原始通信协议的内容。随着执行的命令越来越多，AOF文件的大小也会越来越
大，Redis 可以自动优化AOF文件，每当达到一定条件时Redis就会自动重写AOF文件，这个条件可以在配置文件中设置：

- `auto-aof-rewrite-percentage 100`
	
	当目前的AOF文件大小超过上一次重写时的AOF文件大小的百分之多少时会再次进行重写，如果之前没有重写过，则以启动时的AOF文件大小为依据。

- `auto-aof-rewrite-min-size 64mb`

	限制允许重写的最小AOF文件大小，通常在AOF文件很小的情况下即使其中有很多冗余的命令我们也并不太关心。

除了让Redis自动执行重写外，我们还可以主动使用`BGREWRITEAOF`命令手动执行AOF重写。

**同步硬盘数据**

虽然AOF每次将写命令写入AOF文件中，但由于操作系统的限制，数据并没有真正意义上落盘。在默认情况下系统每30秒会执行一次同步操作，以便将硬盘缓存中的内容真正地写入硬盘，在这30秒的过程中如果系统异常退出则会导致硬盘缓存中的数据丢失。可以通过`appendfsy nc`参数设置同步的时机：

```
# appendfsy nc alway s
appendfsy nc every sec
# appendfsy nc no
```

默认情况下Redis采用every sec规则，即每秒执行一次同步操作。alway s表示每次执行写入都会执行同步，这是最安全也是最慢的方式。no表示不主动进行同步操作，而是完全交由操作系统来做（即每30秒一次），这是最快但最不安全的方式。


## 集群 ##

### 复制 ###

**原理**

复制初始化阶段：

	从数据库启动后，会向主数据库发送SYNC命令。同时主数据库接收到SYNC命令后开始在后台保存快照（即RDB持久化的过程），并将保存快照期间接收到的命令缓存起来。当快照完成后，Redis会将快照文件和所有缓存的命令发送给从数据库。从数据库收到后，载入快照文件并执行收到的缓存的命令。

复制同步阶段：

	复制初始化阶段结束后，主数据库执行的任何会导致数据变化的命令都会异步地传送给从数据库。复制同步阶段会贯穿整个主从同步过程的始终，直到主从关系终止为止。

注意：
	
	当主从数据库之间的连接断开重连后，Redis 2.6以及之前的版本会重新进行复制初始化（即主数据库重新保存快照并传送给从数据库），即使从数据库可以仅有几条命令没有收到，主数据库也必须要将数据库里的所有数据重新传送给从数据库。这使得主从数据库断线重连后的数据恢复过程效率很低下，在网络环境不好的时候这一问题尤其明显。Redis 2.8版的一个重要改进就是断线重连能够支持有条件的增量数据传输，当从数据库重新连接上主数据库后，主数据库只需要将断线期间执行的命令传送给从数据库，从而大大提高Redis复制的实用性。

**乐观复制**

Redis采用了乐观复制（optimistic replication）的复制策略，容忍在一定时间内主从数据库的内容是不同的，但是两者的数据会最终同步。Redis在主从数据库之间复制数据的过程本身是异步的，主数据库执行完客户端请求的命令后会立即将
命令在主数据库的执行结果返回给客户端，并异步地将命令同步给从数据库。Redis 提供了两个配臵选项来限制只有当数据至少同步给指定数量的从数据库时，主数据库才是可写的：

- min-slaves-to-write 3

	表示只有当3个或3个以上的从数据库连接到主数据库时，主数据库才是可写的，否则会返回错误

- min-slaves-max-lag 10

	表示允许从数据库最长失去连接的时间，如果从数据库最后与主数据库联系（即发送REPLCONF ACK命令）的时间小于这个值，则认为从数据库还在保持与主数据库的连接

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
