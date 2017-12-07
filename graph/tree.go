package graph

import "strconv"

type Tree struct {
	Roots []*Vertice
	Name  string
}

type Vertice struct {
	Content  *Node
	Children []*Vertice
	Ancestor *Vertice
	IsRoot   bool
}

type depthTraversing func(vertices []*Vertice)

func (t *Tree) InsertPath(p Path) {
	var v *Vertice
	for i, n := range p.Nodes {
		if i == 0 {
			v = t.FindOrCreateRoot(n)
		} else {
			newNode := v.FindOrCreateChild(n)
			v = newNode
		}
	}
}

func (t *Tree) FindOrCreateRoot(n *Node) *Vertice {
	for _, r := range t.Roots {
		if n.isEqual(*r.Content) {
			return r
		}
	}
	r := createVertice(n)
	r.IsRoot = true
	r.Ancestor = nil
	t.Roots = append(t.Roots, r)
	return r
}

func createVertice(n *Node) *Vertice {
	childs := make([]*Vertice, 0)
	v := Vertice{n, childs, nil, false}
	return &v
}

func (v *Vertice) RightSibling() *Vertice {
	if v.IsRoot == false {
		siblings := v.Ancestor.Children
		for i, s := range siblings {
			if v.Content.isEqual(*s.Content) {
				if i+1 < len(siblings) {
					return siblings[i+1]
				}
			}
		}
	}
	return nil
}

func (v *Vertice) LeftSibling() *Vertice {
	if v.IsRoot == false {
		siblings := v.Ancestor.Children
		for i, s := range siblings {
			if v.Content.isEqual(*s.Content) {
				if i+1 < len(siblings) {
					return siblings[i+1]
				}
			}
		}
	}
	return nil
}

func (v *Vertice) FindOrCreateChild(n *Node) *Vertice {
	for _, c := range v.Children {
		if n.isEqual(*c.Content) {
			return c
		}
	}
	c := createVertice(n)
	c.IsRoot = false
	c.Ancestor = v
	v.Children = append(v.Children, c)
	return c
}

func (t Tree) DepthTraversing(fn depthTraversing) {
	for _, r := range t.Roots {
		ancestors := []*Vertice{r}
		if r.IsLeaf() {
			fn(ancestors)
		} else {
			r.depthTraversing(ancestors, fn)
		}
	}
}

func (v Vertice) IsLeaf() bool {
	return len(v.Children) == 0
}

func (v Vertice) depthTraversing(ancestors []*Vertice, fn depthTraversing) {
	for _, c := range v.Children {
		vertices := append(ancestors, c)
		if c.IsLeaf() {
			fn(vertices)
		} else {
			c.depthTraversing(vertices, fn)
		}
	}
}

func (t Tree) Description() string {
	str := strconv.Itoa(len(t.Roots))
	return str
}
