package sillyname

import (
	"regexp"
	"testing"
)

func TestSillyname(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		name := New()
		if name == "" {
			t.Errorf("name is empty")
		}
		pascalCase := regexp.MustCompile(`[A-Z][^A-Z]*`)
		if !pascalCase.MatchString(name) {
			t.Errorf("name is not in PascalCase: %s", name)
		}
	})
}
