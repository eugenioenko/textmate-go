package textmate

import "testing"

func TestScopeMatchers(t *testing.T) {
	tests := []struct {
		expression string
		input      []string
		want       bool
	}{
		{"foo", []string{"foo"}, true},
		{"foo", []string{"bar"}, false},
		{"- foo", []string{"foo"}, false},
		{"- foo", []string{"bar"}, true},
		{"- - foo", []string{"bar"}, false},
		{"bar foo", []string{"foo"}, false},
		{"bar foo", []string{"bar"}, false},
		{"bar foo", []string{"bar", "foo"}, true},
		{"bar - foo", []string{"bar"}, true},
		{"bar - foo", []string{"foo", "bar"}, false},
		{"bar - foo", []string{"foo"}, false},
		{"bar, foo", []string{"foo"}, true},
		{"bar, foo", []string{"bar"}, true},
		{"bar, foo", []string{"bar", "foo"}, true},
		{"bar, -foo", []string{"bar", "foo"}, true},
		{"bar, -foo", []string{"yo"}, true},
		{"bar, -foo", []string{"foo"}, false},
		{"(foo)", []string{"foo"}, true},
		{"(foo - bar)", []string{"foo"}, true},
		{"(foo - bar)", []string{"foo", "bar"}, false},
		{"foo bar - (yo man)", []string{"foo", "bar"}, true},
		{"foo bar - (yo man)", []string{"foo", "bar", "yo"}, true},
		{"foo bar - (yo man)", []string{"foo", "bar", "yo", "man"}, false},
		{"foo bar - (yo | man)", []string{"foo", "bar", "yo", "man"}, false},
		{"foo bar - (yo | man)", []string{"foo", "bar", "yo"}, false},
		{"R:text.html - (comment.block, text.html source)", []string{"text.html", "bar", "source"}, false},
		{"text.html.php - (meta.embedded | meta.tag), L:text.html.php meta.tag, L:source.js.embedded.html", []string{"text.html.php", "bar", "source.js"}, true},
	}

	for index, test := range tests {
		matched := false
		for _, candidate := range createScopeMatchers(test.expression) {
			if candidate.matcher(test.input) {
				matched = true
				break
			}
		}
		if matched != test.want {
			t.Errorf("case %d createScopeMatchers(%q) = %v, want %v", index, test.expression, matched, test.want)
		}
	}
}

func TestScopeMatcherPriorityAndPrefix(t *testing.T) {
	matchers := createScopeMatchers("L:source.js, R:text.html, comment")
	if len(matchers) != 3 || matchers[0].priority != -1 || matchers[1].priority != 1 || matchers[2].priority != 0 {
		t.Fatalf("priorities = %#v", matchers)
	}
	if !matchers[0].matcher([]string{"source.js.embedded.html"}) {
		t.Fatal("dotted scope prefix did not match")
	}
	if matchers[0].matcher([]string{"source.jsx"}) {
		t.Fatal("partial component matched")
	}
}
