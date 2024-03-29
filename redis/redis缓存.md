# 缓存

在高并发的业务场景下，数据库大多数情况都是用户并发访问最薄弱的环节。所以，就需要使用redis做一个缓冲操作，让请求先访问到redis，而不是直接访问Mysql等数据库。这样可以大大缓解数据库的压力。

使用缓存数据库可能出现以下问题:
- 缓存穿透
- 缓存穿击
- 缓存雪崩
- 缓存污染
- 数据一致性


## 缓存穿透

是指客户端请求的数据在缓存中和数据库中都不存在，这样缓存永远不会生效，这些请求都会打到数据库。

### 示例分析

缓存穿透是指缓存和数据库中都没有的数据，而用户不断发起请求,这将导致这个不存在的数据每次请求都要到存储层去查询，失去了缓存的意义。

在流量大时，可能DB就挂掉了，或者有人利用不存在的key频繁攻击应用，如发起为id为“-1”的数据或id为特别大不存在的数据。这时的用户很可能是攻击者，攻击会导致数据库压力过大。

### 解决方案

1. 缓存null值: 实现简单,维护方便,额外的内存消耗,可能造成短期的不一致
2. 布隆过滤: 内存占用较少,没有多余key,实现复杂,存在误判可能
3. 增强id的复杂度，避免被猜测id规律
4. 做好数据的基础格式校验
5. 加强用户权限校验
6. 做好热点参数的限流

## 缓存击穿

也叫热点Key问题，是指缓存中没有但数据库中有的数据（一般是缓存时间到期），这时由于并发用户特别多，同时读缓存没读到数据，又同时去数据库去取数据，引起数据库压力瞬间增大，造成过大压力

### 示例分析

假设线程1在查询缓存之后，本来应该去查询数据库，然后把这个数据重新加载到缓存的，此时只要线程1走完这个逻辑，其他线程就都能从缓存中加载这些数据了.

假设在线程1没有走完的时候，后续线程2，线程3，线程4同时过来访问同一个数据, 在线程1没完成的情况下缓存未命中，接着同一时间去访问数据库，同时的去执行数据库代码，对数据库访问压力过大

常见的解决方案有两种:
- 热点key永不过期
- 逻辑过期,把过期时间设置在redis的value中通过业务逻辑判断是否过期，如果过期则通过新的子线程去数据库查询更新
- 互斥锁,采用tryLock方法 + double check来
- 接口限流与熔断，降级,防止用户恶意刷接口，同时要降级准备，当接口中的某些服务不可用时候，进行熔断，失败快速返回机制

方案对比:
- 永不过期可能占内存
- 逻辑过期和永不过期实质是一样的,但是具有检查功能
- 互斥锁方案:由于保证了互斥性，所以数据一致，实现简单，仅需加一把锁,没有额外的内存消耗，缺点在于可能存在死锁问题，且只能串行执行性能受到影响
- 逻辑过期方案: 线程读取过程中不需要等待，性能好，有一个额外的线程持有锁去进行重构数据，但是在重构数据完成前，其他的线程只能返回之前的数据，且实现复杂

## 缓存雪崩

是指在同一时段大量的缓存key同时失效或者Redis服务宕机，导致大量请求到达数据库，带来巨大压力。

### 示例分析

缓存中数据大批量到过期时间，而查询数据量巨大，引起数据库压力过大甚至down机。

和缓存击穿不同的是，缓存击穿指并发查同一条数据，缓存雪崩是不同数据都过期了，很多数据都查不到从而查数据库。

### 解决方案

1. 给不同的Key的TTL添加随机值
2. 利用Redis集群提高服务的可用性
3. 给缓存业务添加降级限流策略
4. 给业务添加多级缓存

## 缓存污染

是指缓存中一些只会被访问一次或者几次的的数据，被访问完后，再也不会被访问到，但这部分数据依然留存在缓存中，消耗缓存空间。

存污染会随着数据的持续增加而逐渐显露，随着服务的不断运行，缓存中会存在大量的永远不会再次被访问的数据。缓存空间是有限的，如果缓存空间满了，再往缓存里写数据时就会有额外开销，影响Redis性能。这部分额外开销主要是指写的时候判断淘汰策略，根据淘汰策略去选择要淘汰的数据，然后进行删除操作。

### 缓存淘汰策略

- **内存淘汰：** redis自动进行，当redis内存达到设定的max-memery的时候，会自动触发淘汰机制，淘汰掉一些不重要的数据(可以自己设置策略方式)
- **超时剔除：** 当给redis设置了过期时间ttl之后，redis会将超时的数据进行删除
- **主动更新：** 可以手动调用方法把缓存删掉，通常用于解决缓存和数据库不一致问题

Redis共支持八种淘汰策略,大致分为三类:

1. 不淘汰
   1. noeviction：当内存使用超过配置的时候会返回错误，不会驱逐任何键,无法解决缓存污染问题。一般生产环境不建议使用
