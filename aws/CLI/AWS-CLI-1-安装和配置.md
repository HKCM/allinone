

# AWS CLI

[toc]

AWS CLI有两个版本,V1和V2. V2中添加很多修改,具体看这里:
 https://docs.amazonaws.cn/cli/latest/userguide/cliv2-migration.html

**默认情况下,AWS CLI 版本 2 会对所有输出使用分页程序,在Linux 和 macOS 上是 `less` 程序**

如果不想使用分页程序或者换成其他分页程序,可以配置 `AWS_PAGER` 环境变量或者在`~/.aws/config`文件中添加`cli_pager`。设置为空则表示不使用分页。
```
$ export AWS_PAGER=""

$ aws configure set cli_pager ''

$ cat ~/.aws/config
[default]
cli_pager=
```

所以直接安装V2版本,推荐用官方的方式安装,由于 AWS 不维护第三方存储库,因此不能保证它们包含最新版本的 AWS CLI。

## CLI 安装更新和删除
### Linux安装
```linux X86

# 下载最新版本
$ curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

# 下载指定版本,例如版本 2.0.30
# $ curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.30.zip" -o "awscliv2.zip"

# 解压
$ unzip awscliv2.zip

# 安装至指定路径
# -i - 此选项指定要将所有文件复制到的目录,程序实际所在路径,默认值为 /usr/local/aws-cli
# -b - 此选项指定安装目录中的主 aws 程序通过符号链接指向指定路径中的 aws 文件,默认值为 /usr/local/bin,所以无需再将安装目录添加到用户的 $PATH 变量中
$ sudo ./aws/install -i /usr/local/aws-cli -b /usr/local/bin


# 验证安装
$ aws --version
aws-cli/2.0.47 Python/3.7.4 Linux/4.14.133-113.105.amzn2.x86_64 botocore/2.0.0
```

版本列表: https://github.com/aws/aws-cli/blob/v2/CHANGELOG.rst



### 命令补全

https://docs.amazonaws.cn/cli/latest/userguide/cli-configure-completion.html#cli-command-completion-linux

以zsh为例
```bash
# 找到aws_completer
which aws_completer

# 在~/.zshrc文件中添加以下三行
vim ~/.zshrc
autoload bashcompinit && bashcompinit
autoload -Uz compinit && compinit
complete -C '/usr/local/bin/aws_completer' aws # 这里注意替换实际的aws_completer的位置

source ~/.zshrc
```

### 更新

更新其实就是用新文件覆盖旧文件,和安装基本一致

1. 下载最新版本

   ```
   $ curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
   ```

2. 解压

   ```
   $ unzip awscliv2.zip
   ```

3. 使用 which 命令查找符号链接

   ```
   $ which aws
   /usr/local/bin/aws
   ```

4. 使用 ls 命令查找符号链接指向的目录

   ```
   $ ls -l /usr/local/bin/aws
   lrwxrwxrwx 1 ec2-user ec2-user 49 Oct 22 09:49 /usr/local/bin/aws -> /usr/local/aws-cli/v2/current/bin/aws
   ```

5. 使用带有 --update 参数的 install 命令。

   ```
   # --bin-dir - which命令返回的目录
   # --install-dir - ls显示的链接的实际路径
   $ sudo ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli --update
   ```

   

### 卸载

卸载其实就是删除相关文件和链接

1. 首先使用 which 命令查找符号链接

   ```
   $ which aws
   /usr/local/bin/aws
   ```

   

2. 使用 ls 命令查找符号链接指向的目录

   ```
   $ ls -l /usr/local/bin/aws
   lrwxrwxrwx 1 ec2-user ec2-user 49 Oct 22 09:49 /usr/local/bin/aws -> /usr/local/aws-cli/v2/current/bin/aws
   ```

3. 最后删除链接和目录

   ```
   $ sudo rm /usr/local/bin/aws
   $ sudo rm /usr/local/bin/aws_completer
   $ sudo rm -rf /usr/local/aws-cli
   ```

其他系统安装方式可以参考官方文档,官方文档: https://docs.amazonaws.cn/cli/latest/userguide/install-cliv2-linux.html

## 配置CLI

### 基础配置

1. 默认配置命令 `aws configure` 
  ```
  # 创建默认名称为default的配置文件
  $ aws configure
  AWS Access Key ID [None]: AKIAIOSFODNN7EXAMPLE
  AWS Secret Access Key [None]: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  Default region name [None]: us-west-2
  Default output format [None]: json

  # 创建名为 produser的配置文件
  $ aws configure --profile produser
  ```

  * AWS Access Key ID:源于IAM创建的访问密钥
  * AWS Secret Access Key:源于IAM创建的访问密钥
  * Region:运行CLI命令时,命令发送的区域,取决于服务部署在哪个Region
  * Output:输出格式

  CLI 为使用 `aws configure` 配置的默认配置文件生成的文件看起来类似于以下内容。

  `~/.aws/credentials`
  ```
  [default]
  aws_access_key_id=AKIAIOSFODNN7EXAMPLE
  aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  ```
  `~/.aws/config`
  ```
  [default]
  region=us-east-1
  output=json
  ```

2. 使用`configure`的 `set`, `get` 和 `list`, `list-profiles` 子命令进行单独配置,对于`assume role`操作很有用.

    使用set配置
  ```
$ aws configure set aws_access_key_id AKIAIOSFODNN7EXAMPLE --profile example
$ aws configure set aws_secret_access_key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY --profile example
$ aws configure set region us-east-1 --profile example
$ aws configure set aws_access_key_id AKIAIOSFODNN7EXAMPLE --profile example
  ```

