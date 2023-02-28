package filesystem

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

// trimRight is the equivalent of strings.TrimRight
func trimRight(s string, cutset byte) string {
	lenStr := len(s)
	for lenStr > 0 && s[lenStr-1] == cutset {
		lenStr--
	}
	return s[:lenStr]
}

func getFileExtension(p string) string {
	n := strings.LastIndexByte(p, '.')
	if n < 0 {
		return ""
	}
	return p[n:]
}

func dirList(c *app.RequestContext, f http.File) error {
	fileinfos, err := f.Readdir(-1)
	if err != nil {
		return fmt.Errorf("failed to read dir: %w", err)
	}
	fm := make(map[string]os.FileInfo, len(fileinfos))
	filenames := make([]string, 0, len(fileinfos))
	for _, fi := range fileinfos {
		name := fi.Name()
		fm[name] = fi
		filenames = append(filenames, name)
	}

	basePathEscaped := html.EscapeString(string(c.Path()))
	_, _ = fmt.Fprintf(c, "<html><head><title>%s</title><style>.dir { font-weight: bold }</style></head><body>", basePathEscaped)
	_, _ = fmt.Fprintf(c, "<h1>%s</h1>", basePathEscaped)
	_, _ = fmt.Fprint(c, "<ul>")

	if len(basePathEscaped) > 1 {
		parentPathEscaped := html.EscapeString(trimRight(string(c.Path()), '/') + "/..")
		_, _ = fmt.Fprintf(c, `<li><a href="%s" class="dir">..</a></li>`, parentPathEscaped)
	}
	sort.Strings(filenames)
	for _, name := range filenames {
		pathEscaped := html.EscapeString(path.Join(string(c.Path()) + "/" + name))
		fi := fm[name]
		auxStr := "dir"
		className := "dir"
		if !fi.IsDir() {
			auxStr = fmt.Sprintf("file, %d bytes", fi.Size())
			className = "file"
		}
		_, _ = fmt.Fprintf(c, `<li><a href="%s" class="%s">%s</a>, %s, last modified %s</li>`,
			pathEscaped, className, html.EscapeString(name), auxStr, fi.ModTime())
	}
	_, _ = fmt.Fprint(c, "</ul></body></html>")
	c.Response.Header.SetContentType("text/html")
	return nil
}
