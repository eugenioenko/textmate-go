package textmate

import (
	"reflect"
	"testing"

	"github.com/eugenioenko/textmate-go/oniguruma"
)

func TestReplaceCaptureReferences(t *testing.T) {
	line := []rune("prefix .MiXeD 🙂tail")
	captures := []oniguruma.Capture{
		{Start: 0, End: len(line)},
		{Start: 7, End: 13},
		{Start: 14, End: 15},
		{Start: -1, End: -1},
	}

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "plain reference", value: "entity.$1.example", expected: "entity.MiXeD.example"},
		{name: "downcase", value: "entity.${1:/downcase}.example", expected: "entity.mixed.example"},
		{name: "upcase", value: "entity.${1:/upcase}.example", expected: "entity.MIXED.example"},
		{name: "rune offsets", value: "emoji.$2.example", expected: "emoji.🙂.example"},
		{name: "leading dots stripped", value: "$1", expected: "MiXeD"},
		{name: "unmatched slot", value: "before.$3.after", expected: "before..after"},
		{name: "missing slot remains literal", value: "before.$8.after", expected: "before.$8.after"},
		{name: "all occurrences", value: "$1.${1:/downcase}.${1:/upcase}", expected: "MiXeD.mixed.MIXED"},
		{name: "unsupported command remains literal", value: "${1:/capitalize}", expected: "${1:/capitalize}"},
		{name: "dollar text is unchanged", value: "constant.$name", expected: "constant.$name"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := replaceCaptureReferences(test.value, line, captures); got != test.expected {
				t.Fatalf("replaceCaptureReferences(%q) = %q, want %q", test.value, got, test.expected)
			}
		})
	}
}

func TestHasCaptureReferences(t *testing.T) {
	for _, test := range []struct {
		value string
		want  bool
	}{
		{value: "source.$12", want: true},
		{value: "source.${12:/upcase}", want: true},
		{value: "source.${12:/capitalize}", want: false},
		{value: "source.$name", want: false},
		{value: "source.constant", want: false},
	} {
		if got := hasCaptureReferences(test.value); got != test.want {
			t.Errorf("hasCaptureReferences(%q) = %v, want %v", test.value, got, test.want)
		}
	}
}

func TestCaptureRuneHelpers(t *testing.T) {
	line := []rune("a🙂bc")
	captures := []oniguruma.Capture{
		{Start: 0, End: 4},
		{Start: 1, End: 2},
		{Start: -1, End: -1},
		{Start: 3, End: 2},
		{Start: 3, End: 8},
		{Start: 2, End: 2},
	}

	got := capturedRunes(line, captures)
	want := [][]rune{[]rune("a🙂bc"), []rune("🙂"), nil, nil, nil, {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("capturedRunes() = %#v, want %#v", got, want)
	}

	got[0][0] = 'z'
	if line[0] != 'a' {
		t.Fatal("captureRunes returned an alias of the input line")
	}
}

func TestOrderedRawCaptures(t *testing.T) {
	zero := &RawRule{}
	two := &RawRule{}
	ten := &RawRule{}
	got := orderedRawCaptures(RawCaptures{
		"10":  ten,
		"2":   two,
		"0":   zero,
		"bad": {},
		"-1":  {},
	})

	if len(got) != 11 {
		t.Fatalf("len(orderedRawCaptures()) = %d, want 11", len(got))
	}
	if got[0] != zero || got[2] != two || got[10] != ten {
		t.Fatalf("ordered captures lost numeric placement: %#v", got)
	}
	for _, index := range []int{1, 3, 4, 5, 6, 7, 8, 9} {
		if got[index] != nil {
			t.Errorf("capture slot %d = %#v, want nil", index, got[index])
		}
	}

	if got := orderedRawCaptures(nil); got != nil {
		t.Fatalf("orderedRawCaptures(nil) = %#v, want nil", got)
	}
}
