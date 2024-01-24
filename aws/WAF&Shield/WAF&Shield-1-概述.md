



### WAF-Shield-FirewallManager概述

#### AWS WAF
AWS WAF 是一种 Web 应用程序防火墙,可用于监控转发到 Amazon CloudFront,Amazon API Gateway REST API、ALB 或 AWS AppSync GraphQL API 的 HTTP 和 HTTPS 请求。

* 允许指定的请求以外的所有请求
* 阻止指定的请求之外的所有请求

* 使用指定的条件针对 Web 攻击提供额外保护:
  * 请求源自的 IP 地址.
  * 请求源自的国家/地区。
  * 请求标头中的值.
  * 出现在请求中的字符串 (特定字符串或与正则表达式 (regex) 模式匹配的字符串)。
  * 请求的长度.
  * 存在可能是恶意的 SQL 代码 (称为 SQL 注入).
  * 存在可能是恶意的脚本 (称为跨站点脚本).
* 规则可以允许、阻止或统计满足指定条件的 Web 请求。或者,规则可以阻止或统计不仅满足指定条件,还在任何 5 分钟周期内超过指定请求数的 Web 请求。
* 可以重复用于多个 Web 应用程序的规则.
* 来自 AWS 和 AWS Marketplace 卖家的托管规则组。
* 实时指标和采样的 Web 请求.
* 使用 AWS WAF API 的自动化管理。

[WAF 定价](https://aws.amazon.com/cn/waf/pricing/)

#### AWS Shield

可以使用 AWS WAF web访问控制列表(web ACLs)帮助将分布式拒绝服务(DDoS)攻击的影响降至最低。**为防止 DDoS 攻击, AWS 还提供 AWS Shield Standard 和 AWS Shield Advanced**。 AWS Shield Standard 自动包含在内,无需超出您已经支付的费用 AWS WAF 和你的其他人 AWS 服务。 AWS Shield Advanced 提供扩展的 DDoS 攻击保护 Amazon EC2 实例, Elastic Load Balancing 负载均衡器, CloudFront 分布, Route 53 托管区域,以及 AWS Global Accelerator 加速器。 AWS Shield Advanced 会产生额外费用。

#### AWS Firewall Manager
AWS Firewall Manager 可针对 AWS WAF 规则、AWS Shield Advanced 保护和 Amazon VPC 安全组,简化跨多个账户和多种资源的管理和维护任务。即使您添加新的账户和资源,Firewall Manager 服务也会自动跨账户和资源应用您的规则和其他安全保护。

### Open Web Application Security Project

https://owasp.org/www-project-top-ten/

https://wiki.owasp.org/images/d/dc/OWASP_Top_10_2017_%E4%B8%AD%E6%96%87%E7%89%88v1.3.pdf

