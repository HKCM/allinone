

## 描述: 


### AWS WAF 如何处理 Web ACL 中的规则操作
对于 Web ACL 中的规则,可以选择计数、允许或阻止匹配的 Web 请求:

* 允许和阻止操作是终止操作。它们会停止 Web ACL 对匹配的 Web 请求的所有其他处理。如果在 Web ACL 中包含具有允许或阻止操作的规则,并且该规则找到匹配项,该匹配项将为 Web ACL 确定对 Web 请求的最终处置。AWS WAF 不会处理 Web ACL 中位于匹配项之后的任何其他规则。对于直接添加到 Web ACL 的规则和添加的规则组中的规则,此原理同样适用。

* 计数操作是非终止操作。当具有计数操作的规则与请求匹配时,AWS WAF 会对请求进行计数,然后继续处理 Web ACL 规则集中的后续规则。如果唯一匹配的规则设置了计数操作,AWS WAF 会应用 Web ACL 默认操作设置。

AWS WAF 对 Web 请求应用的操作受规则在 Web ACL 中的相对位置的影响。例如,如果一个 Web 请求匹配允许请求的规则,同时又匹配另一个对请求进行计数的规则,那么如果允许请求的规则先列出,AWS WAF 不会对请求进行计数。

### 预配置的防护

您可以使用我们的预配置模板来快速开始使用 AWS WAF。该模板包含了一组旨在阻止常见的 Web 攻击的 AWS WAF 规则,这些规则可以根据您的需求进行自定义。这些规则可以帮助预防恶意自动程序、SQL 注入、跨站脚本 (XSS)、HTTP 泛洪以及已知攻击者的攻击。在您部署模板后,AWS WAF 将开始阻止发送到您的 CloudFront 分配的与 Web 访问控制列表 (Web ACL) 中预配置的规则匹配的 Web 请求。除您配置的其他 Web ACL 外,您还可以使用此自动化解决方案。

https://aws.amazon.com/cn/answers/security/aws-waf-security-automations/?refid=gs_card

### 阻止超过请求限制的 IP 地址

您可能面临的一个安全挑战是如何防止 Web 服务器受到分布式拒绝服务 (DDoS) 攻击 的影响,这种攻击又常称为 HTTP 泛洪。在此教程中,您将预置一个解决方案,此解决方案可以识别发送超过规定阈值的请求的 IP 地址,并更新您的 AWS WAF 规则以自动阻止来自这些 IP 地址的后续请求。
http://docs.aws.amazon.com/waf/latest/developerguide/tutorials-rate-based-blocking.html?refid=gs_card

### 阻止提交恶意请求的 IP 地址

面向 Internet 的 Web 应用程序经常会被各种来源扫描,除非您进行管理,否则这些来源可能会不怀好意。为了查找漏洞,这些扫描会发出一系列会生成 HTTP 4xx 错误代码的请求,您可以使用这些代码识别和阻止此类扫描。在此教程中,您将创建一个 Lambda 函数,以自动解释 CloudFront 访问日志,统计来自唯一来源(IP 地址)的恶意请求数量,并更新 AWS WAF 以阻止来自这些 IP 地址的进一步扫描。
http://docs.aws.amazon.com/waf/latest/developerguide/tutorials-4xx-blocking.html

### 使用坏人 IP 黑名单来预防 Web 攻击

AWS WAF 可以帮助您保护 Web 应用程序,防止源于由已知坏人操作的 IP 地址的刺探,例如垃圾邮件发送者、恶意软件分发者和僵尸网络等。在此教程中,您将学习如何将 AWS WAF 规则与声誉列表同步,从而阻止不断变化的 Web 攻击 IP 地址列表,掌握不断调换地址和企图逃避检测的坏人动向。
https://blogs.aws.amazon.com/security/post/Tx8GZBDD7HJ6BS/How-to-Import-IP-Address-Reputation-Lists-to-Automatically-Update-AWS-WAF-IP-Bla

