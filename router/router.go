package router

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
	defaultHandlers []Handle
}

type Context struct {
	Params httprouter.Params
	Next   Handle
	Meta   map[string]interface{}
}

type Handle func(w http.ResponseWriter, r *http.Request, c Context) error

func New() *Router {
	return &Router{httprouter.New(), []Handle{}}
}

func (r *Router) Default(handlers ...Handle) {
	r.defaultHandlers = handlers
}

func (r *Router) Handle(method, path string, handlers ...Handle) {
	var nextHandler Handle
	defaultHandlers := r.defaultHandlers

	r.Router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		hs := []Handle{}
		hs = append(hs, defaultHandlers...)
		hs = append(hs, handlers...)

		nextHandler = func(w http.ResponseWriter, r *http.Request, c Context) error {
			if len(hs) == 0 {
				return nil
			}
			// get the next handler
			h := hs[0]
			// remove the next handler from handlers
			hs = hs[1:]
			c.Next = nextHandler
			return h(w, r, c)
		}

		if err := nextHandler(w, r, Context{p, nil, map[string]interface{}{}}); err != nil {
			// log the error to stdout
			log.Println(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
	})
}

// GET is a shortcut for router.Handle("GET", path, handlers...)
func (r *Router) GET(path string, handlers ...Handle) {
	r.Handle("GET", path, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handlers...)
func (r *Router) HEAD(path string, handlers ...Handle) {
	r.Handle("HEAD", path, handlers...)
}

// POST is a shortcut for router.Handle("POST", path, handlers...)
func (r *Router) POST(path string, handlers ...Handle) {
	r.Handle("POST", path, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handlers...)
func (r *Router) PUT(path string, handlers ...Handle) {
	r.Handle("PUT", path, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handlers...)
func (r *Router) PATCH(path string, handlers ...Handle) {
	r.Handle("PATCH", path, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handlers...)
func (r *Router) DELETE(path string, handlers ...Handle) {
	r.Handle("DELETE", path, handlers...)
}
