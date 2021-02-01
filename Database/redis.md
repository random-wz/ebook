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

### 一. redis 数据结构

