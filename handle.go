package ming

import "net/http"

func (r *Router) Use(middlewares ...MiddlewareType) {
	if len(middlewares) > 0 {
		r.middleware = append(r.middleware, middlewares...)
	}
}

func (r *Router) NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	http.NotFound(w, req)
}


func (r *Router) Get(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handlerFn)
}

func (r *Router) Post(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handlerFn)
}

func (r *Router) Put(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handlerFn)
}

func (r *Router) Patch(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handlerFn)
}

func (r *Router) Delete(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handlerFn)
}

func (r *Router) Options(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodOptions, path, handlerFn)
}

func (r *Router) Head(path string, handlerFn http.HandlerFunc) {
	r.Handle(http.MethodHead, path, handlerFn)
}
