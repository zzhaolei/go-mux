package mux

import "net/http"

type Mux struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

func New() *Mux {
	return &Mux{
		mux: http.NewServeMux(),
	}
}

func (m *Mux) Use(mw ...func(http.Handler) http.Handler) {
	m.middleware = append(m.middleware, mw...)
}

func (m *Mux) Group(fn func(*Mux)) {
	mm := *m
	fn(&mm)
}

func (m *Mux) Handle(pattern string, handler http.Handler) {
	m.handle(pattern, handler)
}

func (m *Mux) HandleFunc(pattern string, handler http.HandlerFunc) {
	m.handle(pattern, handler)
}

func (m *Mux) handle(pattern string, handler http.Handler) {
	finalHandler := handler
	for i := len(m.middleware) - 1; i >= 0; i-- {
		finalHandler = m.middleware[i](finalHandler)
	}
	m.mux.Handle(pattern, finalHandler)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.mux.ServeHTTP(w, req)
}
