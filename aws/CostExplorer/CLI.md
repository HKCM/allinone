获取指定日期的 tag为int的cost
```shell
aws ce get-cost-and-usage \
    --time-period Start=2020-01-13,End=2020-01-14 \
    --granularity DAILY \
    --metric BlendedCost \
    --filter '{"Tags":{"Key":"Team","Values":["int"]}}'
```