使用get获取  

```
$ aws configure get aws_access_key_id --profile example
AKIAIOSFODNN7EXAMPLE
```

使用list

```
$ aws configure get aws_access_key_id --profile example
        Name                    Value             Type    Location
        ----                    -----             ----    --------
    profile                  example           manual    --profile
  access_key     ****************MPLE shared-credentials-file    
  secret_key     ****************EKEY shared-credentials-file    
      region                us-east-1      config-file    ~/.aws/config
```

查看当前的profile

```
$ aws configure list-profiles
default
example
```

`assume role`举例

  ```
  KST=($(aws sts assume-role \
  --role-arn "arn:aws:iam::${AWSAccount}:role/${RoleName}" \
  --role-session-name "role-session" \
  --query 'Credentials.[AccessKeyId,SecretAccessKey,SessionToken]' \
  --profile 'bot-aws-profile' \
  --output text))

  echo "Setting AWS configs:"
  aws configure set aws_access_key_id "${KST[0]}" --profile role-profile
  aws configure set aws_secret_access_key "${KST[1]}" --profile role-profile
  aws configure set aws_session_token "${KST[2]}" --profile role-profile
  aws configure set region "${REGION}" --profile role-profile
  ```

3. `import`

  这是最方便的方式,当需要配置本机CLI时,大概率你手里已经有了credentials.csv文件
  必须包含:
  - User Name
  - Access key ID
  - Secret access key

  ```
  $ aws configure import --csv file://./credentials.csv
  ```


### 配置设置

config文件支持大量的字段,下面介绍一些常用的:

* cli_pager
  分页功能,仅适用于 AWS CLI 版本 2,可以被 AWS_PAGER 环境变量覆盖
  使用 --no-cli-pager 命令行选项可禁止在单个命令中使用分页程序。
  将 cli_pager 设置或 AWS_PAGER 变量设置为空字符串。

  export AWS_PAGER=""
  ```shell
  [default]
  cli_pager=less
  # or
  [default]
  cli_pager=
  ```
* duration_seconds
  指定角色会话的最大持续时间(以秒为单位),默认为3600秒
* max_attempts
  指定 AWS CLI 重试处理程序使用的最大重试次数值,其中初始调用计入您提供的 max_attempts 值。可以使用 AWS_MAX_ATTEMPTS 环境变量覆盖此值。
* output 
  CLI命令的输出格式
  * json - 输出采用 JSON 字符串的格式。
  * yaml - 输出采用 YAML 字符串的格式。(仅在 AWS CLI 版本 2 中可用。)
  * yaml-stream - 输出被流式处理并采用 YAML 字符串的格式。流式处理支持更快地处理大型数据类型。(仅在 AWS CLI 版本 2 中可用。)
  * text - 输出采用多个制表符分隔字符串值行的格式。这对于将输出传递到文本处理器(如 grep、sed 或 awk)很有用。
  * table - 输出采用表格形式,使用字符 +|- 以形成单元格边框。它通常以“人性化”格式呈现信息,这种格式比其他格式更容易阅读,但从编程方面来讲不是那么有用。
* region
  对于使用该配置文件请求的命令,指定要将请求发送到的 AWS 区域
* role_arn
  指定要用于运行 AWS CLI 命令的 IAM 角色的 Amazon 资源名称 (ARN)。此外,还必须指定以下参数之一以标识有权代入此角色的凭证:
  * source_profile
  * credential_source
* source_profile
  指定包含长期凭证的命名配置文件,AWS CLI 可使用这些凭证代入通过 role_arn 参数指定的角色。不能在同一配置文件中同时指定 source_profile 和 credential_source。
* credential_source
  在 Amazon EC2 实例或 EC2 容器中使用,指定 AWS CLI 在何处可以找到要用于代入通过 role_arn 参数指定的角色的凭证。不能在同一配置文件中同时指定 source_profile 和 credential_source。  
  此参数具有三个值:
  * Environment - 指定 AWS CLI 从环境变量检索源凭证。
  * Ec2InstanceMetadata - 指定 AWS CLI 将使用附加到 EC2 实例配置文件的 IAM 角色以获取源凭证。
  * EcsContainer - 指定 AWS CLI 将附加到 ECS 容器的 IAM 角色用作源凭证。
* retry_mode 
  重试操作,有三种类型传统重试模式(legacy),标准重试模式(standard)和自适应重试模式(adaptive)
  传统模式:AWS CLI 版本 1 使用的默认模式,重试4次,总共可发出 5 次调用尝试。
  标准模式:AWS CLI 版本 2 使用的默认模式,最大重试次数的默认值为 2,总共可发出 3 次调用尝试。多了一些报错代码。
  自适应模式:自适应模式是一种试验模式,AWSCLI会自己根据情况改变重试次数
* max_attempts
  最大重试次数

### config中的s3设置

```
[profile development]
s3 =
  max_concurrent_requests = 20
  max_queue_size = 10000
  multipart_threshold = 64MB
  multipart_chunksize = 16MB
  max_bandwidth = 50MB/s
  use_accelerate_endpoint = true
  addressing_style = path
```

具体信息:https://docs.amazonaws.cn/cli/latest/userguide/cli-configure-files.html#cli-configure-files-where


### 优先级

AWS CLI 凭证和配置设置的优先顺序如下:

