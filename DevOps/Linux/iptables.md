**iptables 四张表：**

- filter表：负责过滤功能，防火墙；内核模块：iptables_filter

- nat表：network address translation，网络地址转换功能；内核模块：iptable_nat

- mangle表：拆解报文，做出修改，并重新封装 的功能；iptable_mangle

- raw表：关闭nat表上启用的连接追踪机制；iptable_raw

---

**iptables 四种规则：**

- PREROUTING    路由前的规则，可以存在于：raw表，mangle表，nat表。

- INPUT      输入的规则，可以存在于：mangle表，filter表，（centos7中还有nat表，centos6中没有）。

- FORWARD     转发的规则，可以存在于：mangle表，filter表。

- OUTPUT     输出的规则，可以存在于：raw表mangle表，nat表，filter表。

- POSTROUTING    路由后的规则，可以存在于：mangle表，nat表。

---

**处理规则：**

- **ACCEPT**：允许数据包通过。

- **DROP**：直接丢弃数据包，不给任何回应信息，这时候客户端会感觉自己的请求泥牛入海了，过了超时时间才会有反应。

- **REJECT**：拒绝数据包通过，必要时会给数据发送端一个响应的信息，客户端刚请求就会收到拒绝的信息。

- **SNAT**：源地址转换，解决内网用户用同一个公网地址上网的问题。

- **MASQUERADE**：是SNAT的一种特殊形式，适用于动态的、临时会变的ip上。

- **DNAT**：目标地址转换。

- **REDIRECT**：在本机做端口映射。