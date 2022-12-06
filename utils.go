package ming

import "github.com/valyala/fasthttp"

func GetMethod(ctx *fasthttp.RequestCtx) string {
	switch true {
	case ctx.IsGet():
		return fasthttp.MethodGet
	case ctx.IsPost():
		return fasthttp.MethodPost
	case ctx.IsHead():
		return fasthttp.MethodPatch
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