1. 命令行选项 - 覆盖任何其他位置的设置。可以在命令行上指定 --region、--output 和 --profile 作为参数。
2. 环境变量 - 您可以在系统的环境变量中存储值。 环境变量参考:https://docs.amazonaws.cn/cli/latest/userguide/cli-configure-envvars.html
3. CLI 凭证文件 - credentials 文件。credentials 文件位于 ~/.aws/credentials(在 Linux 或 macOS 上)或 C:\Users\USERNAME\.aws\credentials(在 Windows 上)。该文件可以包含 default 配置文件和任何命名配置文件的`凭证详细信息`。
4. CLI 配置文件 - config 文件。config 文件位于 ~/.aws/config(在 Linux 或 macOS 上)或 C:\Users\USERNAME\.aws\config(在 Windows 上)。该文件包含默认配置文件和任何命名配置文件的`配置设置`。
5. 容器凭证 - ECS的IAM role
6. 实例配置文件凭证 - EC2的IAM role

### CLI常用命令行选项

1. **--profile**

  最常用的选项没有之一,当你具有多个配置文件时,用于让AWS命令运行在不同配置下
  ```
  $ aws --profile dev servicename commandname options 
  ```
2. **--region**
    一般般,当你的资源在同一个AWS账号下但是不在同一个region时,在不同region之间控制资源需要用到。

3. **--output <json|yaml|text|table>**
    比较常用,尤其是编程时与`--query`结合

3. **--query <string>**
    比较常用,尤其是编程时与`--output`结合
  
4. **--debug**
    不算常用,但是在CLI出现意料之外的情况时很有用。通常可以将错误输出到错误日志。
  ```
  $ aws servicename commandname options --debug 2> debug.txt
  ```

5. **--generate-cli-skeleton**
    不常用,至少我自己没用过。它可以帮你生成CLI命令的模版,模版中包括该命令的所有参数,你可以修改这个模版,删除你不需要的参数,然后填上你需要的参数。  
    最后运行CLI命令时使用`--cli-input-yaml`选项,就可以自动将你的参数带入命令中。
  ```
  $ aws ec2 run-instances --generate-cli-skeleton yaml-input > ec2runinst.yaml
  $ vim ec2runinst.yaml
  DryRun: true
  ImageId: 'ami-dfc39aef'
  KeyName: 'mykey'
  SecurityGroups:
  - 'my-sg'
  InstanceType: 't2.micro'
  Monitoring: 
    Enabled: true
  $ aws ec2 run-instances --cli-input-yaml file://ec2runinst.yaml
  ```

6. **--no-cli-pager**
    命令输出禁用分页程序

  https://docs.amazonaws.cn/cli/latest/userguide/cli-configure-options.html

### 设置AWS自动补全

  通常AWS安装后会有自动补全功能如果没有,首先找到AWS的命令补全程序的位置(通常和`aws`在同一目录下)
  ```
  $ which aws_completer
  /usr/local/bin/aws_completer
  ```
  `/usr/local/bin`已经在我的$PATH,如果路径不在`$PATH`需要在`~/.bash_profile`文件末尾中添加这一行,并应用生效`source ~/.bash_profile`

  ```
  export PATH=/other/path/folder:$PATH
  ```
  最后`bash`使用内置命令 `complete`。
  ``` 
  $ complete -C '/usr/local/aws/bin/aws_completer' aws
  ```
  https://docs.amazonaws.cn/cli/latest/userguide/cli-configure-completion.html

## CLI 使用



### 获取帮助

首先当然是获取帮助啦,获取帮助用`help`命令

```
$ aws help
$ aws ec2 help
$ aws ec2 describe-instances help
```

### 使用 --debug 选项

当 AWS CLI 报告一个你不能立即理解的错误,或者产生你不期望的结果时,首先要做的是获得关于该错误的更多详细信息。通过再次运行该命令并在命令行末尾包含 --debug 选项,即可完成此操作。这会使 AWS CLI 报告有关完成以下过程所需的每个步骤的详细信息:处理命令、将请求发送到 AWS 服务器、接收响应并将响应处理为你所看到的输出。输出中的详细信息可帮助你确定发生错误的步骤,以及获得有关是什么触发错误的线索的上下文。  

可以将输出发送到一个要捕获它的文本文件以供日后查看,或者按要求将输出发送给 AWS 技术支持。

不使用 `--debug` 选项运行命令的示例

```
aws iam list-groups --profile MyTestProfile
{
    "Groups": [
        {
            "Path": "/",
            "GroupName": "MyTestGroup",
            "GroupId": "AGPA0123456789EXAMPLE",
            "Arn": "arn:aws:iam::123456789012:group/MyTestGroup",
            "CreateDate": "2019-08-12T19:34:04Z"
        }
    ]
}
```

使用 `--debug` 选项运行命令的示例

当包含 --debug 选项时,详细信息包括(或其他):

* 查找凭证
* 解析提供的参数
* 构建发送到 AWS 服务器的请求
* 发送到 AWS 的请求的内容
* 原始响应的内容
* 带格式的输出

