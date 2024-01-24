# EKS



## 安装kubectl

1. 下载二进制文件

- **Kubernetes 1.20:**

  ```
  curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.20.4/2021-04-12/bin/darwin/amd64/kubectl
  ```

- **Kubernetes 1.19:**

  ```
  curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/darwin/amd64/kubectl
  ```

- **Kubernetes 1.18:**

  ```
  curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.9/2020-11-02/bin/darwin/amd64/kubectl
  ```

- **Kubernetes 1.17:**

  ```
  curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.17.12/2020-11-02/bin/darwin/amd64/kubectl
  ```

- **Kubernetes 1.16:**

  ```
  curl -o kubectl https://amazon-eks.s3.us-west-2.amazonaws.com/1.16.15/2020-11-02/bin/darwin/amd64/kubectl
  ```

2. 添加执行权限

   ```
   chmod +x ./kubectl
   ```

3. 将二进制文件复制到您的 `PATH` 中的文件夹。如果您已经安装了某个版本的 `kubectl`,建议您创建一个 `$HOME/bin/kubectl` 并确保 `$HOME/bin` 先出现在您的 `$PATH` 中。

   ```shell
   mkdir -p $HOME/bin && cp ./kubectl $HOME/bin/kubectl && export PATH=$HOME/bin:$PATH
   ```

4. 将 `$HOME/bin` 路径添加到 shell 初始化文件,以便在打开 shell 时配置此路径。

   ```shell
   echo 'export PATH=$PATH:$HOME/bin' >> ~/.bash_profile
   ```

5. 安装 `kubectl` 后,可以使用以下命令验证其版本:

   ```
   kubectl version --short --client
   ```

参考: https://docs.aws.amazon.com/zh_cn/eks/latest/userguide/install-kubectl.html



## 安装eksctl

1. 如果您尚未在 macOS 上安装 Homebrew,请使用以下命令安装它。

   ```
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
   ```

2. 安装 Weaveworks Homebrew tap。

   ```
   brew tap weaveworks/tap
   ```

3. 安装或升级 `eksctl`。

- 使用以下命令安装 `eksctl`:

```shell
brew install weaveworks/tap/eksctl
```

- 如果已安装 `eksctl`,请运行以下命令进行升级:

```shell
brew upgrade eksctl && brew link --overwrite eksctl
```

4. 使用以下命令测试您的安装是否成功。

```shell
eksctl version
```

官网: https://github.com/weaveworks/eksctl

## 启动集群

1. 首先创建 key pair

```shell
aws ec2 --region us-east-2 create-key-pair --key-name myKey-us-east-2 --query "KeyMaterial" --output text > myKey-us-east-2.pemchmod 400 myKey-us-east-2.pem
```

2. 创建集群

默认使用 2 个m5.large机器创建在 us-west-2 区域

```shell
eksctl create cluster
```


