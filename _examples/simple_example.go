package main

import (
	"fmt"
	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	r := ming.New()

	// Basic named parameter test
	r.Get("/hello/{name}", func(ctx *fasthttp.RequestCtx) {
		name := ming.Param(ctx, "name")
		fmt.Fprintf(ctx, "Hello %s!", name)
	})

	// Test with multiple parameters
	r.Get("/user/{id}/post/{postId}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		postId := ming.Param(ctx, "postId")
		fmt.Fprintf(ctx, "User: %s, Post: %s", id, postId)
	})

	// Regex validation test
	r.Get("/product/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		fmt.Fprintf(ctx, "Product ID: %s", id)
	})

	// Catch-all test
	r.Get("/files/{path:*}", func(ctx *fasthttp.RequestCtx) {
		path := ming.Param(ctx, "path")
		fmt.Fprintf(ctx, "File path: %s", path)
	})

	fmt.Println("Server starting on :8080")
	fmt.Println("Test URLs:")
	fmt.Println("  http://localhost:8080/hello/world")
	fmt.Println("  http://localhost:8080/user/123/post/456")
	fmt.Println("  http://localhost:8080/product/789")
	fmt.Println("  http://localhost:8080/files/docs/readme.txt")

	r.Run(":8080")
}
