### 心得

防火墙如果有特殊设置,先允许特殊许可,例如内网IP可以跳过防火墙,其他

#### Template总览
```yaml
Resources:
  # ===========================================WAF==================================
  IPSet:
    Type: AWS::WAFv2::IPSet
    Properties: 
      Addresses: 
        - 110.110.110.110/32
        - 210.210.210.210/32
        - 202.202.202.202/32
      Description: !Sub "XMN VPN IPs, Created by CloudFormation Stack ${Name}"
      IPAddressVersion: IPV4
      Name: !Sub ${Name}-IPset
      Scope: REGIONAL
  
  Regex:
    Type: AWS::WAFv2::RegexPatternSet
    Properties: 
      Description: !Sub "Sentry Regex, Created by CloudFormation Stack ${Name}"
      Name: !Sub ${Name}-Regex
      RegularExpressionList: 
        - ^/api/\d/*
      Scope: REGIONAL

  WAF:
    Type: AWS::WAFv2::WebACL
    Properties: 
      DefaultAction: 
        Block: {} 
      Description: !Sub "A basic WAF, follow basic AWS Managed rules and a rate-limit rule, Created by CloudFormation Stack ${Name}"
      Name: BasicWAF
      Rules: 
        - Action: 
            Allow: {} 
          Name: AllowInternalIP
          Priority: 0
          Statement:
            IPSetReferenceStatement:
              Arn: !GetAtt IPSet.Arn
          VisibilityConfig: 
            CloudWatchMetricsEnabled: true
            MetricName: AllowInternalIP
            SampledRequestsEnabled: true
        - OverrideAction:
            None: {}
          Name: AWSManagedRulesCommonRuleSet
          Priority: 1
          Statement: 
            ManagedRuleGroupStatement:
              VendorName: AWS
              Name: AWSManagedRulesCommonRuleSet
              ExcludedRules: [NoUserAgent_HEADER,EC2MetaDataSSRF_BODY,GenericRFI_BODY]
          VisibilityConfig: 
            CloudWatchMetricsEnabled: true
            MetricName: AWSManagedRulesCommonRuleSet
            SampledRequestsEnabled: true
        - OverrideAction:
            None: {}
          Name: AWSManagedRulesAmazonIpReputationList
          Priority: 2
          Statement: 
            ManagedRuleGroupStatement:
              VendorName: AWS
              Name: AWSManagedRulesAmazonIpReputationList
              ExcludedRules: []
          VisibilityConfig: 
            CloudWatchMetricsEnabled: true
            MetricName: AWSManagedRulesAmazonIpReputationList
            SampledRequestsEnabled: true
        - Action: 
            Block: {} 
          Name: RateLimit
          Priority: 3
          Statement:
            RateBasedStatement:
              AggregateKeyType: IP
              Limit: 100
          VisibilityConfig: 
            CloudWatchMetricsEnabled: true
            MetricName: RateLimit
            SampledRequestsEnabled: true
        # Maybe can add "Regex pattern sets" when client post events
        - Action: 
            Allow: {} 
          Name: AllowPOSTMethod
          Priority: 4
          Statement:
            AndStatement:
              Statements:
                - ByteMatchStatement:
                    FieldToMatch:
                      Method: 
                        POST: {}
                    SearchString: post
                    PositionalConstraint: CONTAINS_WORD
                    TextTransformations:
                      - Priority: 0
                        Type: LOWERCASE
                - ByteMatchStatement:
                    FieldToMatch:
                      SingleHeader: 
                        Name: x-sentry-auth
                    SearchString: sentry
                    PositionalConstraint: CONTAINS_WORD
                    TextTransformations:
                      - Priority: 0
                        Type: LOWERCASE
                - RegexPatternSetReferenceStatement:
                    Arn: !GetAtt Regex
                    FieldToMatch: 
                      UriPath: {}
                    TextTransformations: 
                      - Priority: 0
                        Type: LOWERCASE
          VisibilityConfig: 
            CloudWatchMetricsEnabled: true
            MetricName: AllowPOSTMethod
            SampledRequestsEnabled: true
      Scope: REGIONAL
      VisibilityConfig: 
        CloudWatchMetricsEnabled: true
        MetricName: BasicWAF
        SampledRequestsEnabled: true

  WAFAssociate:
    Type: AWS::WAFv2::WebACLAssociation
    DependsOn: SentryElasticLoadBalancer
    Properties: 
      ResourceArn: !Ref SentryElasticLoadBalancer
      WebACLArn: !GetAtt WAF.Arn

```

#### 引用官方规则
```yaml
- OverrideAction:
    None: {}
  Name: AWSManagedRulesCommonRuleSet
  Priority: 0
  Statement: 
    ManagedRuleGroupStatement:
      VendorName: AWS
      Name: AWSManagedRulesCommonRuleSet
      ExcludedRules: []
  VisibilityConfig: 
    CloudWatchMetricsEnabled: true
    MetricName: AWSManagedRulesCommonRuleSet
    SampledRequestsEnabled: true

```

#### 允许指定IP或POST请求
```yaml
- Action: 
    Allow: {} 
  Name: OnlyPOSTExcludeXMNIP
  Priority: 3
  Statement:
    OrStatement: 
      Statements:
        - ByteMatchStatement:
            FieldToMatch:
              Method: 
                POST: {}
            SearchString: post
            PositionalConstraint: CONTAINS_WORD
            TextTransformations:
              - Priority: 0
                Type: LOWERCASE
        - IPSetReferenceStatement:
            Arn: !GetAtt IPSet.Arn
```

#### 除了指定IP外,其他IP受速率限制
```yaml
- Action: 
    Block: {} 
  Name: RateLimitExcludeXMNIP
  Priority: 2
  Statement:
    RateBasedStatement:
      AggregateKeyType: IP
      Limit: 100
      ScopeDownStatement:
        NotStatement:
          Statement:
            IPSetReferenceStatement:
              Arn: !GetAtt IPSet.Arn
  VisibilityConfig: 
    CloudWatchMetricsEnabled: true
    MetricName: RateLimitExcludeXMNIP
    SampledRequestsEnabled: true
```

#### 指定IP受速率限制
```yaml
- Action: 
    Block: {} 
  Name: RateLimitExcludeXMNIP
  Priority: 2
  Statement:
    RateBasedStatement:
      AggregateKeyType: IP
      Limit: 100
      ScopeDownStatement:
        IPSetReferenceStatement:
          Arn: !GetAtt IPSet.Arn
  VisibilityConfig: 
    CloudWatchMetricsEnabled: true
    MetricName: RateLimitExcludeXMNIP
    SampledRequestsEnabled: true
```

