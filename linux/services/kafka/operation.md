# Installation

https://kafka.apache.org/downloads

```bash
apt update && apt install openjdk-17-jdk curl vim -y

mkdir -p /opt/kafka/data && cd /opt/kafka

curl -O https://downloads.apache.org/kafka/3.5.2/kafka_2.13-3.5.2.tgz

tar -xzvf kafka_2.13-3.5.2.tgz

cd kafka_2.13-3.5.2
# 修改log目录
vim config/kraft/server.properties
node.id
log.dirs=/opt/kafka/data
advertised.listeners=PLAINTEXT://192.168.151.213:9092


# 生成集群uuid
uuid=$(sh bin/kafka-storage.sh random-uuid) # T-Fmq8TVTaqvZwmeWt0s5Q

# data目录下会出现bootstrap.checkpoint 和 meta.properties
bin/kafka-storage.sh format -t ${uuid} -c config/kraft/server.properties

cat ../data/meta.properties 
#Sun Jan 21 06:02:15 GMT 2024
node.id=1
version=1
cluster.id=KnpaXeYaT6qSqbpXGgWQjQ

# 启动
bin/kafka-server-start.sh -daemon config/kraft/server.properties

# 停止
bin/kafka-server-stop.sh
```

# 操作

```bash
# 创建Topic
bin/kafka-topics.sh --create --replication-factor 2 --bootstrap-server localhost:9092 --partitions 3 --topic test2

# 列出Topic
bin/kafka-topics.sh --list --bootstrap-server localhost:9092

# 描述Topic
bin/kafka-topics.sh --describe --bootstrap-server localhost:9092 --topic test1

# 修改Topic保留时间
/bin/kafka-configs.sh --alter \
      --add-config retention.ms=300000 \
      --bootstrap-server=0.0.0.0:9092 \
      --topic ${topic_name}

# 查看topic消费到的offset
bin/kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic test1 --time -1

# 删除Topic
bin/kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic test1
```

# producer

```bash
bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test1

bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic test1 --property parse.key=true
id123   messages # 使用tab分隔
id124   another_massages

# 持续发送消息
bin/kafka-verifiable-producer.sh  --broker-list localhost:9092 --topic test1 --max-messages 64
bin/kafka-verifiable-producer.sh  --broker-list localhost:9092 --topic test1 --max-messages 64000 --throughput 1
```

# consumer

```bash
bin/kafka-topics.sh --bootstrap-server localhost:9092 --list
# 从头开始消费
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test1 --from-beginning

# 描述Topic
bin/kafka-topics.sh --describe --bootstrap-server localhost:9092 --topic test1
# 消费指定分区
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test1 --partition 0 --from-beginning

bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test1 --partition 0 --offset 20
```