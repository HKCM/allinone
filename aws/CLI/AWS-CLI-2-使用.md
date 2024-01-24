

### 获取帮助

首先当然是获取帮助啦,获取帮助用`help`命令

```
$ aws help
$ aws ec2 help
$ aws ec2 describe-instances help
```

### wait命令

`wait` 命令是一个特殊的命令,很多service都有`wait`这个子命令,`wait`之后还要接一个子命令。它也是一个承上启下的命令,尤其是在编写AWS脚本很常用。举个例子,在使用创建EC2的CLI命令后,希望在EC2的`status`处于`ok`状态之后告诉你EC2已经正常运行起来了(EC2有很多状态,启动状态initial,运行状态running,停止状态stopping)

`wait` 命令是在command命令之后
```
$ aws ec2 wait instance-status-ok --instance-ids i-1234567890abcdef0
```

并不是说有服务都有`wait`子命令,通常有状态变化或者需要一段时间才能完成的服务才有`wait`子命令,例如EC2和CloudFormation。

`wait` 每15秒轮询一次,直到达到成功状态为止。40次失败检查后,将以255的返回码退出。

### CLI参数

CLI参数类型不多,但是也不算简单

1. 字符串
    最简单的参数,`--instance-id`后面的参数就是字符串
  ```
  $ aws ec2 describe-instance --instance-id i-1234567890abcdef0
  ```
  如果字符串参数包含空格,则字符串必须用引号引起来

2. 时间戳参数

  时间戳参数适用于类似`--start-time` 或 `--end-time`的参数选项,或者叫`date`参数.
  可接受格式包括
  * `YYYY-MM-DDThh:mm:ss.sssTZD (UTC)`,例如,2014-10-01T20:30:00.000Z
  * `YYYY-MM-DDThh:mm:ss.sssTZD(带偏移量`,例如,2014-10-01T12:30:00.000-08:00
  * `YYYY-MM-DD`,例如,2014-10-01
  * 以秒为单位的 Unix 时间,如 `1412195400`。表示自 1970 年 1 月 1 日午夜 (UTC) 以来经历的秒数。
  ```
  $ aws ec2 describe-spot-price-history --start-time 2014-10-13T19:00:00Z
  ```

3. 列表参数
    以空格分隔的一个或多个字符串。如果任何字符串项目包含空格,则必须用引号括起该项目
  ```
  $ aws ec2 describe-spot-price-history --instance-types m1.xlarge m1.medium
  ```

4. 整数参数
    同样很好理解
  ```
  $ aws ec2 describe-spot-price-history --max-items 5
  ```

5. 二进制参数
    简单来说二进制参数就是一个非文本和数字类型的参数,例如,图片,音频和视频
  ```
  $ aws s3api put-object --bucket my-bucket --key testimage.png --body /tmp/image.png
  ```

6. 映射参数
    映射参数通常比较复杂,像这样
  ```
  $ aws dynamodb get-item --table-name my-table --key '{"id": {"N":"1"}}'

  $ aws ec2 run-instances \
      --image-id ami-12345678 \
      --block-device-mappings '[{"DeviceName":"/dev/sdb","Ebs":{"VolumeSize":20,"DeleteOnTermination":false,"VolumeType":"standard"}}]'
  ```

7. 文件参数
    使用文件参数加载参数值有时候会很有用,例如
  ```
  $ aws sqs create-queue --queue-name my-queue --attributes file://attributes.json
  ```
  `attributes.json`文件内容为,在文件中
  ```
  {
    "RedrivePolicy":{
        "deadLetterTargetArn":"arn:aws-cn:sqs:us-west-2:0123456789012:deadletter",
        "maxReceiveCount":"5"
    }
  }
  ```

  文件参数还包括远程文件和二进制文件,具体看这:https://docs.amazonaws.cn/cli/latest/userguide/cli-usage-parameters-file.html

### 控制命令输出

前面说了CLI大致有四种输出`json`,`yaml`,`text`和`table`,并且可以通过配置文件`~/.aws/config`
```
[default]
output=text
```
环境变量
```
$ export AWS_DEFAULT_OUTPUT="table"
```
以及CLI`--output`选项进行配置,其中CLI的优先级最高
```
$ aws swf list-domains --registration-status REGISTERED --output json
```

使用`json`输出时,可以考虑使用 `jq`进行高级筛选,这是一个命令行 `JSON` 处理器。可以通过 http://mikefarah.github.io/jq/ 下载它并查找文档。

使用`yaml`输出时,可以考虑使用 `yq`进行高级筛选,这是一个命令行 `YAML` 处理器。可以通过 http://mikefarah.github.io/yq/ 下载它并查找文档。

使用`text`输出时,强烈建议使用`--query`,脚本中本人最常用的输出格式,因为不用额外安装 `jq` 或 `yq`

使用`table`输出时,也是和`--query`结合使用,最便于认为阅读,不适合脚本中使用

### 安全

AWS CLI 使用的凭证存储在纯文本文件中,并且`不加密`。

* $HOME/.aws/credentials 文件存储访问 AWS 资源所需的长期凭证。这包括访问密钥 ID 和秘密访问密钥。
* 短期凭证(例如承担的角色或用于 AWS Single Sign-On 服务的角色的凭证)也分别存储在 $HOME/.aws/cli/cache 和 $HOME/.aws/sso/cache 文件夹中。

强烈建议您在 $HOME/.aws 文件夹及其子文件夹和文件上配置文件系统权限,仅限授权用户访问。  
尽可能使用具有临时凭证的角色,以减少凭证泄露时造成损坏的机会。仅使用长期凭证来请求和刷新短期角色凭证。

### 常见问题

1. 系统时间与实际时间不匹配导致验证不通过,需要调整系统时间

2. ~/.aws/credential 文件中格式不正确或者包含错误的字符串.需要生成格式正确的`credential`文件

