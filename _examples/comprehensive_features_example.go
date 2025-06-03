package main

import (
	"fmt"

	"github.com/hypnguyen1209/ming"
	"github.com/valyala/fasthttp"
)

// This example demonstrates all the major Ming router features:
// - High Performance (fasthttp-based)
// - Named Parameters
// - Optional Parameters
// - Regex Validation
// - Catch-All Routes
// - Priority-Based Routing
// - Route Conflict Resolution

// Home handler - basic route
func homeHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Welcome to the Ming Router Features Demo!")
}

// Named parameter handler
func userHandler(ctx *fasthttp.RequestCtx) {
	// Extract the 'name' parameter from the URL
	name := ming.Param(ctx, "name")
	fmt.Fprintf(ctx, "User profile: %s", name)
}

// Multiple parameters handler
func userPostHandler(ctx *fasthttp.RequestCtx) {
	// Extract multiple parameters
	userId := ming.Param(ctx, "userId")
	postId := ming.Param(ctx, "postId")
	fmt.Fprintf(ctx, "User %s, Post %s", userId, postId)
}

// Optional parameter handler
func usersHandler(ctx *fasthttp.RequestCtx) {
	// Check if optional 'id' parameter exists
	id := ming.Param(ctx, "id")
	if id == "" {
		// No ID provided - list all users
		fmt.Fprintf(ctx, "Listing all users")
	} else {
		// ID provided - show specific user
		fmt.Fprintf(ctx, "Details for user: %s", id)
	}
}

// Regex validation handler - only accepts numeric IDs
func productHandler(ctx *fasthttp.RequestCtx) {
	// The router validates that 'id' contains only digits before reaching this handler
	id := ming.Param(ctx, "id")
	fmt.Fprintf(ctx, "Product ID: %s (guaranteed to be numeric)", id)
}

// Category handler with alphabetic-only validation
func categoryHandler(ctx *fasthttp.RequestCtx) {
	// The router validates that category contains only letters
	category := ming.Param(ctx, "category")
	fmt.Fprintf(ctx, "Category: %s (guaranteed to be alphabetic)", category)
}

// Catch-all handler - works with wildcard paths
func fileHandler(ctx *fasthttp.RequestCtx) {
	// Get the entire path after /files/
	path := ming.Param(ctx, "filepath")
	fmt.Fprintf(ctx, "Requested file path: %s", path)
}

// Static route handler - demonstrates route priority
func staticProfileHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "This is the static profile page (high priority)")
}

// Specific API version handler
func apiV1Handler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "API v1 endpoint")
}

// Generic API version handler
func apiVersionHandler(ctx *fasthttp.RequestCtx) {
	version := ming.Param(ctx, "version")
	fmt.Fprintf(ctx, "API %s endpoint (fallback for non-v1)", version)
}

// Custom 404 handler
func notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	fmt.Fprintf(ctx, "Custom 404: Page not found - %s", ctx.Path())
}

// Custom method not allowed handler
func methodNotAllowedHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	fmt.Fprintf(ctx, "Method %s not allowed for path %s", ctx.Method(), ctx.Path())
}

// Custom panic handler
func panicHandler(ctx *fasthttp.RequestCtx, err interface{}) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	fmt.Fprintf(ctx, "Panic recovered: %v", err)
}

// Handler that will cause a panic
func buggyHandler(ctx *fasthttp.RequestCtx) {
	// Intentional panic to demonstrate the panic handler
	var ptr *int
	*ptr = 42 // This will cause a panic (nil pointer dereference)
}

func main() {
	r := ming.New()
	
	// Set up custom error handlers
	r.NotFound = notFoundHandler
	r.MethodNotAllowed = methodNotAllowedHandler
	r.PanicHandler = panicHandler

	// Basic route
	r.Get("/", homeHandler)
	
	// === NAMED PARAMETERS ===
	// Basic named parameter
	r.Get("/user/{name}", userHandler)
	
	// Multiple parameters
	r.Get("/user/{userId}/post/{postId}", userPostHandler)
	
	// === OPTIONAL PARAMETERS ===
	// Optional parameter - matches both with and without the ID
	r.Get("/api/users/{id?}", usersHandler)
	
	// === REGEX VALIDATION ===
	// Numeric validation - only matches if ID consists of digits
	r.Get("/product/{id:[0-9]+}", productHandler)
	
	// Alphabetic validation - only matches if category consists of letters
	r.Get("/category/{category:[a-zA-Z]+}", categoryHandler)
	
	// === CATCH-ALL ROUTES ===
	// Wildcard matching - captures everything after /files/
	r.Get("/files/{filepath:*}", fileHandler)
	
	// === ROUTE PRIORITY & CONFLICT RESOLUTION ===
	// Static route (higher priority)
	r.Get("/user/profile", staticProfileHandler)
	// This route won't match /user/profile because the static route above has higher priority
	// But it will match any other /user/X path
	// r.Get("/user/{name}", userHandler) - already defined above
	
	// Specific API version (higher priority)
	r.Get("/api/v1/status", apiV1Handler)
	// Generic API version (lower priority - will match any version except v1)
	r.Get("/api/{version}/status", apiVersionHandler)
	
	// Route that will cause a panic (to demonstrate panic recovery)
	r.Get("/panic", buggyHandler)
	
	// Start the server
	fmt.Println("Server running on http://127.0.0.1:8080")
	fmt.Println("\nTest the following features:")
	fmt.Println("1. Named Parameters: http://127.0.0.1:8080/user/john")
	fmt.Println("2. Multiple Parameters: http://127.0.0.1:8080/user/42/post/123")
	fmt.Println("3. Optional Parameters: http://127.0.0.1:8080/api/users (without ID)")
	fmt.Println("                        http://127.0.0.1:8080/api/users/42 (with ID)")
	fmt.Println("4. Regex Validation: http://127.0.0.1:8080/product/123 (valid)")
	fmt.Println("                     http://127.0.0.1:8080/product/abc (invalid - will 404)")
	fmt.Println("5. Alphabetic Validation: http://127.0.0.1:8080/category/electronics (valid)")
	fmt.Println("                          http://127.0.0.1:8080/category/123 (invalid - will 404)")
	fmt.Println("6. Catch-All Routes: http://127.0.0.1:8080/files/docs/report.pdf")
	fmt.Println("7. Route Priority: http://127.0.0.1:8080/user/profile (static route)")
	fmt.Println("                   http://127.0.0.1:8080/user/john (parameter route)")
	fmt.Println("8. API Version Priority: http://127.0.0.1:8080/api/v1/status (specific)")
	fmt.Println("                         http://127.0.0.1:8080/api/v2/status (generic)")
	fmt.Println("9. Error Handlers: http://127.0.0.1:8080/unknown (custom 404)")
	fmt.Println("                    POST http://127.0.0.1:8080/ (method not allowed)")
	fmt.Println("10. Panic Recovery: http://127.0.0.1:8080/panic")
	
	r.Run(":8080")
}
