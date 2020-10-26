<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Viper(配置文件管理利器)](#viper%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E7%AE%A1%E7%90%86%E5%88%A9%E5%99%A8)
    - [1. viper 可以用来做什么](#1-viper-%E5%8F%AF%E4%BB%A5%E7%94%A8%E6%9D%A5%E5%81%9A%E4%BB%80%E4%B9%88)
    - [2. 加载本地配置文件](#2-%E5%8A%A0%E8%BD%BD%E6%9C%AC%E5%9C%B0%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
    - [3. 加载命令行参数](#3-%E5%8A%A0%E8%BD%BD%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0)
    - [4. 加载环境变量](#4-%E5%8A%A0%E8%BD%BD%E7%8E%AF%E5%A2%83%E5%8F%98%E9%87%8F)
    - [5. 加载etcd、Firestore或consul中的配置信息](#5-%E5%8A%A0%E8%BD%BDetcdfirestore%E6%88%96consul%E4%B8%AD%E7%9A%84%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF)
    - [6. 从 io.Reader 中加载配置信息](#6-%E4%BB%8E-ioreader-%E4%B8%AD%E5%8A%A0%E8%BD%BD%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF)
    - [7. viper的热更新](#7-viper%E7%9A%84%E7%83%AD%E6%9B%B4%E6%96%B0)
    - [8. 配置信息序列化](#8-%E9%85%8D%E7%BD%AE%E4%BF%A1%E6%81%AF%E5%BA%8F%E5%88%97%E5%8C%96)
    - [9. 保存配置文件](#9-%E4%BF%9D%E5%AD%98%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Viper(配置文件管理利器)

> 在项目开发过程中我们经常需要处理各种配置文件、环境变量、用户输入参数等，今天推荐一个好用的开源项目 github.com/spf13/viper ，它可以很好的帮助我们获取配置信息，希望对你有帮助。

#### 1. viper 可以用来做什么

viper 是适用于 Go 应用程序的完整配置解决方案，它运行在程序内部，并且可以处理所有类型的配置需求和格式。他的主要功能如下：

- 设置默认值
- 从 JSON, TOML, YAML, HCL, envfile and Java properties 配置文件读取配置信息
- 实时读取加载配置文件
- 从环境变量中读取配置信息
- 监听和读取etcd、firestore或consul中的配置信息
- 从命令行读取配置信息
- 从缓存中读取配置信息
- 代码逻辑中显示设置键值

总之 viper 可以解决你软件开发的所有配置需求，它的安装也十分简单：

```bash
go get github.com/spf13/viper
```

使用的时候直接导入即可，下面举一个简单的例子：

```go
func main() {
	// 设置 ip 的值为 127.0.0.1 ,不会覆盖已有的值
	viper.SetDefault("ip", "localhost")
	fmt.Println("SetDefaultIP:", viper.GetString("ip"))
	// 设置 ip 的值为 127.0.0.1 ,会覆盖已有的值
	viper.Set("ip", "127.0.0.1")
	fmt.Println("SetIP:", viper.GetString("ip"))
	// 给ip设置别名为server_ip
	viper.RegisterAlias("server_ip", "ip")
	fmt.Println("ServerIP:", viper.GetString("server_ip"))
}
```

Output:

```bash
$ go run main.go
SetDefaultIP: localhost
SetIP: 127.0.0.1
ServerIP: 127.0.0.1
```

<font color=red>注意：viper 对大小写并不敏感。</font>



#### 2. 加载本地配置文件

加载本地配置文件主要用到下面五种方法：

- SetConfigFile  指定具体的配置文件，将会忽略其他配置文件
- SetConfigName  指定配置文件名称，需要与配置文件名称相同
- SetConfigType  指定配置文件类型
- SetConfigPermissions  配置配置文件权限，一般不用
- AddConfigPath  配置加载配置文件的路径，可以设置多个，viper 会从所有的目录查找配置文件

下面举个例子：

本地创建一个配置文件 config.json，内容如下：

```json
{
  "id": 10,
  "name": "random",
  "country": [
    "China",
    "USA"
  ],
  "student": {
    "random1": "China"
  }
}
```

加载配置文件：

```go
func LoadConfig() {
	//viper.SetConfigFile()  指定具体的配置文件，将会忽略其他配置文件
	viper.SetConfigName("config")  // 指定配置文件名称，需要与配置文件名称相同
	viper.SetConfigType("json")    // 指定配置文件类型
	//viper.SetConfigPermissions() // 配置配置文件权限，一般不用
	viper.AddConfigPath(".")       // 配置加载配置文件的路径，可以设置多个，viper 会从所有的目录查找配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatal("Config File Not Found")
		} else {
			// Config file was found but another error was produced
			log.Fatal("Load Config File Fail", err.Error())
		}
	}
}
```

查看配置信息：

```go
func main() {
	LoadConfig()
	Config := fmt.Sprintf("ID: %d\tName: %s\tCountry: %v\tStudent: %v",
		viper.GetInt("id"),
		viper.GetString("name"),
		viper.GetStringSlice("country"),
		viper.GetStringMapString("student"),
	)
	fmt.Println(Config)
}
```



#### 3. 加载命令行参数

与通过命令行加载配置信息相关的方法有下面四个：

- ```go
  BindPFlag(key string, flag *pflag.Flag) error  使用pflag绑定配置信息，绑定一个参数
  ```

- ```go
  BindPFlags(flags *pflag.FlagSet) error 使用pflag绑定配置信息，绑定多个参数
  ```

- ```
  BindFlagValue(key string, flag FlagValue) error 自定义绑定配置信息的方法，绑定一个参数
  ```

- ```go 
  BindFlagValues(flags FlagValueSet) error 自定义绑定配置信息的方法，绑定多个参数
  ```

按照下面三种场景我们来看一下，如何使用这四个方法：

**场景一:**  使用 pflag 接收命令行参数

这里需要使用一个第三方库`github.com/spf13/pflag`, pflag 和 flag 的功能类似，但是它的功能更强大。

```go
func main() {
	pflag.String("ip", "127.0.0.1", "IP")
	pflag.Int("port", 8080, "Port")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err.Error())
	}
	log.Println(viper.AllSettings())
}
```

另外我们还可以在 cobra 库定义参数后将对应的值绑定到 viper ：

```go
serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
```



**场景二：**使用 flag 接收命令行参数

```go
func main() {
	flag.String("ip", "127.0.0.1", "IP")
	flag.Int("port", 8080, "Port")
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("IP:", viper.GetString("ip"), "Port:", viper.GetInt("port"))
}
```



**场景三：**自定义

viper 提供的另外 两个方法`BindFlagValues`和`BindFlagValue`，这两种方法进行绑定时，被绑定的对象必须实现下面接口的四个方法：

```go
// FlagValue is an interface that users can implement
// to bind different flags to viper.
type FlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}
```

举个例子：

```go
// FlagValue 只是对flag的封装，因为flag本身并没有实现FlagValue接口
type FlagValue struct {
	flag *flag.Flag
}

func (f FlagValue) HasChanged() bool {
	return false
}
func (f FlagValue) Name() string {
	return f.flag.Name
}
func (f FlagValue) ValueString() string {
	return f.flag.Value.String()
}
func (f FlagValue) ValueType() string {
	return reflect.TypeOf(f.flag.Value).Kind().String()
}

func main() {
	flag.String("ip", "127.0.0.1", "ip address")
	flag.Int("port", 8080, "port")
	flag.Parse()
	if err := viper.BindFlagValue("ip", &FlagValue{flag.Lookup("ip")}); err != nil {
		log.Fatal(err.Error())
	}
	if err := viper.BindFlagValue("port", &FlagValue{flag.Lookup("port")}); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("IP:", viper.GetString("ip"), "Port:", viper.GetInt("port"))
}
```

要使用`BindFlagValues`来绑定配置信息，则需要实现 FlagValueSet 接口：

```go
// FlagValueSet is an interface that users can implement
// to bind a set of flags to viper.
type FlagValueSet interface {
	VisitAll(fn func(FlagValue))
}
```

举个例子：

```go
// FlagValue 只是对flag的封装，因为flag本身并没有实现FlagValue接口
type FlagValue struct {
	flag *flag.Flag
}

func (f FlagValue) HasChanged() bool {
	return false
}
func (f FlagValue) Name() string {
	return f.flag.Name
}
func (f FlagValue) ValueString() string {
	return f.flag.Value.String()
}
func (f FlagValue) ValueType() string {
	return reflect.TypeOf(f.flag.Value).Kind().String()
}

type FlagValueSet struct {
	flags *flag.FlagSet
}

//  FlagValueSet 实现FlagValueSet 接口
func (p FlagValueSet) VisitAll(fn func(flag viper.FlagValue)) {
	p.flags.VisitAll(func(flag *flag.Flag) {
		fn(FlagValue{flag})
	})
}

func main() {
	flag.String("ip", "127.0.0.1", "ip address")
	flag.Int("port", 8080, "port")
	flag.Parse()
	if err := viper.BindFlagValues(&FlagValueSet{flag.CommandLine}); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("IP:", viper.GetString("ip"), "Port:", viper.GetInt("port"))
}
```





#### 4. 加载环境变量

加载环境变量的方式很简单，直接调用`viper.AutomaticEnv()`方法即可，读取配置信息方法不变。另外 viper 还提供了下面四个方法，大大提高了加载环境变量的灵活性：

- `BindEnv(string...)  error`  加载配置信息
- `SetEnvPrefix(string) `设置环境变量前缀
- `SetEnvKeyReplacer(string...) *strings.Replacer` 替换环境变量中的字符
- `AllowEmptyEnv(bool)`  允许边境变量为空

下面举个例子：

```go
func main() {
	viper.AutomaticEnv()
	if err := viper.BindEnv("port", "ip"); err != nil {
		log.Fatal(err.Error())
	}
	viper.SetEnvPrefix("server")
	viper.AllowEmptyEnv(true)
	_ = os.Setenv("SERVER_PORT", "8080")
	_ = os.Setenv("SERVER_IP", "127.0.0.1")
	fmt.Println(viper.GetString("ip"), viper.GetInt("port"))
}
```

Output:

```bash
$ go run main.go
127.0.0.1 8080
```

<font color=red>注意：viper 默认环境变量名称为大写，配置的时候写成小写，系统会自动将名称进行大写转换。</font>



#### 5. 加载etcd、Firestore或consul中的配置信息

要通过配置文件系统中绑定配置信息，首先需要匿名导入`viper/remote`包：

```go
import _ "github.com/spf13/viper/remote"
```

viper 可以根据 key 获取到 JSON, TOML, YAML, HCL 或者 envfile类型的配置信息，这些信息可以通过 crypt 进行加密，在获取配置信息的时候指定gpg文件就行，下面罗列了，读取etcd、Firestore和consul中的配置信息的方法：

**加密：**

- etcd:

  ```go
  viper.AddSecureRemoteProvider("etcd","http://127.0.0.1:4001","/config/hugo.json","/etc/secrets/mykeyring.gpg")
  viper.SetConfigType("json")
  err := viper.ReadRemoteConfig()
  ```

  

- consul:

  ```go
  viper.AddSecureRemoteProvider("consul", "localhost:8500", "MY_CONSUL_KEY","/etc/secrets/mykeyring.gpg")
  viper.SetConfigType("json")
  err := viper.ReadRemoteConfig()
  ```

  

- firestore

  ```go
  viper.AddSecureRemoteProvider("firestore", "google-cloud-project-id", "collection/document","/etc/secrets/mykeyring.gpg")
  viper.SetConfigType("json")
  err := viper.ReadRemoteConfig()
  ```

  

**未加密:**

- etcd:

  ```go
  viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001","/config/hugo.json")
  // 支持的类型有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
  viper.SetConfigType("json") 
  err := viper.ReadRemoteConfig()
  ```

  

- consul:

  ```go
  viper.AddRemoteProvider("consul", "localhost:8500", "MY_CONSUL_KEY")
  // 支持的类型有 "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
  viper.SetConfigType("json") 
  err := viper.ReadRemoteConfig()
  ```

  

- firestore

  ```go
  viper.AddRemoteProvider("firestore", "google-cloud-project-id", "collection/document")
  // 支持的类型有  "json", "toml", "yaml", "yml"
  viper.SetConfigType("json") 
  err := viper.ReadRemoteConfig()
  ```

  

#### 6. 从 io.Reader 中加载配置信息

这里使用了`viper.ReadConfig`方法，如下例：

```go
// any approach to require this configuration into your program.
var yamlExample = []byte(`
Hacker: true
name: steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: denim
age: 35
eyes : brown
beard: true
`)

func LoadConfig() {
    viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
    viper.ReadConfig(bytes.NewBuffer(yamlExample))
}
```



#### 7. viper的热更新

viper 支持监听配置文件变化实现配置信息热更新，通过`viper.WatchConfig`开启监听，我们可以通过`viper.OnConfigChange`方法配置当配置文件变更的时候做什么操作，下面举个例子：

```go
func LoadConfig() {
	//viper.SetConfigFile()
	viper.SetConfigName("config2")
	viper.SetConfigType("json")
	//viper.SetConfigPermissions()
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatal("Config File Not Found")
		} else {
			// Config file was found but another error was produced
			log.Fatal("Load Config File Fail", err.Error())
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config File Change: ", in.String())
	})
}
```

如果我们的配置信息是从类似于 etcd 的系统中获取的，我们改如何进行热更新呢？请看下面的例子：

```go
// 创建一个实例
var runtime_viper = viper.New()

runtime_viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/hugo.yml")
runtime_viper.SetConfigType("yaml") 

// 读取配置文件
err := runtime_viper.ReadRemoteConfig()

// 序列化
runtime_viper.Unmarshal(&runtime_conf)

// 通过 Goroutine 监控配置信息变化
go func(){
	for {
	    time.Sleep(time.Second * 5) // 设置时间间隔

	    // 目前只支持 etcd, 重新读取配置信息
	    err := runtime_viper.WatchRemoteConfig()
	    if err != nil {
	        log.Errorf("unable to read remote config: %v", err)
	        continue
	    }
	    runtime_viper.Unmarshal(&runtime_conf)
	}
}()
```



#### 8. 配置信息序列化

viper 提供了序列化的功能，我们可以将读取到的配置信息进行序列化。

**Unmarshal: ** 

如下例，我们有一个配置文件config2.json：

```json
{
  "id": 101,
  "name": "random1",
  "country": [
    "China1",
    "USA1"
  ],
  "student": {
    "random1": "China1"
  }
}
```

加载配置文件并反序列化：

```go
func LoadConfig() {
	//viper.SetConfigFile()
	viper.SetConfigName("config2")
	viper.SetConfigType("json")
	//viper.SetConfigPermissions()
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatal("Config File Not Found")
		} else {
			// Config file was found but another error was produced
			log.Fatal("Load Config File Fail", err.Error())
		}
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config File Change: ", in.String())
	})
}

type Data struct {
	ID      int
	Name    string
	Country []string
	Student map[string]string
}

func main() {
	var d Data
    // 加载配置文件
	LoadConfig()
    // 将配置文件的内容解析到d
	if err := viper.Unmarshal(&d); err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)
}
```

**Marshal:**

如下例：

```go
func main() {
	c := viper.AllSettings()
	bs, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to JSON: %v", err)
	}
	fmt.Println(string(bs))
}
```



#### 9. 保存配置文件

我们可以将 viper 加载的配置文件信息保存到指定的文件中，下面是会用到的四个方法：

- WriteConfig -将配置信息写入之前定义的配置文件中，如果配置文件存在，则会覆盖原文内容，如果不存在则报错。
- SafeWriteConfig - 将配置信息写入之前定义的配置文件中，如果配置文件存在，不会覆盖原文内容，如果不存在则报错。
- WriteConfigAs - 将配置信息写入指定的文件，如果文件存在，则会覆盖原文件内容
- SafeWriteConfigAs - 将配置信息写入指定的文件，如果文件存在，不会覆盖原文件内容

我们可以使用这个功能来生成配置文件，如下例：

```go
func main() {
	var Resp struct {
		IP     string
		Port   int
		Expire time.Duration
	}
	Resp.IP = "127.0.0.1"
	Resp.Port = 8080
	Resp.Expire = time.Second * 10
	viper.Set("data", Resp)
	if err := viper.SafeWriteConfigAs("conf.json"); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Write conf.json success")
}
```

如果运行没有报错则会在当前目录下生成一个 conf.json 文件，文件内容如下：

```json
{
  "data": {
    "IP": "127.0.0.1",
    "Port": 8080,
    "Expire": 10000000000
  }
}
```

其他三个方法就不测试了，感兴趣的可以自己试一下。

