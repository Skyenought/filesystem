package filesystem

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// New creates a new middleware handler.
//
// filesystem does not handle url encoded values (for example spaces)
// on it's own.
func New(urlPrefix string, root http.FileSystem, opts ...Option) app.HandlerFunc {
	cfg := newOption(opts)

	var once sync.Once
	var prefix string
	cacheControlStr := "public, max-age=" + strconv.Itoa(cfg.maxAge)

	return func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Method())

		// Don't execute middleware if method != "GET" OR "HEAD"
		if method != http.MethodGet && method != http.MethodHead {
			c.Next(ctx)
		}

		once.Do(func() {
			prefix = urlPrefix
		})

		path := strings.TrimPrefix(string(c.Path()), prefix)
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		if cfg.pathPrefix != "" {
			// PathPrefix already has a "/" prefix
			path = cfg.pathPrefix + path
		}
		if len(path) > 1 {
			path = trimRight(path, '/')
		}
		file, err := root.Open(path)
		if err != nil && os.IsNotExist(err) && cfg.notFoundFile != "" {
			file, err = root.Open(cfg.notFoundFile)
		}
		if err != nil {
			if os.IsNotExist(err) {
				c.AbortWithStatus(consts.StatusNotFound)
			}
			hlog.Errorf("failed to open: %s", err.Error())
			return
		}

		stat, err := file.Stat()
		if err != nil {
			hlog.Errorf("failed to stat: %s", err)
			return
		}

		// Serve index if urlPrefix is directory
		if stat.IsDir() {
			indexPath := trimRight(path, '/') + cfg.index
			index, err := root.Open(indexPath)
			if err == nil {
				indexStat, err := index.Stat()
				if err == nil {
					file = index
					stat = indexStat
				}
			}
		}

		// Browse directory if no index found and browsing is enabled
		if stat.IsDir() {
			if cfg.browse {
				if err := dirList(c, file); err != nil {
					hlog.Errorf("show dirList fail, err: %s", err)
				}
				return
			}
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

		c.Response.Header.SetContentType(getFileExtension(stat.Name()))
		if !modTime.IsZero() {
			c.Response.Header.Set(consts.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

		if method == consts.MethodGet {
			if cfg.maxAge > 0 {
				c.Response.Header.Set("Cache-Control", cacheControlStr)
			}
			c.Response.SetBodyStream(file, contentLength)
		}

		if method == consts.MethodHead {
			c.Request.ResetBody()
			c.Response.SkipBody = true
			c.Response.Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				hlog.Errorf("failed to close: %s", err)
			}
		}
	}
}
