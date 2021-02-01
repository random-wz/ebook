> gorm 是 golang 的一个 orm 框架，orm 全称 object relation mapping 对象映射关系，目的是解决面向对象和关系数据库之间存在的互不匹配的现象。本文向大家介绍 gorm 的踩坑经历，持续更新，希望对你们有帮助。

#### 1. 表名

在 gorm 中默认采用蛇形小写的方式定义表名和列名，我们在使用的时候其实是可以自定义表名的。

【1】方法一：定义表名规则

```go
// 通过定义DefaultTableNameHandler对默认表名应用任何规则。
gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
    return "prefix_" + defaultTableName;
}
```

【2】方法二：通过TableName方法设置表名

```go
type User struct {} // 默认表名是`users`即结构体名称的复数形式

// 设置User的表名为`student`
func (User) TableName() string {
  return "student"
}
// 不同情况返回不同的表名
func (u User) TableName() string {
    if u.Role == "admin" {
        return "admin"
    } else {
        return "student"
    }
}
```

如果不想让表明为结构体名称的复数形式，我们可以禁用：

```go
// 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
db.SingularTable(true) 
```



#### 2. 更新数据

更新数据的时候有下面几点需要注意：

- Update只能更新单个字段

- 当使用struct更新时，Updates将仅更新具有非空值的字段（""，0，false为空字段）

  这种情况我们可以采用map或者Select来更新我们要更改的字段，或者直接自己写Sql语句。

#### 3. 删除数据

如果模型有`DeletedAt`字段，它将自动获得软删除功能！ 那么在调用`Delete`时不会从数据库中永久删除，而是只将字段`DeletedAt`的值设置为当前时间。因此如果想要永久删除需要调用Unscoped方法。

#### 4. Callback



