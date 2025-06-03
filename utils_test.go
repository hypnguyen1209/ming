package ming

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestGetMethod(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected string
	}{
		{"GET", fasthttp.MethodGet, fasthttp.MethodGet},
		{"POST", fasthttp.MethodPost, fasthttp.MethodPost},
		{"PUT", fasthttp.MethodPut, fasthttp.MethodPut},
		{"DELETE", fasthttp.MethodDelete, fasthttp.MethodDelete},
		{"PATCH", fasthttp.MethodPatch, fasthttp.MethodPatch},
		{"HEAD", fasthttp.MethodHead, fasthttp.MethodHead},
		{"OPTIONS", fasthttp.MethodOptions, fasthttp.MethodOptions},
		{"CONNECT", fasthttp.MethodConnect, fasthttp.MethodConnect},
		{"TRACE", fasthttp.MethodTrace, fasthttp.MethodTrace},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.SetMethod(tc.method)

			result := GetMethod(ctx)
			if result != tc.expected {
				t.Errorf("GetMethod() for %s = %s, want %s", tc.method, result, tc.expected)
			}
		})
	}
}

func TestGetMethodUnknown(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	// Set a custom method that's not in our switch statement
	ctx.Request.Header.SetMethod("CUSTOM")
	result := GetMethod(ctx)
	if result != "" {
		t.Errorf("GetMethod() for unknown method = %s, want empty string", result)
	}
}

func TestGetMethodWithCustomMethod(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("CUSTOM")
	
	result := GetMethod(ctx)
	if result != "" {
		t.Errorf("GetMethod() for custom method = %s, want empty string", result)
	}
}

func TestQuery(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("http://example.com/path?name=john&age=30&city=")

	tests := []struct {
		key      string
		expected string
	}{
		{"name", "john"},
		{"age", "30"},
		{"city", ""},
		{"nonexistent", ""},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			result := string(Query(ctx, tc.key))
			if result != tc.expected {
				t.Errorf("Query(%s) = %s, want %s", tc.key, result, tc.expected)
			}
		})
	}
}

func TestQueryWithMultipleValues(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("http://example.com/path?tags=go&tags=web&tags=router")

	// Query should return the first value when multiple values exist
	result := string(Query(ctx, "tags"))
	if result != "go" {
		t.Errorf("Query(tags) with multiple values = %s, want go", result)
	}
}

func TestSetHeader(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	tests := []struct {
		key   string
		value string
	}{
		{"Content-Type", "application/json"},
		{"X-Custom-Header", "custom-value"},
		{"Authorization", "Bearer token123"},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			SetHeader(ctx, tc.key, tc.value)
			result := string(ctx.Response.Header.Peek(tc.key))
			if result != tc.value {
				t.Errorf("SetHeader(%s, %s) result = %s, want %s", tc.key, tc.value, result, tc.value)
			}
		})
	}
}

func TestSetHeaderOverwrite(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	// Set initial value
	SetHeader(ctx, "Content-Type", "text/plain")
	result1 := string(ctx.Response.Header.Peek("Content-Type"))
	if result1 != "text/plain" {
		t.Errorf("First SetHeader(Content-Type) = %s, want text/plain", result1)
	}

	// Overwrite with new value
	SetHeader(ctx, "Content-Type", "application/json")
	result2 := string(ctx.Response.Header.Peek("Content-Type"))
	if result2 != "application/json" {
		t.Errorf("Second SetHeader(Content-Type) = %s, want application/json", result2)
	}
}

func TestBody(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	tests := []struct {
		name string
		body string
	}{
		{"empty body", ""},
		{"json body", `{"name":"john","age":30}`},
		{"plain text", "Hello, World!"},
		{"binary data", "\x00\x01\x02\x03"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx.Request.SetBody([]byte(tc.body))
			
			result := Body(ctx)
			if string(result) != tc.body {
				t.Errorf("Body() = %s, want %s", string(result), tc.body)
			}
		})
	}
}

