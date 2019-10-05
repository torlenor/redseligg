package httppingplugin

import "testing"

func TestCreateMatrixBot(t *testing.T) {
	result := stripCmd("!CMD test", "CMD")
	if result != "test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = stripCmd("test !CMD", "CMD")
	if result != "test !CMD" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = stripCmd("!CMDtest", "CMD")
	if result != "!CMDtest" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = stripCmd("!TEST test", "CMD")
	if result != "!TEST test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = stripCmd("!CMD2 test", "CMD")
	if result != "!CMD2 test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}
}
