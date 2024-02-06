```bash
docker network create -d bridge redis-network
docker run --network redis-network --name my-redis -d redis:7.0
docker run -it --network redis-network --rm redis:7.0 redis-cli -h my-redis
```

命令查询：https://www.redis.com.cn/commands

## String

- SET 获取存储在给定键中的值
- GET 设置存储在给定键中的值
- DEL 删除存储在给定键中的值
- INCR 将键存储的值加1
- INCRBY 将键存储的值加上整数
- DECR 将键存储的值减1
- DECRBY 将键存储的值减整数


```bash
my-redis:6379> SET hello world # OK
my-redis:6379> GET hello # "world"
my-redis:6379> DEL hello # (integer) 1
my-redis:6379> GET hello # (nil)
my-redis:6379> SET counter 2 #OK
my-redis:6379> GET counter # "2"
my-redis:6379> INCR counter # (integer) 3
my-redis:6379> GET counter # "3"
my-redis:6379> INCRBY counter 100 # (integer) 103
my-redis:6379> GET counter # "103"
my-redis:6379> DECR counter # (integer) 102
my-redis:6379> DECRBY counter 100 # (integer) 102
```

## List

- RPUSH 将给定值推入到列表右端
- LPUSH 将给定值推入到列表左端
- RPOP 从列表的右端弹出一个值，并返回被弹出的值
- LPOP 从列表的左端弹出一个值，并返回被弹出的值
- LRANGE 获取列表在给定范围上的所有值
- LINDEX 通过索引获取列表中的元素。也可以使用负数下标，以-1表示列表的最后一个元素

```bash
my-redis:6379> LPUSH mylist 1 2 ll ls mem # (integer) 5
my-redis:6379> LRANGE mylist 0 -1
1) "mem"
2) "ls"
3) "ll"
4) "2"
5) "1"
my-redis:6379> LINDEX mylist -1 # "1"
my-redis:6379> LINDEX mylist 10 # (nil)
```

## Set

- SADD 向集合添加一个或多个成员
- SCARD 获取集合的总数
- SMEMBERS 返回集合中的所有成员
- SISMEMBER 判断元素是否是集合的成员
- SINTER 返回给定所有集合的交集
- SDIFF 返回第一个集合与其他集合之间的差异

https://www.runoob.com/redis/redis-sets.html

```bash
my-redis:6379> SADD set1 1 2 3 4 5 # (integer) 5
my-redis:6379> SADD set2 3 4 5 6 7 # (integer) 5
my-redis:6379> SCARD set1 # 获取集合的总数 (integer) 5
my-redis:6379> SMEMBERS set1 # 返回集合中的所有成员
1) "1"
2) "2"
3) "3"
4) "4"
5) "5"
my-redis:6379> SISMEMBER set1 5 # 判断元素是否是集合的成员 (integer) 1
my-redis:6379> SISMEMBER set1 6 # 判断元素是否是集合的成员 (integer) 0
my-redis:6379> SINTER set1 set2 # 返回给定所有集合的交集
1) "3"
2) "4"
3) "5"
```

## Hash

- HSET 添加键值对
- HGET 获取指定散列键的值
- HGETALL 获取散列中包含的所有键值对
- HDEL 如果给定键存在于散列中，那么就移除这个键

```bash
my-redis:6379> HSET user name Alice
(integer) 1
my-redis:6379> HSET user email alice@163.com
(integer) 1
my-redis:6379> HGETALL user
1) "name"
2) "Alice"
3) "email"
4) "alice@163.com"
my-redis:6379> HGET user name
"Alice"
my-redis:6379> HGET user email
"alice@163.com"
```

## Zset

- ZADD 将一个带有给定分值的成员添加到有序集合里面
- ZCARD 获取有序集合的成员数
- ZRANGE 根据元素在有序集合中所处的位置，从有序集合中获取多个元素
- ZREM 如果给定元素成员存在于有序集合中，那么就移除这个元素
- ZCOUNT 计算在有序集合中指定区间分数的成员数

https://www.runoob.com/redis/redis-sorted-sets.html

```bash
my-redis:6379> ZADD language 1 go 2 java 3 php 4 c 5 c++ 6 c# # (integer) 6
my-redis:6379> ZCARD language # (integer) 6
my-redis:6379> ZRANGE language 2 3
1) "php"
2) "c"
my-redis:6379> ZCOUNT language 2 4 # (integer) 3
my-redis:6379> ZREM language php
(integer) 1
```

## 特殊类型

https://pdai.tech/md/db/nosql-redis/db-redis-data-type-special.html

## Stream

https://pdai.tech/md/db/nosql-redis/db-redis-data-type-stream.html

