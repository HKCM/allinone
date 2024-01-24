
在parameter group中启用`binlog_format ROW`

1. 查看Blue数据
```sql
select version();
show master status; 
```

2. 查看Green状态

```sql
select version();
show master status; 
show slave status;
show global variables like "log_bin"; # 检查是否开启binlog
show global variables like "binlog_format"; # 检查是否开启binlog格式
show binary logs; # 查看binlog文件大小
call mysql.rds_show_configuration; # 查看binlog retention时间
call mysql.rds_set_configuration('binlog retention hours',72); # 设置binlog retention时间为72小时
select @@GLOBAL.transaction_isolation,@@transaction_isolation;
show VARIABLES like 'character_set%';
show global VARIABLES like 'read_only';
```

通过在CloudWatch的 AuroraBinlogReplicaLag 查看同步滞后的情况

https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/CHAP_Troubleshooting.html#CHAP_Troubleshooting.MySQL.ReplicaLag

https://docs.aws.amazon.com/zh_cn/AmazonRDS/latest/AuroraUserGuide/Aurora.Managing.Clone.html