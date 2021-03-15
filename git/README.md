#### 1. gitignore 文件

我们使用 `git` 做版本管理的时候，通常会使用 `.gitignore` 文件来配置哪些文件我们不希望长传到 git 仓库，比如本地测试或者编辑器在运行的时候会产生一些文件，我们就可以通过配置 `.gitignore` 文件的方式让 git 忽略这些文件。

#### 2. gitignore 文件配置语法

- 以斜杠`/`开头表示目录；
- 以星号`*`通配多个字符；
- 以问号`?`通配单个字符
-  以方括号`[]`包含单个字符的匹配列表；
-  以叹号`!`表示不忽略(跟踪)匹配到的文件或目录；

需要注意的是，`.gitignore` 文件是按行从上到下进行规则匹配的。