> 在Go语言标准库中提供了进行数据库操作的 database/sql 库，需要注意的是在使用 sql 库的时候需要导入数据库驱动。本文记录了 database/sql 标准库的学习笔记，希望对你有帮助。

#### 1. 数据库连接

链接数据库只需要四步即可：

```mermaid
graph LR
A[导入数据库引擎] -- 生成连接语句 --> B[连接数据库]
B --> C[设置连接参数]
C --> D[测试链接]
```

- 导入数据库引擎

  下面列出了常用的数据库引擎，更过请参考 [数据库引擎列表]( https://github.com/golang/go/wiki/SQLDrivers )

  ```go
   github.com/go-sql-driver/mysql/ 
   github.com/mattn/go-sqlite3 
  ```

  

- 连接数据库

  连接数据库前需要设置数据库连接参数，如下例：

  ```go
  type MySqlConfig struct {
  	UserName string
  	Password string
  	IP       string
  	Port     int
  	Database string
  }
  
  func parseMysqlConfig(m MySqlConfig) string {
  	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
  		m.UserName, m.Password, m.IP, m.Port, m.Database)
  }
  ```

  设置完参数，我们通过Open函数连接数据库：

  ```go
  db, err := sql.Open("mysql", parseMysqlConfig(mysql))
  if err != nil {
  	panic(err)
  }
  ```

  

- 设置连接参数

  设置的参数有三种：

  ```go
  // 设置连接的生命周期，超过这个时间他会被Close，一般不会设置这个
  db.SetConnMaxLifetime(time.Hour * 24)
  // 设置最大空闲连接池的大小，默认为2
  db.SetMaxIdleConns(100)
  // 设置最大连接数量
  db.SetMaxOpenConns(20)
  ```

  

- 测试链接

  db 有一个 ping 方法用来测试是否连接成功：

  ```go
  // 通过 ping 测试数据库是否连接成功
  if err := db.Ping(); err != nil {
  	panic(err)
  }
  ```

我们看一下完整代码：

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySqlConfig struct {
	UserName string
	Password string
	IP       string
	Port     int
	Database string
}

func parseMysqlConfig(m MySqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		m.UserName, m.Password, m.IP, m.Port, m.Database)
}

func ConnectMysql() (*sql.DB, error) {
	// Mysql 数据库配置
	mysql := MySqlConfig{
		UserName: "root",
		Password: "root",
		IP:       "127.0.0.1",
		Port:     3306,
		Database: "test",
	}
	db, err := sql.Open("mysql", parseMysqlConfig(mysql))
	if err != nil {
		return db, err
	}
	// 设置连接的生命周期，超过这个时间他会被Close，一般不会设置这个
	db.SetConnMaxLifetime(time.Hour * 24)
	// 设置最大空闲连接池的大小，默认为2
	db.SetMaxIdleConns(100)
	// 设置最大连接数量
	db.SetMaxOpenConns(20)
	// 通过 ping 测试数据库是否连接成功
	if err := db.Ping(); err != nil {
		return db, err
	}
	fmt.Println("Connect Mysql Success....")
	return db, nil
}

```

调用 ConnectMysql 函数，输出`Connect Mysql Success....`表示连接成功。

#### 2. 执行 SQL 语句

- 执行一条 Sql 命令

  ```go
  // Exec执行一次命令（包括查询、删除、更新、插入等），不返回任何执行结果。
  // 参数args表示query中的占位参数。
  // MySQL 占位符为 ?
  // PostgreSQL 占位符为	$1, $2等
  // SQLite 占位符为 ? 和$1
  // Oracle 占位符为 :name
  func (db *DB) Exec(query string, args ...interface{}) (Result, error)
  ```

  

- 执行一次查询返回多条结果

  ```go
  // Query执行一次查询，返回多行结果（即Rows），一般用于执行select命令。
  // 参数args表示query中的占位参数。
  func (db *DB) Query(query string, args ...interface{}) (*Rows, error)
  ```

  

- 执行一次查询期望返回最多一条结果

  ```go
  // QueryRow执行一次查询，并期望返回最多一行结果（即Row）。
  // QueryRow总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误。（如：未找到结果）
  func (db *DB) QueryRow(query string, args ...interface{}) *Row
  ```

  

- Prepare

  ```go
  // Prepare创建一个准备好的状态用于之后的查询和命令。
  // 返回值可以同时执行多个查询和命令。
  func (db *DB) Prepare(query string) (*Stmt, error)
  ```

接下来我们使用上面的方法进行数据库的修改与查询：

```go
func main() {
	db, _ := ConnectMysql()
	// 插入数据
	if _, err := db.Exec("INSERT INTO student (name,age,school,sex) VALUES (?, ?, ?, ?) ",
		"random1", 20, "school1", 1); err != nil {
		fmt.Println(err)
		return
	}
	// 使用Query查询数据
	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("使用Query:\n", "ID", "|", "Name", "|", "Age", "|", "School", "|", "Sex")
	for rows.Next() {
		var (
			ID     int
			Name   string
			Age    int
			School string
			Sex    int
		)
		if err := rows.Scan(&ID, &Name, &Age, &School, &Sex); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ID, Name, Age, School, Sex)
	}
	// 使用
	var (
		ID     int
		Name   string
		Age    int
		School string
		Sex    int
	)
	_ = db.QueryRow("SELECT * FROM student WHERE id = 3").Scan(&ID, &Name, &Age, &School, &Sex)
	fmt.Println("使用QueryRow:\n", ID, Name, Age, School, Sex)
}
```

Output:

```bash
$ go run main.go
Connect Mysql Success....
使用Query:
 ID | Name | Age | School | Sex
