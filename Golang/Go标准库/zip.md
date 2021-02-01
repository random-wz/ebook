> 在 Go 语言标准库中提供了 archive/zip 包用来进行文件的压缩和解压缩，正好最近工作中用到了这个库，在这里向大家介绍 zip 库的使用方法，希望对你有帮助。

#### 1. 主要方法介绍

##### FileHeader 对象描述了 zip 文件中的一个文件信息，相关方法如下：

- FileInfoHeader

  ```go
  // FileInfoHeader返回一个根据fi填写了部分字段的Header。
  // 因为os.FileInfo接口的Name方法只返回它描述的文件的无路径名，有可能需要将返回值的Name字段修改为文件的完整路径名。
  func FileInfoHeader(fi os.FileInfo) (*FileHeader, error)
  ```

- FileInfo

  ```go
  // FileInfo返回一个根据h的信息生成的os.FileInfo。
  func (h *FileHeader) FileInfo() os.FileInfo
  ```

- Mode

  ```go
  // Mode返回h的权限和模式位。
  func (h *FileHeader) Mode() (mode os.FileMode)
  ```

- SetMode

  ```go
  // SetMode修改h的权限和模式位。
  func (h *FileHeader) SetMode(mode os.FileMode)
  ```

- ModTime

  ```go
  // 返回最近一次修改的UTC时间。（精度2s）
  func (h *FileHeader) ModTime() time.Time
  ```

- SetModTime

  ```go
  // 将ModifiedTime和ModifiedDate字段设置为给定的UTC时间。（精度2s）
  func (h *FileHeader) SetModTime(t time.Time)
  ```

##### File 对象继承了 FileHeader 的所有方法，它提供了两个重要方法：

- DataOffset

  ```go
  // DataOffset返回文件的可能存在的压缩数据相对于zip文件起始的偏移量。
  // 大多数调用者应使用Open代替，该方法会主动解压缩数据并验证校验和。
  func (f *File) DataOffset() (offset int64, err error)
  ```

- Open

  ```go
  // Open方法返回一个io.ReadCloser接口，提供读取文件内容的方法。可以同时读取多个文件。
  func (f *File) Open() (rc io.ReadCloser, err error)
  ```

##### Reader 对象提供了一个方法用来读取文件内容：

- NewReader

  ```go
  // NewReader返回一个从r读取数据的*Reader，r被假设其大小为size字节。
  func NewReader(r io.ReaderAt, size int64) (*Reader, error)
  ```

##### ReaderClose 对象提供了两个方法用来读取压缩文件，以及关闭文件：

- OpenReader

  ```go
  // OpenReader会打开name指定的zip文件并返回一个*ReadCloser。
  func OpenReader(name string) (*ReadCloser, error)
  ```

- Close

  ```go
  // Close关闭zip文件，使它不能用于I/O。
  func (rc *ReadCloser) Close() error
  ```

##### 前面都是和解压缩相关的方法，接下来我们看一下 Writer 对象，它提供了压缩文件的相关方法：

- NewWriter

  ```go
  // NewWriter创建并返回一个将zip文件写入w的*Writer。
  func NewWriter(w io.Writer) *Writer
  ```

- CreateHeader

  ```go
  // 使用给出的*FileHeader来作为文件的元数据添加一个文件进zip文件。
  // 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
  // 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
  func (w *Writer) CreateHeader(fh *FileHeader) (io.Writer, error)
  ```

- Create

  ```go
  // 使用给出的文件名添加一个文件进zip文件。
  // 本方法返回一个io.Writer接口（用于写入新添加文件的内容）。
  // 文件名必须是相对路径，不能以设备或斜杠开始，只接受'/'作为路径分隔。
  // 新增文件的内容必须在下一次调用CreateHeader、Create或Close方法之前全部写入。
  func (w *Writer) Create(name string) (io.Writer, error)
  ```

- Close

  ```go
  // Close方法通过写入中央目录关闭该*Writer。
  // 本方法不会也没办法关闭下层的io.Writer接口。
  func (w *Writer) Close() error
  ```

  

#### 2. 压缩文件

##### 场景一：读取文件内容并压缩

这种场景相对第二种常见一点，过程中需要有原文件以及压缩后的文件，示例代码如下：

```go
func CompressedFile(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil || info.IsDir() {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = prefix + "/" + header.Name
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(writer, file); err != nil {
		return err
	}
	return nil
}

func main() {
	f, _ := os.Open("test.txt")
	// 压缩文件
	dst, _ := os.Create("test.zip")
	zipWriter := zip.NewWriter(dst)
	if err := CompressFile(f, "", zipWriter); err != nil {
		log.Fatalln(err.Error())
	}
    // Make sure to check the error on Close.
	if err := zipWriter.Close(); err != nil {
		log.Fatalln(err.Error())
	}
	return
}
```

CompressedFile 用来将文件内容进行压缩，但是如果我们想目录进行压缩改怎么做呢，其实很简单，只需要将目录里面的文件提取出来然后调用CompressedFile将文件压缩并写入zip.Writer即可：

```go
// Compress 压缩文件
func Compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	// 如果是目录调用CompressedDir
	if info.IsDir() {
		return CompressedDir(file, prefix, zw)
	}
	// 如果是文件调用CompressedFile 
	return CompressedFile(file, prefix, zw)
}
// CompressedDir
func CompressedDir(file *os.File, prefix string, zw *zip.Writer) error {
	info, _ := file.Stat()
	prefix = prefix + "/" + info.Name()
	dirInfo, err := file.Readdir(-1)
	if err != nil {
		return err
	}
	for _, f := range dirInfo {
		f, err := os.Open(file.Name() + "/" + f.Name())
		if err != nil {
			return err
		}
		err = Compress(f, prefix, zw)
		if err != nil {
			return err
		}
	}
	return nil
}
```

<font color=red>注意：文件压缩完成后一定要关闭zip.Writer，否则压缩可能失败。</font>

##### 场景二：将数据直接写入压缩文件

近期在工作中遇到这种情况，我需要将程序格式化的数据写入压缩文件让用户下载，但是我并不想在本地生成文件，只对数据进行压缩，示例代码如下：

```go
func CompressedData(data *bytes.Buffer, dest string) error {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	// Create entry in zip file
	zipEntry, err := zipWriter.Create(dest)
	if err != nil {
		return err
	}
	// Write content into zip writer
	if _, err := zipEntry.Write(data.Bytes()); err != nil {
		return err
	}
	// Make sure to check the error on Close.
	if err := zipWriter.Close(); err != nil {
		return err
	}
	return nil
}
```



#### 3. 解压缩文件

在下面的例子中有两个方法，DeCompressed负责读取压缩文件，并调用deCompressed，将读取的内容写入解压缩后的文件：

```go
func DeCompressed(src string) error {
	s, _ := os.Open(src)
	info, _ := s.Stat()
	ZipReader, err := zip.NewReader(s, info.Size())
	if err != nil {
		return err
	}
	for _, f := range ZipReader.File {
		if err := deCompressed(f); err != nil {
			return err
		}
	}
	return nil
}

func deCompressed(f *zip.File) error {
	d, _ := os.Create(f.Name)
	unzipFile, err := f.Open()
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, unzipFile); err != nil {
		return err
	}
	if err := unzipFile.Close(); err != nil {
		return err
	}
	return d.Close()
}
```

