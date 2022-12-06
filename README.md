# ming

Custom HTTP Mux lightweight and high performance 🥗

![](https://i.imgur.com/yCMS1yq.png)


## Examples:

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
	r.Get("/", Home)
	r.Post("/add", PostHandler)
	r.All("/all", AllHandler)
	r.Get("/search", SearchHandler)
	r.Run("127.0.0.1:8000")
	//r.Run(":8000")
}
```

## Test

Source: https://github.com/smallnest/go-web-framework-benchmark

![](https://github.com/smallnest/go-web-framework-benchmark/raw/master/cpubound_benchmark.png)

![](https://github.com/smallnest/go-web-framework-benchmark/raw/master/concurrency.png)


## Base on

+ https://pkg.go.dev/github.com/valyala/fasthttp

+ https://github.com/fasthttp
## 🎊 Inspired by

+ @fasthttp ([router](https://github.com/fasthttp/router))

+ @julienschmidt ([httprouter](https://github.com/julienschmidt/httprouter))

+ @bmizerany ([pat](https://github.com/bmizerany/pat))

