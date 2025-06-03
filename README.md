# ming

Custom HTTP Mux lightweight and high performance 🥗

![](https://i.imgur.com/yCMS1yq.png)

## Features

✅ **High Performance** - Built on fasthttp for maximum speed  
✅ **Named Parameters** - `/user/{name}` with parameter extraction  
✅ **Optional Parameters** - `/api/users/{id?}` for flexible routing  
✅ **Regex Validation** - `/product/{id:[0-9]+}` for parameter validation  
✅ **Catch-All Routes** - `/files/{path:*}` for wildcard matching  
✅ **Priority-Based Routing** - Radix tree with optimized lookups  
✅ **Route Conflict Resolution** - Static routes take precedence over parameter routes  
✅ **Method Support** - GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, TRACE, CONNECT  
✅ **Static File Serving** - Built-in static file handler  
✅ **Middleware Support** - Custom panic handlers and error handling  
✅ **Thread-Safe** - Concurrent request handling  
✅ **Zero Memory Allocations** - Built on fasthttp's zero-allocation philosophy  

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/hypnguyen1209/ming"
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
	
	// Multiple parameters
	r.Get("/user/{id}/post/{postId}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		postId := ming.Param(ctx, "postId")
		fmt.Fprintf(ctx, "User: %s, Post: %s", id, postId)
	})
	
	// Regex validation
	r.Get("/product/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		fmt.Fprintf(ctx, "Product: %s", id)
	})
	
	// Catch-all routes
	r.Get("/files/{path:*}", func(ctx *fasthttp.RequestCtx) {
		path := ming.Param(ctx, "path")
		fmt.Fprintf(ctx, "File: %s", path)
	})
	
	r.Run(":8080")
}
```

## Parameter Types

### Named Parameters
```go
r.Get("/user/{name}", handler)     // Matches: /user/john
```

### Optional Parameters  
```go
r.Get("/api/users/{id?}", handler)  // Matches: /api/users/ and /api/users/123
```

### Regex Validation
```go
r.Get("/product/{id:[0-9]+}", handler)      // Only numeric IDs
r.Get("/category/{name:[a-zA-Z]+}", handler) // Only alphabetic names
```

### Catch-All Parameters
```go
r.Get("/files/{filepath:*}", handler)  // Matches: /files/docs/readme.txt
```

### Complex Routes
```go
r.Get("/api/{version:[v][0-9]+}/users/{userId:[0-9]+}/files/{filepath:*}", handler)
// Matches: /api/v1/users/123/files/documents/report.pdf
```

## API Reference

### Router Methods
```go
r.Get(path, handler)     // GET requests
r.Post(path, handler)    // POST requests  
r.Put(path, handler)     // PUT requests
r.Delete(path, handler)  // DELETE requests
r.Patch(path, handler)   // PATCH requests
r.Head(path, handler)    // HEAD requests
r.Options(path, handler) // OPTIONS requests
r.Trace(path, handler)   // TRACE requests
r.Connect(path, handler) // CONNECT requests
r.All(path, handler)     // All HTTP methods
```

### Parameter Access
```go
// Get parameter value as string
name := ming.Param(ctx, "name")

// Get query parameter
q := string(ming.Query(ctx, "search"))

// Get request body
body := ming.Body(ctx)
```

### Request Data Handling

#### GET Request Query Parameters

Ming provides easy access to URL query parameters through the `ming.Query()` function:

```go
// For URL: /search?term=ming&page=2
func SearchHandler(ctx *fasthttp.RequestCtx) {
    // Get single query parameters
    term := string(ming.Query(ctx, "term"))     // Returns "ming"
    page := string(ming.Query(ctx, "page"))     // Returns "2"
    
    // Check if parameter exists
    if len(ming.Query(ctx, "term")) > 0 {
        // Parameter exists
    }
    
    // Get all query parameters
    queryArgs := ctx.QueryArgs()
    
    // Iterate through all query parameters
    queryArgs.VisitAll(func(key, value []byte) {
        fmt.Printf("Parameter %s = %s\n", string(key), string(value))
    })
    
    // Respond with search results
    fmt.Fprintf(ctx, "Search results for: %s (Page %s)", term, page)
}
```

#### POST Request Body

Ming makes it easy to handle various types of POST request data:

```go
// Handle JSON POST request
func JsonHandler(ctx *fasthttp.RequestCtx) {
    // Get full request body as []byte
    body := ming.Body(ctx)
    
    // Now you can unmarshal the JSON data
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err != nil {
        ctx.SetStatusCode(400)
        fmt.Fprintf(ctx, "Invalid JSON: %s", err.Error())
        return
    }
    
    // Process the data
    fmt.Fprintf(ctx, "Received JSON with %d fields", len(data))
}

