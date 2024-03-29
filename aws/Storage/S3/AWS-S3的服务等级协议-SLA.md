描述: AWS S3 服务等级协议(SLA)


首先给出官方原版文件:

[AWS S3 服务等级协议](https://amazonaws-china.com/s3/sla/?nc1=h_ls) 
[AWS S3 服务等级协议中文版](https://d1.awsstatic-china.com/legal/amazons3service/Amazon-S3-Service-Level-Agreement-Chinese.pdf)

简单解释一下

当AWS的S3服务的每月正常运行时间低于某个值时,你可以申请AWS的S3服务积分,S3服务积分可以抵扣AWS的**S3**账单相当于S3优惠卷.
正常运行时间比|服务积分|大致的达成条件
-|-|-
99.9%-99%|10%|30天内平均每1000个请求至少失败一个或累计44分钟无法正常工作
99%-95%|25%|30天内平均每100个请求至少失败一个或累计7个半小时无法正常工作
<95%|100%|30天内平均每20个请求至少失败一个或累计36个小时无法正常工作


之所以是大致达成条件是因为AWS的条件限制以5分钟内的平均失败率为准

举个例子:

假设,某月(30天)AWS S3持续故障一小时,并且公司A在这一小时内受到影响,导致完全无法使用S3,月内其余时间S3正常运行。
* 情况1:

这一小时内,公司A每5分钟访问S3一或多次,并且都得到error消息
计算: 
在这一个小时内,一共有12次100%的错误率
60 ➗ 5 = 12
整体错误率 = 平均每5分钟的错误率 = 每个5分钟错误率之和 / 30天内5分钟的数量 
(12 * 100% + 0%【其余时间没有错误】)➗ (30 ✖️ 24 ✖️ 60 ➗ 5)= 0.139% 
正常运行时间比
100% - 0.139% = 99.861% 

**结论:可以申请服务积分**

* 情况2:

这一小时内,公司A已知S3故障只有第一个5分钟内访问过S3一或多次,并且得到error消息。其余55分钟没有访问。 
如果5分钟内没有任何请求,则该间隔错误率为0.

计算: 
在这一个小时内,只有第一个5分钟是100%的错误率,其余时间没有访问所以错误率为0. 
5 ➗ 5 = 1   
整体错误率 = 平均每5分钟的错误率 = 每个5分钟错误率之和 / 30天内5分钟的数量 
(1 * 100% + 0%【其余时间没有错误】)➗ (30 ✖️ 24 ✖️ 60 ➗ 5)= 0.0116% 
正常运行时间比 
100% - 0.0116% = 99.9884%  

**结论:无法申请服务积分**

以上只是申请理赔的必要条件,简单来说申请服务积分还需要

1. 必须在故障发生的第二个账单周期结束前通过AWS support center提交申请  
一般理解为事故发生的两个月内,例如1月15日发生的故障必须在3月之前提交申请
2. 申请的主题必须包含"SLA Credit Request"字样
3. 带上发生故障AWS的region,发生的时间日期和发生次数
4. 发生错误的日志 
日志内的敏感信息以*号代替

S3的免责条款大致内容

1. 不可抗力和网络接入 
不可抗力很好理解,个人理解"网络接入"可能是指防火墙之类的
2. 用户或第三方自愿作为或不作为 
作为是指做了禁止的事,不作为指有义务做但是没做的事,例如需要升级CLI才能访问S3但是用户不升级导致无法访问
3. 用户或第三方设备软件或其他技术导致
4. AWS协议中止或终止用户使用S3的权利导致 
类似删号和封号

