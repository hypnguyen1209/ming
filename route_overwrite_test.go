package ming

import (
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
)

// Helper function to get value from a handler for testing
func getHandler(handler fasthttp.RequestHandler) string {
	if handler == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%p", handler) // Print handler pointer address
}

// TestSimpleRouteOverwrite verifies that a new route with the same path overwrites the previous one
func TestSimpleRouteOverwrite(t *testing.T) {
	tree := NewTree("GET")
	
	// First handler
	handler1 := func(ctx *fasthttp.RequestCtx) {}
	tree.addRoute("/test", handler1)
	
	// Get the handler directly from the tree
	h1, _, _ := tree.getValue("/test", "GET")
	fmt.Printf("Handler 1 address: %s\n", getHandler(h1))
	
	// Second handler for the same path
	handler2 := func(ctx *fasthttp.RequestCtx) {}
	tree.addRoute("/test", handler2)
	
	// Get the updated handler from the tree
	h2, _, _ := tree.getValue("/test", "GET")
	fmt.Printf("Handler 2 address: %s\n", getHandler(h2))
	
	// Verify handlers are different (we've overwritten the handler)
	if fmt.Sprintf("%p", h1) == fmt.Sprintf("%p", h2) {
		t.Error("Route was not overwritten")
	}
	
	// Verify the second handler is returned
	if fmt.Sprintf("%p", h2) != fmt.Sprintf("%p", handler2) {
		t.Error("The wrong handler is being returned")
	}
}

// TestParamRouteOverwrite tests overwriting a route with parameters
func TestParamRouteOverwrite(t *testing.T) {
	tree := NewTree("GET")
	
	// First handler with parameter
	handler1 := func(ctx *fasthttp.RequestCtx) {}
	tree.addRoute("/users/{id}", handler1)
	
	// Get the handler directly from the tree
	h1, params1, _ := tree.getValue("/users/123", "GET")
	fmt.Printf("Handler 1 address: %s\n", getHandler(h1))
	fmt.Printf("Parameters 1: %v\n", params1)
	
	// Second handler for the same path
	handler2 := func(ctx *fasthttp.RequestCtx) {}
	tree.addRoute("/users/{id}", handler2)
	
	// Get the updated handler from the tree
	h2, params2, _ := tree.getValue("/users/123", "GET")
	fmt.Printf("Handler 2 address: %s\n", getHandler(h2))
	fmt.Printf("Parameters 2: %v\n", params2)
	
	// Verify handlers are different (we've overwritten the handler)
	if fmt.Sprintf("%p", h1) == fmt.Sprintf("%p", h2) {
		t.Error("Parameter route was not overwritten")
	}
	
	// Verify parameters are still being captured
	if len(params2) == 0 || params2[0].Key != "id" || params2[0].Value != "123" {
		t.Errorf("Parameters not captured correctly: %v", params2)
	}
}

// TestPrintTreeAfterOverwrite prints the tree before and after overwriting a route
func TestPrintTreeAfterOverwrite(t *testing.T) {
	tree := NewTree("GET")
	
	// Add first handler
	handler1 := func(ctx *fasthttp.RequestCtx) {}
	fmt.Println("=== ADDING FIRST HANDLER ===")
	tree.addRoute("/api/users/{id}", handler1)
	fmt.Println("=== TREE AFTER FIRST HANDLER ===")
	tree.PrintTree()
	
	// Add second handler (overwrite)
	handler2 := func(ctx *fasthttp.RequestCtx) {}
	fmt.Println("\n=== ADDING SECOND HANDLER ===")
	tree.addRoute("/api/users/{id}", handler2)
	fmt.Println("=== TREE AFTER SECOND HANDLER ===")
	tree.PrintTree()
	
	// Verify the correct handler is returned
	h, params, _ := tree.getValue("/api/users/123", "GET")
	
	if fmt.Sprintf("%p", h) != fmt.Sprintf("%p", handler2) {
		t.Error("The wrong handler is being returned after overwrite")
	}
	
	fmt.Printf("\nFinal handler address: %s\n", getHandler(h))
	fmt.Printf("Final parameters: %v\n", params)
}