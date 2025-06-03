package main

import (
	"fmt"

	"github.com/hypnguyen1209/ming"
	"github.com/valyala/fasthttp"
)

// This example demonstrates Ming's priority-based routing and conflict resolution.
// When multiple routes could match a request path, Ming makes decisions based on:
// 1. Static routes have higher priority than parameter routes
// 2. The most specific route wins over more general routes
// 3. Earlier declared routes have priority over later ones (when other rules don't apply)

func main() {
	r := ming.New()

	// Setup basic logging middleware
	r.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
		fmt.Printf("Panic recovered: %v\n", p)
		ctx.SetStatusCode(500)
		fmt.Fprintf(ctx, "Internal server error")
	}
	
	// Setup a route matching helper
	setupRouteMatching(r)
	
	// Setup examples of route conflicts and their resolution
	setupConflictExamples(r)
	
	// Setup a visual route tree
	setupRouteTree(r)
	
	// Start the server
	fmt.Println("Ming Router - Priority-Based Routing and Conflict Resolution Example")
	fmt.Println("Server running at http://127.0.0.1:8000")
	fmt.Println("")
	fmt.Println("Test the following endpoints to see routing priority in action:")
	fmt.Println("1. http://127.0.0.1:8000/conflict/static")
	fmt.Println("2. http://127.0.0.1:8000/conflict/param")
	fmt.Println("3. http://127.0.0.1:8000/api/v1/users")
	fmt.Println("4. http://127.0.0.1:8000/api/v2/users")
	fmt.Println("5. http://127.0.0.1:8000/files/report.pdf")
	fmt.Println("6. http://127.0.0.1:8000/files/images/photo.jpg")
	fmt.Println("7. http://127.0.0.1:8000/route-tree (visual representation of routes)")
	fmt.Println("")
	
	r.Run(":8000")
}

// Setup route matching helper
func setupRouteMatching(r *ming.Router) {
	// Home page with basic information
	r.Get("/", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Ming Router - Priority and Conflict Resolution Example\n\n")
		fmt.Fprintf(ctx, "Visit /route-tree to see the routing tree\n")
		fmt.Fprintf(ctx, "Visit /conflict/static to see static vs parameter priority\n")
		fmt.Fprintf(ctx, "Visit /api/v1/users or /api/v2/users to see specific vs general routes\n")
	})
	
	// Route to check how a path would be matched
	r.Get("/match", func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.QueryArgs().Peek("path"))
		if path == "" {
			fmt.Fprintf(ctx, "Please provide a path to match using ?path=/your/path")
			return
		}
		
		fmt.Fprintf(ctx, "Path matching for: %s\n\n", path)
		fmt.Fprintf(ctx, "This would match to a route based on Ming's priority rules:\n")
		fmt.Fprintf(ctx, "1. Static routes have priority over parameter routes\n")
		fmt.Fprintf(ctx, "2. More specific routes have priority over general ones\n")
		fmt.Fprintf(ctx, "3. Earlier registered routes have priority when other rules don't apply\n")
	})
}

