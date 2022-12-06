package ming

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func (r *Router) Handle(method, path string, handler fasthttp.RequestHandler) {
	if !strings.HasPrefix(path, "/") {
		panic("path must begin with \"/\" in \"" + path + "\"")
	}
	r.trees.Add(&Node{
		method:  method,
		path:    path,
		handler: handler,
	})
}

func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}
	path := string(ctx.Path())
	method := GetMethod(ctx)
	if nodeFindByPath := r.trees.FindPath(path); nodeFindByPath.Len() != 0 {
		if node := nodeFindByPath.FindMethod(method); node != nil {
			handler := node.GetHandler()
			handler(ctx)
		} else {
			if node := nodeFindByPath.GetMethodAll(); node != nil {
				handler := node.GetHandler()
				handler(ctx)
			} else {
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed(ctx)
				} else {
					ctx.Error("method not allowed", fasthttp.StatusMethodNotAllowed)
				}
			}
		}
	} else {
		if r.NotFound != nil {
			r.NotFound(ctx)
		} else {
			ctx.Error(fmt.Sprintf("%s %s not found", method, path), fasthttp.StatusNotFound)
		}
	}
}

func (r *Router) Get(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodGet, path, handler)
}

func (r *Router) Head(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodHead, path, handler)
}

func (r *Router) Post(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodDelete, path, handler)
}

func (r *Router) Connect(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodConnect, path, handler)
}

func (r *Router) Options(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodOptions, path, handler)
}

func (r *Router) Trace(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodTrace, path, handler)
}
func (r *Router) All(path string, handler fasthttp.RequestHandler) {
	r.Handle("ALL", path, handler)
}

func (r *Router) Static(rootPath string, IsIndexPage bool) {
	fs := &fasthttp.FS{
		Root:               rootPath,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: IsIndexPage,
	}
	r.NotFound = fs.NewRequestHandler()
}