2. 对设置了过期时间的数据中进行淘汰
   1. volatile-random：从配置了过期时间的键中进行随机删除。可能依然会存在缓存污染现象，无法解决缓存污染问题
   2. volatile-lfu：从配置了过期时间的键中驱逐使用频率最少的键
   3. volatile-ttl：从配置了过期时间的键中驱逐马上就要过期的键
   4. volatile-lru：从配置了过期时间的键中驱逐最久没有使用的键
3. 全部数据进行淘汰
   1. allkeys-random：从所有键值对中随机选择并删除数据,无法解决缓存污染问题
   2. allkeys-lru：从所有键值对中通过LRU算法驱逐最久没有使用的键
   3. allkeys-lfu：从所有键值对中驱逐使用频率最少的键

### LRU算法
LRU 算法的全称是 Least Recently Used,按照最近最少使用的原则来筛选数据。这种模式下会使用 LRU 算法筛选设置了过期时间的键值对。

Redis会记录每个数据的最近一次被访问的时间戳。在Redis在决定淘汰的数据时，第一次会随机选出 N 个数据，把它们作为一个候选集合。

接下来，Redis 会比较这 N 个数据的 lru 字段，把 lru 字段值最小的数据从缓存中淘汰出去。通过随机读取待删除集合，可以让Redis不用维护一个巨大的链表，也不用操作链表，进而提升性能。

Redis 选出的数据个数 N，通过配置参数 `maxmemory-samples` 进行配置。个数N越大，则候选集合越大，选择到的最久未被使用的就更准确，N越小，选择到最久未被使用的数据的概率也会随之减小


### LFU算法

LFU 缓存策略是在 LRU 策略基础上，为每个数据增加了一个计数器，来统计这个数据的访问次数。

当使用 LFU 策略筛选淘汰数据时，首先会根据数据的访问次数进行筛选，把访问次数最低的数据淘汰出缓存。如果两个数据的访问次数相同，LFU 策略再比较这两个数据的访问时效性，把距离上一次访问时间更久的数据淘汰出缓存。 

Redis的LFU算法实现:当 LFU 策略筛选数据时，Redis 会在候选集合中，根据数据 lru 字段的后 8bit 选择访问次数最少的数据进行淘汰。当访问次数相同时，再根据 lru 字段的前 16bit 值大小，选择访问时间最久远的数据进行淘汰。

Redis 只使用了 8bit 记录数据的访问次数，而 8bit 记录的最大值是 255，这样在访问快速的情况下，如果每次被访问就将访问次数加一，很快某条数据就达到最大值255，可能很多数据都是255，那么退化成LRU算法了。

所以Redis为了解决这个问题，实现了一个更优的计数规则，并可以通过配置项，来控制计数器增加的速度。

- lfu-log-factor: 用计数器当前的值乘以配置项 lfu_log_factor 再加 1，再取其倒数，得到一个 p 值；然后，把这个 p 值和一个取值范围在（0，1）间的随机数 r 值比大小，只有 p 值大于 r 值时，计数器才加 1。
- lfu-decay-time: 控制访问次数衰减。LFU 策略会计算当前时间和数据最近一次访问时间的差值，并把这个差值换算成以分钟为单位。然后，LFU 策略再把这个差值除以 lfu_decay_time 值，所得的结果就是数据 counter 要衰减的值。
- lfu-log-factor: 设置越大，递增概率越低，lfu-decay-time设置越大，衰减速度会越慢。

在应用 LFU 策略时，一般可以将 lfu_log_factor 取值为 10。 如果业务应用中有短时高频访问的数据的话，建议把 lfu_decay_time 值设置为 1。可以快速衰减访问次数。

## 数据一致性

当涉及到数据更新时,容易出现缓存(Redis)和数据库（MySQL）间的数据一致性问题

### 示例分析

1. 如果删除了缓存Redis，还没有来得及写库MySQL，另一个线程就来读取，发现缓存为空，则去数据库中读取数据写入缓存，此时缓存中为脏数据。
2. 如果先写了库，在删除缓存前，写库的线程宕机了，没有删除掉缓存，则也会出现数据不一致情况。因为写和读是并发的，没法保证顺序,就会出现缓存和数据库的数据不一致的问题。

### 解决方案

- 策略1：先更新缓存，再更新数据库 如果只有缓存更新成功，会造成脏数据
- 策略2：先更新数据库，再更新缓存 多线程下，A更新数据库4,B更新数据库5, B更新缓存5, A更新缓存4
- 策略3：先删除缓存，再更新数据库 多线程下，A删除缓存,B读取缓存失败, B从数据库读取并更新缓存, A更新数据库（再删除一次，延迟双删）
- 策略4：先更新数据库，再删除缓存 A更新数据库完数据库后，再删除一次，延迟双删