1 random 20 school1 1
3 random1 20 school1 1
使用QueryRow:
 3 random1 20 school1 1
```



#### 3. 获取 Sql 执行结果

- Row ( Row 只有一个方法，用来保存结果到指定的变量 )

  ```go
  // Scan将该行查询结果各列分别保存进dest参数指定的值中。
  // 如果该查询匹配多行，Scan会使用第一行结果并丢弃其余各行。如果没有匹配查询的行，Scan会返回ErrNoRows。
  func (r *Row) Scan(dest ...interface{}) error
  ```

- Rows 

  ```go
  // Columns返回列名。如果Rows已经关闭会返回错误。
  func (rs *Rows) Columns() ([]string, error)
  // Scan将当前行各列结果填充进dest指定的各个值中。
  func (rs *Rows) Scan(dest ...interface{}) error
  // Next准备用于Scan方法的下一行结果。
  // 如果成功会返回真，如果没有下一行或者出现错误会返回假。
  // Err应该被调用以区分这两种情况。
  func (rs *Rows) Next() bool
  // Close关闭Rows，阻止对其更多的列举。 
  func (rs *Rows) Close() error
  // Err返回可能的、在迭代时出现的错误。Err需在显式或隐式调用Close方法后调用。
  func (rs *Rows) Err() error
  ```

  

#### 4. Stmt

Stmt是准备好的状态。Stmt可以安全的被多个go程同时使用，db.Prepare函数返回与当前连接相关的执行 Sql 语句的准备状态，可以进行查询、删除等操作。

```go
// Exec使用提供的参数执行准备好的命令状态，返回Result类型的该状态执行结果的总结。
func (s *Stmt) Exec(args ...interface{}) (Result, error)
// Query使用提供的参数执行准备好的查询状态，返回Rows类型查询结果。
func (s *Stmt) Query(args ...interface{}) (*Rows, error)
// QueryRow使用提供的参数执行准备好的查询状态。
func (s *Stmt) QueryRow(args ...interface{}) *Row
// Close关闭状态。
func (s *Stmt) Close() error
```

举个例子：

```go
func main(){
    db, _ := ConnectMysql()
	stmt, _ := db.Prepare("SELECT * FROM student WHERE id = ?")
	var g sync.WaitGroup
	g.Add(3)
	for i := 1; i < 4; i++ {
        // 创建三个协程只从sql语句
		go func(i int) {
			var (
				ID     int
				Name   string
				Age    int
				School string
				Sex    int
			)
			if err := stmt.QueryRow(i).Scan(&ID, &Name, &Age, &School, &Sex); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(ID, Name, Age, School, Sex)
			}
			g.Done()
		}(i)
	}
	g.Wait()
}
```

Output: 

```go
$ go run main.go
Connect Mysql Success....
3 random1 20 school1 1
2 random 20 school1 1
1 random 20 school1 1
```

我们可以看到三个协程都执行成功了。

#### 5. 事务

通过 db.Begin() 可以生成一个事务 Tx ，相信大家对对数据库的事务是有一定了解的，一次事务必须以对Commit或Rollback的调用结束。下面是 Tx 相关的方法：

```go
// Exec执行命令，但不返回结果。例如执行insert和update。
func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)
// Query执行查询并返回零到多行结果（Rows），一般执行select命令。
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)
// QueryRow执行查询并期望返回最多一行结果（Row）。
// QueryRow总是返回非nil的结果，查询失败的错误会延迟到在调用该结果的Scan方法时释放。
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row
// Prepare准备一个专用于该事务的状态。
// 返回的该事务专属状态操作在Tx递交会回滚后不能再使用。
// 要在该事务中使用已存在的状态，参见Tx.Stmt方法。
func (tx *Tx) Prepare(query string) (*Stmt, error)
// Stmt使用已存在的状态生成一个该事务特定的状态。
func (tx *Tx) Stmt(stmt *Stmt) *Stmt
// Commit递交事务。
func (tx *Tx) Commit() error
// Rollback放弃并回滚事务。
func (tx *Tx) Rollback() error
```

举个例子：

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	db, _ := ConnectMysql()
	tx, _ := db.Begin()
	// 写入数据
	result, err := tx.Exec("INSERT INTO student (name,age,school,sex) VALUES (?, ?, ?, ?) ",
		"random_w", 18, "school2", 2)
	if err != nil {
		_ = tx.Rollback()
		fmt.Println(err)
		return
	}
	// 获取写入数据的ID
	id, _ := result.LastInsertId()
	var (
		ID     int
		Name   string
		Age    int
		School string
		Sex    int
	)
	// 通过ID检索数据
	if err := tx.QueryRow("SELECT * FROM student WHERE id = ?", id).Scan(&ID, &Name, &Age, &School, &Sex); err != nil {
		_ = tx.Rollback()
		fmt.Println(err)
		return
	}
	fmt.Println(ID, Name, Age, School, Sex)
	_ = tx.Commit()
	return
}
```

Output:

```bash
$ go run main.go
Connect Mysql Success....
8 random_w 18 school2 2
```

