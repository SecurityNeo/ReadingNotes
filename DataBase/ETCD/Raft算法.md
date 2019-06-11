# Raft算法 #
[https://www.cnblogs.com/hzmark/p/raft.html](https://www.cnblogs.com/hzmark/p/raft.html "https://www.cnblogs.com/hzmark/p/raft.html")

Raft是一种共识算法，旨在替代Paxos。它的首要设计目的就是易于理解，所以在选主的冲突处理等方式上它都选择了非常简单明了的解决方案。

**三种节点角色**

- Leader: 处理所有客户端交互，日志复制等，一般一次只有一个Leader.
- Follower: 类似选民，完全被动
- Candidate候选人: 类似Proposer律师，可以被选为一个新的领导人。、

角色之间的转换关系：

![](img/Raft_state.png)

- 所有节点初始状态都是Follower角色
- 超时时间内没有收到Leader的请求则转换为Candidate进行选举
- Candidate收到大多数节点的选票则转换为Leader；发现Leader或者收到更高任期的请求则转换为Follower
- Leader在收到更高任期的请求后转换为Follower

**任期**

![](img/Raft_Election.png)
![](img/Raft_Election_Term.png)

每个任期都由一次选举开始，若选举失败则这个任期内没有Leader；如果选举出了Leader则这个任期内有Leader负责集群状态管理。Splite Vote是因为如果同时有两个候选人向大家邀票，这时通过类似加时赛来解决，两个候选者在一段timeout比如300ms互相不服气的等待以后，因为双方得到的票数是一样的，一半对一半，那么在300ms以后，再由这两个候选者发出邀票，这时同时的概率大大降低，那么首先发出邀票的的候选者得到了大多数同意，成为领导者Leader，而另外一个候选者后来发出邀票时，那些Follower选民已经投票给第一个候选者，不能再投票给它，它就成为落选者了，最后这个落选者也成为普通Follower一员了。

**心跳（heartbeats）和超时机制（timeout）**

在Raft算法中，有两个timeout机制来控制领导人选举：

- 选举定时器（eletion timeout）：Follower等待成为Candidate状态的等待时间，这个时间被随机设定为150ms~300ms之间
- headrbeat timeout： 在某个节点成为Leader以后，它会发送Append Entries消息给其他节点，这些消息就是通过heartbeat timeout来传送，Follower接收到Leader的心跳包的同时也重置选举定时器。

## 日志复制 ##

- Client向Leader提交指令，Leader收到命令后，将命令追加到本地日志中。此时，这个命令处于“uncomitted”状态，复制状态机不会执行该命令。
- Leader将命令并发复制给其他节点，并等待其他其他节点将命令写入到日志中，如果此时有些节点失败或者比较慢，Leader节点会一直重试，直到所有节点都保存了命令到日志中。之后Leader节点就提交命令，并将结果返回给Client节点。
- Leader节点在提交命令后，下一次的心跳包中就带有通知其他节点提交命令的消息，其他节点收到Leader的消息后，就将命令应用到状态机中（State Machine），最终每个节点的日志都保持了一致性。

