
## SFTP选择

### Serverlsee SFTP

https://amazonaws-china.com/aws-transfer-family/pricing/?nc1=h_ls
Time SFTP is enabled on your endpoint	$0.30 per hour
SFTP data uploads	                    $0.04 per gigabyte (GB) transferred
SFTP data downloads	                    $0.04 per gigabyte (GB) transferred

Example 1: Single protocol use

Assume you set up your endpoint and enable it for SFTP access. You then configure 20 end users to access data stored in your Amazon S3 bucket over SFTP. In total, your users download 1 GB/day of data from your Amazon S3 bucket. We calculate your monthly AWS Transfer Family costs using pricing in the US-East-1 Region as follows:

SFTP enabled on your endpoint: At $0.30 hourly rate, your monthly charge for SFTP is:

$0.30 * 24 hours * 30 days = $216

SFTP data upload and download: At $0.04/GB, your monthly charge for data uploads and downloads over SFTP is:

$0.04 * 1 GB * 30 days = $1.20

Adding the charges above, your total monthly bill for the AWS Transfer Family would be:

$216 + $1.20 = $217.20.

**总结:**
简单总结就是AWS官方的SFTP或者是说transfer family价格太高单纯服务器每个月就是216刀

### EC2 SFTP Server

Type t3.micro machines are enough

https://amazonaws-china.com/cn/ec2/pricing/on-demand/

```
t3.micro	                            $0.0104 per hour
SFTP data uploads	                    $0.00 per gigabyte (GB) transferred
SFTP data downloads	                    $0.09 per gigabyte (GB) transferred
```

SFTP enabled on your EC2: At $0.0104 hourly rate, your monthly charge for EC2 is:

$0.0104 * 24 hours * 30 days = $7.488

SFTP data upload and download: At $0.09/GB, your monthly charge for data uploads and downloads over SFTP is:
```
$0.09 * 1 GB * 30 days = $2.70
```
Adding the charges above, your total monthly bill for the EC2 SFTP would be:
```
$7.488 + $2.70 = $10.188
```
**总结:**
简单总结就是使用EC2搭建SFTP会更加便宜

### 简单比较优缺点

AWS SFTP优点
1. 利用S3作为数据存储,可靠性高
2. Serverless架构,无需管理
3. 数据传输双向收费但单价低
4. 支持与API Gateway结合进行身份认证
5. 无存储上限,可以无限存储

AWS SFTP缺点
1. $0.3每小时月单价达到$216

EC2优点
1. $0.0104每小时月单价达到$7.28
2. EC2随起随停,可以进一步降低费用
3. 数据传输单向收费(流出EC2流量)但单价高
4. 可以SSH public key和用户名密码认证

EC2 SFTP缺点
1. EBS块存储可用性不如S3
2. 需要管理底层EC2
3. 存储容量有限不如S3

## Ubuntu搭建SFTP

Refer:

https://www.digitalocean.com/community/tutorials/how-to-enable-sftp-without-shell-access-on-ubuntu-18-04

### Introduction
SFTP stands for SSH File Transfer Protocol. As its name suggests, it’s a secure way to transfer files between machines using an encrypted SSH connection. Despite the name, it’s a completely different protocol than FTP (File Transfer Protocol), though it’s widely supported by modern FTP clients.

SFTP is available by default with no additional configuration on all servers that have SSH access enabled. It’s secure and easy to use, but comes with a disadvantage: in a standard configuration, the SSH server grants file transfer access and terminal shell access to all users with an account on the system.

In some cases, you might want only certain users to be allowed file transfers and no SSH access. In this tutorial, we’ll set up the SSH daemon to limit SFTP access to one directory with no SSH access allowed on per-user basis.

#### Prerequisites
To follow this tutorial, you will need access to an Ubuntu 18.04 server. This server should have a non-root user with sudo privileges, as well as a firewall enabled. For help with setting this up, follow our Initial Server Setup Guide for Ubuntu 18.04.

#### Step 1 — Creating a New User
First, create a new user who will be granted only file transfer access to the server. Here, we’re using the username sammyfiles, but you can use any username you like.

```
sudo adduser sammyfiles
```
You’ll be prompted to create a password for the account, followed by some information about the user. The user information is optional, so you can press ENTER to leave those fields blank.

You have now created a new user that will be granted access to the restricted directory. In the next step we will create the directory for file transfers and set up the necessary permissions.

#### Step 2 — Creating a Directory for File Transfers
In order to restrict SFTP access to one directory, we first have to make sure the directory complies with the SSH server’s permissions requirements, which are very particular.

Specifically, the directory itself and all directories above it in the filesystem tree must be owned by root and not writable by anyone else. Consequently, it’s not possible to simply give restricted access to a user’s home directory because home directories are owned by the user, not root.

Note: Some versions of OpenSSH do not have such strict requirements for the directory structure and ownership, but most modern Linux distributions (including Ubuntu 18.04) do.

There are a number of ways to work around this ownership issue. In this tutorial, we’ll create and use /var/sftp/uploads as the target upload directory. /var/sftp will be owned by root and will not be writable by other users; the subdirectory /var/sftp/uploads will be owned by sammyfiles, so that user will be able to upload files to it.

First, create the directories.

sudo mkdir -p /var/sftp/uploads
Set the owner of /var/sftp to root.

sudo chown root:root /var/sftp
Give root write permissions to the same directory, and give other users only read and execute rights.

sudo chmod 755 /var/sftp
Change the ownership on the uploads directory to sammyfiles.

sudo chown sammyfiles:sammyfiles /var/sftp/uploads
Now that the directory structure is in place, we can configure the SSH server itself.

#### Step 3 — Restricting Access to One Directory
In this step, we’ll modify the SSH server configuration to disallow terminal access for sammyfiles but allow file transfer access.

