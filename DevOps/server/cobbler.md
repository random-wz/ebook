### 1. cobbler 基本概念

`Cobbler`的配置结构基于一组注册的对象。每个对象表示一个与另一个实体相关联的实体。当一个对象指向另一个对象时，它就继承了被指向对象的数据，并可覆盖或添加更多特定信息。

- 发行版(`distros`)： 表示一个操作系统。它承载了内核和`initrd`的信息，以及内核参数等其他数据。
- 配置文件(`profiles`)：包含一个发行版、一个`kickstart`文件以及可能的存储库，还包括更多特定的内核参数等其他数据。
- 系统(`systems`)：表示要配给的机器。它包括一个配置文件或一个镜像、`IP`和`MAC`地址、电源管理（地址、凭据、类型）以及更为专业的数据等信息。
- 镜像(`images`)：可以替换一个包含不屑于此类别的文件的发行版对象（例如，无法分为内核和`initrd`的对象）。

### 2. cobbler集成的服务

- PXE服务支持
- DHCP服务管理
- DNS服务管理
- 电源管理
- Kickstart服务支持
- YUM仓库管理
- TFTP
- Apache

### 3. 安装 cobbler 的两种方式

#### 3.1 rpm 包安装

环境准备：

```bash
# 关闭防火墙、selinux等
[root@cobbler ~]# systemctl stop firewalld
[root@cobbler ~]# systemctl disable firewalld
[root@cobbler ~]# setenforce 0
[root@cobbler ~]# sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/sysconfig/selinux
```

#### 安装 cobbler：

```bash
# 配置epel源
[root@cobbler ~]# yum -y install epel-release

# 安装cobbler及dhcp httpd xinetd cobbler-web
[root@cobbler ~]# yum -y install cobbler cobbler-web tftp-server dhcp httpd xinetd

# 启动cobbler及httpd并加入开机启动
[root@cobbler ~]# systemctl start httpd cobblerd
[root@cobbler ~]# systemctl enable httpd cobblerd
```

配置 cobbler：检查`Cobbler`的配置，如果看不到下面的结果，再次重启`cobbler`

```bash
[root@cobbler ~]# cobbler check
The following are potential configuration items that you may want to fix:

1 : The 'server' field in /etc/cobbler/settings must be set to something other than localhost, or kickstarting features will not work.  This should be a resolvable hostname or IP for the boot server as reachable by all machines that will use it.
2 : For PXE to be functional, the 'next_server' field in /etc/cobbler/settings must be set to something other than 127.0.0.1, and should match the IP of the boot server on the PXE network.
3 : change 'disable' to 'no' in /etc/xinetd.d/tftp
4 : Some network boot-loaders are missing from /var/lib/cobbler/loaders, you may run 'cobbler get-loaders' to download them, or, if you only want to handle x86/x86_64 netbooting, you may ensure that you have installed a *recent* version of the syslinux package installed and can ignore this message entirely.  Files in this directory, should you want to support all architectures, should include pxelinux.0, menu.c32, elilo.efi, and yaboot. The 'cobbler get-loaders' command is the easiest way to resolve these requirements.
5 : enable and start rsyncd.service with systemctl
6 : debmirror package is not installed, it will be required to manage debian deployments and repositories
7 : ksvalidator was not found, install pykickstart
8 : The default password used by the sample templates for newly installed machines (default_password_crypted in /etc/cobbler/settings) is still set to 'cobbler' and should be changed, try: "openssl passwd -1 -salt 'random-phrase-here' 'your-password-here'" to generate new one
9 : fencing tools were not found, and are required to use the (optional) power management features. install cman or fence-agents to use them

Restart cobblerd and then run 'cobbler sync' to apply changes.
```

看到上面出现的问题，然后一个一个的进行解决，可以动态配置，也可以直接更改配置文件。

```bash
1. server 注意：server 后面的IP地址必须是cobbler所在服务器的IP地址
[root@cobbler ~]# cobbler setting edit --name=server --value=192.168.2.128

2. next_server 注意：next_server 后面的IP地址必须是cobbler所在服务器的IP地址
[root@cobbler ~]# cobbler setting edit --name=next_server --value=192.168.2.128

3. tftp_server
[root@cobbler ~]# sed -ri '/disable/c\disable = no' /etc/xinetd.d/tftp
[root@cobbler ~]# systemctl enable xinetd
[root@cobbler ~]# systemctl restart xinetd

4. boot-loaders
[root@cobbler ~]# cobbler get-loaders

5. rsyncd
[root@cobbler ~]# systemctl start rsyncd
[root@cobbler ~]# systemctl enable rsyncd

6. debmirror [optional]
# 这个是可选项的，可以忽略。这里就忽略了

7. pykickstart
[root@cobbler ~]# yum -y install pykickstart

8. default_password_crypted  #注意：这里设置的密码，也就是后面安装完系统的初始化登录密码
[root@cobbler ~]# openssl passwd -1 -salt `openssl rand -hex 4` 'admin'
$1$675f1d08$oJoAMVxdbdKHjQXbGqNTX0
[root@cobbler ~]# cobbler setting edit --name=default_password_crypted --value='$1$675f1d08$oJoAMVxdbdKHjQXbGqNTX0'

9. fencing tools [optional]
[root@cobbler ~]# yum -y install fence-agents
```

同步 cobbler 配置：

```bash
[root@cobbler ~]# cobbler sync
```

#### 3.2 通过 docker 容器的方式安装



### 4. cobbler CLI

| 命令             | 说明                                       |
| ---------------- | ------------------------------------------ |
| cobbler check    | 核对当前设置是否有问题                     |
| cobbler list     | 列出所有的cobbler元素                      |
| cobbler report   | 列出元素的详细信息                         |
| cobbler sync     | 同步配置到数据目录，更改配置最好都执行一下 |
| cobbler reposync | 同步yum仓库                                |
| cobbler distro   | 查看导入的发行版系统信息                   |
| cobbler system   | 查看添加的系统信息                         |
| cobbler profile  | 查看配置信息                               |