```shell
aws iam list-groups --profile MyTestProfile --debug
2019-08-12 12:36:18,305 - MainThread - awscli.clidriver - DEBUG - CLI version: aws-cli/1.16.215 Python/3.7.3 Linux/4.14.133-113.105.amzn2.x86_64 botocore/1.12.205
2019-08-12 12:36:18,305 - MainThread - awscli.clidriver - DEBUG - Arguments entered to CLI: ['iam', 'list-groups', '--debug']
2019-08-12 12:36:18,305 - MainThread - botocore.hooks - DEBUG - Event session-initialized: calling handler <function add_scalar_parsers at 0x7fdf173161e0>
2019-08-12 12:36:18,305 - MainThread - botocore.hooks - DEBUG - Event session-initialized: calling handler <function register_uri_param_handler at 0x7fdf17dec400>
2019-08-12 12:36:18,305 - MainThread - botocore.hooks - DEBUG - Event session-initialized: calling handler <function inject_assume_role_provider_cache at 0x7fdf17da9378>
2019-08-12 12:36:18,307 - MainThread - botocore.credentials - DEBUG - Skipping environment variable credential check because profile name was explicitly set.
2019-08-12 12:36:18,307 - MainThread - botocore.hooks - DEBUG - Event session-initialized: calling handler <function attach_history_handler at 0x7fdf173ed9d8>
2019-08-12 12:36:18,308 - MainThread - botocore.loaders - DEBUG - Loading JSON file: /home/ec2-user/venv/lib/python3.7/site-packages/botocore/data/iam/2010-05-08/service-2.json
2019-08-12 12:36:18,317 - MainThread - botocore.hooks - DEBUG - Event building-command-table.iam: calling handler <function add_waiters at 0x7fdf1731a840>
2019-08-12 12:36:18,320 - MainThread - botocore.loaders - DEBUG - Loading JSON file: /home/ec2-user/venv/lib/python3.7/site-packages/botocore/data/iam/2010-05-08/waiters-2.json
2019-08-12 12:36:18,321 - MainThread - awscli.clidriver - DEBUG - OrderedDict([('path-prefix', <awscli.arguments.CLIArgument object at 0x7fdf171ac780>), ('marker', <awscli.arguments.CLIArgument object at 0x7fdf171b09e8>), ('max-items', <awscli.arguments.CLIArgument object at 0x7fdf171b09b0>)])
2019-08-12 12:36:18,322 - MainThread - botocore.hooks - DEBUG - Event building-argument-table.iam.list-groups: calling handler <function add_streaming_output_arg at 0x7fdf17316510>
2019-08-12 12:36:18,322 - MainThread - botocore.hooks - DEBUG - Event building-argument-table.iam.list-groups: calling handler <function add_cli_input_json at 0x7fdf17da9d90>
2019-08-12 12:36:18,322 - MainThread - botocore.hooks - DEBUG - Event building-argument-table.iam.list-groups: calling handler <function unify_paging_params at 0x7fdf17328048>
2019-08-12 12:36:18,326 - MainThread - botocore.loaders - DEBUG - Loading JSON file: /home/ec2-user/venv/lib/python3.7/site-packages/botocore/data/iam/2010-05-08/paginators-1.json
2019-08-12 12:36:18,326 - MainThread - awscli.customizations.paginate - DEBUG - Modifying paging parameters for operation: ListGroups
2019-08-12 12:36:18,326 - MainThread - botocore.hooks - DEBUG - Event building-argument-table.iam.list-groups: calling handler <function add_generate_skeleton at 0x7fdf1737eae8>
2019-08-12 12:36:18,326 - MainThread - botocore.hooks - DEBUG - Event before-building-argument-table-parser.iam.list-groups: calling handler <bound method OverrideRequiredArgsArgument.override_required_args of <awscli.customizations.cliinputjson.CliInputJSONArgument object at 0x7fdf171b0a58>>
2019-08-12 12:36:18,327 - MainThread - botocore.hooks - DEBUG - Event before-building-argument-table-parser.iam.list-groups: calling handler <bound method GenerateCliSkeletonArgument.override_required_args of <awscli.customizations.generatecliskeleton.GenerateCliSkeletonArgument object at 0x7fdf171c5978>>
2019-08-12 12:36:18,327 - MainThread - botocore.hooks - DEBUG - Event operation-args-parsed.iam.list-groups: calling handler functools.partial(<function check_should_enable_pagination at 0x7fdf17328158>, ['marker', 'max-items'], {'max-items': <awscli.arguments.CLIArgument object at 0x7fdf171b09b0>}, OrderedDict([('path-prefix', <awscli.arguments.CLIArgument object at 0x7fdf171ac780>), ('marker', <awscli.arguments.CLIArgument object at 0x7fdf171b09e8>), ('max-items', <awscli.customizations.paginate.PageArgument object at 0x7fdf171c58d0>), ('cli-input-json', <awscli.customizations.cliinputjson.CliInputJSONArgument object at 0x7fdf171b0a58>), ('starting-token', <awscli.customizations.paginate.PageArgument object at 0x7fdf171b0a20>), ('page-size', <awscli.customizations.paginate.PageArgument object at 0x7fdf171c5828>), ('generate-cli-skeleton', <awscli.customizations.generatecliskeleton.GenerateCliSkeletonArgument object at 0x7fdf171c5978>)]))
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.path-prefix: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.marker: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.max-items: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.cli-input-json: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.starting-token: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.page-size: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,328 - MainThread - botocore.hooks - DEBUG - Event load-cli-arg.iam.list-groups.generate-cli-skeleton: calling handler <awscli.paramfile.URIArgumentHandler object at 0x7fdf1725c978>
2019-08-12 12:36:18,329 - MainThread - botocore.hooks - DEBUG - Event calling-command.iam.list-groups: calling handler <bound method CliInputJSONArgument.add_to_call_parameters of <awscli.customizations.cliinputjson.CliInputJSONArgument object at 0x7fdf171b0a58>>
2019-08-12 12:36:18,329 - MainThread - botocore.hooks - DEBUG - Event calling-command.iam.list-groups: calling handler <bound method GenerateCliSkeletonArgument.generate_json_skeleton of <awscli.customizations.generatecliskeleton.GenerateCliSkeletonArgument object at 0x7fdf171c5978>>
2019-08-12 12:36:18,329 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: assume-role
2019-08-12 12:36:18,329 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: assume-role-with-web-identity
2019-08-12 12:36:18,329 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: shared-credentials-file
2019-08-12 12:36:18,329 - MainThread - botocore.credentials - INFO - Found credentials in shared credentials file: ~/.aws/credentials
2019-08-12 12:36:18,330 - MainThread - botocore.loaders - DEBUG - Loading JSON file: /home/ec2-user/venv/lib/python3.7/site-packages/botocore/data/endpoints.json
2019-08-12 12:36:18,334 - MainThread - botocore.hooks - DEBUG - Event choose-service-name: calling handler <function handle_service_name_alias at 0x7fdf1898eb70>
2019-08-12 12:36:18,337 - MainThread - botocore.hooks - DEBUG - Event creating-client-class.iam: calling handler <function add_generate_presigned_url at 0x7fdf18a028c8>
2019-08-12 12:36:18,337 - MainThread - botocore.regions - DEBUG - Using partition endpoint for iam, us-west-2: aws-global
2019-08-12 12:36:18,337 - MainThread - botocore.args - DEBUG - The s3 config key is not a dictionary type, ignoring its value of: None
2019-08-12 12:36:18,340 - MainThread - botocore.endpoint - DEBUG - Setting iam timeout as (60, 60)
2019-08-12 12:36:18,341 - MainThread - botocore.loaders - DEBUG - Loading JSON file: /home/ec2-user/venv/lib/python3.7/site-packages/botocore/data/_retry.json
2019-08-12 12:36:18,341 - MainThread - botocore.client - DEBUG - Registering retry handlers for service: iam
2019-08-12 12:36:18,342 - MainThread - botocore.hooks - DEBUG - Event before-parameter-build.iam.ListGroups: calling handler <function generate_idempotent_uuid at 0x7fdf189b10d0>
2019-08-12 12:36:18,342 - MainThread - botocore.hooks - DEBUG - Event before-call.iam.ListGroups: calling handler <function inject_api_version_header_if_needed at 0x7fdf189b2a60>
2019-08-12 12:36:18,343 - MainThread - botocore.endpoint - DEBUG - Making request for OperationModel(name=ListGroups) with params: {'url_path': '/', 'query_string': '', 'method': 'POST', 'headers': {'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8', 'User-Agent': 'aws-cli/1.16.215 Python/3.7.3 Linux/4.14.133-113.105.amzn2.x86_64 botocore/1.12.205'}, 'body': {'Action': 'ListGroups', 'Version': '2010-05-08'}, 'url': 'https://iam.amazonaws.com/', 'context': {'client_region': 'aws-global', 'client_config': <botocore.config.Config object at 0x7fdf16e9a4a8>, 'has_streaming_input': False, 'auth_type': None}}
2019-08-12 12:36:18,343 - MainThread - botocore.hooks - DEBUG - Event request-created.iam.ListGroups: calling handler <bound method RequestSigner.handler of <botocore.signers.RequestSigner object at 0x7fdf16e9a470>>
2019-08-12 12:36:18,343 - MainThread - botocore.hooks - DEBUG - Event choose-signer.iam.ListGroups: calling handler <function set_operation_specific_signer at 0x7fdf18996f28>
2019-08-12 12:36:18,343 - MainThread - botocore.auth - DEBUG - Calculating signature using v4 auth.
2019-08-12 12:36:18,343 - MainThread - botocore.auth - DEBUG - CanonicalRequest:
POST
/

content-type:application/x-www-form-urlencoded; charset=utf-8
host:iam.amazonaws.com
x-amz-date:20190812T193618Z

content-type;host;x-amz-date
5f776d91EXAMPLE9b8cb5eb5d6d4a787a33ae41c8cd6eEXAMPLEca69080e1e1f
2019-08-12 12:36:18,344 - MainThread - botocore.auth - DEBUG - StringToSign:
AWS4-HMAC-SHA256
20190812T193618Z
20190812/us-east-1/iam/aws4_request
ab7e367eEXAMPLE2769f178ea509978cf8bfa054874b3EXAMPLE8d043fab6cc9
2019-08-12 12:36:18,344 - MainThread - botocore.auth - DEBUG - Signature:
d85a0EXAMPLEb40164f2f539cdc76d4f294fe822EXAMPLE18ad1ddf58a1a3ce7
2019-08-12 12:36:18,344 - MainThread - botocore.endpoint - DEBUG - Sending http request: <AWSPreparedRequest stream_output=False, method=POST, url=https://iam.amazonaws.com/, headers={'Content-Type': b'application/x-www-form-urlencoded; charset=utf-8', 'User-Agent': b'aws-cli/1.16.215 Python/3.7.3 Linux/4.14.133-113.105.amzn2.x86_64 botocore/1.12.205', 'X-Amz-Date': b'20190812T193618Z', 'Authorization': b'AWS4-HMAC-SHA256 Credential=AKIA01234567890EXAMPLE-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=d85a07692aceb401EXAMPLEa1b18ad1ddf58a1a3ce7EXAMPLE', 'Content-Length': '36'}>
2019-08-12 12:36:18,344 - MainThread - urllib3.util.retry - DEBUG - Converted retries value: False -> Retry(total=False, connect=None, read=None, redirect=0, status=None)
2019-08-12 12:36:18,344 - MainThread - urllib3.connectionpool - DEBUG - Starting new HTTPS connection (1): iam.amazonaws.com:443
2019-08-12 12:36:18,664 - MainThread - urllib3.connectionpool - DEBUG - https://iam.amazonaws.com:443 "POST / HTTP/1.1" 200 570
2019-08-12 12:36:18,664 - MainThread - botocore.parsers - DEBUG - Response headers: {'x-amzn-RequestId': '74c11606-bd38-11e9-9c82-559da0adb349', 'Content-Type': 'text/xml', 'Content-Length': '570', 'Date': 'Mon, 12 Aug 2019 19:36:18 GMT'}
2019-08-12 12:36:18,664 - MainThread - botocore.parsers - DEBUG - Response body:
b'<ListGroupsResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">\n  <ListGroupsResult>\n    <IsTruncated>false</IsTruncated>\n    <Groups>\n      <member>\n        <Path>/</Path>\n        <GroupName>MyTestGroup</GroupName>\n        <Arn>arn:aws:iam::123456789012:group/MyTestGroup</Arn>\n        <GroupId>AGPA1234567890EXAMPLE</GroupId>\n        <CreateDate>2019-08-12T19:34:04Z</CreateDate>\n      </member>\n    </Groups>\n  </ListGroupsResult>\n  <ResponseMetadata>\n    <RequestId>74c11606-bd38-11e9-9c82-559da0adb349</RequestId>\n  </ResponseMetadata>\n</ListGroupsResponse>\n'
2019-08-12 12:36:18,665 - MainThread - botocore.hooks - DEBUG - Event needs-retry.iam.ListGroups: calling handler <botocore.retryhandler.RetryHandler object at 0x7fdf16e9a780>
2019-08-12 12:36:18,665 - MainThread - botocore.retryhandler - DEBUG - No retry needed.
2019-08-12 12:36:18,665 - MainThread - botocore.hooks - DEBUG - Event after-call.iam.ListGroups: calling handler <function json_decode_policies at 0x7fdf189b1d90>
{
    "Groups": [
        {
            "Path": "/",
            "GroupName": "MyTestGroup",
            "GroupId": "AGPA123456789012EXAMPLE",
            "Arn": "arn:aws:iam::123456789012:group/MyTestGroup",
            "CreateDate": "2019-08-12T19:34:04Z"
        }
    ]
}
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
{  "RedrivePolicy":{      "deadLetterTargetArn":"arn:aws-cn:sqs:us-west-2:0123456789012:deadletter",      "maxReceiveCount":"5"  }}
  ```

  文件参数还包括远程文件和二进制文件,具体看这:https://docs.amazonaws.cn/cli/latest/userguide/cli-usage-parameters-file.html

