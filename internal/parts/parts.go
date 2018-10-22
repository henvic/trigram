package parts

import (
	"math/rand"
)

// Level is the position of the word.
type Level int

const (
	// Root element
	Root = iota

	// First word
	First

	// Second word
	Second

	// Third word
	Third
)

// Part of the trigram.
type Part struct {
	Word string

	appears int
	level   Level

	children map[string]*Part
}

// Add child part.
func (p *Part) Add(w string) *Part {
	if p.level >= Third {
		panic("the depth of a trigram is exactly 3")
	}

	child, ok := p.children[w]

	if ok {
		child.appears++
		return child
	}

	child = &Part{
		Word: w,

		appears: 1,
		level:   p.level + 1,
	}

	if p.children == nil {
		p.children = map[string]*Part{}
	}

	p.children[w] = child
	return child
}

// RandChild returns a random (weighted) child.
func (p *Part) RandChild() *Part {
	var total = 1

	for _, c := range p.children {
		total += c.appears
	}

	r := rand.Intn(total)

	for _, c := range p.children {
		r -= total

		if r <= 0 {
			return c
		}
	}

	return nil
}
