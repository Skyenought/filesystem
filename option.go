package filesystem

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// option defines the config for middleware.
type option struct {
	next         app.HandlerFunc
	pathPrefix   string
	browse       bool
	index        string
	maxAge       int
	notFoundFile string
}

type Option func(o *option)

func newOption(opts []Option) *option {
	cfg := new(option)
	cfg.index = "index.html"
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

// WithCallback defines a function to skip this middleware when returned true.
func WithCallback(next app.HandlerFunc) Option {
	return func(o *option) {
		o.next = next
	}
}
