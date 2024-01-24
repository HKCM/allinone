

## 描述: WAF&Shield 常用的CLI命令


### WAF

#### 列出托管的的规则组
```shell
aws wafv2 list-available-managed-rule-groups --scope REGIONAL

aws wafv2 list-available-managed-rule-groups --scope=CLOUDFRONT --region=us-east-1
```

#### 获取 Amazon CloudFront 分配中使用的基于速率的规则阻止的 IP 地址列表
```shell
aws wafv2 get-rate-based-statement-managed-keys \
  --scope=CLOUDFRONT \
  --region=us-east-1 \
  --web-acl-name=my-test \
  --web-acl-id=1234abcd-d4e3-486c-a95d-fd25c4dbabcd \
  --rule-name=rate-limit
```

#### 获取基于API Gateway REST API、ALB 或AppSync GraphQL 基于速率的规则阻止的 IP 地址列表
```shell
aws wafv2 get-rate-based-statement-managed-keys \
  --scope=REGIONAL \
  --region=region \
  --web-acl-name=WebACLName \
  --web-acl-id=WebACLId \
  --rule-name=RuleName
```
