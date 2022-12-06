package main

import (
	"fmt"

	"github.com/hypnguyen1209/ming"
	"github.com/valyala/fasthttp"
)

func Home(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Home")
}

func AllHandler(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("123")
}

func SearchHandler(ctx *fasthttp.RequestCtx) {
	q := string(ming.Query(ctx, "name"))
	fmt.Fprintf(ctx, "Hello %s", q)
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	ctx.Write(ming.Body(ctx))
}

func main() {
	r := ming.New()
	r.Get("/", Home)
	r.Post("/add", PostHandler)
	r.All("/all", AllHandler)
	r.Get("/search", SearchHandler)
	r.Run("127.0.0.1:8000")
}
