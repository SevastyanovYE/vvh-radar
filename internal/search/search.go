package search

import "strings"

var synonyms = map[string][]string{
	"зачет":            {"зачет", "зач", "зачёт"},
	"экзамен":          {"экзамен", "экз"},
	"пересдача":        {"пересдача", "хвост"},
	"билеты":           {"билеты", "билет"},
	"программирование": {"программирование", "прога", "matlab", "матлаб"},
}

func Normalize(s string) string {
	s = strings.ToLower(strings.ReplaceAll(s, "ё", "е"))
	return strings.Join(strings.Fields(s), " ")
}

func Expand(query string) []string {
	n := Normalize(query)
	set := map[string]struct{}{}
	for _, tok := range strings.Fields(n) {
		set[tok] = struct{}{}
		for _, group := range synonyms {
			match := false
			for _, item := range group {
				if item == tok {
					match = true
					break
				}
			}
			if match {
				for _, item := range group {
					set[Normalize(item)] = struct{}{}
				}
			}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	return out
}
