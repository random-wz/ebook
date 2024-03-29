> 这篇笔记记录重学Python的相关学习记录。

### 一、基本语法

#### 1. 中文编码

在Python中默认使用的是ASCII格式，因此需要打印汉语的时候需要指定编码格式为UTF-8，语法如下：

```python
# 第一种
# —*— coding: UTF-8 —*—
# 第二种
# coding=utf-8
```

#### 2. Python标识符

在Python中标识符由字母、数字、下划线组成，但是不能以数字开头，Python中的标识符是区分大小写的，在Python中以下划线开头的标识符都有特殊含义，如`_foo`代表不能直接访问的类属性，需要通过类的接口进行访问，不能通过import导入，以双下划线开头的`__foo`代表类的私有成员，以双下划线开头和结尾的`__foo__`代表Python里特殊方法的专用标识，如`__init__()`代表类的构造函数。

#### 3. Python保留字符

| and      | exec    | not    |
| -------- | ------- | ------ |
| assert   | finally | or     |
| break    | for     | pass   |
| class    | from    | print  |
| continue | global  | raise  |
| def      | if      | return |
| del      | import  | try    |
| elif     | in      | while  |
| else     | is      | with   |
| except   | lambda  | yield  |

#### 4. 行与缩进

在Python中对缩进的要求很高，代码的缩进空白数量是可变的，但是所有代码块语句必须包含相同的缩进空白数量。

#### 5. 引号

Python可以使用单引号(')、(")、三引号(""")来表示字符串，引号的开始和结束必须是相同类型，其中三引号可以由多行组成，编写多行文本时使用。

#### 6. 注释

Python中单行注释适用`#`开头，多行注释使用三个单引号或者三个双引号组成：

```python
#!/usr/bin/python
# -*- coding: UTF-8 -*-
# 文件名：test.py

# 第一个注释
'''
这是多行注释，使用单引号。
这是多行注释，使用单引号。
这是多行注释，使用单引号。
'''
```

#### 7. 空行

函数之间和类的方法之间使用空行分隔，表示一段新的代码的开始。类和函数入口之间也用一行空行分隔，以突出函数入口的开始。需要注意的是空行并不会影响代码编译，使用空行主要是提高代码的可读性，便于维护或重构。

#### 8. 同一行显示多条语句

Python中可以在同一行中使用多条语句，语句之间使用分号分隔。

#### 9. 输入输出





### 二、数据类型

### 三、条件语句

### 四、循环语句

### 五、函数