package sillyname

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	name := New()
	if name == "" {
		t.Errorf("name is empty string")
	}
	pattern := regexp.MustCompile(`[A-Z][^A-Z]*`)
	if len(pattern.FindAllString(name, -1)) != 2 {
		t.Errorf("name does not match pattern: %s", name)
	}
}
