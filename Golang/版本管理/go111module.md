> Go开发环境中有一个环境变量 GO111MODULE，默认变量值为 auto，在项目开发过程中我们可能会因为项目版本控制等原因，需要修改该变量的值，本文记录了 GO111MOODULE 环境变量的详细介绍，希望对你有帮助。

我们可以用环境变量 `GO111MODULE` 开启或关闭模块支持，它有三个可选值：off、on、auto，默认值是 auto。

- `GO111MODULE=off` 无模块支持，go 会从 GOPATH 和 vendor 文件夹寻找包。
- `GO111MODULE=on` 模块支持，go 会忽略 GOPATH 和 vendor 文件夹，只根据 go.mod 下载依赖。
- `GO111MODULE=auto` 在 $GOPATH/src 外面且根目录有 go.mod 文件时，开启模块支持。

在使用模块的时候，GOPATH 是无意义的，不过它还是会把下载的依赖储存在 GOPATH/pkg/mod中，也会把goinstall的结果放在GOPATH/pkg/mod中，也会把goinstall的结果放在GOPATH/bin 中。

