# go-mux

## Usage

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zzhaolei/go-mux"
)

func main() {
	app := mux.New()
	app.Use(globalLogger)

	app.Group(func(m *mux.Mux) {
		m.Use(logger)

		m.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write(fmt.Appendf(nil, `{"user_id": "%s"}`, r.PathValue("id")))
		})
	})
	if err := http.ListenAndServe(":8080", app); err != nil {
		panic(err)
	}
}

func globalLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("all: %s - %s", r.Method, r.URL.Path)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s", r.Method, r.URL.Path)

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
```