### 控制命令输出

前面说了CLI大致有四种输出`json`,`yaml`,`text`和`table`,并且可以通过配置文件`~/.aws/config`

```
[default]output=text
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

## 示例

演示一些常见命令

### DynamoDB

#### 创建table

```
$ aws dynamodb create-table \
    --table-name MusicCollection \
    --attribute-definitions AttributeName=Artist,AttributeType=S AttributeName=SongTitle,AttributeType=S \
    --key-schema AttributeName=Artist,KeyType=HASH AttributeName=SongTitle,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
```

#### 插入项目

```
$ aws dynamodb put-item \
    --table-name MusicCollection \
    --item '{
        "Artist": {"S": "No One You Know"},
        "SongTitle": {"S": "Call Me Today"} ,
        "AlbumTitle": {"S": "Somewhat Famous"} 
      }' \
    --return-consumed-capacity TOTAL
{
    "ConsumedCapacity": {
        "CapacityUnits": 1.0,
        "TableName": "MusicCollection"
    }
}

$ aws dynamodb put-item \
    --table-name MusicCollection \
    --item '{ 
        "Artist": {"S": "Acme Band"}, 
        "SongTitle": {"S": "Happy Day"} , 
        "AlbumTitle": {"S": "Songs About Life"} 
      }' \
    --return-consumed-capacity TOTAL