以Spot实例启用eks
```shell
# with-oidc Kubernetes 服务账户可以访问Amazon资源
$ eksctl create cluster \
--profile hkc_admin \
--version=1.18 \
--tags environment=staging \
--node-volume-size=20 \
--node-volume-type=gp2 \
--name my-cluster \
--with-oidc \
--region us-east-2 \
--ssh-access \
--ssh-public-key myKey-us-east-2 \
--spot \
--instance-types=t3.medium \
--managed
2021-09-26 09:26:35 [ℹ]  eksctl version 0.67.0
2021-09-26 09:26:35 [ℹ]  using region us-east-2
2021-09-26 09:26:37 [ℹ]  setting availability zones to [us-east-2a us-east-2b us-east-2c]
2021-09-26 09:26:37 [ℹ]  subnets for us-east-2a - public:192.168.0.0/19 private:192.168.96.0/19
2021-09-26 09:26:37 [ℹ]  subnets for us-east-2b - public:192.168.32.0/19 private:192.168.128.0/19
2021-09-26 09:26:37 [ℹ]  subnets for us-east-2c - public:192.168.64.0/19 private:192.168.160.0/19
2021-09-26 09:26:37 [ℹ]  nodegroup "ng-2195232e" will use "" [AmazonLinux2/1.18]
2021-09-26 09:26:37 [ℹ]  using EC2 key pair %!q(*string=<nil>)
2021-09-26 09:26:37 [ℹ]  using Kubernetes version 1.18
2021-09-26 09:26:37 [ℹ]  creating EKS cluster "my-cluster" in "us-east-2" region with managed nodes
2021-09-26 09:26:37 [ℹ]  will create 2 separate CloudFormation stacks for cluster itself and the initial managed nodegroup
2021-09-26 09:26:37 [ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-east-2 --cluster=my-cluster'
2021-09-26 09:26:37 [ℹ]  CloudWatch logging will not be enabled for cluster "my-cluster" in "us-east-2"
2021-09-26 09:26:37 [ℹ]  you can enable it with 'eksctl utils update-cluster-logging --enable-types={SPECIFY-YOUR-LOG-TYPES-HERE (e.g. all)} --region=us-east-2 --cluster=my-cluster'
2021-09-26 09:26:37 [ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "my-cluster" in "us-east-2"
2021-09-26 09:26:37 [ℹ]  2 sequential tasks: { create cluster control plane "my-cluster", 3 sequential sub-tasks: { 5 sequential sub-tasks: { wait for control plane to become ready, tag cluster, associate IAM OIDC provider, 2 sequential sub-tasks: { create IAM role for serviceaccount "kube-system/aws-node", create serviceaccount "kube-system/aws-node" }, restart daemonset "kube-system/aws-node" }, 1 task: { create addons }, create managed nodegroup "ng-2195232e" } }
2021-09-26 09:26:37 [ℹ]  building cluster stack "eksctl-my-cluster-cluster"
2021-09-26 09:26:41 [ℹ]  deploying stack "eksctl-my-cluster-cluster"
2021-09-26 09:27:11 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:27:43 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:28:45 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:29:46 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:30:47 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:31:49 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:32:51 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:33:52 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:34:53 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:35:54 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:36:56 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:37:58 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:39:00 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:40:01 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:41:02 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-cluster"
2021-09-26 09:43:14 [✔]  tagged EKS cluster (environment=staging)
2021-09-26 09:45:18 [ℹ]  building iamserviceaccount stack "eksctl-my-cluster-addon-iamserviceaccount-kube-system-aws-node"
2021-09-26 09:45:20 [ℹ]  deploying stack "eksctl-my-cluster-addon-iamserviceaccount-kube-system-aws-node"
2021-09-26 09:45:20 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-iamserviceaccount-kube-system-aws-node"
2021-09-26 09:45:38 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-addon-iamserviceaccount-kube-system-aws-node"
2021-09-26 09:45:41 [ℹ]  serviceaccount "kube-system/aws-node" already exists
2021-09-26 09:45:41 [ℹ]  updated serviceaccount "kube-system/aws-node"
2021-09-26 09:45:42 [ℹ]  daemonset "kube-system/aws-node" restarted
2021-09-26 09:47:47 [ℹ]  building managed nodegroup stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:47:49 [ℹ]  deploying stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:47:49 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:48:07 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:48:28 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:48:48 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:49:05 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:49:24 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:49:43 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:50:01 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:50:19 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:50:37 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:50:53 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:51:11 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:51:32 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:51:53 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:52:14 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:52:34 [ℹ]  waiting for CloudFormation stack "eksctl-my-cluster-nodegroup-ng-2195232e"
2021-09-26 09:52:37 [ℹ]  waiting for the control plane availability...
2021-09-26 09:52:37 [✔]  saved kubeconfig as "/Users/USER/.kube/config"
2021-09-26 09:52:37 [ℹ]  no tasks
2021-09-26 09:52:37 [✔]  all EKS cluster resources for "my-cluster" have been created
2021-09-26 09:52:38 [ℹ]  nodegroup "ng-2195232e" has 2 node(s)
2021-09-26 09:52:38 [ℹ]  node "ip-192-168-31-108.us-east-2.compute.internal" is ready
2021-09-26 09:52:38 [ℹ]  node "ip-192-168-83-220.us-east-2.compute.internal" is ready
2021-09-26 09:52:38 [ℹ]  waiting for at least 2 node(s) to become ready in "ng-2195232e"
2021-09-26 09:52:38 [ℹ]  nodegroup "ng-2195232e" has 2 node(s)
2021-09-26 09:52:38 [ℹ]  node "ip-192-168-31-108.us-east-2.compute.internal" is ready
2021-09-26 09:52:38 [ℹ]  node "ip-192-168-83-220.us-east-2.compute.internal" is ready
2021-09-26 09:54:49 [ℹ]  kubectl command should work with "/Users/USER/.kube/config", try 'kubectl get nodes'
2021-09-26 09:54:49 [✔]  EKS cluster "my-cluster" in "us-east-2" region is ready
```

