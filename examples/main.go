package main

import (
	"embed"
	"github.com/Skyenought/filesystem"
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
	h.Use(filesystem.New("/fs", http.FS(fs), filesystem.WithBrowse(true)))
	h.Spin()
}
