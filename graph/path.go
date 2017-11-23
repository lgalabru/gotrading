package graph

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Path struct {
	ContextualNodes []*ContextualNode `json:"node"`
	Id              *string           `json:"id"`
	Name            *string           `json:"description"`
}

func (p *Path) encode() {
	desc := p.Description()
	p.Name = &desc

	h := sha1.New()
	h.Write([]byte(desc))
	enc := hex.EncodeToString(h.Sum(nil))
	p.Id = &enc
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
