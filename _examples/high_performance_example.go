package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/hypnguyen1209/ming/v2"
	"github.com/valyala/fasthttp"
)

// High Performance example demonstrating:
// 1. Ming's fasthttp foundation
// 2. Simple benchmarking comparison with net/http
// 3. Concurrent request handling
// 4. Low overhead routing

// Simple response handler for Ming
func mingHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.WriteString(`{"status":"success","message":"Ming router response"}`)
}

// Simple response handler for standard net/http
func stdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Standard net/http response"}`))
}

// Generate random path for testing
func randomPath() string {
	return "/test/" + strconv.Itoa(rand.Intn(100000))
}

// Simple benchmark function for Ming router
func benchmarkMing(routeCount, requestCount int) time.Duration {
	r := ming.New()
	
	// Register routes
	for i := 0; i < routeCount; i++ {
		path := "/test/" + strconv.Itoa(i)
		r.Get(path, mingHandler)
	}
	
	// Prepare request
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://localhost:8080/test/50")
	req.Header.SetMethod("GET")
	defer fasthttp.ReleaseRequest(req)
	
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	
	// Run benchmark
	start := time.Now()
	for i := 0; i < requestCount; i++ {
		client := &fasthttp.Client{}
		handler := r.Handler
		
		err := fasthttp.ServeConn(nil, req, resp, handler)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
	
	return time.Since(start)
}

// Simple concurrent load test for Ming
func concurrentLoadTest(routeCount, requestCount, concurrency int) {
	r := ming.New()
	
	// Register routes
	for i := 0; i < routeCount; i++ {
		path := "/test/" + strconv.Itoa(i)
		r.Get(path, mingHandler)
	}
	
	// Set up a simple server
	go func() {
		fmt.Println("Starting Ming server on :8080 for concurrent load testing...")
		if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
			log.Fatalf("Error in ListenAndServe: %v", err)
		}
	}()
	
	// Wait for server to start
	time.Sleep(500 * time.Millisecond)
	
	// Create a client with timeout
	client := &fasthttp.Client{
		MaxConnsPerHost: concurrency,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
	}
	
	var wg sync.WaitGroup
	requestsPerGoroutine := requestCount / concurrency
	
	start := time.Now()
	
	// Launch concurrent goroutines to send requests
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			
			for j := 0; j < requestsPerGoroutine; j++ {
				// Choose a random route
				routeNum := rand.Intn(routeCount)
				uri := "http://localhost:8080/test/" + strconv.Itoa(routeNum)
				
				// Create request
				req := fasthttp.AcquireRequest()
				req.SetRequestURI(uri)
				req.Header.SetMethod("GET")
				
				// Create response object
				resp := fasthttp.AcquireResponse()
				
				// Send request
				err := client.Do(req, resp)
				
				// Clean up
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(resp)
				
				if err != nil {
					fmt.Printf("Error in goroutine %d: %v\n", num, err)
				}
			}
		}(i)
	}
	
	// Wait for all requests to complete
	wg.Wait()
	
	elapsed := time.Since(start)
	rps := float64(requestCount) / elapsed.Seconds()
	
	fmt.Printf("\n=== Concurrent Load Test Results ===\n")
	fmt.Printf("Total requests:    %d\n", requestCount)
	fmt.Printf("Concurrency level: %d\n", concurrency)
	fmt.Printf("Time taken:        %v\n", elapsed)
	fmt.Printf("Requests per sec:  %.2f\n", rps)
}

func main() {
	// Print runtime information
	fmt.Printf("Ming Router - High Performance Example\n")
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
	
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Check if we should run the concurrent test
	if len(os.Args) > 1 && os.Args[1] == "concurrent" {
		routeCount := 100      // Number of routes to register
		requestCount := 10000  // Total number of requests to make
		concurrency := 10      // Number of concurrent goroutines
		
		// Allow overriding values from command line
		if len(os.Args) > 2 {
			if val, err := strconv.Atoi(os.Args[2]); err == nil {
				concurrency = val
			}
		}
		if len(os.Args) > 3 {
			if val, err := strconv.Atoi(os.Args[3]); err == nil {
				requestCount = val
			}
		}
		
		fmt.Printf("\nRunning concurrent load test with %d routes, %d requests, %d concurrency...\n", 
			routeCount, requestCount, concurrency)
		
		concurrentLoadTest(routeCount, requestCount, concurrency)
		return
	}
	
	// Simple Ming demonstration with essential features
	r := ming.New()
	
	// Register some routes
	r.Get("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("Ming Router - High Performance Example")
	})
	
	r.Get("/json", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		ctx.WriteString(`{"status":"success","message":"High-performance JSON response"}`)
	})
	
	r.Get("/user/{id:[0-9]+}", func(ctx *fasthttp.RequestCtx) {
		id := ming.Param(ctx, "id")
		ctx.SetContentType("application/json")
		fmt.Fprintf(ctx, `{"id":"%s","name":"User %s"}`, id, id)
	})
	
	// Run simple benchmark
	routeCount := 100
	requestCount := 1000
	
	fmt.Printf("\nRunning simple benchmark with %d routes and %d requests...\n", routeCount, requestCount)
	duration := benchmarkMing(routeCount, requestCount)
	
	fmt.Printf("\n=== Benchmark Results ===\n")
	fmt.Printf("Routes:         %d\n", routeCount)
	fmt.Printf("Requests:       %d\n", requestCount)
	fmt.Printf("Total time:     %v\n", duration)
	fmt.Printf("Avg time/req:   %v\n", duration/time.Duration(requestCount))
	fmt.Printf("Requests/sec:   %.2f\n\n", float64(requestCount)/duration.Seconds())
	
	// Print instructions
	fmt.Println("To start the server, run the program with no arguments.")
	fmt.Println("To run a concurrent load test, use: go run high_performance_example.go concurrent [concurrency] [requests]")
	fmt.Println("\nStarting the server on http://127.0.0.1:8080...")
	
	// Start the server
	r.Run(":8080")
}
