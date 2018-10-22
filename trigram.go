// Package trigram can be used to create trigrams from books, articles, and more.
package trigram

import (
	"bufio"
	"context"
	"errors"
	"io"
	"math/rand"
	"strings"
	"sync"
	"unicode"

	"github.com/henvic/trigram/internal/parts"
)

// Store for trigrams.
type Store struct {
	p *parts.Part
	m sync.RWMutex
}

// NewStore for the trigrams.
func NewStore() *Store {
	return &Store{
		p: &parts.Part{},
	}
}

// ErrStoreNotFound is returned when the trigram database is not found.
var ErrStoreNotFound = errors.New("trigram store not found")

// ErrTooShort is returned when text has less than three words.
var ErrTooShort = errors.New("text is too short to extract trigram")

// Learn text.
func (s *Store) Learn(ctx context.Context, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	var f, m, l string
	var v = parts.Root

	for scanner.Scan() {
		if v != parts.Third {
			v++
		}

		if f, m, l = m, l, scanner.Text(); v != parts.Third {
			continue
		}

		s.add(f, m, l)
		// log.Println("Adding:", f, m, l)
	}

	return scanner.Err()
}

func (s *Store) add(f, m, l string) {
	var root = s.p
	s.m.Lock()
	defer s.m.Unlock()

	current := root.Add(f)
	current = current.Add(m)
	current.Add(l)
}

// Generate trigrams.
func (s *Store) Generate() (string, error) {
	if s.p == nil {
		return "", ErrStoreNotFound
	}

	var sa []string

	// Notice average sentences close to 14 words are better to understand.
	// Reference: http://prsay.prsa.org/2009/01/14/how-to-make-your-copy-more-readable-make-sentences-shorter/
	for left := normal(28, 5); left > 0; left-- {
		st, err := s.trigram()

		if err != nil {
			return "", err
		}

		sa = append(sa, st)
	}

	return normalize(strings.Join(sa, " ")), nil
}

func (s *Store) trigram() (string, error) {
	var root = s.p
	s.m.RLock()
	defer s.m.RUnlock()

	var current = root

	var sa []string

	for c := 0; c <= 2; c++ {
		current = current.RandChild()

		if current == nil {
			return "", ErrTooShort
		}

		sa = append(sa, current.Word)
	}

	return strings.Join(sa, " "), nil
}

func normalize(s string) string {
	nr := []rune(s)

	if !unicode.IsUpper(nr[0]) {
		nr[0] = unicode.ToUpper(nr[0])
	}

	switch last := nr[len(nr)-1]; {
	case unicode.IsNumber(last) || unicode.IsLetter(last):
		nr = append(nr, '.')
	case last == ',':
		nr[len(nr)-1] = '.'
	}

	return string(nr)
}

func normal(desiredMean, desiredStdDev int) int {
	var sample = rand.NormFloat64()*float64(desiredStdDev) + float64(desiredMean)
	return int(sample)
}
