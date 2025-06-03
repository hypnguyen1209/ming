package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

// Simple in-memory user store
// In a real application, you would use a database
var users = map[string]string{
	"admin": hashPassword("admin123"),
	"john":  hashPassword("password123"),
	"jane":  hashPassword("securepwd456"),
}

// Store for refresh tokens
var refreshTokens = map[string]string{}

// Simple helper to hash passwords
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Basic auth middleware
func authMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Get token from header
		token := string(ctx.Request.Header.Peek("Authorization"))
		
		// Check if token starts with "Bearer "
		if !strings.HasPrefix(token, "Bearer ") {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			fmt.Fprintf(ctx, `{"error":"Unauthorized: Missing or invalid token"}`)
			return
		}
		
		// Validate token (this is a simple example, use JWT for production)
		token = token[7:] // Remove "Bearer " prefix
		username, valid := validateToken(token)
		if !valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			fmt.Fprintf(ctx, `{"error":"Unauthorized: Invalid token"}`)
			return
		}
		
		// Set username in context for handlers to use
		ctx.SetUserValue("username", username)
		
		// Call the next handler
		next(ctx)
	}
}

// Admin access middleware - adds additional access control layer
func adminOnlyMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		username := ctx.UserValue("username")
		if username != "admin" {
			ctx.SetStatusCode(fasthttp.StatusForbidden)
			fmt.Fprintf(ctx, `{"error":"Forbidden: Admin access required"}`)
			return
		}
		next(ctx)
	}
}

// Simple token generation (for demo purposes)
// In a real application, use JWT with proper signing
func generateToken(username string) string {
	// This is just a demo token format: username + timestamp + simple hash
	// In production, use JWT with proper signing
	timestamp := time.Now().Unix()
	tokenData := fmt.Sprintf("%s:%d", username, timestamp)
	hash := sha256.Sum256([]byte(tokenData + "secret-key"))
	token := base64.StdEncoding.EncodeToString([]byte(tokenData)) + "." + 
		hex.EncodeToString(hash[:8])
	return token
}

// Simple token validation
func validateToken(token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", false
	}
	
	// Decode the token data
	decoded, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	
	// Parse data
	tokenParts := strings.Split(string(decoded), ":")
	if len(tokenParts) != 2 {
		return "", false
	}
	
	username := tokenParts[0]
	
	// Re-compute the hash to verify token integrity
	hash := sha256.Sum256([]byte(string(decoded) + "secret-key"))
	expectedHash := hex.EncodeToString(hash[:8])
	
	// Verify hash
	if parts[1] != expectedHash {
		return "", false
	}
	
	return username, true
}

// Generate refresh token
func generateRefreshToken(username string) string {
	tokenBytes := make([]byte, 32)
	for i := 0; i < len(tokenBytes); i++ {
		tokenBytes[i] = byte(i + 42) // Simple deterministic value for demo
	}
	token := hex.EncodeToString(tokenBytes)
	refreshTokens[token] = username
	return token
}

