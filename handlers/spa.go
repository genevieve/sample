package handlers

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SPA struct {
	content embed.FS
}

func NewSPA(content embed.FS) *SPA {
	return &SPA{content: content}
}

func (h SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join("web/build", path)

	// check whether a file exists at the given path
	_, err = fs.Stat(h.content, path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		b, err := h.content.ReadFile("web/build/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(b))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dir, err := fs.Sub(h.content, "web/build")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the build dir
	http.FileServer(http.FS(dir)).ServeHTTP(w, r)
}