{
    "ConsumedCapacity": {
        "CapacityUnits": 1.0,
        "TableName": "MusicCollection"
    }
}
```

#### 使用json文件

json文件`expression-attributes.json`的内容

```
{
  ":v1": {"S": "No One You Know"},
  ":v2": {"S": "Call Me Today"}
}
```

```
$ aws dynamodb query --table-name MusicCollection \
    --key-condition-expression "Artist = :v1 AND SongTitle = :v2" \
    --expression-attribute-values file://expression-attributes.json
{
    "Count": 1,
    "Items": [
        {
            "AlbumTitle": {
                "S": "Somewhat Famous"
            },
            "SongTitle": {
                "S": "Call Me Today"
            },
            "Artist": {
                "S": "No One You Know"
            }
        }
    ],
    "ScannedCount": 1,
    "ConsumedCapacity": null
}
```

### EC2

#### 创建keypair

```
$ aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text > MyKeyPair.pem
```

#### 创建安全组

```
$ aws ec2 create-security-group --group-name my-sg --description "My security group" --vpc-id vpc-1a2b3c4d
{
    "GroupId": "sg-903004f8"
}
```

#### 添加入站规则

```
$ aws ec2 authorize-security-group-ingress --group-id sg-903004f8 --protocol tcp --port 22 --cidr 203.0.113.0/24
```

#### 创建EC2

```
$ aws ec2 run-instances --image-id ami-xxxxxxxx --count 1 --instance-type t2.micro --key-name MyKeyPair --security-group-ids sg-903004f8 --subnet-id subnet-6e7f829e
```

#### 列出EC2

列出类型为t2.micro的实例

```
$ aws ec2 describe-instances --filters "Name=instance-type,Values=t2.micro" --query "Reservations[].Instances[].InstanceId"
```


### IAM

#### 创建用户组

使用 `create-group` 命令创建组

```
$ aws iam create-group --group-name MyIamGroup
{
    "Group": {
        "GroupName": "MyIamGroup",
        "CreateDate": "2018-12-14T03:03:52.834Z",
        "GroupId": "AGPAJNUJ2W4IJVEXAMPLE",
        "Arn": "arn:aws-cn:iam::123456789012:group/MyIamGroup",
        "Path": "/"
    }
}
```

#### 创建用户

使用 `create-user` 命令创建用户

```
$ aws iam create-user --user-name MyUser
{
    "User": {
        "UserName": "MyUser",
        "Path": "/",
        "CreateDate": "2018-12-14T03:13:02.581Z",
        "UserId": "AIDAJY2PE5XUZ4EXAMPLE",
        "Arn": "arn:aws-cn:iam::123456789012:user/MyUser"
    }
}
```

#### 将用户添加到组中

使用 `add-user-to-group` 命令将用户添加到组中

```
$ aws iam add-user-to-group --user-name MyUser --group-name MyIamGroup
```

#### 将 IAM 托管策略附加到 IAM 用户

1. 确定要附加的策略的 Amazon 资源名称 (ARN)。以下命令使用 `list-policies` 查找具有名称 `PowerUserAccess` 的策略的 ARN。然后,它会将该 ARN 存储在环境变量中。

  ```
