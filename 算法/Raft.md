# Raft


https://raft.github.io/
https://thesecretlivesofdata.com/raft/

## 角色

Raft将系统中的角色分为领导者(Leader)、跟从者(Follower)和候选人(Candidate):
- Leader: 接受客户端请求，并向Follower同步请求日志，当日志同步到大多数节点上后告诉Follower提交日志。
- Follower: 接受并持久化Leader同步的日志，在Leader告之日志可以提交之后，提交日志。
- Candidate: Leader选举过程中的临时角色。

Raft要求系统在任意时刻最多只有一个Leader，正常工作期间只有Leader和Followers。

Follower只响应其他服务器的请求。如果Follower超时没有收到Leader的消息，它会成为一个Candidate并且开始一次Leader选举。收到大多数服务器投票的Candidate会成为新的Leader。Leader在宕机之前会一直保持Leader的状态