// Handle form POST request
func FormHandler(ctx *fasthttp.RequestCtx) {
    // Access form values
    username := string(ctx.FormValue("username"))
    email := string(ctx.FormValue("email"))
    
    // Check if a specific form field was provided
    if len(ctx.FormValue("username")) == 0 {
        ctx.SetStatusCode(400)
        ctx.WriteString("Username is required")
        return
    }
    
    // Get all form values
    ctx.PostArgs().VisitAll(func(key, value []byte) {
        fmt.Printf("Form field %s = %s\n", string(key), string(value))
    })
    
    // Access uploaded files (for multipart/form-data)
    ctx.Request.SetBodyStream(ctx.RequestBodyStream(), int(ctx.Request.Header.ContentLength()))
    form, err := ctx.MultipartForm()
    if err == nil {
        // Get uploaded files by form field name
        if files, ok := form.File["upload"]; ok && len(files) > 0 {
            filename := files[0].Filename
            // Process the uploaded file
        }
    }
    
    fmt.Fprintf(ctx, "Form processed for user: %s", username)
}
```

#### Working with Different Content Types

Ming handles various content types transparently:

- **application/json**: Use `ming.Body(ctx)` and then unmarshal the JSON
- **application/x-www-form-urlencoded**: Use `ctx.FormValue(key)` to access form fields
- **multipart/form-data**: Use `ctx.MultipartForm()` to access form fields and uploaded files
- **text/plain** or other raw data: Use `ming.Body(ctx)` to access the raw body

### Static Files
```go
r.Static("./public", true)  // Serve static files with directory listing
```

### Error Handling
```go
r.NotFound = func(ctx *fasthttp.RequestCtx) {
    ctx.SetStatusCode(404)
    ctx.WriteString("Page not found!")
}

r.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
    ctx.SetStatusCode(405) 
    ctx.WriteString("Method not allowed!")
}
```

## Examples

```go
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
	r.Static("./", true)
	r.Get("/", Home)
	r.Post("/add", PostHandler)
	r.All("/all", AllHandler)
	r.Get("/search", SearchHandler)
	r.Run("127.0.0.1:8000")
    // r.Run(":8000")
}
```

See more examples in the [_examples](_examples/) directory:
- [basic example](_examples/main.go) - Basic routing and static files
- [parameters example](_examples/parameters_example.go) - Comprehensive parameter features
- [simple example](_examples/simple_example.go) - Quick parameter demonstration

## Feature Details

### High Performance

Ming is built on top of [fasthttp](https://github.com/valyala/fasthttp), which provides exceptional performance compared to the standard `net/http` package. The router is designed to be:

- **Ultra-fast**: Up to 10x faster than standard Go net/http routing
- **Memory-efficient**: Minimizes allocations to reduce GC pressure
- **Optimized routing**: Uses a fast radix tree implementation for route matching
- **Connection reuse**: Takes advantage of fasthttp's connection pooling

Performance benchmarks show Ming handling thousands of requests per second with minimal resource consumption. For high-load applications, this translates to better response times and lower infrastructure costs.

### Named Parameters

Ming provides a simple way to extract variables from URL paths:

```go
// Define a route with a named parameter
r.Get("/user/{name}", func(ctx *fasthttp.RequestCtx) {
    // Extract the parameter value
    name := ming.Param(ctx, "name")
    fmt.Fprintf(ctx, "Hello, %s!", name)
})
```

The parameter values are extracted at request time and can be accessed using `ming.Param(ctx, "paramName")`. Named parameters are perfect for RESTful APIs or any scenario where you need to capture parts of the URL path.

### Optional Parameters

Ming supports optional parameters with the `?` syntax, allowing a single route to handle multiple URL patterns:

```go
// This route handles both /api/users and /api/users/42
r.Get("/api/users/{id?}", func(ctx *fasthttp.RequestCtx) {
    id := ming.Param(ctx, "id")
    if id == "" {
        // Handle case without ID (list all users)
        fmt.Fprintf(ctx, "List all users")
    } else {
        // Handle case with ID (get specific user)
        fmt.Fprintf(ctx, "Get user with ID: %s", id)
    }
})
```

This feature helps reduce code duplication by handling related endpoints in a single handler.

### Regex Validation

Ming supports parameter validation using regex patterns, ensuring that parameters match specific formats:

```go
// Only match numeric IDs
r.Get("/product/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
    id := ming.Param(ctx, "id")
    // ID is guaranteed to contain only digits
    fmt.Fprintf(ctx, "Product ID: %s", id)
})

