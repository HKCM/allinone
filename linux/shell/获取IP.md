# 获取IP

```shell
curl icanhazip.com
curl ifconfig.me
curl http://checkip.amazonaws.com
wget http://ipecho.net/plain -O - -q

# EC2
curl http://169.254.169.254/latest/meta-data/local-ipv4 # Get private IPv4
curl http://169.254.169.254/latest/meta-data/public-ipv4 # Get public IPv4
```

# 查询域名

```shell
dig @8.8.8.8 www.baidu.com +short
```