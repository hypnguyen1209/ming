package ming

import (
	"github.com/valyala/fasthttp"
)

type Tree []*Node
type Node struct {
	method  string
	path    string
	handler fasthttp.RequestHandler
}

func (t *Tree) Add(n *Node) {
	*t = append(*t, n)
}

func (t *Tree) FindMethod(method string) *Node {
	for _, v := range *t {
		if v.method == method {
			return v
		}
	}
	return nil
}

func (t *Tree) FindPath(path string) *Tree {
	result := &Tree{}
	for _, v := range *t {
		if v.path == path {
			result.Add(v)
		}
	}
	return result
}

func (t *Tree) Len() int {
	result := 0
	for i := 0; i < len(*t); i++ {
		result += 1
	}
	return result
}

func (n *Node) GetHandler() fasthttp.RequestHandler {
	return n.handler
}

func (t *Tree) GetMethodAll() *Node {
	if result := t.FindMethod("ALL"); result != nil {
		return result
	} else {
		return nil
	}
}