func main() {
	r := ming.New()
	
	// Set up custom 404 and panic handlers
	r.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetContentType("application/json")
		fmt.Fprintf(ctx, `{"error":"Not Found"}`)
	}
	
	r.PanicHandler = func(ctx *fasthttp.RequestCtx, err interface{}) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetContentType("application/json")
		fmt.Fprintf(ctx, `{"error":"Internal Server Error"}`)
		fmt.Printf("Panic recovered: %v\n", err)
	}
	
	// Public routes
	r.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		fmt.Fprintf(ctx, `{"message":"Welcome to Ming Auth Example API"}`)
	})
	
	// Authentication endpoint
	r.Post("/auth/login", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		// Parse credentials from request body
		username := string(ctx.FormValue("username"))
		password := string(ctx.FormValue("password"))
		
		// Validate credentials
		storedHash, exists := users[username]
		if !exists || storedHash != hashPassword(password) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			fmt.Fprintf(ctx, `{"error":"Invalid username or password"}`)
			return
		}
		
		// Generate access and refresh tokens
		accessToken := generateToken(username)
		refreshToken := generateRefreshToken(username)
		
		// Return tokens
		fmt.Fprintf(ctx, `{
			"access_token": "%s",
			"refresh_token": "%s",
			"token_type": "Bearer",
			"expires_in": 3600
		}`, accessToken, refreshToken)
	})
	
	// Token refresh endpoint
	r.Post("/auth/refresh", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		refreshToken := string(ctx.FormValue("refresh_token"))
		username, exists := refreshTokens[refreshToken]
		
		if !exists {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			fmt.Fprintf(ctx, `{"error":"Invalid refresh token"}`)
			return
		}
		
		// Generate new access token
		newAccessToken := generateToken(username)
		
		// Return new access token
		fmt.Fprintf(ctx, `{
			"access_token": "%s",
			"token_type": "Bearer",
			"expires_in": 3600
		}`, newAccessToken)
	})
	
	// User registration endpoint
	r.Post("/auth/register", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		username := string(ctx.FormValue("username"))
		password := string(ctx.FormValue("password"))
		
		// Simple validation
		if username == "" || password == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			fmt.Fprintf(ctx, `{"error":"Username and password are required"}`)
			return
		}
		
		// Check if user already exists
		if _, exists := users[username]; exists {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			fmt.Fprintf(ctx, `{"error":"Username already exists"}`)
			return
		}
		
		// Add user to the store
		users[username] = hashPassword(password)
		
		fmt.Fprintf(ctx, `{"message":"User registered successfully"}`)
	})
	
	// Protected API endpoints
	
	// User profile - requires authentication
	r.Get("/api/profile", authMiddleware(func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		username := ctx.UserValue("username").(string)
		
		// Return profile data
		fmt.Fprintf(ctx, `{
			"username": "%s",
			"profile": {
				"name": "User %s",
				"email": "%s@example.com",
				"role": "user"
			}
		}`, username, username, username)
	}))
	
	// User data - requires authentication
	r.Get("/api/data", authMiddleware(func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		username := ctx.UserValue("username").(string)
		
		// Return some user data
		fmt.Fprintf(ctx, `{
			"username": "%s",
			"data": [
				{"id": 1, "value": "Item 1"},
				{"id": 2, "value": "Item 2"},
				{"id": 3, "value": "Item 3"}
			]
		}`, username)
	}))
	
	// Admin endpoint - requires admin role
	r.Get("/api/admin", authMiddleware(adminOnlyMiddleware(func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		
		// List all users (admin only function)
		userList := []string{}
		for username := range users {
			userList = append(userList, username)
		}
		
		// Return admin data
		fmt.Fprintf(ctx, `{
			"message": "Admin access granted",
			"users": %q
		}`, userList)
	})))
	
	// Start server
	fmt.Println("Ming Auth Example Server")
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("\nAvailable routes:")
	fmt.Println("- GET  / (Public welcome endpoint)")
	fmt.Println("- POST /auth/register (Register new user)")
	fmt.Println("- POST /auth/login (Login and get tokens)")
	fmt.Println("- POST /auth/refresh (Refresh access token)")
	fmt.Println("- GET  /api/profile (Protected - requires auth)")
	fmt.Println("- GET  /api/data (Protected - requires auth)")
	fmt.Println("- GET  /api/admin (Protected - requires admin role)")
	fmt.Println("\nTest users:")
	fmt.Println("- admin / admin123 (has admin access)")
	fmt.Println("- john / password123 (regular user)")
	fmt.Println("- jane / securepwd456 (regular user)")
	fmt.Println("\nExample curl commands:")
	fmt.Println("1. Login:")
	fmt.Println("   curl -X POST http://localhost:8080/auth/login -d \"username=admin&password=admin123\"")
	fmt.Println("\n2. Access protected endpoint:")
	fmt.Println("   curl -H \"Authorization: Bearer YOUR_ACCESS_TOKEN\" http://localhost:8080/api/profile")
	
	r.Run(":8080")
}
