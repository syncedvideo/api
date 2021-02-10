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
	pascalCase := regexp.MustCompile(`[A-Z][^A-Z]*`)
	if !pascalCase.MatchString(name) {
		t.Errorf("name is not in PascalCase: %s", name)
	}
}
