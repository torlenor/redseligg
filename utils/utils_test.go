package utils

import "testing"

func TestCreateMatrixBot(t *testing.T) {
	result := StripCmd("!CMD test", "CMD")
	if result != "test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = StripCmd("test !CMD", "CMD")
	if result != "test !CMD" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = StripCmd("!CMDtest", "CMD")
	if result != "!CMDtest" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = StripCmd("!TEST test", "CMD")
	if result != "!TEST test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}

	result = StripCmd("!CMD2 test", "CMD")
	if result != "!CMD2 test" {
		t.Fatalf("Stripping resulted in wrong string: %s", result)
	}
}
