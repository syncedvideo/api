package youtube

import "testing"

func TestExtractVideoID(t *testing.T) {
	case1 := ExtractVideoID("https://www.youtube.com/watch?v=LOTf6EXI-Uw")
	if case1 != "LOTf6EXI-Uw" {
		t.Errorf("failed case 1: %s", case1)
	}

	case2 := ExtractVideoID("LOTf6EXI-Uw")
	if case2 != "LOTf6EXI-Uw" {
		t.Errorf("failed case 2: %s", case2)
	}

	case3 := ExtractVideoID("")
	if case3 != "" {
		t.Errorf("failed case 3: %s", case3)
	}
}
