package graph

import (
	"gotrading/core"
	"strconv"
)

type Tree struct {
	Roots []*Vertice
	Name  string
}

type Vertice struct {
	Hit      *core.Hit
	Children []*Vertice
	Ancestor *Vertice
	IsRoot   bool
}

type depthTraversing func(hits []*core.Hit)

func (t *Tree) InsertPath(p Path) {
	var v *Vertice
	for i, n := range p.Hits {
		if i == 0 {
			v = t.FindOrCreateRoot(n)
		} else {
			newNode := v.FindOrCreateChild(n)
			v = newNode
		}
	}
}

func (t *Tree) FindOrCreateRoot(h *core.Hit) *Vertice {
	for _, r := range t.Roots {
		if h.IsEqual(*r.Hit) {
			return r
		}
	}
	r := createVertice(h)
	r.IsRoot = true
	r.Ancestor = nil
	t.Roots = append(t.Roots, r)
	return r
}

func createVertice(h *core.Hit) *Vertice {
	childs := make([]*Vertice, 0)
	v := Vertice{h, childs, nil, false}
	return &v
}

func (v *Vertice) RightSibling() *Vertice {
	if v.IsRoot == false {
		siblings := v.Ancestor.Children
		for i, s := range siblings {
			if v.Hit.IsEqual(*s.Hit) {
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
			if v.Hit.IsEqual(*s.Hit) {
				if i+1 < len(siblings) {
					return siblings[i+1]
				}
			}
		}
	}
	return nil
}

func (v *Vertice) FindOrCreateChild(h *core.Hit) *Vertice {
	for _, c := range v.Children {
		if h.IsEqual(*c.Hit) {
			return c
		}
	}
	c := createVertice(h)
	c.IsRoot = false
	c.Ancestor = v
	v.Children = append(v.Children, c)
	return c
}

func (t Tree) DepthTraversing(fn depthTraversing) {
	for _, r := range t.Roots {
		ancestors := []*Vertice{r}
		hits := []*core.Hit{r.Hit}
		if r.IsLeaf() {
			fn(hits)
		} else {
			r.depthTraversing(ancestors, hits, fn)
		}
	}
}

func (v Vertice) IsLeaf() bool {
	return len(v.Children) == 0
}

func (v Vertice) depthTraversing(ancestors []*Vertice, hits []*core.Hit, fn depthTraversing) {
	for _, c := range v.Children {
		vertices := append(ancestors, c)
		contents := append(hits, c.Hit)
		if c.IsLeaf() {
			fn(contents)
		} else {
			c.depthTraversing(vertices, contents, fn)
		}
	}
}

func (t Tree) Description() string {
	str := strconv.Itoa(len(t.Roots))
	return str
}
