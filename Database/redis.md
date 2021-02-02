#### 使用 Docker 安装一个 redis 实例

` docker run -itd --name redis-test -p 6379:6379 redis`

进入容器：

`docker exec -it redis-test /bin/bash`

```bash
root@f891ff0e29db:/data# redis-cli
127.0.0.1:6379> help
redis-cli 6.0.9
To get help about Redis commands type:
      "help @<group>" to get a list of commands in <group>
      "help <command>" for help on <command>
      "help <tab>" to get a list of possible help topics
      "quit" to exit

To set redis-cli preferences:
      ":set hints" enable online hints
      ":set nohints" disable online hints
Set your preferences in ~/.redisclirc
```

### 一. redis 数据类型

> redis 共支持五种数据类型：string（字符串），hash（哈希），list（列表），set（集合）及zset(sorted set：有序集合)。

#### 1.  字符串 (string)

string 是 redis 最基本的类型，一个 key 对应一个 value。在 redis 中 string 类型是二进制安全的，也就是说 string 可以包含任何数据，比如 jpg 图片或者序列化的对象，string 类型的值最大能存储 512MB。

举个例子：

```bash
127.0.0.1:6379> SET age 18
OK
127.0.0.1:6379> GET age
"18"
```



#### 2. 哈希 (hash)

哈希是一个键值对集合，也是一个 string 类型的 field 和 value 的映射表，适合用来存储对象。每个 hash 可以存储 232 -1 键值对（40多亿）。

举个例子：

```bash
127.0.0.1:6379> HMSET random age 18 addr china
OK
127.0.0.1:6379> HGET random age
"18"
127.0.0.1:6379> HGET random addr
"china"
```



#### 3. 列表 (List)

列表是一个简单的字符串列表，按照插入顺序排序，你可以插入一个元素到列表的头部或者尾部。列表最多可以存储 2<sup>32</sup> - 1 个元素（4292967295，每个列表可存储40多亿）。

举个例子：

```bash
127.0.0.1:6379> LPUSH student wang
(integer) 1
127.0.0.1:6379> LPUSH student li
(integer) 2
127.0.0.1:6379> LPUSH student zhang
(integer) 3
127.0.0.1:6379> LRANGE student 0 2
1) "zhang"
2) "li"
3) "wang"
```



#### 4. 集合 (Set)

集合是 string 类型的无序集合，集合通常通过哈希表实现，所以增删改查的复杂度都是O(1)。

举个例子：

```bash
127.0.0.1:6379> SADD food noodle rice potato
(integer) 2
127.0.0.1:6379> SMEMBERS food
1) "rice"
2) "potato"
3) "noodle"
```



#### 5. 有序集合 (zset)

有序集合和集合一样都是 string 类型元素的集合，且不允许重复的成员。不同的是有序集合中每个元素都会关联一个 double 类型的分数，redis 通过这个分数来为集合中的成员进行从小到大的排序，zset 成员是唯一的，但分数却可以重复。

举个例子：

```bash
127.0.0.1:6379> ZADD date 1 day
(integer) 1
127.0.0.1:6379> ZADD date 2 month
(integer) 1
127.0.0.1:6379> ZADD date 2 year
(integer) 1
127.0.0.1:6379> ZRANGEBYSCORE date 1 2
1) "day"
2) "month"
3) "year"
```

### 二、基本命令

#### 1. redis 命令

**redis-cli 命令：**



**redis-server 命令：**

#### 2. 键

#### 3. 字符串

#### 4. 哈希

#### 5. 列表

#### 6. 集合

#### 7. 有序集合

#### 8. HperLogLog

#### 9. 发布订阅

#### 10. 事务

#### 11. 脚本



