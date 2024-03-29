




# 安全性

当使用主从,哨兵或集群模式时,不开启持久化是危险的,例如

1. 设置节点A为主服务器，关闭持久化，节点B和C从节点A复制数据。
2. 此时出现了一个崩溃，但Redis具有自动重启系统，重启了进程，因为关闭了持久化，节点重启后只有一个空的数据集。
3. 由于节点B和C从节点A进行复制，现在节点A是空的，所以节点B和C上的复制数据也会被删除。

当在高可用系统中使用Redis Sentinel，关闭了主服务器的持久化，并且允许自动重启，这种情况是很危险的。

比如主服务器可能在很短的时间就完成了重启，以至于Sentinel都无法检测到这次失败，那么上面说的这种失败的情况就发生了。

如果数据比较重要，并且在使用主从复制时关闭了主服务器持久化功能的场景中，都应该禁止实例自动重启