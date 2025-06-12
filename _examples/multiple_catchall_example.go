package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/hypnguyen1209/ming/v2"
)

// FileHandler handles requests for serving static files
func FileHandler(ctx *fasthttp.RequestCtx) {
	// Extract the filepath from the request
	filepath := ming.Param(ctx, "filepath")
	
	// In a real application, you would validate and serve the file
	// For this example, we'll just print the filepath
	fmt.Fprintf(ctx, "File Handler: Serving file at path: %s\n", filepath)
}

// DocumentHandler handles requests for document files
func DocumentHandler(ctx *fasthttp.RequestCtx) {
	// Extract the document path from the request
	docpath := ming.Param(ctx, "docpath")
	
	fmt.Fprintf(ctx, "Document Handler: Accessing document at path: %s\n", docpath)
}

// APIProxyHandler handles proxy requests to external APIs
func APIProxyHandler(ctx *fasthttp.RequestCtx) {
	// Extract the URL from the request
	url := ming.Param(ctx, "url")
	
	// In a real application, you would forward the request to the URL
	// For this example, we'll just print the URL
	fmt.Fprintf(ctx, "API Proxy: Forwarding request to: %s\n", url)
	
	// NOTE: Due to how URL paths are parsed, if a URL contains "https://",
	// it will be captured as "https://" (with one slash instead of two).
	// You would need to handle this in your application:
	if strings.HasPrefix(url, "http:/") && !strings.HasPrefix(url, "http://") {
		url = strings.Replace(url, "http:/", "http://", 1)
	} else if strings.HasPrefix(url, "https:/") && !strings.HasPrefix(url, "https://") {
		url = strings.Replace(url, "https:/", "https://", 1)
	}
	
	fmt.Fprintf(ctx, "Corrected URL: %s\n", url)
}

// MediaHandler handles streaming media content
func MediaHandler(ctx *fasthttp.RequestCtx) {
	// Extract the media path from the request
	mediapath := ming.Param(ctx, "mediapath")
	
	fmt.Fprintf(ctx, "Media Handler: Streaming media from: %s\n", mediapath)
}

// Error handling middleware
func errorHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				fmt.Fprintf(ctx, "Internal Server Error")
			}
		}()
		next(ctx)
	}
}

// Logging middleware
func logger(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		
		log.Printf("%s %s - Start", method, path)
		next(ctx)
		log.Printf("%s %s - End (Status: %d)", method, path, ctx.Response.StatusCode())
	}
}

func main() {
	// Create a new Ming router
	router := ming.New()
	
	// Set up middleware
	router.PanicHandler = func(ctx *fasthttp.RequestCtx, err interface{}) {
		log.Printf("Panic: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString("Internal Server Error")
	}
	
	// Define static routes first (higher priority)
	router.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Welcome to the Multiple Catch-All Routes Example Server!")
		ctx.WriteString("\n\nTry these routes:\n")
		ctx.WriteString("1. /files/example.txt\n")
		ctx.WriteString("2. /documents/report.pdf\n")
		ctx.WriteString("3. /api/v1/proxy/https://example.com\n")
		ctx.WriteString("4. /media/videos/demo.mp4\n")
	})
	
	// Handle requests for empty directories with specific handlers
	router.Get("/files/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("File directory root. Browse files by adding a path.")
	})
	
	router.Get("/documents/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Document directory root. Browse documents by adding a path.")
	})
	
	router.Get("/media/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Media directory root. Browse media by adding a path.")
	})
	
	// Register multiple catch-all routes with different prefixes
	// Each catch-all captures everything after its prefix
	router.Get("/files/{filepath:*}", errorHandler(logger(FileHandler)))
	router.Get("/documents/{docpath:*}", errorHandler(logger(DocumentHandler)))
	router.Get("/api/v1/proxy/{url:*}", errorHandler(logger(APIProxyHandler)))
	router.Get("/media/{mediapath:*}", errorHandler(logger(MediaHandler)))
	
	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Println("Try accessing different catch-all routes such as:")
	log.Println("- http://localhost:8080/files/example.txt")
	log.Println("- http://localhost:8080/documents/report.pdf")
	log.Println("- http://localhost:8080/api/v1/proxy/https://example.com")
	log.Println("- http://localhost:8080/media/videos/demo.mp4")
	
	// Use standard library to start HTTP server
	err := http.ListenAndServe(":8080", router.Handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