func TestBodyWithLargePayload(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	// Create a large body (1MB)
	largeBody := make([]byte, 1024*1024)
	for i := range largeBody {
		largeBody[i] = byte(i % 256)
	}
	
	ctx.Request.SetBody(largeBody)
	
	result := Body(ctx)
	if len(result) != len(largeBody) {
		t.Errorf("Body() length = %d, want %d", len(result), len(largeBody))
	}
	
	// Verify content
	for i := 0; i < len(result); i++ {
		if result[i] != largeBody[i] {
			t.Errorf("Body() content mismatch at index %d: got %d, want %d", i, result[i], largeBody[i])
			break
		}
	}
}

func TestUserValue(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	tests := []struct {
		key   string
		value interface{}
	}{
		{"string_param", "test_value"},
		{"int_param", 42},
		{"bool_param", true},
		{"nil_param", nil},
		{"struct_param", struct{ Name string }{"test"}},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			// Set the user value
			ctx.SetUserValue(tc.key, tc.value)
			
			// Get it back using UserValue
			result := UserValue(ctx, tc.key)
			if result != tc.value {
				t.Errorf("UserValue(%s) = %v, want %v", tc.key, result, tc.value)
			}
		})
	}
}

func TestUserValueNonexistent(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	result := UserValue(ctx, "nonexistent")
	if result != nil {
		t.Errorf("UserValue(nonexistent) = %v, want nil", result)
	}
}

func TestParam(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	tests := []struct {
		key   string
		value string
	}{
		{"id", "123"},
		{"name", "john"},
		{"category", "electronics"},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			// Set the parameter as user value
			ctx.SetUserValue(tc.key, tc.value)
			
			// Get it back using Param
			result := Param(ctx, tc.key)
			if result != tc.value {
				t.Errorf("Param(%s) = %s, want %s", tc.key, result, tc.value)
			}
		})
	}
}

func TestParamNonexistent(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	result := Param(ctx, "nonexistent")
	if result != "" {
		t.Errorf("Param(nonexistent) = %s, want empty string", result)
	}
}

func TestParamWithNonStringValue(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	// Set a non-string value
	ctx.SetUserValue("number", 42)
	
	// This should cause a panic due to type assertion
	defer func() {
		if r := recover(); r == nil {
			t.Error("Param() should panic when value is not a string")
		}
	}()
	
	Param(ctx, "number")
}

func TestUtilityFunctionsIntegration(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	
	// Setup request
	ctx.Request.SetRequestURI("http://example.com/users/123?format=json&detailed=true")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetBody([]byte(`{"name":"john","email":"john@example.com"}`))
	
	// Set some parameters
	ctx.SetUserValue("id", "123")
	ctx.SetUserValue("action", "update")
	
	// Test all utility functions together
	method := GetMethod(ctx)
	if method != fasthttp.MethodPost {
		t.Errorf("GetMethod() = %s, want %s", method, fasthttp.MethodPost)
	}
	
	format := string(Query(ctx, "format"))
	if format != "json" {
		t.Errorf("Query(format) = %s, want json", format)
	}
	
	detailed := string(Query(ctx, "detailed"))
	if detailed != "true" {
		t.Errorf("Query(detailed) = %s, want true", detailed)
	}
	
	id := Param(ctx, "id")
	if id != "123" {
		t.Errorf("Param(id) = %s, want 123", id)
	}
	
	action := UserValue(ctx, "action")
	if action != "update" {
		t.Errorf("UserValue(action) = %v, want update", action)
	}
	
	body := Body(ctx)
	expectedBody := `{"name":"john","email":"john@example.com"}`
	if string(body) != expectedBody {
		t.Errorf("Body() = %s, want %s", string(body), expectedBody)
	}
	
	// Test setting headers
	SetHeader(ctx, "Content-Type", "application/json")
	SetHeader(ctx, "X-Request-ID", "req-123")
	
	contentType := string(ctx.Response.Header.Peek("Content-Type"))
	if contentType != "application/json" {
		t.Errorf("Response Content-Type = %s, want application/json", contentType)
	}
	
	requestId := string(ctx.Response.Header.Peek("X-Request-ID"))
	if requestId != "req-123" {
		t.Errorf("Response X-Request-ID = %s, want req-123", requestId)
	}
}
