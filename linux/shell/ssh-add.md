# ssh-add

`ssh-add`命令用来将私钥加入`ssh-agent`

```bash
ssh-add -l # 列出所有已经添加的私钥
ssh-add -d name-of-key-file # 从内存中删除指定的私钥
ssh-add -D # 从内存中删除所有已经添加的私钥

SSH_KEY_ALREADY_ADDED="`ssh-add -l | grep -v grep | grep 'ec2.pem' | wc -l`"
echo ""
if [ "$SSH_KEY_ALREADY_ADDED" = "0" ]; then
    echo "add ec2 ssh key to agent..."
    ssh-add ~/.ssh/ec2.pem
else
    echo "prod ec2 ssh key already added to ssh agent"
fi
```