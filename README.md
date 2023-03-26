# filesystem

[English](./README_EN.md)

[Hertz](https://github.com/cloudwego/hertz) 的文件系统中间件，使用户可以直接使用原生的 `http.Dir`, `http.FS` 等进行静态文件的映射。

注意:

⚠️ **不支持前缀路径中的 `:params`**

⚠️ **此版本尚不完善, 只可以 `h.Use` 使用一次该中间件**

安装:

```shell
go get -u github.com/Skyenought/filesystem
```

简易使用示例:

```go
package main

// ...

func main() {
	h := server.Default()
	h.Use(filesystem.New("/", http.Dir("./testdata"))) // 需要访问的文件夹的相对路径
	h.Spin()
}
```
完整使用示例

```go
package main

func main() {
	h := server.Default()
	h.Use(filesystem.New("/dir", http.Dir("./testdata"), // 支持使用 embed.FS, 即 http.FS
		filesystem.WithBrowse(true),     // 开启浏览器预览文件, 默认为 false
		filesystem.WithPathPrefix(""),   // PathPrefix定义了一个前缀，当从FileSystem读取文件时, 会添加到文件路径中, 在使用Go 1.16 embed.FS时使用
		filesystem.WithNotFoundFile(""), // 设置未访问到相应文件的自定义页面或数据
		filesystem.WithIndexFile(""),    // 设置访问设置目录的主页内容的路径
		filesystem.WithMaxAge(0),        // 设置文件响应中的Cache-Control HTTP头的值。MaxAge以秒为单位定义
	))
	h.Spin()
}

```
使用 embed.FS 的示例[示例](./examples/main.go)

