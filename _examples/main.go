// filepath: /root/debug_ming/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

// User represents a user in our system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Created  string `json:"created"`
}

// APIResponse is a standard response format for our API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Global in-memory user database for demonstration
var users = make(map[string]User)

// Middleware to log requests
func loggerMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()
		
		// Call the next handler
		next(ctx)
		
		// Log after the request is processed
		fmt.Printf(
			"[%s] %s - %s - %s - %d - %v\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			ctx.Method(),
			ctx.Path(),
			ctx.RemoteAddr(),
			ctx.Response.StatusCode(),
			time.Since(startTime),
		)
	}
}

// Middleware to check API key
func authMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Skip auth for certain paths
		path := string(ctx.Path())
		if path == "/" || path == "/ping" {
			next(ctx)
			return
		}

		// Check for API key in header
		apiKey := string(ctx.Request.Header.Peek("X-API-Key"))
		if apiKey != "dev-secret-key" {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			response := APIResponse{
				Success: false,
				Message: "Unauthorized - Invalid API Key",
			}
			json.NewEncoder(ctx).Encode(response)
			return
		}

		next(ctx)
	}
}

// Handler for the home page
func homeHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	fmt.Fprintf(ctx, "<html><body>")
	fmt.Fprintf(ctx, "<h1>Welcome to ming Web Server</h1>")
	fmt.Fprintf(ctx, "<p>A high-performance HTTP router for Go</p>")
	fmt.Fprintf(ctx, "<h2>Available Routes:</h2>")
	fmt.Fprintf(ctx, "<ul>")
	fmt.Fprintf(ctx, "<li>GET /ping - Health check</li>")
	fmt.Fprintf(ctx, "<li>GET /api/users - List all users</li>")
	fmt.Fprintf(ctx, "<li>GET /api/users/{id} - Get user by ID</li>")
	fmt.Fprintf(ctx, "<li>POST /api/users - Create a new user</li>")
	fmt.Fprintf(ctx, "<li>PUT /api/users/{id} - Update a user</li>")
	fmt.Fprintf(ctx, "<li>DELETE /api/users/{id} - Delete a user</li>")
	fmt.Fprintf(ctx, "<li>GET /static/{filepath:*} - Serve static files</li>")
	fmt.Fprintf(ctx, "</ul>")
	fmt.Fprintf(ctx, "<p>Note: All API endpoints require X-API-Key header with value 'dev-secret-key'</p>")
	fmt.Fprintf(ctx, "</body></html>")
}

