// Package ming is a high-performance HTTP router based on fasthttp.
//
// Ming provides a priority-based routing system with Radix Tree for efficient request routing.
// It supports named parameters, optional parameters, regex validation, and catch-all routes.
//
// For more details and examples, see https://github.com/hypnguyen1209/ming
package ming

import (
	"fmt"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	DefaultContentType = []byte("text/plain; charset=utf-8")
)

type Router struct {
	trees            map[string]*Tree
	PanicHandler     func(*fasthttp.RequestCtx, interface{})
	NotFound         fasthttp.RequestHandler
	MethodNotAllowed fasthttp.RequestHandler
}

// New creates a new Router instance.
// The router implements fasthttp.RequestHandler and can be directly passed to fasthttp server.
func New() *Router {
	return &Router{
		trees: make(map[string]*Tree),
	}
}

// HostSwitch is a map of host names to request handlers.
// It can be used to implement virtual hosting functionality
// where different hosts are handled by different handlers.
type HostSwitch map[string]fasthttp.RequestHandler

// CheckHost implements the fasthttp.RequestHandler interface.
// It checks the Host header to select the appropriate handler for the request.
// If no handler is found for the host, a 403 Forbidden response is returned.
func (hs HostSwitch) CheckHost(ctx *fasthttp.RequestCtx) {
	if handler := hs[string(ctx.Host())]; handler != nil {
		handler(ctx)
	} else {
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
	}
}

// Run starts a fasthttp server with the router as handler.
// The addr parameter can be either ":8080" for all interfaces or "127.0.0.1:8080" for specific interface.
// This is a convenience function to start the server with default configuration.
func (r *Router) Run(addr string) {
	if strings.HasPrefix(addr, ":") {
		log.Fatal(fasthttp.ListenAndServe(addr, r.Handler))
	} else {
		port := ":" + strings.Split(addr, ":")[1]
		hs := make(HostSwitch)
		hs[addr] = r.Handler
		log.Fatal(fasthttp.ListenAndServe(port, hs.CheckHost))
	}
}

// Query returns the query parameter value for the given key.
// If no value is found, an empty byte slice is returned.
func Query(ctx *fasthttp.RequestCtx, str string) []byte {
	return ctx.QueryArgs().Peek(str)
}

// SetHeader sets the response header with the given key and value.
func SetHeader(ctx *fasthttp.RequestCtx, key string, value string) {
	ctx.Response.Header.Set(key, value)
}

// Body returns the raw request body.
func Body(ctx *fasthttp.RequestCtx) []byte {
	return ctx.Request.Body()
}

func (r *Router) recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

// UserValue retrieves a user value from the request context by key.
// This is a convenience wrapper around fasthttp.RequestCtx.UserValue
// and is primarily used to access route parameters.
func UserValue(ctx *fasthttp.RequestCtx, key string) interface{} {
	return ctx.UserValue(key)
}

// Param is a convenience function to get a parameter value as string
// Param returns the value of the URL parameter from the request context.
// For example, if the route is defined as "/user/{id}" and the request path is "/user/123",
// then Param(ctx, "id") would return "123".
// If the parameter is not found, an empty string is returned.
func Param(ctx *fasthttp.RequestCtx, key string) string {
	if value := ctx.UserValue(key); value != nil {
		return value.(string)
	}
	return ""
}

// GetTree returns the routing tree for the specified HTTP method.
// This is primarily used for testing and debugging purposes.
func (r *Router) GetTree(method string) *Tree {
	return r.trees[method]
}

// GetValue looks up a handler for the given path and method in the routing tree.
// It returns the handler function, any extracted parameters, and a boolean indicating if
// a trailing slash redirect should occur.
// This is primarily used for testing and debugging purposes.
func (t *Tree) GetValue(path, method string) (fasthttp.RequestHandler, Parameters, bool) {
	return t.getValue(path, method)
}

// AddRoute adds a new route to the tree.
// This is primarily used for testing and debugging purposes.
func (t *Tree) AddRoute(path string, handler fasthttp.RequestHandler) {
	t.addRoute(path, handler)
}

// GetRoot returns the root node of the routing tree.
// This is primarily used for testing and debugging purposes.
func (t *Tree) GetRoot() *Node {
	return t.root
}

// PrintTree prints the tree structure to standard output.
// This is useful for debugging the router's internal structure.
func (t *Tree) PrintTree() {
	if t.root != nil {
		printNode(t.root, 0)
	}
}

func printNode(n *Node, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	fmt.Printf("%sNode: path='%s', type=%v, wildChild=%v, indices='%s'\n",
		indent, n.path, n.nType, n.wildChild, n.indices)

	if len(n.handlers) > 0 {
		fmt.Printf("%s  Handlers: %d\n", indent, len(n.handlers))
	}
	if len(n.paramNames) > 0 {
		fmt.Printf("%s  ParamNames: %v\n", indent, n.paramNames)
	}

	for _, child := range n.children {
		printNode(child, depth+1)
	}
}
