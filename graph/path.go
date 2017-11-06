package graph

import "fmt"

type Path struct {
	Nodes []Node
}

func (p Path) contains(n Node) bool {
	found := false
	for _, m := range p.Nodes {
		found = n.isEqual(m)
	}
	return found
}

func (p Path) Description() string {
	str := ""
	for _, n := range p.Nodes {
		str += n.Description() + " /"
	}
	return str
}

func (p Path) Display() {
	fmt.Println(p.Description())
}
