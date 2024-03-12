# ping

```shell
ping -c 6 -i 3 url # ping6次每次间隔3秒

# 探测192.168.103.0/24网段中有多少可以通信的主机
ping -b -c 3 192.168.103.255
WARNING: pinging broadcast address
PING 192.168.103.255 (192.168.103.255) 56(84) bytes of data.
64 bytes from 192.168.103.199: icmp_seq=1 ttl=64 time=1.95 ms
64 bytes from 192.168.103.168: icmp_seq=1 ttl=64 time=1.97 ms (DUP! )
64 bytes from 192.168.103.252: icmp_seq=1 ttl=64 time=2.29 ms (DUP! )
```