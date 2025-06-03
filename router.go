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

func New() *Router {
	return &Router{
		trees: make(map[string]*Tree),
	}
}

type HostSwitch map[string]fasthttp.RequestHandler

func (hs HostSwitch) CheckHost(ctx *fasthttp.RequestCtx) {
	if handler := hs[string(ctx.Host())]; handler != nil {
		handler(ctx)
	} else {
		ctx.Error("Forbidden", fasthttp.StatusForbidden)
	}
}

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

func Query(ctx *fasthttp.RequestCtx, str string) []byte {
	return ctx.QueryArgs().Peek(str)
}

func SetHeader(ctx *fasthttp.RequestCtx, key string, value string) {
	ctx.Response.Header.Set(key, value)
}

func Body(ctx *fasthttp.RequestCtx) []byte {
	return ctx.Request.Body()
}

func (r *Router) recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

// UserValue gets a user value from the context for parameter access
func UserValue(ctx *fasthttp.RequestCtx, key string) interface{} {
	return ctx.UserValue(key)
}

// Param is a convenience function to get a parameter value as string
func Param(ctx *fasthttp.RequestCtx, key string) string {
	if value := ctx.UserValue(key); value != nil {
		return value.(string)
	}
	return ""
}

// Debug methods for testing
func (r *Router) GetTree(method string) *Tree {
	return r.trees[method]
}

func (t *Tree) GetValue(path, method string) (fasthttp.RequestHandler, Parameters, bool) {
	return t.getValue(path, method)
}

// Debug methods for testing
func (t *Tree) AddRoute(path string, handler fasthttp.RequestHandler) {
	t.addRoute(path, handler)
}

func (t *Tree) GetRoot() *Node {
	return t.root
}

// PrintTree prints the tree structure for debugging
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