// Only match alphabetic categories
r.Get("/category/{name:[a-zA-Z]+}", func(ctx *fasthttp.RequestCtx) {
    name := ming.Param(ctx, "name")
    // Name is guaranteed to contain only letters
    fmt.Fprintf(ctx, "Category: %s", name)
})

// Match specific date format
r.Get("/date/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", func(ctx *fasthttp.RequestCtx) {
    date := ming.Param(ctx, "date")
    // Date is guaranteed to be in YYYY-MM-DD format
    fmt.Fprintf(ctx, "Date: %s", date)
})
```

Regex validation helps ensure data integrity and proper API usage, preventing malformed requests from reaching your handlers.

### Catch-All Routes

For situations where you need to capture an entire path segment, including slashes, Ming provides catch-all parameters:

```go
// Match any path under /files/
r.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
    filepath := ming.Param(ctx, "filepath")
    fmt.Fprintf(ctx, "Requested file: %s", filepath)
})
```

This is particularly useful for:
- File servers
- Proxy routes
- Backend APIs that preserve URL structure
- Documentation systems

### Priority-Based Routing

Ming uses a sophisticated radix tree to store and match routes. The routing algorithm follows specific priority rules:

1. **Static routes** have the highest priority
2. **More specific routes** take precedence over general ones
3. **Parameter routes** have lower priority than static routes

For example:

```go
r.Get("/user/profile", specificHandler)  // Higher priority - matches exactly /user/profile
r.Get("/user/{name}", generalHandler)    // Lower priority - matches /user/john but not /user/profile
```

This priority system ensures that the most appropriate handler is always selected for each request.

### Route Conflict Resolution

When multiple routes could match the same URL pattern, Ming resolves conflicts using these rules:

1. Static routes always win over parameter routes
2. The most specific route wins (fewer parameters/wildcards)
3. Routes defined earlier win when other rules don't apply

For example:

```go
// This will match /api/v1/users exactly
r.Get("/api/v1/users", listUsersV1)

// This will match any other version (v2, v3, etc.)
r.Get("/api/{version}/users", listUsersGeneric)

// This won't match any of the above, because it's less specific
r.Get("/{service}/{version}/users", generalHandler)
```

This conflict resolution system makes routing predictable and intuitive, while still providing flexibility.

## Advanced Usage

### Static File Serving

Ming provides a simple way to serve static files from your filesystem. The `Static` method configures the router to serve files from a specified directory:

```go
r.Static(rootPath string, isIndexPage bool)
```

Parameters:
- `rootPath`: The root directory path from which to serve files
- `isIndexPage`: Boolean flag to enable/disable automatic generation of directory listing pages

#### Basic Static File Server

```go
r := ming.New()

// Serve files from the current directory with directory listing enabled
r.Static("./", true)

// Now run the server
r.Run(":8080")
```

#### Serving from a Specific Directory

```go
r := ming.New()

// Serve files from the "public" directory without directory listings
r.Static("./public", false)

r.Run(":8080")
```

#### How Static File Serving Works

Ming's static file serving uses the `fasthttp.FS` functionality behind the scenes. When you call `r.Static()`:

1. It configures the `NotFound` handler of the router to serve static files
2. Any request that doesn't match a defined route will attempt to serve a file from the specified directory
3. If a matching file is found, it's served with the appropriate content-type
4. If directory listings are enabled and a directory is requested, an index page is generated
5. If no matching file exists, a 404 response is returned

#### Example with Mixed Dynamic Routes and Static Files

```go
r := ming.New()

// Define your API routes first
r.Get("/api/users", apiUsersHandler)
r.Post("/api/login", apiLoginHandler)

// Then configure static file serving
// This will handle any routes not matched by the API endpoints
r.Static("./public", true)

