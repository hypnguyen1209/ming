# Named Parameters Documentation

Ming Router now supports advanced routing with named parameters, optional parameters, regex validation, and catch-all parameters.

## Features

### 1. Named Parameters

Named parameters are defined using curly braces `{name}` and match a single path segment.

```go
r.Get("/user/{name}", func(ctx *fasthttp.RequestCtx) {
    name := ming.Param(ctx, "name")
    fmt.Fprintf(ctx, "Hello, %s!", name)
})
```

**Examples:**
- `/user/john` ✅ matches
- `/user/jane` ✅ matches
- `/user/john/profile` ❌ no match (contains additional segment)
- `/user/` ❌ no match (empty parameter)

### 2. Optional Parameters

Optional parameters are defined by adding `?` after the parameter name: `{name?}`

```go
r.Get("/api/users/{id?}", func(ctx *fasthttp.RequestCtx) {
    id := ming.Param(ctx, "id")
    if id == "" {
        ctx.WriteString("List all users")
    } else {
        fmt.Fprintf(ctx, "Get user: %s", id)
    }
})
```

**Examples:**
- `/api/users/` ✅ matches (id is empty)
- `/api/users/123` ✅ matches (id is "123")

### 3. Regex Validation

Parameters can include regex validation using the syntax `{name:[regex]}`:

```go
// Only numeric IDs
r.Get("/product/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
    id := ming.Param(ctx, "id")
    fmt.Fprintf(ctx, "Product ID: %s", id)
})

// Only alphabetic names
r.Get("/category/{name:[a-zA-Z]+}", func(ctx *fasthttp.RequestCtx) {
    name := ming.Param(ctx, "name")
    fmt.Fprintf(ctx, "Category: %s", name)
})
```

**Examples:**
- `/product/123` ✅ matches
- `/product/abc` ❌ no match (not numeric)
- `/category/electronics` ✅ matches
- `/category/123` ❌ no match (not alphabetic)

### 4. Optional Parameters with Regex

Combine optional parameters with regex validation: `{name?:[regex]}`

```go
r.Get("/article/{slug?:[a-z0-9-]+}", func(ctx *fasthttp.RequestCtx) {
    slug := ming.Param(ctx, "slug")
    if slug == "" {
        ctx.WriteString("List all articles")
    } else {
        fmt.Fprintf(ctx, "Article: %s", slug)
    }
})
```

### 5. Catch-All Parameters

Catch-all parameters use `{name:*}` and match everything including slashes. They must be at the end of the path.

```go
r.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
    filepath := ming.Param(ctx, "filepath")
    fmt.Fprintf(ctx, "File: %s", filepath)
})
```

**Examples:**
- `/files/` ✅ matches (filepath is empty)
- `/files/readme.txt` ✅ matches
- `/files/docs/readme.txt` ✅ matches
- `/files/path/to/deep/file.pdf` ✅ matches

### 6. Named Parameters with Suffixes

Parameters can have suffixes:

```go
r.Get("/admin/{name}_profile", func(ctx *fasthttp.RequestCtx) {
    name := ming.Param(ctx, "name")
    fmt.Fprintf(ctx, "Admin profile: %s", name)
})
```

**Examples:**
- `/admin/john_profile` ✅ matches
- `/admin/jane_profile` ✅ matches
- `/admin/john` ❌ no match (missing suffix)

## API Reference

### Getting Parameter Values

```go
// Get parameter value as string
name := ming.Param(ctx, "paramName")

// Get raw user value (returns interface{})
value := ming.UserValue(ctx, "paramName")
```

### Complex Route Example

```go
r.Get("/api/{version:[v][0-9]+}/users/{userId:[0-9]+}/files/{filepath:*}", 
    func(ctx *fasthttp.RequestCtx) {
        version := ming.Param(ctx, "version")
        userId := ming.Param(ctx, "userId")
        filepath := ming.Param(ctx, "filepath")
        
        fmt.Fprintf(ctx, "API %s, User %s, File %s", 
            version, userId, filepath)
    })
```

This matches: `/api/v1/users/123/files/documents/report.pdf`

## Priority and Matching

The router uses a radix tree with priority-based matching:

1. Static routes have highest priority
2. Routes with parameters are evaluated by registration order
3. Catch-all routes have lowest priority

## Error Handling

- Invalid regex patterns will panic during route registration
- Catch-all parameters not at the end will panic
- Parameters that don't match their regex validation will not match the route
