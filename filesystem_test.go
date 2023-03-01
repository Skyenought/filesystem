package filesystem

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"net/http"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	fsFiles = "examples/testdata/fs"
)

// go test -run TestFileSystem
func TestFileSystem(t *testing.T) {
	t.Parallel()
	h := server.Default(server.WithHostPorts(":3000"))
	h.Use(New("/test", http.Dir(fsFiles)))
	h.Use(New("/dir", http.Dir(fsFiles), WithBrowse(true)))
	h.GET("/", func(_ context.Context, c *app.RequestContext) {
		c.String(200, "Hello world!")
	})
	h.Use(New("/spatest", http.Dir(fsFiles), WithIndexFile("index.html"), WithNotFoundFile("index.html")))
	h.Use(New("/prefix", http.Dir(fsFiles), WithPathPrefix("img")))
	go h.Spin()
	time.Sleep(1 * time.Second)

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
			testClient, _ := client.NewClient()
			req, resp := &protocol.Request{}, &protocol.Response{}
			req.SetRequestURI(tt.url)
			req.SetMethod(consts.MethodGet)
			_ = testClient.Do(context.Background(), req, resp)
			assert.DeepEqual(t, tt.statusCode, resp.StatusCode())
			assert.DeepEqual(t, tt.contentType, resp.Header.Get("Content-Type"))
		})
	}
}
