#### 获取集群

```
kubectl config get-clusters
NAME
aws94-c01-kbm20
aws94-c01-kbm10
```

#### 删除集群配置

```
kubectl config delete-cluster aws94-c01-kbm20
```

#### 查看当前使用的context

```shell
# 查看当前使用的context
kubectl config current-context

# 查看所有context
kubectl config get-contexts
CURRENT   NAME    CLUSTER           AUTHINFO                NAMESPACE
          itl     aws94-c01-kbm20   aws94-c01-kbm20-admin   
          xmnup   aws94-c01-kbm10   aws94-c01-kbm10-admin
          
# 切换context
kubectl config use-context itl
Switched to context "itl".

# 查看详细信息
kubectl config view

# 删除 context
kubectl config unset contexts.aws_cluster1-kubernetes
```



#### 删除 user

```
kubectl config unset users.xxx
Property "users.xxxx" unset.
```

#### 更新配置

```shell
aws eks update-kubeconfig \
--region ap-southeast-1 \
--name aws94-c01-kbm10 \
--profile int-xmn \
--role-arn arn:aws:iam::1234567890:role/cops-eks-admin-aws94-c01-kbm10
```

#### 查看所有pods

```shell
# 查看所有pods
kubectl get pods -A

# 查看pod事件
kubectl describe pod aws94-abe-8f9cd4697-m85lx -n test-namespace
```

#### AWS EKS 示例config文件

通过特定的role更新kube config文件
```shell
aws eks update-kubeconfig \
--region ap-southeast-1 \
--name clusterName \
--profile myprofile \
--role-arn arn:aws:iam::123456789012:role/eks-admin-role-name
```
```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: EXAMPLE--URVJUSUZ...JSUN5RENDQWJDZ0F3SUJL3ZMUmJ
    server: https://XXXXC93234XXX5B30B624FXXXXX.gr7.ap-southeast-1.eks.amazonaws.com
  name: arn:aws:eks:ap-southeast-1:123456789012:cluster/example
contexts:
- context:
    cluster: arn:aws:eks:ap-southeast-1:123456789012:cluster/example
    user: arn:aws:eks:ap-southeast-1:123456789012:cluster/example
  name: xmnup
current-context: xmnup
kind: Config
preferences: {}
users:
- name: arn:aws:eks:ap-southeast-1:123456789012:cluster/example
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - --region
      - ap-southeast-1
      - eks
      - get-token
      - --cluster-name
      - example
      - --role
      - arn:aws:iam::123456789012:role/eks-admin-example
      command: aws
      env:
      - name: AWS_PROFILE
        value: myprofile
```

该文件是从`arn:aws:iam::123456789012:role/eks-admin-example` Role中获取对EKS的控制权, 又是通过本地的`myprofile` 获取该Role的token

#### 更新配置

```shell
aws eks update-kubeconfig \
--region ap-southeast-1 \
--name aws94-c01-kbm10 \
--profile int-xmn \
--role-arn arn:aws:iam::123456789:role/cops-eks-admin-aws94-c01-kbm10
```


