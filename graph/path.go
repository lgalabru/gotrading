package graph

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"gotrading/core"
)

type Path struct {
	Hits []*core.Hit `json:"hits"`
	Id   *string     `json:"id"`
	Name *string     `json:"description"`
	USD  float64
}

func (p *Path) Encode() {
	desc := p.Description()
	p.Name = &desc

	h := sha1.New()
	h.Write([]byte(desc))
	enc := hex.EncodeToString(h.Sum(nil))
	p.Id = &enc
}

func (p Path) contains(h core.Hit) bool {
	found := false
	for _, m := range p.Hits {
		found = h.IsEqual(*m)
	}
	return found
}

func (p Path) Description() string {
	str := ""
	for _, n := range p.Hits {
		str += n.Description() + " -> "
	}
	return str
}

func (p Path) Display() {
	fmt.Println(p.Description())
}
