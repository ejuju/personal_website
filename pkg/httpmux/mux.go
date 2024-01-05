package httpmux

import (
	"io"
	"net/http"
)

type Endpoints []Endpoint

type Endpoint struct {
	Handler  http.HandlerFunc
	Matchers []RequestMatcher
}

type RequestMatcher func(r *http.Request) (matched bool)

func (mux Endpoints) Append(h http.HandlerFunc, matchers ...RequestMatcher) Endpoints {
	return append(mux, Endpoint{Handler: h, Matchers: matchers})
}

func (mux Endpoints) Route(path, method string, h http.HandlerFunc) Endpoints {
	return mux.Append(h, MatchPath(path), MatchMethod(method))
}

func (mux Endpoints) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, endpoint := range mux {
		if endpoint.Match(r) {
			endpoint.Handler.ServeHTTP(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, http.StatusText(http.StatusNotFound))
}

func (ep Endpoint) Match(r *http.Request) (matched bool) {
	for _, matcher := range ep.Matchers {
		if !matcher(r) {
			return false
		}
	}
	return true
}

func MatchPath(p string) RequestMatcher {
	return func(r *http.Request) (matched bool) { return r.URL.Path == p }
}

func MatchGet(r *http.Request) (matched bool)     { return r.Method == http.MethodGet }
func MatchPost(r *http.Request) (matched bool)    { return r.Method == http.MethodPost }
func MatchPut(r *http.Request) (matched bool)     { return r.Method == http.MethodPut }
func MatchPatch(r *http.Request) (matched bool)   { return r.Method == http.MethodPatch }
func MatchDelete(r *http.Request) (matched bool)  { return r.Method == http.MethodDelete }
func MatchHead(r *http.Request) (matched bool)    { return r.Method == http.MethodHead }
func MatchOptions(r *http.Request) (matched bool) { return r.Method == http.MethodOptions }
func MatchMethod(method string) RequestMatcher {
	return func(r *http.Request) (matched bool) { return r.Method == method }
}
