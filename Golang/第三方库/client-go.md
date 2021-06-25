### 通过 go module 引入`cleint-go`的正确姿势

#### 1. 问题

在项目上我喜欢通过`go module`做包管理，在使用的过程中会遇到各种导入包但因为版本问题无法正常使用的问题，这次遇到的问题如下：

> 项目中引用了`client-go`包，执行完 go mod tidy 及 go mod vendor 后，运行程序发现程序无法运行，报错如下： cannot load k8s.io/api/auditregistration/v1alpha1.....

#### 2. 刨根问底

第一时间检查vendor目录下报错信息中的库是否存在，发现真的不存在，于是查看`go.mod`文件：

```bash
require (
	github.com/imdario/mergo v0.3.12 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.21.0 // indirect
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10 // indirect
)
```

查看 [cleint-go v11.0.0](https://github.com/kubernetes/client-go/tree/release-11.0) 版本的源码，发现这个版本并没有使用 go module 做版本管理，在README文档中说了下面的一句话：

> We currently recommend using the v10.0.0 tag. See [INSTALL.md](https://github.com/kubernetes/client-go/blob/v11.0.0/INSTALL.md) for detailed installation instructions. `go get k8s.io/client-go/...` works, but will build `master`, which doesn't handle the dependencies well.

意思大概就是说这个版本不能很好的处理依赖项，建议使用 v10.0.0 版本，找到了产生上面报错的根本原因了：<font color=red>go module 下载的依赖库与 client-go 不匹配导致，程序无法正常运行。</font>

#### 3. 解决问题

那就找一个稳定版本，然后手动修改 go.mod 文件，举例如下：

```bash
replace (
	k8s.io/api => k8s.io/api v0.15.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.10
	k8s.io/client-go => k8s.io/client-go v0.15.10
)
```

当然你也可以根据自己的需要选择满足自己需求的版本，只需要更换示例中的版本号就可以了。

#### 4. 总结

在博客中我省略了谷歌的过程，因为在搜索的过程中，我看到了很多博客直接给出了修改 go.mod 来解决这个问题的办法，但是我想知道为什么那么做，经过反向探索，一步步找出了产生问题的根本原因及解决办法，这里总结一下就是：<font color=red>go module 做版本管理如果报找不到依赖之类的问题，排除网络问题后可以直接去源码找问题，源码的 README 文档中一般都能找到解决思路。</font>



