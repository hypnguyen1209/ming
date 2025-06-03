package ming

import "github.com/valyala/fasthttp"

// GetMethod returns the HTTP method of the request as a string.
// This is a utility function that converts fasthttp method detection
// into a standardized string format.
func GetMethod(ctx *fasthttp.RequestCtx) string {
	switch true {
	case ctx.IsGet():
		return fasthttp.MethodGet
	case ctx.IsPost():
		return fasthttp.MethodPost
	case ctx.IsHead():
		return fasthttp.MethodHead
	case ctx.IsPut():
		return fasthttp.MethodPut
	case ctx.IsPatch():
		return fasthttp.MethodPatch
	case ctx.IsDelete():
		return fasthttp.MethodDelete
	case ctx.IsConnect():
		return fasthttp.MethodConnect
	case ctx.IsOptions():
		return fasthttp.MethodOptions
	case ctx.IsTrace():
		return fasthttp.MethodTrace
	default:
		return ""
	}
}