// Setup examples of route conflicts and their resolution
func setupConflictExamples(r *ming.Router) {
	// Example 1: Static vs Parameter routes
	// The static route has higher priority
	r.Get("/conflict/static", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "This is the /conflict/static route (STATIC ROUTE)\n\n")
		fmt.Fprintf(ctx, "This path matched the static route because static routes have\n")
		fmt.Fprintf(ctx, "higher priority than parameter routes in Ming's router.\n\n")
		fmt.Fprintf(ctx, "There's also a conflicting route with pattern /conflict/{param}\n")
		fmt.Fprintf(ctx, "but it will never be matched for this specific URL.")
	})
	
	// This route won't be matched for /conflict/static but will match other paths
	r.Get("/conflict/{param}", func(ctx *fasthttp.RequestCtx) {
		param := ming.Param(ctx, "param")
		fmt.Fprintf(ctx, "This is the /conflict/{param} route (PARAMETER ROUTE)\n\n")
		fmt.Fprintf(ctx, "Param value: %s\n\n", param)
		fmt.Fprintf(ctx, "This path matched the parameter route because there's no static route\n")
		fmt.Fprintf(ctx, "that directly matches this path. Static routes always take\n")
		fmt.Fprintf(ctx, "precedence over parameter routes for the same path segment.")
	})
	
	// Example 2: Specific vs General API version routes
	// The specific route has higher priority
	r.Get("/api/v1/users", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "API v1 Users Endpoint (SPECIFIC ROUTE)\n\n")
		fmt.Fprintf(ctx, "This matched the specific /api/v1/users route instead of\n")
		fmt.Fprintf(ctx, "the more general /api/{version}/users route because Ming's\n")
		fmt.Fprintf(ctx, "router gives priority to more specific routes.")
	})
	
	// This route won't be matched for /api/v1/users but will match other API versions
	r.Get("/api/{version}/users", func(ctx *fasthttp.RequestCtx) {
		version := ming.Param(ctx, "version")
		fmt.Fprintf(ctx, "API %s Users Endpoint (GENERAL ROUTE)\n\n", version)
		fmt.Fprintf(ctx, "This matched the general /api/{version}/users route because\n")
		fmt.Fprintf(ctx, "there's no specific route defined for this API version.\n")
		fmt.Fprintf(ctx, "Note that requests to /api/v1/users would match the specific route instead.")
	})
	
	// Example 3: Multiple route parameters with different levels of specificity
	// Direct file route - highest priority for direct matches
	r.Get("/files/report.pdf", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Direct file route: /files/report.pdf (STATIC ROUTE)\n\n")
		fmt.Fprintf(ctx, "This matched the static route specifically defined for this file.\n")
		fmt.Fprintf(ctx, "This has higher priority than any parameter routes.")
	})
	
	// Files in a specific folder - higher priority than general catchall
	r.Get("/files/images/{filename}", func(ctx *fasthttp.RequestCtx) {
		filename := ming.Param(ctx, "filename")
		fmt.Fprintf(ctx, "Image file route: /files/images/%s (PARAMETER ROUTE)\n\n", filename)
		fmt.Fprintf(ctx, "This matched the specific folder parameter route.\n")
		fmt.Fprintf(ctx, "This has higher priority than the general catchall route.")
	})
	
	// Catchall files route - lowest priority
	r.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		filepath := ming.Param(ctx, "filepath")
		fmt.Fprintf(ctx, "Catchall file route: /files/%s (CATCHALL ROUTE)\n\n", filepath)
		fmt.Fprintf(ctx, "This matched the general catchall route because no more\n")
		fmt.Fprintf(ctx, "specific route was found for this path.")
	})
}

// Setup a visual representation of the route tree
func setupRouteTree(r *ming.Router) {
	r.Get("/route-tree", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Ming Router - Route Priority Tree\n\n")
		fmt.Fprintf(ctx, "The Ming router uses a sophisticated radix tree to store and match routes.\n")
		fmt.Fprintf(ctx, "Routes are stored in a tree structure with the following priority rules:\n\n")
		fmt.Fprintf(ctx, "1. Static routes have higher priority than parameter routes\n")
		fmt.Fprintf(ctx, "2. More specific routes have higher priority than general routes\n")
		fmt.Fprintf(ctx, "3. Earlier registered routes have priority when other rules don't apply\n\n")
		fmt.Fprintf(ctx, "Visual representation of our current route tree:\n\n")
		
		// ASCII representation of route tree
		fmt.Fprintf(ctx, "/ (ROOT)\n")
		fmt.Fprintf(ctx, "├── / (GET)\n")
		fmt.Fprintf(ctx, "├── match (GET)\n")
		fmt.Fprintf(ctx, "├── conflict/\n")
		fmt.Fprintf(ctx, "│   ├── static (GET) [STATIC - HIGHER PRIORITY]\n")
		fmt.Fprintf(ctx, "│   └── {param} (GET) [PARAMETER - LOWER PRIORITY]\n")
		fmt.Fprintf(ctx, "├── api/\n")
		fmt.Fprintf(ctx, "│   ├── v1/\n")
		fmt.Fprintf(ctx, "│   │   └── users (GET) [SPECIFIC - HIGHER PRIORITY]\n")
		fmt.Fprintf(ctx, "│   └── {version}/\n")
		fmt.Fprintf(ctx, "│       └── users (GET) [GENERAL - LOWER PRIORITY]\n")
		fmt.Fprintf(ctx, "├── files/\n")
		fmt.Fprintf(ctx, "│   ├── report.pdf (GET) [STATIC - HIGHEST PRIORITY]\n")
		fmt.Fprintf(ctx, "│   ├── images/\n")
		fmt.Fprintf(ctx, "│   │   └── {filename} (GET) [PARAMETER - MEDIUM PRIORITY]\n")
		fmt.Fprintf(ctx, "│   └── {filepath:*} (GET) [CATCHALL - LOWEST PRIORITY]\n")
		fmt.Fprintf(ctx, "└── route-tree (GET)\n")
	})
}
