package main

import (
	"context"
	"embed"
	"github.com/Skyenought/filesystem"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

//go:embed testdata/*
var fs embed.FS

func main() {
	h := server.Default(
		server.WithExitWaitTime(500*time.Millisecond),
		server.WithHostPorts(":3000"),
	)
	h.GET("/", func(_ context.Context, c *app.RequestContext) {
		c.String(200, "Hello World!")
	})
	h.Use(filesystem.New("/dir", http.FS(fs),
		filesystem.WithBrowse(true),
		filesystem.WithPreHandler(func(ctx context.Context, c *app.RequestContext) bool {
			get := c.Request.Header.Get("token")
			if get != "123" {
				return false
			}
			return true
		}),
	))
	h.Spin()
}
