[toc]

> InfluxDB是一个时间序列数据库，旨在处理较高的写入和查询负载，因此它常被用来存储监控数据，IoT行业的实时数据等。influxDB是用Go语言写的开源项目 ，点击[这里](git@github.com:influxdata/influxdb.git)查看源码。

#### 1. 简介及安装

influxDB默认开通的端口号是`8086`，我们可以通过访问API接口进行数据库的写入和查询，目前有三种途径：

- CLI 本地安装的数据库可以直接通过influx命令行工具连接数据库（类似于mysql命令），连接之后我们可以通过influxQL进行数据读写。
- curl 直接通过HTTP请求进行数据库的读写。
- sdk influxDB提供了多种语言的sdk，在Gthub上可以找到，其实本质上也是通过HTTP请求的方式进行数据库的读写。

influxDB数据库的安装也是十分方便，官方文档中给出了各种Linux操作系统的[安装手册](https://docs.influxdata.com/influxdb/v1.8/introduction/install/)，这里只介绍如何在Centos&Red Hat 操作系统中安装influxDB数据库。

因为influx数据库通过`8086`端口提供API访问，通过`8088`端口提供RPC服务执行备份和还原操作，所以安装前需要确保这两个端口可用，或者修改influxDB数据库的配置文件：`/etc/influxdb/influxdb.conf`。另外因为influxDB数据库是时序数据库，因此对服务器的时间准确性要求很高，默认情况下influxDB使用服务器的UTC时间给数据打时间戳，多台服务器通过NTP服务同步时间，因此在influxDB数据库集群中一定要确保NTP服务器正常运行，确保时间准确。

上面两个条件OK后我们开始安装，首先配置influxDB数据库的yum源：

```bash
cat <<EOF | sudo tee /etc/yum.repos.d/influxdb.repo
[influxdb]
name = InfluxDB Repository - RHEL \$releasever
baseurl = https://repos.influxdata.com/rhel/\$releasever/\$basearch/stable
enabled = 1
gpgcheck = 1
gpgkey = https://repos.influxdata.com/influxdb.key
EOF
```

yum源配置成功后就可以安装influxDB了：

```bash
 yum install -y influxdb
```

启动服务：

```bash
systemctl start influxdb
```

influxDB的默认配置文件为：`/etc/influxdb/influxdb.conf`，如果不想使用默认的配置文件，我们就可以使用`-config`命令指定配置文件，如下例：

```bash
[root@random ~]# influxd -config /etc/influxdb/influxdb.conf

 8888888           .d888 888                   8888888b.  888888b.
   888            d88P"  888                   888  "Y88b 888  "88b
   888            888    888                   888    888 888  .88P
   888   88888b.  888888 888 888  888 888  888 888    888 8888888K.
   888   888 "88b 888    888 888  888  Y8bd8P' 888    888 888  "Y88b
   888   888  888 888    888 888  888   X88K   888    888 888    888
   888   888  888 888    888 Y88b 888 .d8""8b. 888  .d88P 888   d88P
 8888888 888  888 888    888  "Y88888 888  888 8888888P"  8888888P"

2020-10-26T14:29:37.285419Z     info    InfluxDB starting       {"log_id": "0Q5KseXG000", "version": "1.8.3", "branch": "1.8", "commit": "563e6c3d1a7a2790763c6289501095dbec19244e"}
```

我们也可以通过配置环境变量的方式指定配置文件：

```bash
[root@random ~]# echo $INFLUXDB_CONFIG_PATH
/etc/influxdb/influxdb.conf
[root@random ~]# influxd

 8888888           .d888 888                   8888888b.  888888b.
   888            d88P"  888                   888  "Y88b 888  "88b
   888            888    888                   888    888 888  .88P
   888   88888b.  888888 888 888  888 888  888 888    888 8888888K.
   888   888 "88b 888    888 888  888  Y8bd8P' 888    888 888  "Y88b
   888   888  888 888    888 888  888   X88K   888    888 888    888
   888   888  888 888    888 Y88b 888 .d8""8b. 888  .d88P 888   d88P
 8888888 888  888 888    888  "Y88888 888  888 8888888P"  8888888P"

2020-10-26T14:29:37.285419Z     info    InfluxDB starting       {"log_id": "0Q5KseXG000", "version": "1.8.3", "branch": "1.8", "commit": "563e6c3d1a7a2790763c6289501095dbec19244e"}
```

到此单机版influxDB就安装完成了。

<font color=red size=3>注意：influxDB会首先检查`-config`命令，然后是环境变量。</font>

#### 2. 数据格式定义

indluxDB数据库名词解释：

| 名称        | 含义                                                         |
| ----------- | ------------------------------------------------------------ |
| database    | 数据库，与Mysql中的Database类似。                            |
| measurement | 类似于数据库中的表，influxDB通过时间序列保存数据，每个时间点的数据由时间戳、measurement、至少一组field以及零到多个tags组成。我们可以将measurement看成数据库的一张表，并且他的主键是时间戳time。 |
| field       | 键值对，influxDB中保存的数据中至少包含一组键值对，这里是指被测量的值本身，比如“value=0.64”、 “temperature=21.2”，需要注意的是不能通过field进行检索。 |
| tag         | 标签，tag中保存的都是一些元数据，比如：host=server01”, “region=EMEA”, “dc=Frankfurt”，可以通过tag进行检索。 |

了解了专业名词含义，我们来看一下在influxDB中是如何保存数据的，在influxDB中定义了一套数据保存规则，其中规定influxDB每一行数据的格式如下：

`<measurement>[,<tag-key>=<tag-value>...] <field-key>=<field-value>[,<field2-key>=<field2-value>...] [unix-nano-timestamp]`



#### 3.  使用指导

安装完数据库后，我们就可以通过influx命令连接数据库了，我们看一下influx常用命令：

```bash
[root@localhost ~]# influx --help
Usage of influx:
  -version   查看当前版本
  -path-prefix 'url path'  指定host和port的访问前缀，比如域名前缀

  -host 'host name' 指定访问服务器的主机名
  -port 'port #' 指定连接的端口，默认8086
  -socket 'unix domain socket' 通过socket连接数据库
  -database 'database name' 指定连接的数据库名称
  -password 'password' 指定数据库密码，默认为空
  -username 'username' 指定用户名
  -ssl 使用https连接
  -execute 'command' 执行命令后退出
  -pretty json格式打印时以更易读的方式打印。
  -format 'json|csv|column' 指定数据格式化类型json, csv, or column.
  -precision 'rfc3339|h|m|s|ms|u|ns' 时间戳:  rfc3339, h, m, s, ms, u or ns.
  -import 从文件中导入数据库
Examples:

    # Use influx in a non-interactive mode to query the database "metrics" and pretty print json:
    $ influx -database 'metrics' -execute 'select * from cpu' -format 'json' -pretty

    # Connect to a specific database on startup and set database context:
    $ influx -database 'metrics' -host 'localhost' -port '8086'
```

**（1）连接数据库：**

如果没有配置用户名密码，直接使用influx命令连接数据库：

```bash
[root@localhost ~]# influx
Connected to http://localhost:8086 version 1.8.3
InfluxDB shell version: 1.8.3
>
```

**（2）创建数据库：**

```bash
> CREATE DATABASE dbname
```

如果执行上面的命令没有任何错误信息输出则说明数据库已经创建成功了，我们可以通过下面的命令查看数据库：

```bash
> CREATE DATABASE test
> SHOW DATABASES
name: databases
name
----
_internal
test
```

我们可以看到除了我们创建的test数据库外，还有一个_internal数据库，这个数据库是用来保存influxDB运行时的metrics信息。

创建完数据库，我们可以通过下面的命令使用数据库：

```bash
> USE test
Using database test
```



**（3）写入数据：**

前面我们已经了解了influxDB是如何保存数据的，现在我们来实际操练一下，我们可以通过`INSERT`命令写入数据，数据格式必须满足[influxDB line protocol](https://docs.influxdata.com/influxdb/v1.8/write_protocols/line_protocol_reference/#line-protocol-syntax)。

```bash
> INSERT cpu,host=serverA,region=us_west value=0.64
```



**（4）查询数据：**

```bash
# 只检索需要的字段：
> SELECT "host", "region", "value" FROM "cpu"
name: cpu
time                host    region  value
----                ----    ------  -----
1603696637517565428 serverA us_west 0.64
# 查询所有字段
> SELECT * FROM "cpu"
name: cpu
time                host    region  value
----                ----    ------  -----
1603696637517565428 serverA us_west 0.64
# 带条件查询
> SELECT * FROM "cpu" WHERE "value">0.5
name: cpu
time                host    region  value
----                ----    ------  -----
1603696637517565428 serverA us_west 0.64
# 设置limit
> SELECT * FROM "cpu" LIMIT 1
name: cpu
time                host    region  value
----                ----    ------  -----
1603696637517565428 serverA us_west 0.64
```

<font color=red size=2>注意：influxQL对大小写不敏感，因此上面的命令用小写也是OK的。</font>

#### 5. 数据库配置

#### 6. influxQL

#### 7. influx

#### 8. 支持的协议

#### 9. 高可用

#### 10. 排查故障





参考文档：

- https://docs.influxdata.com/influxdb/v2.0/