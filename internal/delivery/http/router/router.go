package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	GET(uri string, f func(w http.ResponseWriter, r *http.Request))
	POST(uri string, f func(w http.ResponseWriter, r *http.Request))
	PUT(uri string, f func(w http.ResponseWriter, r *http.Request))
	DELETE(uri string, f func(w http.ResponseWriter, r *http.Request))
	GETWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc)
	POSTWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc)
	PUTWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc)
	DELETEWithMiddleware(uri string, f func(w http.ResponseWriter, r *http.Request), middlewares ...mux.MiddlewareFunc)
	OPTIONS(uri string)
	Mux() *mux.Router
	SERVE(port string)
}
