package filesystem

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"net/http"
	"testing"
)

func TestFileSystem(t *testing.T) {
	t.Parallel()

	h := server.New()
	NewFSHandler(h, "/test", http.Dir("./examples/testdata/fs"))
	NewFSHandler(h, "/dir", http.Dir("./examples/testdata/fs"))
	h.GET("/", func(ctx context.Context, c *app.RequestContext) { c.String(200, "Hello World!") })
	NewFSHandler(h, "/spatest", http.Dir("./examples/testdata/fs"), WithIndexFile("index.html"), WithNotFoundFile("index.html"))
	NewFSHandler(h, "/prefix", http.Dir("./examples/testdata/fs"), WithPathPrefix("img"))
	tests := []struct {
		name         string
		url          string
		statusCode   int
		contentType  string
		modifiedTime string
	}{
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test/index.html",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200 with suitable content-type",
			url:         "/test/css/style.css",
			statusCode:  200,
			contentType: "text/css",
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile.js",
			statusCode: 404,
		},
		{
			name:       "Should be returns status 404",
			url:        "/test/nofile",
			statusCode: 404,
		},
		{
			name:        "Should be returns status 200",
			url:         "/",
			statusCode:  200,
			contentType: "text/plain; charset=utf-8",
		},
		{
			name:       "Should be returns status 403",
			url:        "/test/img",
			statusCode: 403,
		},
		{
			name:        "Should list the directory contents",
			url:         "/dir/img",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should list the directory contents",
			url:         "/dir/img/",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "Should be returns status 200",
			url:         "/dir/img/fiber.png",
			statusCode:  200,
			contentType: "image/png",
		},
		{
			name:        "Should be return status 200",
			url:         "/spatest/doesnotexist",
			statusCode:  200,
			contentType: "text/html",
		},
		{
			name:        "PathPrefix should be applied",
			url:         "/prefix/fiber.png",
			statusCode:  200,
			contentType: "image/png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := ut.PerformRequest(h.Engine, consts.MethodGet, tt.url, nil)
			response := w.Result()
			assert.DeepEqual(t, tt.statusCode, response.StatusCode())

			if tt.contentType != "" {
				ct := response.Header.Get("Content-Type")
				assert.DeepEqual(t, tt.contentType, ct)
			}
		})
	}
}

func TestPreHandler(t *testing.T) {
	t.Parallel()

	h := server.New()
	NewFSHandler(h, "/", http.Dir("./examples/testdata/fs"),
		WithPreHandler(func(_ context.Context, c *app.RequestContext) (func(), bool) {
			return nil, false
		}),
	)

	w := ut.PerformRequest(h.Engine, consts.MethodGet, "/", nil)
	response := w.Result()
	assert.DeepEqual(t, 401, response.StatusCode())
}
