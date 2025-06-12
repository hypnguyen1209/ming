package ming

import (
	"testing"
	"github.com/valyala/fasthttp"
)

func TestMultipleCatchAllRoutes(t *testing.T) {
	router := New()
	var capturedPath string
	var handlerCounter int
	
	// Define multiple catch-all routes with different prefixes
	router.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		capturedPath = Param(ctx, "filepath")
		handlerCounter = 1
	})
	
	router.Get("/documents/{docpath:*}", func(ctx *fasthttp.RequestCtx) {
		capturedPath = Param(ctx, "docpath")
		handlerCounter = 2
	})
	
	router.Get("/api/v1/proxy/{url:*}", func(ctx *fasthttp.RequestCtx) {
		capturedPath = Param(ctx, "url")
		handlerCounter = 3
	})
	
	// Test different paths
	testCases := []struct {
		name           string
		path           string
		expectedPath   string
		expectedHandler int
		shouldMatch    bool
	}{
		{"files route - single file", "/files/document.txt", "document.txt", 1, true},
		{"files route - nested path", "/files/folder/subfolder/file.pdf", "folder/subfolder/file.pdf", 1, true},
		{"documents route - single file", "/documents/report.pdf", "report.pdf", 2, true},
		{"documents route - nested path", "/documents/2024/quarterly/q1.docx", "2024/quarterly/q1.docx", 2, true},
		{"api proxy route", "/api/v1/proxy/https:/example.com", "https:/example.com", 3, true}, // Note: URL has only one slash
		{"api proxy nested path", "/api/v1/proxy/service/endpoint", "service/endpoint", 3, true},
		{"non-existent route", "/images/logo.png", "", 0, false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capturedPath = ""
			handlerCounter = 0
			
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI(tc.path)
			ctx.Request.Header.SetMethod(fasthttp.MethodGet)
			
			router.Handler(ctx)
			
			t.Logf("Test %s: Request path %s, captured path: %q", tc.name, tc.path, capturedPath)
			
			if tc.shouldMatch {
				if handlerCounter == 0 {
					t.Errorf("Expected route to match for path %s, but no handler was called", tc.path)
				}
				
				if handlerCounter != tc.expectedHandler {
					t.Errorf("Expected handler %d to be called for path %s, but handler %d was called", 
					         tc.expectedHandler, tc.path, handlerCounter)
				}
				
				if capturedPath != tc.expectedPath {
					t.Errorf("Expected captured path %q, got %q", tc.expectedPath, capturedPath)
				}
			} else {
				if handlerCounter != 0 {
					t.Errorf("Expected no handler to be called for path %s, but handler %d was called", 
					         tc.path, handlerCounter)
				}
			}
		})
	}
}

func TestMultipleCatchAllRoutesWithSameBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for duplicate catch-all route with same prefix, but none occurred")
		}
	}()
	
	router := New()
	
	// Define first catch-all route
	router.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		// First handler
	})
	
	// Define second catch-all route with the same base
	// This should panic
	router.Get("/files/{otherpath:*}", func(ctx *fasthttp.RequestCtx) {
		// Second handler
	})
}

func TestMultipleCatchAllRoutesWithTrailingSlash(t *testing.T) {
	router := New()
	var capturedPath string
	var handlerCounter int
	
	// Define routes with trailing slash and catch-all 
	router.Get("/files/", func(ctx *fasthttp.RequestCtx) {
		handlerCounter = 1
	})
	
	router.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		capturedPath = Param(ctx, "filepath")
		handlerCounter = 2
	})
	
	router.Get("/documents/", func(ctx *fasthttp.RequestCtx) {
		handlerCounter = 3
	})
	
	router.Get("/documents/{docpath:*}", func(ctx *fasthttp.RequestCtx) {
		capturedPath = Param(ctx, "docpath")
		handlerCounter = 4
	})
	
	// Test cases
	testCases := []struct {
		name           string
		path           string
		expectedPath   string
		expectedHandler int
	}{
		{"files empty path", "/files/", "", 1}, // Direct match with /files/
		{"files with content", "/files/hello.txt", "hello.txt", 2},
		{"documents empty path", "/documents/", "", 3}, // Direct match with /documents/
		{"documents with content", "/documents/report.pdf", "report.pdf", 4},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			capturedPath = ""
			handlerCounter = 0
			
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.SetRequestURI(tc.path)
			ctx.Request.Header.SetMethod(fasthttp.MethodGet)
			
			router.Handler(ctx)
			
			if handlerCounter != tc.expectedHandler {
				t.Errorf("Expected handler %d to be called for path %s, but handler %d was called", 
				         tc.expectedHandler, tc.path, handlerCounter)
			}
			
			if capturedPath != tc.expectedPath {
				t.Errorf("Expected captured path %q, got %q", tc.expectedPath, capturedPath)
			}
		})
	}
}
