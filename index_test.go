package genddl

import (
	"strings"
	"testing"
)

func TestJoinAndStripName(t *testing.T) {
	safeSize := []string{
		"aaaa",
		"bbbb",
		"cccc",
		"dddd",
		"eeee",
		"ffff",
		"gg",
	}
	ss1 := joinAndStripName(strings.Join(safeSize, "_"))
	if ss1 != strings.Join(safeSize, "_") {
		t.Errorf("return string is invalid: %s", ss1)
	}

	overSize := append(safeSize, "h")
	expectedIndexName := "aaaa_bbbb_cccc_dddd_eeeee98dee81"
	ss2 := joinAndStripName(strings.Join(overSize, "_"))
	if ss2 != expectedIndexName {
		t.Errorf("return string is invalid: %s", ss2)
	}
}