`eksctl`创建`kubectl` `config`在中的文件`~/.kube`或将新集群的配置添加到现有`config`在中的文件`~/.kube`

实际行为会创建三个 cloudformation:

1. cluster cloudformation分别创建以下资源:
   1. VPC
   2. Subnet
   3. SecurityGroup
   4. Gateway(NatGateway, InternetGateway)
   5. RouteTable
   6. EIP
   7. Cluster
2. iamserviceaccount cloudformation分别创建以下资源:
   1. IAM Role
3. NodeGroup cloudformation分别创建以下资源:
   1. SecurityGroup
   2. Nodegroup
   3. LaunchTemplate
   4. EC2 IAM Role

3. (可选)启用日志

```shell
eksctl utils update-cluster-logging --enable-types=all --region=us-east-2 --cluster=my-cluster
```

   

4. 查看资源

为集群创建或更新 `kubeconfig` 文件

```shell
aws eks update-kubeconfig --region us-east-2 --name my-cluster
```

   

```shell
kubectl get nodes -o wide
kubectl get nodes -o wide -l eks.amazonaws.com/capacityType=SPOT
kubectl get nodes -o wide -l eks.amazonaws.com/capacityType=ONDEMAND
nodeselector
eks.amazonaws.com/capacityType=SPOT
eks.amazonaws.com/nodegroup=ng-9c0e7ac2
```

**Amazon EC2 节点输出**

```shell
NAME                                           STATUS   ROLES    AGE   VERSION               INTERNAL-IP      EXTERNAL-IP      OS-IMAGE         KERNEL-VERSION                  CONTAINER-RUNTIMEip-192-168-31-108.us-east-2.compute.internal   Ready    <none>   74m   v1.18.20-eks-c9f1ce   192.168.31.108   18.191.234.232   Amazon Linux 2   4.14.243-185.433.amzn2.x86_64   docker://19.3.13ip-192-168-83-220.us-east-2.compute.internal   Ready    <none>   74m   v1.18.20-eks-c9f1ce   192.168.83.220   18.188.6.113     Amazon Linux 2   4.14.243-185.433.amzn2.x86_64   docker://19.3.13
```

5. 安装CNI插件

