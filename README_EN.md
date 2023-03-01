# filesystem

[中文](./README.md)

The file system middleware for [Hertz](https://github.com/cloudwego/hertz), which enables you to serve files from a directory.

Note:

⚠️ **Do not support `:params` in the prefix path**

⚠️ **This version is not perfect, and the middleware can only be used once with `h.Use`.**

Installation:

```shell
go get -u github.com/Skyenought/filesystem
```

Simple usage example:

```go
package main

// ...

func main() {
	h := server.Default()
	h.Use(filesystem.New("/", http.Dir("./testdata"))) // The relative path of the folder to be accessed
	h.Spin()
}
```

Complete usage example:

```go
package main

func main() {
	h := server.Default()
	h.Use(filesystem.New("/dir", http.Dir("./testdata"), // Supports using embed.FS, that is, http.FS
		filesystem.WithBrowse(true),     // Enable browsing files in the directory, default is false
		filesystem.WithPathPrefix(""),   // PathPrefix defines a prefix to be added to a filepath when reading a file from the FileSystem. Use when using Go 1.16 embed.FS.
		filesystem.WithNotFoundFile(""), // Set custom page or data for the file that has not been accessed
		filesystem.WithIndexFile(""),    // Set the path to the home page content of the accessed setting directory
		filesystem.WithMaxAge(0),        // Set the value for the Cache-Control HTTP-header that is set on the file response. MaxAge is defined in seconds.
	))
	h.Spin()
}
```

Example using embed.FS [example](./examples/main.go)
