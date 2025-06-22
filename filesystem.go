package filesystem

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// NewFSHandler creates a new middleware handler.
func NewFSHandler(engine *server.Hertz, relpath string, root http.FileSystem, opts ...Option) {
	var prefix string
	prefix = relpath
	cfg := newOption(root, opts)
	if cfg.pathPrefix != "" && !strings.HasPrefix(cfg.pathPrefix, "/") {
		prefix = "/" + cfg.pathPrefix
	}
	cacheControlStr := "public, max-age=" + strconv.Itoa(cfg.maxAge)

	logicFunc := func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Method())

		// Check that the request has the correct headers, and that the request is well-formed.
		// If the request does not have the correct headers, or is malformed, return an error.
		// Otherwise, return nil.
		if cfg.preHandler != nil {
			customFallback, ok := cfg.preHandler(ctx, c)
			if !ok {
				if customFallback == nil {
					c.AbortWithStatus(consts.StatusUnauthorized)
					return
				}
				customFallback()
				return
			}
		}

		path := c.Param("filepath")
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		if cfg.pathPrefix != "" {
			// PathPrefix already has a "/" relpath
			path = cfg.pathPrefix + path
		}
		if len(path) > 1 {
			path = trimRight(path, '/')
		}

		file, err := cfg.root.Open(path)
		if err != nil && os.IsNotExist(err) && cfg.notFoundFile != "" {
			file, err = cfg.root.Open(cfg.notFoundFile)
		}
		if err != nil {
			if os.IsNotExist(err) {
				hlog.SystemLogger().Errorf("Cannot open file or Directory, path: %s, err = %s", path, err)
				c.AbortWithMsg("Cannot open file or Directory", consts.StatusNotFound)
				return
			}
			hlog.SystemLogger().Errorf("Failed to open: %s", err)
			c.AbortWithMsg("Cannot open file or Directory", consts.StatusNotFound)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			hlog.SystemLogger().Errorf("failed to stat: %s", err)
			c.AbortWithMsg("failed to stat", consts.StatusInternalServerError)
			return
		}

		// Serve index if relpath is directory
		if stat.IsDir() {
			indexPath := trimRight(path, '/') + cfg.index
			index, err := cfg.root.Open(indexPath)
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
					c.String(consts.StatusInternalServerError, err.Error())
					hlog.Errorf("show dirList fail, err: %s", err)
				}
				return
			}
			c.AbortWithStatus(consts.StatusForbidden)
			return
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

		c.Response.Header.SetContentType(getMIME(getFileExtension(stat.Name())))
		if !modTime.IsZero() {
			c.Response.Header.Set(consts.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

		if method == consts.MethodGet {
			if cfg.maxAge > 0 {
				c.Response.Header.Set("Cache-Control", cacheControlStr)
			}
			c.Response.SetBodyStream(file, contentLength)
			return
		}

		if method == consts.MethodHead {
			c.Request.ResetBody()
			c.Response.SkipBody = true
			c.Response.Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				hlog.SystemLogger().Errorf("failed to close: %s", err.Error())
				c.AbortWithMsg("fail to close file", consts.StatusInternalServerError)
				return
			}
			return
		}
		c.Next(ctx)
	}
	engine.GET(prefix+"/*filepath", logicFunc)
	engine.HEAD(prefix+"/*filepath", logicFunc)
}

// New creates a new middleware handler.
//
// filesystem does not handle url encoded values (for example spaces)
// on its own.
//
// Deprecated: use NewFSHandler instead.
func New(urlPrefix string, root http.FileSystem, opts ...Option) app.HandlerFunc {
	cfg := newOption(root, opts)

	var once sync.Once
	var prefix string
	cacheControlStr := "public, max-age=" + strconv.Itoa(cfg.maxAge)

	return func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Method())

		// Don't execute middleware if method != "GET" OR "HEAD"
		if method != http.MethodGet && method != http.MethodHead {
			c.Next(ctx)
			return
		}

		// Check that the request has the correct headers, and that the request is well-formed.
		// If the request does not have the correct headers, or is malformed, return an error.
		// Otherwise, return nil.
		if cfg.preHandler != nil {
			customFallbackFunc, ok := cfg.preHandler(ctx, c)
			if !ok {
				if customFallbackFunc != nil {
					customFallbackFunc()
					return
				}
				c.AbortWithStatus(consts.StatusUnauthorized)
				return
			}
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
		file, err := cfg.root.Open(path)
		if err != nil && os.IsNotExist(err) && cfg.notFoundFile != "" {
			file, err = cfg.root.Open(cfg.notFoundFile)
		}
		if err != nil {
			if os.IsNotExist(err) {
				hlog.SystemLogger().Errorf("Cannot open file or Directory, path: %s, err = %s", path, err)
				c.AbortWithMsg("Cannot open file or Directory", consts.StatusNotFound)
				return
			}
			hlog.SystemLogger().Errorf("Failed to open: %s", err)
			c.AbortWithMsg("Cannot open file or Directory", consts.StatusNotFound)
			return
		}

		stat, err := file.Stat()
		if err != nil {
			hlog.SystemLogger().Errorf("failed to stat: %s", err)
			c.AbortWithMsg("failed to stat", consts.StatusInternalServerError)
			return
		}

		// Serve index if urlPrefix is directory
		if stat.IsDir() {
			indexPath := trimRight(path, '/') + cfg.index
			index, err := cfg.root.Open(indexPath)
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
					c.String(consts.StatusInternalServerError, err.Error())
					hlog.Errorf("show dirList fail, err: %s", err)
				}
				return
			}
			c.AbortWithStatus(consts.StatusForbidden)
			return
		}

		modTime := stat.ModTime()
		contentLength := int(stat.Size())

		c.Response.Header.SetContentType(getMIME(getFileExtension(stat.Name())))
		if !modTime.IsZero() {
			c.Response.Header.Set(consts.HeaderLastModified, modTime.UTC().Format(http.TimeFormat))
		}

		if method == consts.MethodGet {
			if cfg.maxAge > 0 {
				c.Response.Header.Set("Cache-Control", cacheControlStr)
			}
			c.Response.SetBodyStream(file, contentLength)
			return
		}

		if method == consts.MethodHead {
			c.Request.ResetBody()
			c.Response.SkipBody = true
			c.Response.Header.SetContentLength(contentLength)
			if err := file.Close(); err != nil {
				hlog.SystemLogger().Errorf("failed to close: %s", err.Error())
				c.AbortWithMsg("fail to close file", consts.StatusInternalServerError)
				return
			}
			return
		}
		c.Next(ctx)
	}
}