Open the SSH server configuration file using nano or your favorite text editor.

sudo nano /etc/ssh/sshd_config
Scroll to the very bottom of the file and append the following configuration snippet:

/etc/ssh/sshd_config
. . .

Match User sammyfiles
ForceCommand internal-sftp
PasswordAuthentication yes
ChrootDirectory /var/sftp
PermitTunnel no
AllowAgentForwarding no
AllowTcpForwarding no
X11Forwarding no
Then save and close the file.

Here’s what each of those directives do:

Match User tells the SSH server to apply the following commands only to the user specified. Here, we specify sammyfiles.
ForceCommand internal-sftp forces the SSH server to run the SFTP server upon login, disallowing shell access.
PasswordAuthentication yes allows password authentication for this user.
ChrootDirectory /var/sftp/ ensures that the user will not be allowed access to anything beyond the /var/sftp directory.
AllowAgentForwarding no, AllowTcpForwarding no. and X11Forwarding no disables port forwarding, tunneling and X11 forwarding for this user.
This set of commands, starting with Match User, can be copied and repeated for different users too. Make sure to modify the username in the Match User line accordingly.

Note: You can omit the PasswordAuthentication yes line and instead set up SSH key access for increased security. Follow the Copying your Public SSH Key section of the SSH Essentials: Working with SSH Servers, Clients, and Keys tutorial to do so. Make sure to do this before you disable shell access for the user.

In the next step, we’ll test the configuration by SSHing locally with password access, but if you set up SSH keys, you’ll instead need access to a computer with the user’s keypair.

To apply the configuration changes, restart the service.

sudo systemctl restart sshd
You have now configured the SSH server to restrict access to file transfer only for sammyfiles. The last step is testing the configuration to make sure it works as intended.

#### Step 4 — Verifying the Configuration
Let’s ensure that our new sammyfiles user can only transfer files.

Logging in to the server as sammyfiles using normal shell access should no longer be possible. Let’s try it:

ssh sammyfiles@localhost
You’ll see the following message before being returned to your original prompt:

Error message
This service allows sftp connections only.
Connection to localhost closed.
This means that sammyfiles can no longer can access the server shell using SSH.

Next, let’s verify if the user can successfully access SFTP for file transfer.

sftp sammyfiles@localhost
Instead of an error message, this command will show a successful login message with an interactive prompt.

SFTP prompt
Connected to localhost.
sftp>
You can list the directory contents using ls in the prompt:

ls
This will show the uploads directory that was created in the previous step and return you to the sftp> prompt.

SFTP file list output
uploads
To verify that the user is indeed restricted to this directory and cannot access any directory above it, you can try changing the directory to the one above it.

cd ..
This command will not give an error, but listing the directory contents as before will show no change, proving that the user was not able to switch to the parent directory.

You have now verified that the restricted configuration works as intended. The newly created sammyfiles user can access the server only using the SFTP protocol for file transfer and has no ability to access the full shell.

Conclusion
You’ve restricted a user to SFTP-only access to a single directory on a server without full shell access. While this tutorial uses only one directory and one user for brevity, you can extend this example to multiple users and multiple directories.

The SSH server allows more complex configuration schemes, including limiting access to groups or multiple users at once, or even limited access to certain IP addresses. You can find examples of additional configuration options and explanation of possible directives in the OpenSSH Cookbook. If you run into any issues with SSH, you can debug and fix them with this troubleshooting SSH series.


### CentOS系统记录sftp log
```shell
vim /etc/ssh/sshd.conf

#这行指定使用sftp服务使用系统自带的internal-sftp
Subsystem sftp internal-sftp -l info -f local5
#这行用来匹配sftpusers组的用户,如果要匹配多个组,多个组之间用逗号分割； 
Match Group sftp     
    #用chroot将用户的根目录指定到%h,%h代表用户home目录,这样用户就只能在用户目录下活动。也可用%u,%u代表用户名。
    ChrootDirectory /data/sftp/%u  
    #指定sftp命令       
    ForceCommand internal-sftp    
    AllowTcpForwarding no
    X11Forwarding no

#匹配用户了,多个用户名之间也是用逗号分割
Match user sftp-user        
    ChrootDirectory /data/sftp
    ForceCommand internal-sftp
    AllowTcpForwarding no
    X11Forwarding no
```

```shell
vim /etc/rsyslog.conf
auth,authpriv.*,local5.*      /var/log/sftp.log
```

## SFTP命令

Linux SFTP 子系统支持以下功能,它们都默认启用。执行以下命令以查看所有支持的协议功能列表。
```shell
# /usr/libexec/openssh/sftp-server -Q requests
open
close
read
write
lstat
fstat
setstat
fsetstat
opendir
readdir
remove
mkdir
rmdir
realpath
stat
rename
readlink
symlink
posix-rename
statvfs
fstatvfs
hardlink
fsync
```

在 SFTP 服务器系统的文件中`/etc/ssh/sshd_config`添加以下配置,以禁用所有 SFTP 用户的rename和delete权限,然后重新启动 sshd 服务。

```shell
# For normal SFTP(Without Chroot) add the following configuration in sshd_config file. 
Subsystem       sftp    /usr/libexec/openssh/sftp-server -P remove,rmdir,rename,posix-rename

# For Chroot SFTP  add the following configuration in sshd_config file.
Subsystem sftp internal-sftp -P remove,rmdir,rename,posix-rename

# chroot匹配块的示例配置
Match Group sftp_users
        ChrootDirectory /data/sftp
        ForceCommand internal-sftp  -P remove,rmdir,rename,posix-rename
        X11Forwarding no
        AllowTcpForwarding no
```

参考
```shell
$ man sftp
$ man sftp-server
```


