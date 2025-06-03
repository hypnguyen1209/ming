package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

// API route handler
func apiHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "API response: %s", ctx.Path())
}

// Custom 404 handler
func custom404Handler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetContentType("text/html; charset=utf-8")
	fmt.Fprintf(ctx, "<html><body><h1>Custom 404 Page</h1><p>Path not found: %s</p></body></html>", ctx.Path())
}

func main() {
	// Create a demo directory structure for static files
	setupStaticFiles()

	// Create a new router
	r := ming.New()

	// Define API routes first
	r.Get("/api/info", apiHandler)
	r.Get("/api/status", apiHandler)

	// Set custom 404 handler (optional)
	// Note: This won't be used for static file serving if r.Static() is called
	r.NotFound = custom404Handler

	// Serve static files from the "static_demo" directory (with directory listing enabled)
	// Any request not matching a defined route will try to serve a static file
	r.Static("./static_demo", true)

	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("Try these URLs:")
	fmt.Println("- http://localhost:8080/           (Serves index.html)")
	fmt.Println("- http://localhost:8080/styles.css (Serves CSS file)")
	fmt.Println("- http://localhost:8080/images/    (Shows directory listing)")
	fmt.Println("- http://localhost:8080/api/info   (API endpoint)")
	fmt.Println("- http://localhost:8080/not-found  (404 for non-existent files)")

	// Start the server
	r.Run(":8080")
}

// Helper function to set up demo static files
func setupStaticFiles() {
	// Create directories
	os.RemoveAll("./static_demo") // Clean up any existing directory
	os.MkdirAll("./static_demo/images", os.ModePerm)
	os.MkdirAll("./static_demo/js", os.ModePerm)

	// Create index.html
	indexHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Ming Static File Server</title>
    <link rel="stylesheet" href="/styles.css">
    <script src="/js/app.js"></script>
</head>
<body>
    <h1>Welcome to Ming Static File Server</h1>
    <p>This is a demonstration of the static file serving capability.</p>
    <ul>
        <li><a href="/images/">Browse Images Directory</a></li>
        <li><a href="/api/info">API Example</a></li>
        <li><a href="/not-found">Test 404 Page</a></li>
    </ul>
</body>
</html>`
	os.WriteFile("./static_demo/index.html", []byte(indexHTML), 0644)

	// Create CSS file
	cssContent := `body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    color: #333;
}

h1 {
    color: #0066cc;
    border-bottom: 1px solid #eee;
    padding-bottom: 10px;
}

a {
    color: #0066cc;
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

ul {
    background-color: #f5f5f5;
    padding: 15px 15px 15px 30px;
    border-radius: 5px;
}`
	os.WriteFile("./static_demo/styles.css", []byte(cssContent), 0644)

	// Create JavaScript file
	jsContent := `document.addEventListener('DOMContentLoaded', function() {
    console.log('Static file server example loaded');
    
    // Add a timestamp to the page
    const footer = document.createElement('footer');
    footer.style.marginTop = '30px';
    footer.style.borderTop = '1px solid #eee';
    footer.style.paddingTop = '10px';
    footer.style.fontSize = '12px';
    footer.style.color = '#666';
    footer.textContent = 'Page rendered at: ' + new Date().toLocaleString();
    document.body.appendChild(footer);
});`
	os.WriteFile("./static_demo/js/app.js", []byte(jsContent), 0644)

	// Create a sample image (1x1 pixel transparent GIF)
	transparentGif := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00,
		0x00, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x21, 0xF9, 0x04, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
		0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3B,
	}
	os.WriteFile("./static_demo/images/sample.gif", transparentGif, 0644)

	// Create a README file in the images directory
	readmeContent := `# Image Directory
This directory contains sample images for the static file server example.

Current files:
- sample.gif - A 1x1 transparent GIF image
`
	os.WriteFile("./static_demo/images/README.md", []byte(readmeContent), 0644)

	// Print the created directory structure
	fmt.Println("Created static file demo directory structure:")
	filepath.Walk("./static_demo", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Printf("Directory: %s/\n", path)
		} else {
			fmt.Printf("File: %s (%d bytes)\n", path, info.Size())
		}
		return nil
	})
	fmt.Println()
}
