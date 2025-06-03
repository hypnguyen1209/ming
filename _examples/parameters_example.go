package main

import (
	"fmt"
	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

// This example demonstrates all the parameter types supported by Ming router:
// 1. Named Parameters - {name}
// 2. Optional Parameters - {name?}
// 3. Regex Validation - {name:[pattern]}
// 4. Catch-All Parameters - {name:*}

func main() {
	r := ming.New()

	// Basic route
	r.Get("/", func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Ming Router Parameter Examples\n\n")
		fmt.Fprintf(ctx, "This example demonstrates all parameter types supported by Ming.\n\n")
		fmt.Fprintf(ctx, "Try the following routes:\n")
		fmt.Fprintf(ctx, "1. Named Parameters:\n")
		fmt.Fprintf(ctx, "   - /user/john\n")
		fmt.Fprintf(ctx, "   - /user/43/post/789\n\n")
		fmt.Fprintf(ctx, "2. Optional Parameters:\n")
		fmt.Fprintf(ctx, "   - /api/users (list all users)\n")
		fmt.Fprintf(ctx, "   - /api/users/42 (get specific user)\n")
		fmt.Fprintf(ctx, "   - /blog (all posts)\n")
		fmt.Fprintf(ctx, "   - /blog/golang (posts about golang)\n\n")
		fmt.Fprintf(ctx, "3. Regex Validation:\n")
		fmt.Fprintf(ctx, "   - /product/123 (numeric id)\n")
		fmt.Fprintf(ctx, "   - /product/abc (will 404 - not numeric)\n")
		fmt.Fprintf(ctx, "   - /category/electronics (alphabetic only)\n")
		fmt.Fprintf(ctx, "   - /user-info/john42 (alphanumeric only)\n")
		fmt.Fprintf(ctx, "   - /dates/2023-09-25 (date format only)\n\n")
		fmt.Fprintf(ctx, "4. Catch-All Parameters:\n")
		fmt.Fprintf(ctx, "   - /files/documents/report.pdf\n")
		fmt.Fprintf(ctx, "   - /files/images/vacations/beach.jpg\n\n")
		fmt.Fprintf(ctx, "5. Combined Example:\n")
		fmt.Fprintf(ctx, "   - /api/v1/report/sales/2023/quarterly/Q3.pdf\n")
	})

	//
	// 1. NAMED PARAMETERS
	//

	// Simple named parameter
	r.Get("/user/{name}", func(ctx *fasthttp.RequestCtx) {
		name := ming.Param(ctx, "name")
		fmt.Fprintf(ctx, "User Profile: %s\n\n", name)
		fmt.Fprintf(ctx, "This route uses a named parameter: /user/{name}\n")
		fmt.Fprintf(ctx, "The parameter value is extracted using: ming.Param(ctx, \"name\")")
	})

	// Multiple named parameters
	r.Get("/user/{id}/post/{postId}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		postId := ming.Param(ctx, "postId")
		fmt.Fprintf(ctx, "User ID: %s, Post ID: %s\n\n", id, postId)
		fmt.Fprintf(ctx, "This route uses multiple named parameters: /user/{id}/post/{postId}\n")
		fmt.Fprintf(ctx, "Multiple parameters can be mixed in any order in the URL pattern.")
	})

	//
	// 2. OPTIONAL PARAMETERS
	//

	// Optional parameter example - handles both /api/users and /api/users/123
	r.Get("/api/users/{id?}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		if id == "" {
			fmt.Fprintf(ctx, "List of all users\n\n")
			fmt.Fprintf(ctx, "This route handles /api/users without an ID\n")
			fmt.Fprintf(ctx, "The optional parameter is defined as: /api/users/{id?}\n")
			fmt.Fprintf(ctx, "When the parameter is not provided, ming.Param(ctx, \"id\") returns an empty string.")
		} else {
			fmt.Fprintf(ctx, "Details for user ID: %s\n\n", id)
			fmt.Fprintf(ctx, "This route handles /api/users/{id} with a specific ID\n")
			fmt.Fprintf(ctx, "The same handler processes both cases with and without the ID parameter.")
		}
	})

	// Another optional parameter example
	r.Get("/blog/{topic?}", func(ctx *fasthttp.RequestCtx) {
		topic := ming.Param(ctx, "topic")
		if topic == "" {
			fmt.Fprintf(ctx, "Blog home - All topics\n\n")
			fmt.Fprintf(ctx, "This route handles /blog without a topic\n")
		} else {
			fmt.Fprintf(ctx, "Blog articles about: %s\n\n", topic)
			fmt.Fprintf(ctx, "This route handles /blog/{topic} with a specific topic\n")
		}
	})

	//
	// 3. REGEX VALIDATION PARAMETERS
	//
		id := ming.Param(ctx, "id")
			// Numeric validation - only digits
	r.Get("/product/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		fmt.Fprintf(ctx, "Product ID: %s (numeric only)\n\n", id)
		fmt.Fprintf(ctx, "This route uses regex validation: /product/{id:[0-9]+}\n")
		fmt.Fprintf(ctx, "The pattern [0-9]+ ensures the ID contains only digits\n")
		fmt.Fprintf(ctx, "If the ID contained non-digits, this route would not match.")
	})

	// Alphabetic validation - only letters
	r.Get("/category/{name:[a-zA-Z]+}", func(ctx *fasthttp.RequestCtx) {
		name := ming.Param(ctx, "name")
		fmt.Fprintf(ctx, "Category: %s (alphabetic only)\n\n", name)
		fmt.Fprintf(ctx, "This route uses regex validation: /category/{name:[a-zA-Z]+}\n")
		fmt.Fprintf(ctx, "The pattern [a-zA-Z]+ ensures the name contains only letters.")
	})

	// Alphanumeric validation
	r.Get("/user-info/{username:[a-zA-Z0-9]+}", func(ctx *fasthttp.RequestCtx) {
		username := ming.Param(ctx, "username")
		fmt.Fprintf(ctx, "Username: %s (alphanumeric only)\n\n", username)
		fmt.Fprintf(ctx, "This route uses regex validation: /user-info/{username:[a-zA-Z0-9]+}")
	})

	// Date format validation
	r.Get("/dates/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", func(ctx *fasthttp.RequestCtx) {
		date := ming.Param(ctx, "date")
		fmt.Fprintf(ctx, "Date: %s (YYYY-MM-DD format)\n\n", date)
		fmt.Fprintf(ctx, "This route uses regex validation: /dates/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}\n")
		fmt.Fprintf(ctx, "The pattern ensures the date follows the YYYY-MM-DD format.")
	})

	// Combined optional parameter with regex
	r.Get("/article/{slug?:[a-z0-9-]+}", func(ctx *fasthttp.RequestCtx) {
		slug := ming.Param(ctx, "slug")
		if slug == "" {
			fmt.Fprintf(ctx, "All Articles\n\n")
			fmt.Fprintf(ctx, "This route combines optional parameter with regex validation\n")
			fmt.Fprintf(ctx, "Pattern: /article/{slug?:[a-z0-9-]+}")
		} else {
			fmt.Fprintf(ctx, "Article with slug: %s\n\n", slug)
			fmt.Fprintf(ctx, "The slug is validated to contain only lowercase letters, digits, and hyphens.")
		}
	})

	//
	// 4. CATCH-ALL PARAMETERS
	//

	// Catch-all parameter example
	r.Get("/files/{filepath:*}", func(ctx *fasthttp.RequestCtx) {
		filepath := ming.Param(ctx, "filepath")
		fmt.Fprintf(ctx, "File path: %s\n\n", filepath)
		fmt.Fprintf(ctx, "This route uses a catch-all parameter: /files/{filepath:*}\n")
		fmt.Fprintf(ctx, "The :* syntax captures the entire remaining path, including slashes\n")
		fmt.Fprintf(ctx, "This is perfect for file paths, proxy routes, or any nested resources.")
	})

	//
	// 5. COMBINED EXAMPLES
	//

	// Complex route with mixed parameters
	r.Get("/api/{version:[v][0-9]+}/report/{type:[a-z]+}/{year:[0-9]{4}}/{format}/{filename:*}", func(ctx *fasthttp.RequestCtx) {
		version := ming.Param(ctx, "version")
		reportType := ming.Param(ctx, "type")
		year := ming.Param(ctx, "year")
		format := ming.Param(ctx, "format")
		filename := ming.Param(ctx, "filename")
		
		fmt.Fprintf(ctx, "Complex route with multiple parameter types:\n\n")
		fmt.Fprintf(ctx, "API Version: %s (validates format v followed by digits)\n", version)
		fmt.Fprintf(ctx, "Report Type: %s (validates lowercase letters only)\n", reportType)
		fmt.Fprintf(ctx, "Year: %s (validates exactly 4 digits)\n", year)
		fmt.Fprintf(ctx, "Format: %s (no validation)\n", format)
		fmt.Fprintf(ctx, "Filename: %s (catch-all for remaining path)\n\n", filename)
		
		fmt.Fprintf(ctx, "This example demonstrates combining multiple parameter types in one route pattern.")
	})
	})

	// POST example with named parameters
	r.Post("/user/{id}/update", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		body := string(ming.Body(ctx))
		fmt.Fprintf(ctx, "Updating user %s with data: %s", id, body)
	})

	// Error handling
	r.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString("Page not found!")
	}

	r.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.WriteString("Method not allowed!")
	}

	fmt.Println("Server starting on :8080")
	fmt.Println("\nTry these routes:")
	fmt.Println("GET  /user/john")
	fmt.Println("GET  /user/123/post/456")
	fmt.Println("GET  /api/v1/users/123")
	fmt.Println("GET  /api/v1/users/")
	fmt.Println("GET  /product/123")
	fmt.Println("GET  /category/electronics")
	fmt.Println("GET  /article/my-awesome-post")
	fmt.Println("GET  /files/documents/readme.txt")
	fmt.Println("GET  /api/v1/users/123/files/documents/report.pdf")
	fmt.Println("GET  /admin/john_profile")
	fmt.Println("POST /user/123/update")

	r.Run(":8080")
}
