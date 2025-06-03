package ming

import (
	"testing"

	"github.com/valyala/fasthttp"
)

// TestRouteOverwrite verifies that when a route with the same path is added,
// it overwrites the previous handler.
func TestRouteOverwrite(t *testing.T) {
	router := New()

	// Define our test variables
	path := "/test"
	var handlerCalled string

	// Define the first handler
	firstHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "first"
	}

	// Define the second handler
	secondHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "second"
	}

	// Add the first route
	router.Get(path, firstHandler)

	// Add the second route with the same path (should overwrite)
	router.Get(path, secondHandler)

	// Create a test request context
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI(path)

	// Process the request
	router.Handler(ctx)

	// Check that the second handler was called, not the first
	if handlerCalled != "second" {
		t.Errorf("Expected second handler to be called, but got %s", handlerCalled)
	}
}

// TestRouteOverwriteDifferentMethods verifies that routes with the same path
// but different methods don't overwrite each other.
func TestRouteOverwriteDifferentMethods(t *testing.T) {
	router := New()

	// Define our test variables
	path := "/api/users"
	var getHandlerCalled bool
	var postHandlerCalled bool

	// Define handlers for different methods
	getHandler := func(ctx *fasthttp.RequestCtx) {
		getHandlerCalled = true
	}

	postHandler := func(ctx *fasthttp.RequestCtx) {
		postHandlerCalled = true
	}

	// Add routes for different methods
	router.Get(path, getHandler)
	router.Post(path, postHandler)

	// Create a test GET request context
	getCtx := &fasthttp.RequestCtx{}
	getCtx.Request.Header.SetMethod("GET")
	getCtx.Request.SetRequestURI(path)

	// Process the GET request
	router.Handler(getCtx)

	// Check that only GET handler was called
	if !getHandlerCalled {
		t.Error("GET handler was not called")
	}
	if postHandlerCalled {
		t.Error("POST handler was unexpectedly called during GET request")
	}

	// Reset flags
	getHandlerCalled = false
	postHandlerCalled = false

	// Create a test POST request context
	postCtx := &fasthttp.RequestCtx{}
	postCtx.Request.Header.SetMethod("POST")
	postCtx.Request.SetRequestURI(path)

	// Process the POST request
	router.Handler(postCtx)

	// Check that only POST handler was called
	if !postHandlerCalled {
		t.Error("POST handler was not called")
	}
	if getHandlerCalled {
		t.Error("GET handler was unexpectedly called during POST request")
	}
}

// TestParameterRouteOverwrite tests that parameter routes can also be overwritten
func TestParameterRouteOverwrite(t *testing.T) {
	router := New()

	// Define our test variables
	paramPath := "/users/{id}"
	var handlerCalled string
	var paramValue string

	// Define the first handler
	firstHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "first"
		paramValue = Param(ctx, "id")
	}

	// Define the second handler
	secondHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "second"
		paramValue = Param(ctx, "id")
	}

	// Add the first route with parameter
	router.Get(paramPath, firstHandler)

	// Add the second route with the same path (should overwrite)
	router.Get(paramPath, secondHandler)

	// Create a test request context
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/users/123")

	// Process the request
	router.Handler(ctx)

	// Check that the second handler was called, not the first
	if handlerCalled != "second" {
		t.Errorf("Expected second handler to be called, but got %s", handlerCalled)
	}

	// Verify parameter was still properly captured
	if paramValue != "123" {
		t.Errorf("Expected parameter 'id' to be '123', but got '%s'", paramValue)
	}
}

// TestAllMethodOverwrite tests that the ALL method can be overwritten
func TestAllMethodOverwrite(t *testing.T) {
	router := New()

	// Define our test variables
	path := "/wildcard"
	var handlerCalled string

	// Define the first ALL handler
	firstHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "first"
	}

	// Define the second ALL handler
	secondHandler := func(ctx *fasthttp.RequestCtx) {
		handlerCalled = "second"
	}

	// Add the first ALL route
	router.All(path, firstHandler)

	// Add the second ALL route with the same path (should overwrite)
	router.All(path, secondHandler)

	// Create a test request context
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI(path)

	// Process the request
	router.Handler(ctx)

	// Check that the second handler was called, not the first
	if handlerCalled != "second" {
		t.Errorf("Expected second handler to be called, but got %s", handlerCalled)
	}

	// Reset and try with POST
	handlerCalled = ""
	ctx = &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI(path)

	// Process the POST request
	router.Handler(ctx)

	// Check that the second handler was still called for POST
	if handlerCalled != "second" {
		t.Errorf("Expected second handler to be called for POST, but got %s", handlerCalled)
	}
}