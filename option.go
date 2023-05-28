package filesystem

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"strings"
)

// option defines the config for middleware.
type option struct {
	preHandler
	root         http.FileSystem
	pathPrefix   string
	browse       bool
	index        string
	maxAge       int
	notFoundFile string
}

type Option func(o *option)

type preHandler func(c context.Context, ctx *app.RequestContext) (func(), bool)

func newOption(root http.FileSystem, opts []Option) *option {
	cfg := &option{
		root:  root,
		index: "index.html",
	}
	for _, optionFuc := range opts {
		optionFuc(cfg)
	}
	if !strings.HasPrefix(cfg.index, "/") {
		cfg.index = "/" + cfg.index
	}

	if cfg.notFoundFile != "" && !strings.HasPrefix(cfg.notFoundFile, "/") {
		cfg.notFoundFile = "/" + cfg.notFoundFile
	}

	if cfg.pathPrefix != "" && !strings.HasPrefix(cfg.pathPrefix, "/") {
		cfg.pathPrefix = "/" + cfg.pathPrefix
	}
	return cfg
}

// WithPathPrefix PathPrefix defines a prefix to be added to a filepath when
// reading a file from the FileSystem.
//
// Use when using Go 1.16 embed.FS
func WithPathPrefix(prefix string) Option {
	return func(o *option) {
		o.pathPrefix = prefix
	}
}

// WithBrowse Enable directory browsing.
func WithBrowse(enabled bool) Option {
	return func(o *option) {
		o.browse = enabled
	}
}

// WithIndexFile Index file for serving a directory.
func WithIndexFile(index string) Option {
	return func(o *option) {
		o.index = index
	}
}

// WithMaxAge The value for the Cache-Control HTTP-header
// that is set on the file response. MaxAge is defined in seconds.
func WithMaxAge(age int) Option {
	return func(o *option) {
		o.maxAge = age
	}
}

// WithNotFoundFile File to return if path is not found. Useful for SPA's.
func WithNotFoundFile(path string) Option {
	return func(o *option) {
		o.notFoundFile = path
	}
}

// WithPreHandler PreHandler is executed before the filesystem middleware.
// If the handler returns false, the middleware will abort with a 401 status by default.
//
// If the handler returns false and a custom fallback function is returned,
// the middleware will abort with the custom fallback function.
func WithPreHandler(handler preHandler) Option {
	return func(o *option) {
		o.preHandler = handler
	}
}
