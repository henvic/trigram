package trigram

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestStoreNotFound(t *testing.T) {
	s := Store{}

	if _, err := s.Generate(); err == nil {
		t.Errorf("Expected store not to be found, got %v instead", err)
	}
}

func TestStoreTooShort(t *testing.T) {
	s := NewStore()

	if _, err := s.Generate(); err != ErrTooShort {
		t.Errorf("Expected error to be %v, got %v instead", ErrTooShort, err)
	}
}

func TestStore(t *testing.T) {
	s := NewStore()

	if err := s.Learn(context.Background(), bytes.NewBuffer(book)); err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	for n := 0; n < 20; n++ {
		testGenerate(t, s, book)
	}
}

func TestStoreConcurrency(t *testing.T) {
	s := NewStore()

	var w sync.WaitGroup
	var repeat = 30

	w.Add(repeat * len(sophisms))

	for n := 0; n < repeat; n++ {
		for _, b := range sophisms {
			go testSophism(t, s, b, w.Done)
		}
	}

	w.Wait()
}

func testSophism(t *testing.T, s *Store, b []byte, done func()) {
	defer done()

	if err := s.Learn(context.Background(), bytes.NewBuffer(b)); err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	testGenerate(t, s, b)
}

func testGenerate(t *testing.T, s *Store, b []byte) {
	var got, err = s.Generate()

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if len(got) < 30 {
		t.Errorf("Sentence %s appears to be smaller than expected", got)
	}

	// TODO(henvic): improve test by checking if trigram pattern is contained on input.
	// Denormalize before checking to avoid false negatives.
}

func ExampleStore() {
	s := NewStore()

	if err := s.Learn(context.Background(), bytes.NewBuffer(book)); err != nil {
		panic(err)
	}

	sentences, err := s.Generate()

	if err != nil {
		panic(err) // on the real world, you really want to treat your errors instead of panicking.
	}

	fmt.Println(sentences)
}

type normCase struct {
	in   string
	want string
}

var normCases = []normCase{
	normCase{
		in:   "a little dog,",
		want: "A little dog.",
	},
	normCase{
		in:   "a little dog",
		want: "A little dog.",
	},
	normCase{
		in:   "a little dog is 10",
		want: "A little dog is 10.",
	},
	normCase{
		in:   "A little dog.",
		want: "A little dog.",
	},
	normCase{
		in:   "a little dog!",
		want: "A little dog!",
	},
}

func TestNormalize(t *testing.T) {
	for _, c := range normCases {
		var got = normalize(c.in)

		if c.want != got {
			t.Errorf("Expected %v, got %v instead", c.want, got)
		}
	}
}
