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
	
	if r.trees == nil {
		r.trees = make(map[string]*Tree)
	}
	
	tree := r.trees[method]
	if tree == nil {
		tree = NewTree(method)
		r.trees[method] = tree
	}
	
	tree.addRoute(path, handler)
}

func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}
	
	path := string(ctx.Path())
	method := GetMethod(ctx)
	
	if tree := r.trees[method]; tree != nil {
		if handler, params, tsr := tree.getValue(path, method); handler != nil {
			// Set parameters in context
			for _, param := range params {
				ctx.SetUserValue(param.Key, param.Value)
			}
			handler(ctx)
			return
		} else if tsr && method != fasthttp.MethodConnect {
			// Handle trailing slash redirect
			var redirectPath string
			if len(path) > 1 && path[len(path)-1] == '/' {
				redirectPath = path[:len(path)-1]
			} else {
				redirectPath = path + "/"
			}
			ctx.Response.Header.Set("Location", redirectPath)
			ctx.SetStatusCode(fasthttp.StatusMovedPermanently)
			return
		}
	}
	
	// Try ALL method as fallback
	if tree := r.trees["ALL"]; tree != nil {
		if handler, params, _ := tree.getValue(path, "ALL"); handler != nil {
			// Set parameters in context
			for _, param := range params {
				ctx.SetUserValue(param.Key, param.Value)
			}
			handler(ctx)
			return
		}
	}
	
	// Check if method is allowed for this path
	allowed := make([]string, 0, len(r.trees))
	for m, tree := range r.trees {
		if m == method || m == "ALL" {
			continue
		}
		if handler, _, _ := tree.getValue(path, m); handler != nil {
			allowed = append(allowed, m)
		}
	}
	
	if len(allowed) > 0 {
		allowHeader := strings.Join(allowed, ", ")
		if r.MethodNotAllowed != nil {
			ctx.Response.Header.Set("Allow", allowHeader)
			r.MethodNotAllowed(ctx)
		} else {
			ctx.Error("Method Not Allowed", fasthttp.StatusMethodNotAllowed)
			ctx.Response.Header.Set("Allow", allowHeader)
		}
		return
	}
	
	// Not found
	if r.NotFound != nil {
		r.NotFound(ctx)
	} else {
		ctx.Error(fmt.Sprintf("%s %s not found", method, path), fasthttp.StatusNotFound)
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
