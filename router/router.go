package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/renom/golang-test/response"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	r.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)
	return r
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.RespondMessage(w, http.StatusNotFound, "Page not found")
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	response.RespondMessage(w, http.StatusMethodNotAllowed, "Method isn't allowed")
}
