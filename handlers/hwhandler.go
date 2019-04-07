package handlers

import (
	"fmt"
	"net/http"
)

type HelloWorldHandler struct {}

func (s *HelloWorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello World")
}