[管理 Amazon VPC CNI 附加组件](https://docs.aws.amazon.com/zh_cn/eks/latest/userguide/managing-vpc-cni.html)

Role-arn: addon-iamserviceaccount-kube-system-aws-node output

```shell
eksctl create addon \
--profile hkc_admin \
--region us-east-2 \
--name vpc-cni \
--version latest \
--cluster my-cluster \
--service-account-role-arn arn:aws:iam::734871910852:role/eksctl-my-cluster-addon-iamserviceaccount-ku-Role1-1ICXJNFOS0K2H \
--force
```

删除插件

```
eksctl delete addon --cluster my-cluster --name vpc-cni --profile hkc_admin --region us-east-2
```

查看适用插件,`True`表示版本可用

```shell
aws eks describe-addon-versions \
--profile hkc_admin \
--region us-east-2 \
--addon-name vpc-cni \
--kubernetes-version 1.19 \
--query "addons[].addonVersions[].[addonVersion, compatibilities[].defaultVersion]" \
--output text
...
v1.7.6-eksbuild.1False
v1.7.5-eksbuild.2True
v1.7.5-eksbuild.1False
...
   ```

安装插件

```shell
aws eks --profile hkc_admin --region us-east-2 create-addon --cluster-name my-cluster  --addon-name vpc-cni
{    
   "addon": {
      "addonName": "vpc-cni",
      "clusterName": "my-cluster", 
      "status": "CREATING",
      "addonVersion": "v1.7.5-eksbuild.2",
      "health": {
         "issues": []
      },
      "addonArn": "arn:aws:eks:us-east-2:734871910852:addon/my-cluster/vpc-cni/96be1005-fd41-56ad-009d-753d14d17c5d",
      "createdAt": "2021-09-26T11:00:28.566000+08:00",
      "modifiedAt": "2021-09-26T11:00:28.584000+08:00",
      "tags": {}
   }
}
```

查看插件

```shell
aws eks describe-addon \
--profile hkc_admin \
--region us-east-2 \
--cluster-name my-cluster \
--addon-name vpc-cni \
--query "addon.addonVersion" \
--output text
v1.7.5-eksbuild.2
```

更新插件并指定版本

```shell
aws eks update-addon \
--profile hkc_admin \
--region us-east-2 \
--cluster-name my-cluster \
--addon-name vpc-cni \
--addon-version 1.7.10.eksbuild-1 \
--resolve-conflicts
```

删除插件

```shell
aws eks delete-addon \
--profile hkc_admin \
--region us-east-2 \
--cluster-name my-cluster \
--addon-name vpc-cni
```

6. CoreDNS插件

```shell
aws eks --profile hkc_admin --region us-east-2 create-addon --cluster-name my-cluster --addon-name coredns
{
   "addon": {
      "addonName": "coredns",
      "clusterName": "my-cluster",
      "status": "CREATING",
      "addonVersion": "v1.7.0-eksbuild.1",
      "health": {
         "issues": []
      },
      "addonArn": "arn:aws:eks:us-east-2:734871910852:addon/my-cluster/coredns/36be1072-ab0d-c2c3-fa1f-8c12763e095b",
      "createdAt": "2021-09-26T14:57:53.326000+08:00",
      "modifiedAt": "2021-09-26T14:57:53.342000+08:00",
      "tags": {}
   }
}
```

   

## 删除集群

```shell
eksctl delete cluster \
--profile hkc_admin \
--region us-east-2 \
--name my-cluster 
```

## 更新集群

使用以下命令获取集群控制层面的 Kubernetes 版本。

```
kubectl version --short
```

使用以下命令获取节点的 Kubernetes 版本。此命令会返回所有自我管理和托管的 Amazon EC2 和 Fargate 节点。每个 Fargate Pod 都作为其自身的节点列出。

```shell
kubectl get nodes
```



### 更新集群

```shell
aws eks update-cluster-version --profile hkc_admin --region us-east-2 --name my-cluster --kubernetes-version 1.19
{
   "update": {
      "id": "6251aa95-8084-435a-b181-369dfc753b8a",
      "status": "InProgress",
      "type": "VersionUpdate",
      "params": [
         {
            "type": "Version",
            "value": "1.19"
         },
         {
            "type": "PlatformVersion",
            "value": "eks.6"
         }
      ],
      "createdAt": "2021-10-08T00:36:51.882000+08:00",
      "errors": []
   }
}
aws eks describe-update --profile hkc_admin --region us-east-2 --name my-cluster --update-id 6251aa95-8084-435a-b181-369dfc753b8a
{
   "update": {
      "id": "6251aa95-8084-435a-b181-369dfc753b8a",
      "status": "InProgress",
      "type": "VersionUpdate",
      "params": [
         {
            "type": "Version",
            "value": "1.19"
         },
         {
            "type": "PlatformVersion",
            "value": "eks.6"
         }
      ],
   "createdAt": "2021-10-08T00:36:51.882000+08:00",
   "errors": []
   }
}
```

### 更新节点

```shell
eksctl upgrade nodegroup \
--profile hkc_admin \
--region us-east-2 \
--name=ng-e4d742ae \
--cluster=my-cluster \
--kubernetes-version=1.19
```

## 部署服务

###  添加负载均衡器控制器

1. 下载AWS负载均衡器控制器的 IAM 策略,该策略允许负载均衡器代表您调用 AWS API。您可以查看 GitHub 上的[策略文档](https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.2.0/docs/install/iam_policy.json)。

```shell
curl -o iam_policy.json https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.2.0/docs/install/iam_policy.json
```

IAM Policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:CreateServiceLinkedRole",
                "ec2:DescribeAccountAttributes",
                "ec2:DescribeAddresses",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeInternetGateways",
                "ec2:DescribeVpcs",
                "ec2:DescribeSubnets",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeInstances",
                "ec2:DescribeNetworkInterfaces",
                "ec2:DescribeTags",
                "ec2:GetCoipPoolUsage",
                "ec2:DescribeCoipPools",
                "elasticloadbalancing:DescribeLoadBalancers",
                "elasticloadbalancing:DescribeLoadBalancerAttributes",
                "elasticloadbalancing:DescribeListeners",
                "elasticloadbalancing:DescribeListenerCertificates",
                "elasticloadbalancing:DescribeSSLPolicies",
                "elasticloadbalancing:DescribeRules",
                "elasticloadbalancing:DescribeTargetGroups",
                "elasticloadbalancing:DescribeTargetGroupAttributes",
                "elasticloadbalancing:DescribeTargetHealth",
                "elasticloadbalancing:DescribeTags"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "cognito-idp:DescribeUserPoolClient",
                "acm:ListCertificates",
                "acm:DescribeCertificate",
                "iam:ListServerCertificates",
                "iam:GetServerCertificate",
                "waf-regional:GetWebACL",
                "waf-regional:GetWebACLForResource",
                "waf-regional:AssociateWebACL",
                "waf-regional:DisassociateWebACL",
                "wafv2:GetWebACL",
                "wafv2:GetWebACLForResource",
                "wafv2:AssociateWebACL",
                "wafv2:DisassociateWebACL",
                "shield:GetSubscriptionState",
                "shield:DescribeProtection",
                "shield:CreateProtection",
                "shield:DeleteProtection"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:RevokeSecurityGroupIngress"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:CreateSecurityGroup"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:CreateTags"
            ],
            "Resource": "arn:aws:ec2:*:*:security-group/*",
            "Condition": {
                "StringEquals": {
                    "ec2:CreateAction": "CreateSecurityGroup"
                },
                "Null": {
                    "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:CreateTags",
                "ec2:DeleteTags"
            ],
            "Resource": "arn:aws:ec2:*:*:security-group/*",
            "Condition": {
                "Null": {
                    "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
                    "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:RevokeSecurityGroupIngress",
                "ec2:DeleteSecurityGroup"
            ],
            "Resource": "*",
            "Condition": {
                "Null": {
                    "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:CreateLoadBalancer",
                "elasticloadbalancing:CreateTargetGroup"
            ],
            "Resource": "*",
            "Condition": {
                "Null": {
                    "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:CreateListener",
                "elasticloadbalancing:DeleteListener",
                "elasticloadbalancing:CreateRule",
                "elasticloadbalancing:DeleteRule"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:AddTags",
                "elasticloadbalancing:RemoveTags"
            ],
            "Resource": [
                "arn:aws:elasticloadbalancing:*:*:targetgroup/*/*",
                "arn:aws:elasticloadbalancing:*:*:loadbalancer/net/*/*",
                "arn:aws:elasticloadbalancing:*:*:loadbalancer/app/*/*"
            ],
            "Condition": {
                "Null": {
                    "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
                    "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:AddTags",
                "elasticloadbalancing:RemoveTags"
            ],
            "Resource": [
                "arn:aws:elasticloadbalancing:*:*:listener/net/*/*/*",
                "arn:aws:elasticloadbalancing:*:*:listener/app/*/*/*",
                "arn:aws:elasticloadbalancing:*:*:listener-rule/net/*/*/*",
                "arn:aws:elasticloadbalancing:*:*:listener-rule/app/*/*/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:ModifyLoadBalancerAttributes",
                "elasticloadbalancing:SetIpAddressType",
                "elasticloadbalancing:SetSecurityGroups",
                "elasticloadbalancing:SetSubnets",
                "elasticloadbalancing:DeleteLoadBalancer",
                "elasticloadbalancing:ModifyTargetGroup",
                "elasticloadbalancing:ModifyTargetGroupAttributes",
                "elasticloadbalancing:DeleteTargetGroup"
            ],
            "Resource": "*",
            "Condition": {
                "Null": {
                    "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
                }
            }
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:RegisterTargets",
                "elasticloadbalancing:DeregisterTargets"
            ],
            "Resource": "arn:aws:elasticloadbalancing:*:*:targetgroup/*/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:SetWebAcl",
                "elasticloadbalancing:ModifyListener",
                "elasticloadbalancing:AddListenerCertificates",
                "elasticloadbalancing:RemoveListenerCertificates",
                "elasticloadbalancing:ModifyRule"
            ],
            "Resource": "*"
        }
    ]
}
```

   

2. 创建策略,记下返回的策略 ARN。

```shell
aws iam create-policy \--profile hkc_admin \--policy-name AWSLoadBalancerControllerIAMPolicy \--policy-document file://iam_policy.json
```

3. 创建 role,这里会创建一个 cloudformation

```shell
eksctl create iamserviceaccount \--profile hkc_admin \--region us-east-2 \--cluster=my-cluster \--namespace=kube-system \--name=aws-load-balancer-controller \--attach-policy-arn=arn:aws:iam::734871910852:policy/AWSLoadBalancerControllerIAMPolicy \--override-existing-serviceaccounts \--approve  
```

   

4. 安装 `cert-manager` 以将证书配置注入到 Webhook 中。

```shell
kubectl apply \--validate=false \-f https://github.com/jetstack/cert-manager/releases/download/v1.1.1/cert-manager.yaml
```

   

5. 安装控制器。下载控制器规范。有关控制器的更多信息,请参阅 GitHub 上的[文档](https://kubernetes-sigs.github.io/aws-load-balancer-controller/)。

```shell
curl -o v2_2_0_full.yaml https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.2.0/docs/install/v2_2_0_full.yaml
```

对 `v2_2_0_full.yaml` 文件进行以下编辑。

删除文件的 `ServiceAccount` 部分。删除此部分可防止在部署控制器时覆盖带有 IAM 角色的注释,并在删除控制器时保留您在第 4 步中创建的服务账户。

将文件 `Deploymentspec` 部分的 `your-cluster-name` 替换为自己的集群的名称。应用文件。

```shell
kubectl apply -f v2_2_0_full.yaml
```

6. 验证控制器是否已安装。

```shell
kubectl get deployment -n kube-system aws-load-balancer-controller
```

输出

```
NAME                           READY   UP-TO-DATE   AVAILABLE   AGEaws-load-balancer-controller   1/1     1            1           14s
```


7. 参考

https://docs.aws.amazon.com/zh_cn/eks/latest/userguide/aws-load-balancer-controller.html



###  部署2048服务

部署一个示例小游戏服务

```shell
kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.2.0/docs/examples/2048/2048_full.yaml
```

2048_full.yaml

```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: game-2048
---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: game-2048
  name: deployment-2048
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: app-2048
  replicas: 5
  template:
    metadata:
      labels:
        app.kubernetes.io/name: app-2048
    spec:
      containers:
      - image: alexwhen/docker-2048
        imagePullPolicy: Always
        name: app-2048
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  namespace: game-2048
  name: service-2048
