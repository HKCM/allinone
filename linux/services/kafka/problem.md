# 如何防止消息丢失

- 发送方: ack是1或者-1/al可以防止消息丢失,如果要做到99.9999%,ack设成all,把replicas >= 2
- 消费方: 把自动提交改为手动提交.

# 如何防治重复消费

一条消息被消费者消费多次.如果为了消息的不重复消费,而把生产端的重试机制关闭、消费端的手动提交改成自动提交,这样反而会出现消息丢失.

可以直接在防治消息丢失的手段上再加上消费消息时的幂等性保证,就能解决消息的重复消费问题.

幂等性如何保证: 
- mysql插入业务id作为主键,主键是唯一的,所以一次只能插入一条
- 使用redis或zk的分布式锁（主流的方案）

# 如何做到顺序消费

- Topic使用单partition(会牺牲性能)
- 发送方: 确保消息是顺序发送的,在发送时将ack不能设置0,使用同步发送,等到发送成功再发送下一条.
- 接收方: 消息是发送到一个分区中,只能有一个消费组的消费者来接收消息.

因此,kafka的顺序消费会牺牲掉性能.

# 如何解决消息积压

消息积压会导致很多问题,比如磁盘被打满、生产端发消息导致kafka性能过慢,就容易出现服务雪崩,就需要有相应的手段:

- 方案一: 提升一个消费者的消费能力,在一个消费者中启动多个线程,让多个线程同时消费
- 方案二: 充分利用服务器的cpu资源,可以启动多个消费者,多个消费者部署在相同或不同的服务器上
- 方案三: 让一个消费者去把收到的消息往另外一个topic上发,另一个topic设置多个分区和多个消费者,进行具体的业务消费.

# 延迟队列

延迟队列的应用场景: 在订单创建成功后如果超过30分钟没有付款,则需要取消订单,此时可用延时队列来实现.创建多个topic,每个topic表示延时的间隔

- topic_5s: 延时5s执行的队列
- topic_1m: 延时1分钟执行的队列
- topic_30m: 延时30分钟执行的队列

1. 消息发送者发送消息到相应的topic,并带上消息的发送时间
2. 消费者订阅相应的topic,消费时轮询消费整个topic中的消息
3. 如果消息的发送时间,和消费的当前时间超过预设的值,则------
4. 如果消息的发送时间,和消费的当前时间没有超过预设的值,则------及之后的offset的所有消息都-----
5. 下次继续消费该offset处的消息,判断时间是否已满足预设值

# 死信队列和重试队列

死信可以看作消费者不能处理收到的消息，也可以看作消费者不想处理收到的消息，还可以看作不符合处理要求的消息。比如消息内包含的消息内容无法被消费者解析，为了确保消息的可靠性而不被随意丢弃，故将其投递到死信队列中，这里的死信就可以看作消费者不能处理的消息。再比如超过既定的重试次数之后将消息投入死信队列，这里就可以将死信看作不符合处理要求的消息。

重试队列其实可以看作一种回退队列，具体指消费端消费消息失败时，为了防止消息无故丢失而重新将消息回滚到 broker 中。与回退队列不同的是，重试队列一般分成多个重试等级，每个重试等级一般也会设置重新投递延时，重试次数越多投递延时就越大。

理解了他们的概念之后我们就可以为每个主题设置重试队列，消息第一次消费失败入重试队列 Q1，Q1 的重新投递延时为5s，5s过后重新投递该消息；如果消息再次消费失败则入重试队列 Q2，Q2 的重新投递延时为10s，10s过后再次投递该消息。

然后再设置一个主题作为死信队列，重试越多次重新投递的时间就越久，并且需要设置一个上限，超过投递次数就进入死信队列。重试队列与延时队列有相同的地方，都需要设置延时级别。