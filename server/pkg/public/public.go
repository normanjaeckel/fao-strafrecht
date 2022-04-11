/*
Package public provides a http handler for the the public client files like
index.html and all assets (JS, CSS, ...).
*/
package public

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed files
var files embed.FS

func Files() http.Handler {
	root, err := fs.Sub(files, "files")
	if err != nil {
		panic("Error when getting subtree of embedded filesystem. " +
			"This should never ever happen.")
	}
	return http.FileServer(http.FS(root))
}