spec:
  ports:
    - port: 80
      targetPort: 80
      protocol: TCP
  type: NodePort
  selector:
    app.kubernetes.io/name: app-2048
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  namespace: game-2048
  name: ingress-2048
  annotations:
    kubernetes.io/ingress.class: alb
    alb.ingress.kubernetes.io/scheme: internet-facing
    alb.ingress.kubernetes.io/target-type: ip
spec:
  rules:
    - http:
        paths:
          - path: /*
            backend:
              serviceName: service-2048
              servicePort: 80
```

### 部署 dashboard

recommended.yaml
关于 ingress 的 ann字段说明文档: https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.1/guide/ingress/annotations/

```yaml
# Copyright 2017 The Kubernetes Authors.#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at#
#     http://www.apache.org/licenses/LICENSE-2.0#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
apiVersion: v1
kind: Namespace
metadata:  
   name: kubernetes-dashboard
---
apiVersion: v1
kind: ServiceAccount
metadata:  
   labels:    
      k8s-app: kubernetes-dashboard
      name: kubernetes-dashboard
      namespace: kubernetes-dashboard
---
kind: ServiceapiVersion: v1metadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard  namespace: kubernetes-dashboardspec:  ports:    - port: 8443      targetPort: 8443  selector:    k8s-app: kubernetes-dashboard---apiVersion: networking.k8s.io/v1beta1kind: Ingressmetadata:  name: kubernetes-dashboard  namespace: kubernetes-dashboard  annotations:    kubernetes.io/ingress.class: alb    alb.ingress.kubernetes.io/scheme: internet-facing    alb.ingress.kubernetes.io/target-type: ip    alb.ingress.kubernetes.io/backend-protocol: HTTPSspec:  rules:    - http:        paths:          - path: /*            backend:              serviceName: kubernetes-dashboard              servicePort: 8443---apiVersion: v1kind: Secretmetadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard-certs  namespace: kubernetes-dashboardtype: Opaque---apiVersion: v1kind: Secretmetadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard-csrf  namespace: kubernetes-dashboardtype: Opaquedata:  csrf: ""---apiVersion: v1kind: Secretmetadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard-key-holder  namespace: kubernetes-dashboardtype: Opaque---kind: ConfigMapapiVersion: v1metadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard-settings  namespace: kubernetes-dashboard---kind: RoleapiVersion: rbac.authorization.k8s.io/v1metadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard  namespace: kubernetes-dashboardrules:  # Allow Dashboard to get, update and delete Dashboard exclusive secrets.  - apiGroups: [""]    resources: ["secrets"]    resourceNames: ["kubernetes-dashboard-key-holder", "kubernetes-dashboard-certs", "kubernetes-dashboard-csrf"]    verbs: ["get", "update", "delete"]    # Allow Dashboard to get and update 'kubernetes-dashboard-settings' config map.  - apiGroups: [""]    resources: ["configmaps"]    resourceNames: ["kubernetes-dashboard-settings"]    verbs: ["get", "update"]    # Allow Dashboard to get metrics.  - apiGroups: [""]    resources: ["services"]    resourceNames: ["heapster", "dashboard-metrics-scraper"]    verbs: ["proxy"]  - apiGroups: [""]    resources: ["services/proxy"]    resourceNames: ["heapster", "http:heapster:", "https:heapster:", "dashboard-metrics-scraper", "http:dashboard-metrics-scraper"]    verbs: ["get"]---kind: ClusterRoleapiVersion: rbac.authorization.k8s.io/v1metadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboardrules:  # Allow Metrics Scraper to get metrics from the Metrics server  - apiGroups: ["metrics.k8s.io"]    resources: ["pods", "nodes"]    verbs: ["get", "list", "watch"]---apiVersion: rbac.authorization.k8s.io/v1kind: RoleBindingmetadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard  namespace: kubernetes-dashboardroleRef:  apiGroup: rbac.authorization.k8s.io  kind: Role  name: kubernetes-dashboardsubjects:  - kind: ServiceAccount    name: kubernetes-dashboard    namespace: kubernetes-dashboard---apiVersion: rbac.authorization.k8s.io/v1kind: ClusterRoleBindingmetadata:  name: kubernetes-dashboardroleRef:  apiGroup: rbac.authorization.k8s.io  kind: ClusterRole  name: kubernetes-dashboardsubjects:  - kind: ServiceAccount    name: kubernetes-dashboard    namespace: kubernetes-dashboard---kind: DeploymentapiVersion: apps/v1metadata:  labels:    k8s-app: kubernetes-dashboard  name: kubernetes-dashboard  namespace: kubernetes-dashboardspec:  replicas: 1  revisionHistoryLimit: 10  selector:    matchLabels:      k8s-app: kubernetes-dashboard  template:    metadata:      labels:        k8s-app: kubernetes-dashboard    spec:      containers:        - name: kubernetes-dashboard          image: kubernetesui/dashboard:v2.0.5          imagePullPolicy: Always          ports:            - containerPort: 8443              protocol: TCP          args:            - --auto-generate-certificates            - --namespace=kubernetes-dashboard            # Uncomment the following line to manually specify Kubernetes API server Host            # If not specified, Dashboard will attempt to auto discover the API server and connect            # to it. Uncomment only if the default does not work.            # - --apiserver-host=http://my-address:port          volumeMounts:            - name: kubernetes-dashboard-certs              mountPath: /certs              # Create on-disk volume to store exec logs            - mountPath: /tmp              name: tmp-volume          livenessProbe:            httpGet:              scheme: HTTPS              path: /              port: 8443            initialDelaySeconds: 30            timeoutSeconds: 30          securityContext:            allowPrivilegeEscalation: false            readOnlyRootFilesystem: true            runAsUser: 1001            runAsGroup: 2001      volumes:        - name: kubernetes-dashboard-certs          secret:            secretName: kubernetes-dashboard-certs        - name: tmp-volume          emptyDir: {}      serviceAccountName: kubernetes-dashboard      nodeSelector:        "kubernetes.io/os": linux      # Comment the following tolerations if Dashboard must not be deployed on master      tolerations:        - key: node-role.kubernetes.io/master          effect: NoSchedule---kind: ServiceapiVersion: v1metadata:  labels:    k8s-app: dashboard-metrics-scraper  name: dashboard-metrics-scraper  namespace: kubernetes-dashboardspec:  ports:    - port: 8000      targetPort: 8000  selector:    k8s-app: dashboard-metrics-scraper---kind: DeploymentapiVersion: apps/v1metadata:  labels:    k8s-app: dashboard-metrics-scraper  name: dashboard-metrics-scraper  namespace: kubernetes-dashboardspec:  replicas: 1  revisionHistoryLimit: 10  selector:    matchLabels:      k8s-app: dashboard-metrics-scraper  template:    metadata:      labels:        k8s-app: dashboard-metrics-scraper      annotations:        seccomp.security.alpha.kubernetes.io/pod: 'runtime/default'    spec:      containers:        - name: dashboard-metrics-scraper          image: kubernetesui/metrics-scraper:v1.0.6          ports:            - containerPort: 8000              protocol: TCP          livenessProbe:            httpGet:              scheme: HTTP              path: /              port: 8000            initialDelaySeconds: 30            timeoutSeconds: 30          volumeMounts:          - mountPath: /tmp            name: tmp-volume          securityContext:            allowPrivilegeEscalation: false            readOnlyRootFilesystem: true            runAsUser: 1001            runAsGroup: 2001      serviceAccountName: kubernetes-dashboard      nodeSelector:        "kubernetes.io/os": linux      # Comment the following tolerations if Dashboard must not be deployed on master      tolerations:        - key: node-role.kubernetes.io/master          effect: NoSchedule      volumes:        - name: tmp-volume          emptyDir: {}
```

```shell
kubectl apply -f recommended.yaml
```

获取 Dashboard token

```shell
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep eks-admin | awk '{print $1}')
```




## Blog

### [如何解决连接到 Amazon EKS API 服务器时出现的未经授权服务器错误](https://aws.amazon.com/cn/premiumsupport/knowledge-center/eks-api-server-unauthorized-error/)

检查您是否已经应用了 `aws-auth` ConfigMap。

```
kubectl describe configmap -n kube-system aws-auth
```

如果您收到错误指示“`Error from server (NotFound): configmaps "aws-auth" not found`”,则下载配置映射。

```
curl -o aws-auth-cm.yaml https://s3.us-west-2.amazonaws.com/amazon-eks/cloudformation/2020-10-29/aws-auth-cm.yaml
```

编辑配置

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
data:
  mapRoles: |
    - rolearn: <ARN of instance role (not instance profile)>
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
  mapUsers: |    
    - userarn: arn:aws:iam::734871910852:user/hkc_admin2      
      username: hkc_admin2      
      groups:        
        - system:masters
```

应用配置。此命令可能需要几分钟才能完成。

```
kubectl apply -f aws-auth-cm.yaml
```

在其他用户主机上再次更新`kubeconfig`

```
aws eks update-kubeconfig --region us-east-2 --profile hkc_admin2 --name my-cluster
```



### LoadBanlencer Ingress Annotations

[https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.1/guide/ingress/annotations/#annotations](https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.1/guide/ingress/annotations/#annotations)
