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

