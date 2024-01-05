package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Service) runHTTPServer() {
	s.HTTPEndpoints = s.HTTPEndpoints.Route("/", http.MethodGet, s.handleGetHomeHTTP)
	s.HTTPEndpoints = s.HTTPEndpoints.Append(s.handleNotFoundHTTP)

	s.Logger.Log("Starting HTTP server...")
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(s.Config.HTTPPort),
		Handler:           s,
		MaxHeaderBytes:    1024 * 1024, // = 1 MiB
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ErrorLog:          log.New(s.Logger.Writer, "HTTP-SERVER: ", 0),
	}
	defer server.Shutdown(context.Background())
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Panic(fmt.Errorf("listen and serve HTTP: %w", err))
	}
}

func (s Service) handleNotFoundHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "\""+r.URL.String()+"\" not found...")
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respw := HTTPResponseWriter{ResponseWriter: w, StartedAt: time.Now()}
	s.HTTPEndpoints.ServeHTTP(&respw, r)
	respw.Duration = time.Since(respw.StartedAt)
	respw.Method = r.Method
	respw.Path = r.URL.Path
	s.Logger.Log(fmt.Sprintf("HTTP " + respw.String()))
}

type HTTPResponseWriter struct {
	http.ResponseWriter
	StartedAt  time.Time
	Duration   time.Duration
	Method     string
	Path       string
	StatusCode int
	Written    int
}

func (l HTTPResponseWriter) String() string {
	return l.StartedAt.Format(time.DateTime) + " " +
		strconv.Itoa(l.StatusCode) + " " +
		l.Method + " " +
		l.Path + " " +
		l.Duration.String() + " " +
		strconv.Itoa(l.Written) + "B"
}

func (l *HTTPResponseWriter) WriteHeader(statusCode int) { l.StatusCode = statusCode }

func (l *HTTPResponseWriter) Write(b []byte) (n int, err error) {
	n, err = l.ResponseWriter.Write(b)
	l.Written += n
	return n, err
}
