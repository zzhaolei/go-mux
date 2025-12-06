package mux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write(fmt.Appendf(nil, `{"id": %s}`, r.PathValue("id")))
}

func TestHandle(t *testing.T) {
	app := New()
	app.Handle("GET /users/{id}", http.HandlerFunc(handle))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected := `{"id": 666}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func TestHandleFunc(t *testing.T) {
	app := New()
	app.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write(fmt.Appendf(nil, `{"id": %s}`, r.PathValue("id")))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected := `{"id": 666}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func TestMuxGroup(t *testing.T) {
	app := New()

	app.Group(func(m *Mux) {
		m.Use(logger)

		m.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write(fmt.Appendf(nil, `{"user_id": "%s"}`, r.URL.Query().Get("test")))
		})
	})

	app.Group(func(m *Mux) {
		m.Use(logger)

		m.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write(fmt.Appendf(nil, `{"order_id": "%s"}`, r.URL.Query().Get("test")))
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected := `{"user_id": "/users/666"}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/orders/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected = `{"order_id": "/orders/666"}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	app := New()
	app.Use(logger)
	app.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write(fmt.Appendf(nil, `{"id": "%s"}`, r.URL.Query().Get("test")))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/asdf", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected := `{"id": "/users/asdf"}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func TestMuxMixinGroup(t *testing.T) {
	app := New()

	app.Group(func(m *Mux) {
		m.Use(logger)

		m.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			w.Write(fmt.Appendf(nil, `{"user_id": "%s"}`, r.URL.Query().Get("test")))
		})
	})

	app.HandleFunc("GET /users/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write(fmt.Appendf(nil, `{"order_id": "%s"}`, r.URL.Query().Get("test")))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected := `{"user_id": "/users/666"}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/test/666", nil)
	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	expected = `{"order_id": ""}`
	if w.Body.String() != expected {
		t.Fatalf("expected body %s, got %s", expected, w.Body.String())
	}
}

func logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		query.Add("test", r.URL.Path)
		r.URL.RawQuery = query.Encode()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
