# dd

keyword: 生成数据 生成文件 生成大文件

在当前目录下会生成一个1000M的test文件，文件内容为全0（因从/dev/zero中读取，/dev/zero为0源），但是这样为实际写入硬盘，文件产生速度取决于硬盘读写速度，如果欲产生超大文件，速度很慢。
```bash
dd if=/dev/zero of=test bs=1M count=1000
```

在某种场景下，只想让文件系统认为存在一个超大文件在此，但是并不实际写入硬盘则可以
```bash
dd if=/dev/zero of=test bs=1M count=0 seek=100000
```

此时创建的文件在文件系统中的显示大小为100000MB，但是并不实际占用block，因此创建速度与内存速度相当，seek的作用是跳过输出文件中指定大小的部分，这就达到了创建大文件，但是并不实际写入的目的。当然，因为不实际写入硬盘，所以在容量只有10G的硬盘上创建100G的此类文件都是可以的。

```bash
dd if=/dev/random of=random_file bs=10k count=1
dd if=/dev/urandom of=random_file bs=10k count=1
dd if=/dev/zero of=random_file bs=10k count=1
```