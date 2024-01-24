[toc]

## 常见问题

### EC2设置密码
1. 设置密码
```shell
sudo passwd ubuntu
Changing password for user ubuntu.
New password:
Retype new password:
``` 

2. 更改配置
```shell
sudo vim /etc/ssh/sshd_config

PasswordAuthentication yes
```

3. 重启服务
```shell
sudo systemctl restart sshd.service
```

### Please login as the user "ec2-user" rather than the user "root".

```shell
[ec2-user@ip-10-4-7-51 ~]$ sudo cat /root/.ssh/authorized_keys
no-port-forwarding,no-agent-forwarding,no-X11-forwarding,command="echo 'Please login as the user \"ec2-user\" rather than the user \"root\".';echo;sleep 10" ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFw38SgXXnz+/2AEXS3Iyl04Sq55ersJfmDoor7gdqDVLV5yZ5YGHaSiu/kP/pHHBU7jMhjm6JpjR9xJ9iFp57tNl2eW8Sym02j8z/DI1qUrvpT9I1/sWGWlB80O0c6zfnLq9jEMNCj/0oNtqKluqYD28gD9o56staaj15VWxT9lPI/cjNL/eWoHOh++9DJwBmjtkOutMjjiTeHuM275n8rQkMDTaIfnXitnLGlqvHiQqpyt/6CTxPJLS5Or3HrLkyf3YgPVmhqCdzoBp+KStBnFmGk2aejuKwposPsOR6299CZgQdG8oBHHClQFFRo5IIXJWeKHGkeeDbBahQ7HiB myselfkey
```

解决方法:

```shell
[ec2-user@ip-10-4-7-51 ~]$ sudo cp /root/.ssh/authorized_keys /root/.ssh/authorized_keys.bak
[ec2-user@ip-10-4-7-51 ~]$ sudo cp .ssh/authorized_keys /root/.ssh/authorized_keys
```

### type: <class 'str'>, valid types: <class 'dict'>
```
aws ec2 modify-security-group-rules \
> --group-id sg-0c77d0b4e1f8ccdda \
> --security-group-rules SecurityGroupRuleId=sgr-077eed5ab7ba61131,SecurityGroupRule={IpProtocol=tcp,FromPort=22,ToPort=22,CidrIpv4=111.111.111.111/32,Description=user-home}

Parameter validation failed:
Invalid type for parameter SecurityGroupRules[0].SecurityGroupRule, value: IpProtocol=tcp, type: <class 'str'>, valid types: <class 'dict'>
Invalid type for parameter SecurityGroupRules[1].SecurityGroupRule, value: FromPort=22, type: <class 'str'>, valid types: <class 'dict'>
Invalid type for parameter SecurityGroupRules[2].SecurityGroupRule, value: ToPort=22, type: <class 'str'>, valid types: <class 'dict'>
Invalid type for parameter SecurityGroupRules[3].SecurityGroupRule, value: CidrIpv4=111.111.111.111/32, type: <class 'str'>, valid types: <class 'dict'>
Invalid type for parameter SecurityGroupRules[4].SecurityGroupRule, value: Description=user-home, type: <class 'str'>, valid types: <class 'dict'>
```

解决方案: 对`SecurityGroupRule`加引号
```shell
aws ec2 modify-security-group-rules \
> --group-id sg-0c77d0b4e1f8ccdda \
> --security-group-rules SecurityGroupRuleId=sgr-077eed5ab7ba61131,SecurityGroupRule="{IpProtocol=tcp,FromPort=22,ToPort=22,CidrIpv4=111.111.111.111/32}"
{
    "Return": true
}
```

