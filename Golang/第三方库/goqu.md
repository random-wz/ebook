> [gopu](https://github.com/doug-martin/goqu) 库可以用来生成及执行 SQL 语句

### 1. 安装

如果启用了 `go module`安装包的时候需要加上版本号：

```bash
go get -u github.com/doug-martin/goqu/v9
```

如果没有使用 `go module`：

```bash
go get -u github.com/doug-martin/goqu
```

### 2. Dialect

我们可以通过 Dialect 指定数据库类型，`goqu` 会根据 Dialect 生成对应的 SQL 语句，下面是 `goqu` 提供的四种 Dialect ：

- mysql - `import _ "github.com/doug-martin/goqu/v9/dialect/mysql"`
- postgres - `import _ "github.com/doug-martin/goqu/v9/dialect/postgres"`
- sqlite3 - `import _ "github.com/doug-martin/goqu/v9/dialect/sqlite3"`
- sqlserver - `import _ "github.com/doug-martin/goqu/v9/dialect/sqlserver"`

