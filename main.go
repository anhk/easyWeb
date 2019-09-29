package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
)

var (
	lsn  = flag.String("l", ":80", "listen address")
	path = flag.String("x", "./", "path")
)

func getPath(path string) string {
	f, _ := exec.LookPath(path)
	p, _ := filepath.Abs(f)
	return p
}

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		etagHeaders := []string{
			"ETag",
			"If-Modified-Since",
			"If-Match",
			"If-None-Match",
			"If-Range",
			"If-Unmodified-Since",
		}

		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	fmt.Println("Path: ", getPath(*path))
	fmt.Println("Access: ", *lsn)

	http.Handle("/", NoCache(http.FileServer(http.Dir(getPath(*path)))))
	http.ListenAndServe(*lsn, nil)
}
