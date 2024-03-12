netstat命令各个参数说明如下:
```
-a : 列出所有网络状态，包括Socket程序
-c : 指定每隔几秒刷新一次网络状态
-r : 显示路由表
-t : 指明显示TCP端口
-u : 指明显示UDP端口
-l : 仅显示监听套接字(所谓套接字就是使应用程序能够读写与收发通讯协议(protocol)与资料的程序)
-p : 显示进程标识符和程序名称,每一个套接字/端口都属于一个程序。
-n : 不进行DNS轮询(可以加速操作),使用IP地址和端口号显示，不使用域名与服务名
```


- Proto：网络连接的协议，一般就是TCP协议或者UDP协议。
- Recv-Q：表示接收到的数据，已经在本地的缓冲中，但是还没有被进程取走。
- Send-Q：表示从本机发送，对方还没有收到的数据，依然在本地的缓冲中，一般是不具备ACK标志的数据包。
- Local Address：本机的IP地址和端口号。
- Foreign Address：远程主机的IP地址和端口号。
- State：状态。常见的状态主要有以下几种:
  - LISTEN：监听状态，只有TCP协议需要监听，而UDP协议不需要监听。
  - ESTABLISHED：已经建立连接的状态。如果使用“-l”选项，则看不到已经建立连接的状态。
  - SYN_SENT:SYN发起包，就是主动发起连接的数据包。
  - SYN_RECV：接收到主动连接的数据包。- FIN_WAIT1：正在中断的连接
  - FIN_WAIT2：已经中断的连接，但是正在等待对方主机进行确认
  - TIME_WAIT：连接已经中断，但是套接字依然在网络中等待结束
  - CLOSED：套接字没有被使用


即可显示当前服务器上所有端口及进程服务,于grep结合可查看某个具体端口及服务情况··
```bash
[root@localhost ~]# netstat -nlp |grep LISTEN   //查看当前所有监听端口·
[root@localhost ~]# netstat -nlp |grep 80   //查看所有80端口使用情况·
[root@localhost ~]# netstat -an | grep 3306   //查看所有3306端口使用情况·
```

ss -lnt