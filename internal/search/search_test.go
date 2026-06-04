package search

import "testing"

func TestNormalize(t *testing.T) {
	got := Normalize("  ЗачЁт   по   ОММ ")
	if got != "зачет по омм" {
		t.Fatalf("got %q", got)
	}
}

func TestExpand(t *testing.T) {
	items := Expand("зач")
	m := map[string]bool{}
	for _, v := range items {
		m[v] = true
	}
	if !m["зачет"] || !m["зач"] {
		t.Fatalf("unexpected %v", items)
	}
}
