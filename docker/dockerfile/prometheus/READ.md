## Command

```bash
docker-compose up -d
docker-compose restart # 可以重新加载配置
docker-compose down # 删除所有container
```

## 注意

### Grafana

在Grafana Web界面上配置中,Prometheus地址应写为 `http://prometheus:9090`