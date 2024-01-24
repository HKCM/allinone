

### Spot 实例请求

使用 [run-instances](https://docs.aws.amazon.com/cli/latest/reference/ec2/run-instances.html) 命令并在 `--instance-market-options` 参数中指定 Spot 实例选项。

```shell
aws ec2 run-instances \
--profile hkc_admin \
--image-id ami-00e87074e52e6c9f9 \
--instance-type t2.medium \
--count 1 \
--subnet-id subnet-026fa303b2e5f164d \
--key-name myselfkey \
--security-group-ids sg-0a3761914ac096713 \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":20,"DeleteOnTermination":true,"VolumeType":"gp2"}}]' \
--instance-market-options 'MarketType=spot,SpotOptions={MaxPrice=0.014,SpotInstanceType=one-time}' \
--private-ip-address 10.4.7.51
```

以下是要在 JSON 文件中为 `--instance-market-options` 指定的数据结构。您还可以指定 `ValidUntil` 和 `InstanceInterruptionBehavior`。如果未在数据结构中指定字段,则将使用默认值。此示例创建一个 `one-time` 请求,并指定 `0.014` 作为您愿意为 Spot 实例支付的最高价。

```json
{
  "MarketType": "spot",
  "SpotOptions": {
    "MaxPrice": "0.014",
    "SpotInstanceType": "one-time"
  }
}
```


创建多个具有不同私有地址的EC2
```shell
#!/bin/bash
for i in 52 53 54 55 59
do
aws ec2 run-instances \
--profile hkc_admin \
--image-id ami-087c17d1fe0178315 \
--instance-type t2.medium \
--count 1 \
--subnet-id subnet-026fa303b2e5f164d \
--key-name myselfkey \
--security-group-ids sg-0a3761914ac096713 \
--block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":20,"DeleteOnTermination":true,"VolumeType":"gp2"}}]' \
--tag-specifications 'ResourceType=instance,Tags=[{Key=webserver,Value=production}]' 'ResourceType=volume,Tags=[{Key=cost-center,Value=cc123}]' \
--instance-market-options 'MarketType=spot,SpotOptions={MaxPrice=0.014,SpotInstanceType=one-time}' \
--private-ip-address 10.4.7.$i
done
```

