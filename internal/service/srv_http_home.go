package service

import (
	"io"
	"net/http"
)

func (s Service) handleGetHomeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `
	<!DOCTYPE html>
	<html>
		<head>
			<title>Home - Julien Sellier</title>
			<style>
				body {
					background-color: black;
					color: white;
					font-family: monospace;
				}
			</style>
		</head>

		<body>
			<h1>Coucou</h1>
			<p>Hello world!</p>
		</body>
	</html>
	`)
}