r.Run(":8080")
```

#### Best Practices for Static File Serving

1. **Security**: Be cautious about which directories you expose. Never serve files from sensitive system directories.

2. **Directory Listings**: Use `isIndexPage: false` in production to prevent directory structure exposure.

3. **Performance**: For high-traffic sites, consider using a dedicated static file server like Nginx or a CDN.

4. **Route Priority**: Define all your dynamic routes before calling `r.Static()` to ensure proper routing.

5. **Hidden Files**: Note that Ming will serve all files in the specified directory, including hidden files, unless handled by other routes.

#### Caching Configuration

When using Ming's static file server, you may want to implement HTTP caching headers yourself in a custom middleware for optimal performance.

### Route Priority and Conflict Resolution

Ming router resolves route conflicts with these rules:

1. Static routes always take precedence over parameter routes
   ```go
   r.Get("/user/profile", profileHandler)  // This takes precedence
   r.Get("/user/{id}", userHandler)        // This matches any other /user/X
   ```

2. When routes conflict, the most specific one wins
   ```go
   r.Get("/api/v1/users", listUsersHandler)      // Specific to /api/v1/users
   r.Get("/api/{version}/users", versionHandler) // For any other version
   ```

### Panic Recovery

Ming provides built-in panic recovery to prevent crashes:

```go
r := ming.New()

// Custom panic handler
r.PanicHandler = func(ctx *fasthttp.RequestCtx, p interface{}) {
    ctx.SetStatusCode(500)
    fmt.Fprintf(ctx, "Internal error: %v", p)
    // Log the error, notify admins, etc.
}
```

### Creating Public APIs with Static File Serving

Ming is ideal for creating public APIs that also serve static assets for documentation. Here's a common pattern:

```go
r := ming.New()

// API routes with versioning
r.Get("/api/v1/users", apiV1GetUsers)
r.Post("/api/v1/users", apiV1CreateUser)
r.Get("/api/v1/products", apiV1GetProducts)

// API documentation routes
r.Get("/docs/api", apiDocsHandler)

// Serve static files (API docs, frontend assets, etc.)
r.Static("./public", false)
```

#### Organizing Static Assets

For production applications, consider organizing your static assets with this directory structure:

```
public/
  ├── assets/
  │    ├── css/
  │    ├── images/
  │    └── js/
  ├── docs/
  │    ├── api/
  │    └── guides/
  └── index.html
```

#### Security Considerations

When serving static files in production:

1. Never expose sensitive configuration files or data
2. Disable directory listings with `r.Static("./public", false)`
3. Consider adding a middleware for rate limiting static file access
4. Set proper cache headers for static assets
5. For larger applications, consider a dedicated CDN or static file server

See the [static file example](_examples/static_file_example.go) for a complete working example.

### Authentication and Middleware

Ming router supports middleware patterns for authentication, authorization, and other cross-cutting concerns. Here's how to implement authentication:

```go
// Authentication middleware
func authMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
    return func(ctx *fasthttp.RequestCtx) {
        // Get token from request header
        token := string(ctx.Request.Header.Peek("Authorization"))
        
        // Validate token
        if !validateToken(token) {
            ctx.SetStatusCode(fasthttp.StatusUnauthorized)
            ctx.WriteString("Unauthorized")
            return
        }
        
        // Token is valid, continue to the handler
        next(ctx)
    }
}

// Protected route with middleware
r.Get("/api/profile", authMiddleware(func(ctx *fasthttp.RequestCtx) {
    // This handler is only executed if auth is successful
    ctx.WriteString("Protected data")
}))
```

#### User Context

You can store user information in the request context for use in handlers:

```go
// In auth middleware
ctx.SetUserValue("userId", userId)

// In handler
userId := ctx.UserValue("userId").(string)
```

#### Role-Based Access Control

For role-based authorization, you can stack middleware:

```go
// Admin access middleware
func adminOnly(next fasthttp.RequestHandler) fasthttp.RequestHandler {
    return func(ctx *fasthttp.RequestCtx) {
        role := ctx.UserValue("role").(string)
        if role != "admin" {
            ctx.SetStatusCode(fasthttp.StatusForbidden)
            ctx.WriteString("Forbidden")
            return
        }
        next(ctx)
    }
}

// Endpoint requiring both authentication and admin role
r.Get("/admin/dashboard", authMiddleware(adminOnly(dashboardHandler)))
```

See the [auth example](_examples/auth_example.go) for a complete authentication implementation.

## Performance

Ming is built on fasthttp for maximum performance. See benchmarks:

Source: https://github.com/smallnest/go-web-framework-benchmark

![](https://github.com/smallnest/go-web-framework-benchmark/raw/master/cpubound_benchmark.png)

![](https://github.com/smallnest/go-web-framework-benchmark/raw/master/concurrency.png)