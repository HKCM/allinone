# passwd

```bash
passwd username # 修改用户密码
```

```bash
echo -e "badpass\nbadpass" | passwd ttt # 可以在docker中使用

echo "username:password" | chpasswd # 在/bin/sh和/bin/bash中都可以使用
```

```bash
#!/bin/bash
read -p "Please Enter Your Real Name: " REAL_NAME 
read -p "Please Enter Your User Name: " USER_NAME 
useradd -c "${COMMENT}" -m ${USER_NAME} 
read -p "Please Enter Your Password: " PASSWORD
echo -e "$PASSWORD\n$PASSWORD" |passwd "$USER_NAME"
passwd -e ${USER_NAME}
```