$ export POLICYARN=$(aws iam list-policies --query 'Policies[?PolicyName==`PowerUserAccess`].{ARN:Arn}' --output text)       ~
$ echo $POLICYARN
arn:aws-cn:iam::aws:policy/PowerUserAccess
  ```

2. 要附加策略,请使用 `attach-user-policy` 命令,并引用存放策略 ARN 的环境变量。

  ```
$ aws iam attach-user-policy --user-name MyUser --policy-arn $POLICYARN
  ```

3. 通过运行 `list-attached-user-policies` 命令验证策略已附加到此用户。

  ```
$ aws iam list-attached-user-policies --user-name MyUser
{
    "AttachedPolicies": [
        {
            "PolicyName": "PowerUserAccess",
            "PolicyArn": "arn:aws-cn:iam::aws:policy/PowerUserAccess"
        }
    ]
}
  ```

#### 为 IAM 用户设置初始密码

```
$ aws iam create-login-profile --user-name MyUser --password My!User1Login8P@ssword --password-reset-required
{
    "LoginProfile": {
        "UserName": "MyUser",
        "CreateDate": "2018-12-14T17:27:18Z",
        "PasswordResetRequired": true
    }
}
```

#### 更改 IAM 用户的密码

可以使用 `update-login-profile` 命令更改 IAM 用户的密码。

```
$ aws iam update-login-profile --user-name MyUser --password My!User1ADifferentP@ssword
```

#### 创建访问密钥

使用 `create-access-key` 命令为 IAM 用户创建访问密钥。访问密钥是一组安全凭证,由访问密钥 ID 和私有密钥组成。

IAM 用户一次只能创建两个访问密钥。如果您尝试创建第三组,则命令返回 `LimitExceeded` 错误。

```
$ aws iam create-access-key --user-name MyUser
{
    "AccessKey": {
        "UserName": "MyUser",
        "AccessKeyId": "AKIAIOSFODNN7EXAMPLE",
        "Status": "Active",
        "SecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
        "CreateDate": "2018-12-14T17:34:16Z"
    }
}
```

#### 删除访问密钥

使用 `delete-access-key` 命令为 IAM 用户删除访问密钥。使用访问密钥 ID 指定要删除的访问密钥。

```
$ aws iam delete-access-key --user-name MyUser --access-key-id AKIAIOSFODNN7EXAMPLE
```


### S3

#### 创建桶

```
$ aws s3 mb s3://bucket-name
```

#### 复制

```
$ aws s3 cp filename.txt s3://bucket-name
$ aws s3 cp s3://bucket-name/filename.txt ./

# 使用 cat 文本编辑器,将文本“hello world”流式传输到 s3://bucket-name/filename.txt 文件
$ cat "hello world" | aws s3 cp - s3://bucket-name/filename.txt

# 将 s3://bucket-name/filename.txt 文件流式传输到 stdout,并将内容输出到控制台
$ aws s3 cp s3://bucket-name/filename.txt -
hello world

# 将 s3://bucket-name/pre 的内容流式传输到 stdout,使用 bzip2 命令压缩文件,并将名为 key.bz2 的新压缩文件上传到 s3://bucket-name
$ aws s3 cp s3://bucket-name/pre - | bzip2 --best | aws s3 cp - s3://bucket-name/key.bz2
```

#### 移动

```
# 将对象从 s3://bucket-name/example 移动到 s3://my-bucket/
$ aws s3 mv s3://bucket-name/example s3://my-bucket/

# 将本地文件从当前工作目录移动到 Amazon S3 存储桶
$ aws s3 mv filename.txt s3://bucket-name

# 将文件从 Amazon S3 存储桶移动到当前工作目录,其中 ./ 指定当前的工作目录。
$ aws s3 mv s3://bucket-name/filename.txt ./
```

#### 同步

s3 sync 命令同步一个存储桶与一个目录中的内容,或者同步两个存储桶中的内容。通常,s3 sync 在源和目标之间复制缺失或过时的文件或对象。不过,您还可以提供 --delete 选项来从目标中删除源中不存在的文件或对象。

```
$ aws s3 sync . s3://my-bucket/path
upload: MySubdirectory\MyFile3.txt to s3://my-bucket/path/MySubdirectory/MyFile3.txt
upload: MyFile2.txt to s3://my-bucket/path/MyFile2.txt
upload: MyFile1.txt to s3://my-bucket/path/MyFile1.txt

# Delete local file
$ rm ./MyFile1.txt


# Sync with deletion - object is deleted from bucket
$ aws s3 sync . s3://my-bucket/path --delete
delete: s3://my-bucket/path/MyFile1.txt