// Handler for health check
func pingHandler(ctx *fasthttp.RequestCtx) {
	response := APIResponse{
		Success: true,
		Message: "pong",
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		},
	}
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Handler to list all users
func listUsersHandler(ctx *fasthttp.RequestCtx) {
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	
	response := APIResponse{
		Success: true,
		Data:    userList,
	}
	
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Handler to get a user by ID
func getUserHandler(ctx *fasthttp.RequestCtx) {
	id := ming.Param(ctx, "id")
	
	user, exists := users[id]
	if !exists {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("User with ID %s not found", id),
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Data:    user,
	}
	
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Handler to create a new user
func createUserHandler(ctx *fasthttp.RequestCtx) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		response := APIResponse{
			Success: false,
			Message: "Invalid JSON payload",
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	// Validate input
	if input.Username == "" || input.Email == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		response := APIResponse{
			Success: false,
			Message: "Username and email are required",
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	// Create new user
	id := fmt.Sprintf("user_%d", len(users)+1)
	user := User{
		ID:       id,
		Username: input.Username,
		Email:    input.Email,
		Created:  time.Now().Format(time.RFC3339),
	}
	
	users[id] = user
	
	ctx.SetStatusCode(fasthttp.StatusCreated)
	response := APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	}
	
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Handler to update a user
func updateUserHandler(ctx *fasthttp.RequestCtx) {
	id := ming.Param(ctx, "id")
	
	// Check if user exists
	user, exists := users[id]
	if !exists {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("User with ID %s not found", id),
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	// Parse input
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		response := APIResponse{
			Success: false,
			Message: "Invalid JSON payload",
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	// Update user fields if provided
	if input.Username != "" {
		user.Username = input.Username
	}
	
	if input.Email != "" {
		user.Email = input.Email
	}
	
	// Save updated user
	users[id] = user
	
	response := APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	}
	
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Handler to delete a user
func deleteUserHandler(ctx *fasthttp.RequestCtx) {
	id := ming.Param(ctx, "id")
	
	// Check if user exists
	_, exists := users[id]
	if !exists {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		response := APIResponse{
			Success: false,
			Message: fmt.Sprintf("User with ID %s not found", id),
		}
		json.NewEncoder(ctx).Encode(response)
		return
	}
	
	// Delete the user
	delete(users, id)
	
	response := APIResponse{
		Success: true,
		Message: fmt.Sprintf("User with ID %s deleted successfully", id),
	}
	
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
}

// Create a static directory if it doesn't exist
func ensureStaticDir() {
	staticDir := filepath.Join(".", "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		os.Mkdir(staticDir, 0755)
		
		// Create a sample index.html file
		indexPath := filepath.Join(staticDir, "index.html")
		indexContent := `<!DOCTYPE html>
<html>
<head>
    <title>Ming Static File Server</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #2c3e50; }
        .info { background-color: #f8f9fa; padding: 15px; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>Static File Server</h1>
    <div class="info">
        <p>This is a sample static file served by ming web server.</p>
        <p>You can place your static files in the <code>static</code> directory.</p>
    </div>
</body>
</html>`
		os.WriteFile(indexPath, []byte(indexContent), 0644)
	}
}

func main() {
	// Create static directory with sample file
	ensureStaticDir()
	
	// Add some example users
	users["user_1"] = User{
		ID:       "user_1",
		Username: "admin",
		Email:    "admin@example.com",
		Created:  time.Now().Format(time.RFC3339),
	}
	
	users["user_2"] = User{
		ID:       "user_2",
		Username: "test",
		Email:    "test@example.com",
		Created:  time.Now().Format(time.RFC3339),
	}
	
	// Initialize ming router
	r := ming.New()
	
	// Create middleware chain by wrapping handlers
	withLogger := func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return loggerMiddleware(handler)
	}
	
	withAuth := func(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
		return authMiddleware(handler)
	}
	
	// Define routes
	r.Get("/", withLogger(homeHandler))
	r.Get("/ping", withLogger(pingHandler))
	
	// API routes
	r.Get("/api/users", withLogger(withAuth(listUsersHandler)))
	r.Get("/api/users/{id}", withLogger(withAuth(getUserHandler)))
	r.Post("/api/users", withLogger(withAuth(createUserHandler)))
	r.Put("/api/users/{id}", withLogger(withAuth(updateUserHandler)))
	r.Delete("/api/users/{id}", withLogger(withAuth(deleteUserHandler)))
	
	// Example route with regex validation
	r.Get("/product/{id:[0-9]+}", withLogger(withAuth(func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		response := APIResponse{
			Success: true,
			Message: fmt.Sprintf("Product ID %s is valid", id),
		}
		ctx.SetContentType("application/json")
		json.NewEncoder(ctx).Encode(response)
	})))
	
	// Optional parameter example
	r.Get("/optional/{param?}", withLogger(func(ctx *fasthttp.RequestCtx) {
		param := ming.Param(ctx, "param")
		if param == "" {
			param = "not provided"
		}
		
		response := APIResponse{
			Success: true,
			Message: fmt.Sprintf("Optional parameter: %s", param),
		}
		
		ctx.SetContentType("application/json")
		json.NewEncoder(ctx).Encode(response)
	}))
	
	// Custom 404 handler for non-static paths
	r.NotFound = withLogger(func(ctx *fasthttp.RequestCtx) {
		// Handle the static files directory first
		path := string(ctx.Path())
		if strings.HasPrefix(path, "/static/") {
			// Strip the "/static/" prefix
			filePath := path[8:]
			fullPath := filepath.Join("./static", filePath)
			
			// Check if file exists
			if _, err := os.Stat(fullPath); err == nil {
				fasthttp.ServeFile(ctx, fullPath)
				return
			}
		}
		
		// Otherwise, return a 404 response
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetContentType("application/json")
		json.NewEncoder(ctx).Encode(APIResponse{
			Success: false,
			Message: fmt.Sprintf("Path not found: %s", path),
		})
	})
	
	// Start the server
	fmt.Println("Server starting on :8080")
	fmt.Println("Home page: http://localhost:8080/")
	fmt.Println("Health check: http://localhost:8080/ping")
	fmt.Println("API (requires X-API-Key: dev-secret-key header):")
	fmt.Println("  GET http://localhost:8080/api/users")
	fmt.Println("  GET http://localhost:8080/api/users/user_1")
	fmt.Println("  POST http://localhost:8080/api/users")
	fmt.Println("  PUT http://localhost:8080/api/users/user_1")
	fmt.Println("  DELETE http://localhost:8080/api/users/user_1")
	fmt.Println("Static files: http://localhost:8080/static/index.html")
	
	r.Run(":8080")
}