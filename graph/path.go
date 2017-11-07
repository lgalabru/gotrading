package graph

import "fmt"

type Path struct {
	ContextualNodes []*ContextualNode
}

func (p Path) contains(n ContextualNode) bool {
	found := false
	for _, m := range p.ContextualNodes {
		found = n.isEqual(*m)
	}
	return found
}

func (p Path) Description() string {
	str := ""
	for _, n := range p.ContextualNodes {
		str += n.Description() + " -> "
	}
	return str
}

func (p Path) Display() {
	fmt.Println(p.Description())
}
