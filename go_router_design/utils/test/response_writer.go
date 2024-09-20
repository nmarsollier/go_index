package test

import (
	"net/http"
	"testing"
)

func ResponseWriter(t *testing.T) *FakeResponseWriter {
	return &FakeResponseWriter{
		t:       t,
		headers: make(http.Header),
	}
}

type FakeResponseWriter struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func (r *FakeResponseWriter) Header() http.Header {
	return r.headers
}

func (r *FakeResponseWriter) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *FakeResponseWriter) WriteHeader(status int) {
	r.status = status
}

func (r *FakeResponseWriter) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %#v to equal %#v", string(r.body), body)
	}
}
