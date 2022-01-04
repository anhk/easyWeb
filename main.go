package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	lsn  = flag.String("l", ":80", "listen address")
	path = flag.String("x", "./", "path")
)

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access:", r.URL.String())
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
	fmt.Println("Path: ", *path)
	fmt.Println("Access: ", *lsn)

	http.Handle("/", NoCache(http.FileServer(http.Dir(*path))))
	if err := http.ListenAndServe(*lsn, nil); err != nil {
		fmt.Println("Error:", err.Error())
	}
}
