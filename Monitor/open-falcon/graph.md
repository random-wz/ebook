> graph 是 open-falcon 监控系统中比较重要的一块，graph进程接收从transfer推送来的指标数据，操作rrd文件存储监控数据，graph也为API进程提供查询接口，处理query组件的查询请求、返回绘图数据，本文向大家介绍 graph 的实现，有不足的地方欢迎补充。

