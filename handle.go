package ming

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE and CONNECT requests,
// the respective shortcut methods can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
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

// Handler handles all requests for the router.
// It implements the fasthttp.RequestHandler interface.
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

// Get registers a route for GET requests to the specified path with the provided handler.
func (r *Router) Get(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodGet, path, handler)
}

// Head registers a route for HEAD requests to the specified path with the provided handler.
func (r *Router) Head(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodHead, path, handler)
}

// Post registers a route for POST requests to the specified path with the provided handler.
func (r *Router) Post(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPost, path, handler)
}

// Put registers a route for PUT requests to the specified path with the provided handler.
func (r *Router) Put(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPut, path, handler)
}

// Patch registers a route for PATCH requests to the specified path with the provided handler.
func (r *Router) Patch(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodPatch, path, handler)
}

// Delete registers a route for DELETE requests to the specified path with the provided handler.
func (r *Router) Delete(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodDelete, path, handler)
}

// Connect registers a route for CONNECT requests to the specified path with the provided handler.
func (r *Router) Connect(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodConnect, path, handler)
}

// Options registers a route for OPTIONS requests to the specified path with the provided handler.
func (r *Router) Options(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodOptions, path, handler)
}

// Trace registers a route for TRACE requests to the specified path with the provided handler.
func (r *Router) Trace(path string, handler fasthttp.RequestHandler) {
	r.Handle(fasthttp.MethodTrace, path, handler)
}
// All registers a route that matches all HTTP methods for the specified path with the provided handler.
// This is useful for routes that need to respond to multiple HTTP methods with the same logic.
func (r *Router) All(path string, handler fasthttp.RequestHandler) {
	r.Handle("ALL", path, handler)
}

// Static configures the router to serve static files from the specified directory.
// The rootPath parameter specifies the root directory from which to serve files.
// When IsIndexPage is true, directory listings will be generated for directories that don't have an index file.
func (r *Router) Static(rootPath string, IsIndexPage bool) {
	fs := &fasthttp.FS{
		Root:               rootPath,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: IsIndexPage,
	}
	r.NotFound = fs.NewRequestHandler()
}