# Delete object from bucket
$ aws s3 rm s3://my-bucket/path/MySubdirectory/MyFile3.txt
delete: s3://my-bucket/path/MySubdirectory/MyFile3.txt

# Sync with deletion - local file is deleted
$ aws s3 sync s3://my-bucket/path . --delete
delete: MySubdirectory\MyFile3.txt

# Sync with Infrequent Access storage class
$ aws s3 sync . s3://my-bucket/path --storage-class STANDARD_IA
```

#### 删除桶以及桶内所有有内容

```
$ aws s3 rb s3://bucket-name --force
```

#### 删除桶中对象

从 s3://bucket-name/example 中删除所有对象

```shell
$ aws s3 rm s3://bucket-name/example
```



## 常见问题


### 运行 aws 时收到“找不到命令”错误

**可能的原因:安装期间未更新操作系统“路径”**

此错误意味着操作系统找不到 AWS CLI 程序。安装可能不完整。

如果使用 pip 安装 AWS CLI,可能需要将包含 aws 程序的文件夹添加到操作系统的 PATH 环境变量,或更改其模式以使其可执行。

需要将 aws 可执行文件添加到操作系统的 PATH 环境变量中。按照相应过程中的步骤操作:

* Windows - [将 AWS CLI 版本 1 可执行文件添加到命令行路径](https://docs.aws.amazon.com/zh_cn/cli/latest/userguide/install-windows.html#awscli-install-windows-path)
* macOS - [将 AWS CLI 版本 1 可执行文件添加到 macOS 命令行路径](https://docs.aws.amazon.com/zh_cn/cli/latest/userguide/install-macos.html#awscli-install-osx-path)
* Linux - [将 AWS CLI 版本 1 可执行文件添加到命令行路径](https://docs.aws.amazon.com/zh_cn/cli/latest/userguide/install-linux.html#install-linux-path)

### “拒绝访问”错误

**可能的原因:AWS CLI 程序文件没有“运行”权限**

在 Linux 或 macOS 上,确保 aws 程序具有调用用户的运行权限。通常,权限设置为 755。

要添加用户的运行权限,请运行以下命令,并将 ~/.local/bin/aws 替换为计算机上指向此程序的路径。

```
$ chmod +x ~/.local/bin/aws
```

### “凭证无效”错误

**可能的原因:AWS CLI 从意外位置读取了凭证**

AWS CLI 读取凭证的位置可能与预期不同。可以运行 aws configure list 以确认使用哪些凭证。

以下示例说明如何检查用于默认配置文件的凭证。

```
$ aws configure list
      Name                    Value             Type    Location
      ----                    -----             ----    --------
   profile                <not set>             None    None
access_key     ****************XYVA shared-credentials-file
secret_key     ****************ZAGY shared-credentials-file
    region                us-west-2      config-file    ~/.aws/config
```

以下示例说明如何检查命名配置文件的凭证。

```
$ aws configure list --profile saanvi
      Name                    Value             Type    Location
      ----                    -----             ----    --------
   profile                    saanvi           manual    --profile
access_key         **************** shared-credentials-file
secret_key         **************** shared-credentials-file
    region                us-west-2      config-file    ~/.aws/config
```

**可能的原因:计算机的时钟不同步**

如果使用有效的凭证,则时钟可能不同步导致了这个问题。在 Linux 或 macOS 上,运行 date 以检查时间。

```
$ date
```

如果系统时钟在几分钟内不正确,则使用 ntpd 进行同步。

```
$ sudo service ntpd stop
$ sudo ntpdate time.nist.gov
$ sudo service ntpd start
$ ntpstat
```

在 Windows 上,使用控制面板中的日期和时间选项来配置系统时钟。

### “签名不匹配”错误。

当 AWS CLI 运行命令时,它会向 AWS 服务器发送加密请求以执行适当的 AWS 服务操作。凭证(访问密钥和私有密钥)参与了加密过程,使 AWS 能够对发出请求的人员进行身份验证。有多种因素可能会干扰此过程的正常执行,如下所示。

**可能的原因:时钟与 AWS 服务器不同步**

为了帮助防范重播攻击,在加密/解密过程中可能会使用当前时间。如果客户端和服务器的时间不一致超出允许的时间量,该过程可能会失败,并且请求会被拒绝。当在时钟与主机时钟不同步的虚拟机中运行命令时,也可能发生此错误。一个可能原因是,当虚拟机休眠时,唤醒后需要一些时间才能将时钟与主机重新同步。

在 Linux 或 macOS 上,运行 date 以检查时间。

```
$ date
```

如果系统时钟在几分钟内不正确,则使用 ntpd 进行同步。

```
$ sudo service ntpd stop
$ sudo ntpdate time.nist.gov
$ sudo service ntpd start
$ ntpstat
```

在 Windows 上,使用控制面板中的日期和时间选项来配置系统时钟。

**可能的原因:操作系统错误地处理了包含某些特殊字符的 AWS 私有密钥**

如果 AWS 私有密钥包含某些特殊字符(例如 - 、+、/ 或 %),则某些操作系统变体会不正确地处理该字符串,并导致对该私有密钥字符串的解释不正确。huoz

如果使用其他工具或脚本(例如,在创建新实例期间在其上构建凭证文件的工具)处理访问密钥和私有密钥,则这些工具和脚本可能有自己对特殊字符的处理,这会使它们转换为 AWS 不再识别的内容。

简单的解决方案是重新生成私有密钥,以获得一个不包含特殊字符的私有密钥。

**可能的原因:配置文件包含额外内容,导致配置文件格式不正确**

重新生成配置文件或手动改正配置文件



