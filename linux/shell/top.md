# Top

- PID: 进程的ID
- USER: 进程属主的名字
- PR: 进程的优先级
- NI: 进程的谦让度值
- VIRT: 进程占用的虚拟内存总量
- RES: 进程占用的物理内存总量
- SHR: 进程和其他进程共享的内存总量
- S: 进程的状态（D代表可中断的休眠状态，R代表在运行状态，S代表休眠状态，T代表跟踪状态或停止状态，Z代表僵化状态）
- %CPU: 进程使用的CPU时间比例
- %MEM: 进程使用的内存占可用内存的比例
- TIME+: 自进程启动到目前为止的CPU时间总量
- COMMAND: 进程所对应的命令行名称，也就是启动的程序名

```bash
top         #每5秒显示所有进程的资源占用情况
top -d 2    #每2秒显示所有进程的资源占用情况
top -c      #每5秒显示进程的资源占用情况,并显示命令行参数(默认只有进程名)
top -p 12345 -p 6789    #每隔5秒显示pid是12345和pid是6789的两个进程的资源占用情况
top -d 2 -c -p 123456   #每隔2秒显示pid的进程的资源使用情况,并显示该进程命令行参数

# M 根据驻留内存大小进行排序 
# P 根据CPU使用百分比大小进行排序
```