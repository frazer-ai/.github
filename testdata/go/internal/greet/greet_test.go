package greet

import (
	"strings"
	"testing"
)

func TestHello(t *testing.T) {
	t.Parallel()
	got := Hello("Frazer")
	if !strings.HasPrefix(got, "hello Frazer ") {
		t.Fatalf("Hello() = %q", got)
	}
}
