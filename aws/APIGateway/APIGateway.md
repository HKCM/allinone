### API Gateway Pricing

https://aws.amazon.com/cn/api-gateway/pricing/


# Private API Gateway Demo
https://aws.amazon.com/cn/blogs/compute/introducing-amazon-api-gateway-private-endpoints/

### 为API Gateway配置域名

1. 需要在API gateway上创建一个自定义域名(只是创建,只是一个声明作用),这里需要有域名证书
https://us-east-1.console.aws.amazon.com/apigateway/main/publish/domain-names
2. 然后在这个自定义域名中添加API 映射 就是你创建的API
https://us-east-1.console.aws.amazon.com/apigateway/main/publish/domain-names/api-mappings
1. 在Route53中就可以使用A记录加Alias 锁定到API Gateway中新创建的域名了