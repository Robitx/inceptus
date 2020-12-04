package route

import (
	"net/http"
	"path/filepath"
	"strings"
)

type filesSystem struct {
	fs http.FileSystem
}

func (fs filesSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := fs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// FilesServer gives you a (sub)router that serves only files
// from a directory (it hides your direcory structure) - accessing
// a (sub)directory will lead to index.html or StatusNotFound
//
// A full path to served dir is required.
//
// If mounted as a subrouter under certain path XY,
// specify stripPrefix equal to the mounting path XY.
//
// Here is an example usage:
//	 router := route.New()
//	 ...
//	 workDir, _ := os.Getwd()
//	 dir := filepath.Join(workDir, "relative_path_to_dir")
//	 router.Mount("/static", route.FileServer("/static", dir))
func FilesServer(stripPrefix string, dir string) Router {
	if strings.ContainsAny(stripPrefix, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}
	router := New()

	fileServer := http.FileServer(filesSystem{http.Dir(dir)})

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix(stripPrefix, fileServer)
		fs.ServeHTTP(w, r)
	})

	return router
}
