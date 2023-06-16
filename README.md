# filesystem

[English](./README_EN.md)

[Hertz](https://github.com/cloudwego/hertz) 的文件系统中间件，使用户可以直接使用原生的 `http.Dir`, `http.FS` 等进行静态文件的映射。

## 安装:

```shell
go get -u github.com/Skyenought/filesystem
```

## 简易使用示例:

```go
package main

// ...

func main() {
	h := server.Default()
	// 已过时, 请使用 filesystem.NewFSHandler 替换
	h.Use(filesystem.New("/", http.Dir("./testdata"))) // 需要访问的文件夹的相对路径
	filesystem.NewFSHandler(h, "/dir",  http.Dir("./testdata"),
		filesystem.WithBrowse(true),
	)
	h.Spin()
}
```
## 完整使用示例

```go
package main

func main() {
	h := server.Default()
	filesystem.NewFSHandler(h, "/test", http.FS(fs),
		filesystem.WithBrowse(true),
	)
	h.Use(filesystem.New("/dir", http.Dir("./testdata"), // 支持使用 embed.FS, 即 http.FS
		filesystem.WithBrowse(true),     // 开启浏览器预览文件, 默认为 false
		filesystem.WithPathPrefix(""),   // PathPrefix定义了一个前缀，当从FileSystem读取文件时, 会添加到文件路径中, 在使用Go 1.16 embed.FS时使用
		filesystem.WithNotFoundFile(""), // 设置未访问到相应文件的自定义页面或数据
		filesystem.WithIndexFile(""),    // 设置访问设置目录的主页内容的路径
		filesystem.WithMaxAge(0),        // 设置文件响应中的Cache-Control HTTP头的值。MaxAge以秒为单位定义
		filesystem.WithPreHandler(func(c context.Context, ctx *app.RequestContext) (func(), bool) {
			if ctx.Request.Header.Get("token") != "123" {
				return func() {
					ctx.String(http.StatusUnauthorized, "Authorize Fail!")
				}, false
			}
			return nil, true
		}), // 设置一个预处理函数, 用于在访问文件之前进行一些操作, 如果返回 false, 则不会继续访问文件
	))
	h.Spin()
}

```
使用 embed.FS 的示例[示例](./examples/main.go)

## 原理

使用 any 节点劫持访问路径, 并实现 FS 接口使得 hertz 直接支持直接使用原生的 `http.Dir`, `http.FS` 等实现 FS 接口的方法进行文件管理或访问
