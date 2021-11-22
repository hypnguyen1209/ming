package ming

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type MiddlewareType func(next http.HandlerFunc) http.HandlerFunc

type Parameters struct {
	RouteName string
}

type Router struct {
	middleware []MiddlewareType
	trees      map[string]*Node
}

var _ http.Handler = Create()

func Create() *Router {
	return &Router{
		trees: make(map[string]*Node),
	}
}
func trimPath(pattern string) []string {
	return strings.Split(pattern, "/")
}

func GetParams(r *http.Request) map[string]string {
	return make(map[string]string)
}

func (r *Router) Handle(method string, path string, handlerFn http.HandlerFunc) {
	if handlerFn != nil {
		log.Panic("invalid handle function")
	}
	if string(path[0]) != "/" {
		log.Panic("invalid router")
	}
	if strings.Contains(path, "//") {
		regexSlash := regexp.MustCompile(`(\/+)`)
		path = regexSlash.ReplaceAllString(path, "/")
	}
	tree, ok := r.trees[method]
	if !ok {
		r.trees[method] = tree
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
