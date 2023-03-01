package main

import (
	"context"
	"github.com/Skyenought/filesystem"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default(
		server.WithExitWaitTime(500*time.Millisecond),
		server.WithHostPorts(":3000"),
	)
	h.GET("/", func(_ context.Context, c *app.RequestContext) {
		c.String(200, "Hello World!")
	})
	h.Use(filesystem.New("/dir", http.Dir("./testdata"),
		filesystem.WithBrowse(true)),
	)
	h.Spin()
}
