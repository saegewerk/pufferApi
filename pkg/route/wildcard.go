package route

import "strings"

type WildcardRoot struct {
	Nodes   map[string]WildcardNode
	IsRoute bool
}

type WildcardNode struct {
	Nodes   map[string]WildcardNode
	IsRoute bool
}

func (root *WildcardRoot) Add(path string) {
	if path == "*" {
		root.IsRoute = true
		return
	}
	tokens := strings.Split(path, "/")[1:]

	if _, ok := root.Nodes[tokens[0]]; !ok {
		root.Nodes[tokens[0]] = WildcardNode{
			Nodes: make(map[string]WildcardNode),
		}
	}

	depth := 1
	root.Nodes[tokens[0]].Add(&tokens, &depth)

	return
}

func (node WildcardNode) Add(tokens *[]string, depth *int) {

	if _, ok := node.Nodes[(*tokens)[*depth]]; !ok {
		node.Nodes[(*tokens)[*depth]] = WildcardNode{
			Nodes: make(map[string]WildcardNode),
		}
	}
	*depth++
	if *depth < len(*tokens) {
		node.Nodes[(*tokens)[*depth-1]].Add(tokens, depth)
	} else {
		node.IsRoute = true
		return
	}
	return
}

func CreateWildcards() WildcardRoot {
	return WildcardRoot{
		Nodes: make(map[string]WildcardNode),
	}
}

func (root *WildcardRoot) Find(path string) string {
	tokens := strings.Split(path, "/")[1:]
	depth := 1
	actRoute := "/"
	lastRoute := actRoute

	if node, ok := root.Nodes[tokens[0]]; ok {
		actRoute += tokens[0]
		node.Find(&tokens, &actRoute, &lastRoute, &depth)
	}
	if root.IsRoute {
		return "*"
	}
	return lastRoute
}

func (node *WildcardNode) Find(tokens *[]string, actRoute, lastRoute *string, depth *int) {
	*depth++
	if *depth < len(*tokens) {
		if node, ok := node.Nodes[(*tokens)[*depth]]; ok {
			*actRoute += (*tokens)[0]
			if node.IsRoute {
				*lastRoute = *actRoute
			}
			node.Find(tokens, actRoute, lastRoute, depth)
		}
	} else {
		return
	}
	return
}
