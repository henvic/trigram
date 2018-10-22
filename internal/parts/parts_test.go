package parts

import (
	"testing"
)

func TestPartTooBig(t *testing.T) {
	var current = &Part{}

	current = current.Add("var")
	current = current.Add("foo")
	current = current.Add("=")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected trigram error, got %v instead", r)
		}
	}()

	current.Add("bar")
}
