/*
Package ming is a high-performance HTTP router based on fasthttp.

Ming provides a priority-based routing system with a Radix Tree for efficient request routing.
It is designed to be fast, flexible, and easy to use for building web applications and APIs.

Key Features:

  - High Performance: Built on fasthttp for maximum speed
  - Named Parameters: `/user/{name}` with parameter extraction
  - Optional Parameters: `/api/users/{id?}` for flexible routing
  - Regex Validation: `/product/{id:[0-9]+}` for parameter validation
  - Catch-All Routes: `/files/{path:*}` for wildcard matching
  - Priority-Based Routing: Radix tree with optimized lookups for O(k) matching complexity
  - Route Conflict Resolution: Static routes take precedence over parameter routes
  - Method Support: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT
  - Static File Serving: Built-in static file handler
  - Middleware Support: Custom panic handlers and error handling
  - Thread-Safe: Concurrent request handling
  - Zero Memory Allocations: Built on fasthttp's zero-allocation philosophy

Quick Start:

	package main
	
	import (
		"fmt"
		"github.com/hypnguyen1209/ming/v2"
		"github.com/valyala/fasthttp"
	)
	
	func main() {
		r := ming.New()
		
		// Basic route
		r.Get("/", func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("Hello World!")
		})
		
		// Named parameters
		r.Get("/user/{name}", func(ctx *fasthttp.RequestCtx) {
			name := ming.Param(ctx, "name")
			fmt.Fprintf(ctx, "Hello %s!", name)
		})
		
		// Run the server
		r.Run(":8080")
	}

For more details and examples, see https://github.com/hypnguyen1209/ming/v2
*/
package ming
