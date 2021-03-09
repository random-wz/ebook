再我们平时使用 RHEL 系列 Linux 操作系统的时候，安装软件包通常需要安装一个 `epel-release` 的软件包，那这个包有什么用呢？

EPEL是由 Fedora 社区打造，为 RHEL 及衍生发行版如 CentOS、Scientific Linux 等提供高质量软件包的项目。装上了 EPEL之后，就相当于添加了一个第三方源。官方的rpm repository提供的rpm包也不够丰富，很多时候需要自己编译那太辛苦了，而EPEL恰恰可以解决这两方面的问题。

EPEL 安装方法：

```bash
[root@random ~]# yum -y install epel-release
```

我们平时有找不到的安装包也可以直接去 EPEL rpm包下载官网找：[epel-release rpm包下载地址](https://pkgs.org/download/epel-release) 。

阿里云也提供了丰富的 yum 源：[阿里云 rpm 包下载地址](https://developer.aliyun.com/packageSearch) 。