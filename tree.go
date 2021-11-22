package ming

import "net/http"

type Node struct {
	depth     uint64
	path      string
	handlerFn http.HandlerFunc
	children  []*Node
}

/* import (
	"net/http"
	"strings"
)

type Node struct {
	Key        string
	Path       string
	Handler    http.HandlerFunc
	Count      uint64
	List       map[string]*Node
	isPattern  bool
	Middleware []MiddlewareT
}

type Tree struct {
	Routes     map[string]*Node
	Parameters Parameters
	Root       *Node
}

func CreateNode(key string, count uint64) *Node {
	return &Node{
		Key:   key,
		Count: count,
		List:  make(map[string]*Node),
	}
}

func CreateTree() *Tree {
	return &Tree{
		Root:   CreateNode("/", 1),
		Routes: make(map[string]*Node),
	}
}

func trimPrefixRoute(pattern string) string {
	return strings.TrimPrefix(pattern, "/")
}

func splitRoute(pattern string) []string {
	return strings.Split(pattern, "/")
}

func (t *Tree) Add(handleFn http.HandlerFunc, pattern string, middlewares ...MiddlewareT) {
	curNode := t.Root
	curNode.Handler = handleFn
	curNode.isPattern = true
	curNode.Path = pattern
	if pattern != curNode.Key {
		newPattern := trimPrefixRoute(pattern)
		listPatterns := splitRoute(newPattern)
		for _, v := range listPatterns {
			node, ok := curNode.List[v]
			if !ok {
				newNode := CreateNode(v, curNode.Count+1)
				if len(middlewares) != 0 {
					newNode.Middleware = append(newNode.Middleware, middlewares...)
				}
				curNode.List[v] = newNode
			}
			curNode = node
		}
	}
	if len(middlewares) != 0 && curNode.Count == 1 {
		curNode.Middleware = append(curNode.Middleware, middlewares...)
	}
	if routeName := t.Parameters.RouteName; routeName != "" {
		t.Routes[routeName] = curNode
	}
}

func (t *Tree) Find(pattern string, isMatchRegex bool) []*Node {
	var queue, nodes []*Node
	node := t.Root
	if pattern == node.Path {
		nodes = append(nodes, node)
		return nodes
	}
	if !isMatchRegex {
		pattern = trimPrefixRoute(pattern)
	}
	listPatterns := splitRoute(pattern)
	for _, v := range listPatterns {
		child, ok := node.List[v]
		if !ok {
			if !isMatchRegex {
				return nodes
			} else {
				break
			}
		}
		if !isMatchRegex && pattern == child.Path {
			nodes = append(nodes, child)
			return nodes
		}
		node = child
	}
	queue = append(queue, node)
	for len(queue) > 0 {
		var temp []*Node
		for _, v := range queue {
			if v.isPattern {
				nodes = append(nodes, v)
			}
			for _, val := range v.List {
				temp = append(temp, val)
			}
		}
		queue = temp
	}
	return nodes
}
*/
