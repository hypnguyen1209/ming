package ming

import (
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

var (
	DefaultContentType = []byte("text/plain; charset=utf-8")
)

type Router struct {
	trees *Tree
}

func New() *Router {
	tree := new(Tree)
	return &Router{
		trees: tree,
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
