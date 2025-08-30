package trie

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrieInsertAndSearch(t *testing.T) {
	tr := NewTrie()

	words := []string{"hello", "helium", "he", "hero"}
	for _, w := range words {
		tr.Insert(w)
	}

	for _, w := range words {
		if !tr.Search(w) {
			t.Errorf("Search(%q) = false; want true", w)
		}
	}

	nonWords := []string{"hey", "her", ""}
	for _, w := range nonWords {
		if tr.Search(w) {
			t.Errorf("Search(%q) = true; want false", w)
		}
	}
}

func TestTrieStartsWith(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")
	tr.Insert("helium")
	tr.Insert("he")
	tr.Insert("hero")

	tests := []struct {
		prefix   string
		expected bool
	}{
		{"he", true},
		{"hel", true},
		{"hero", true},
		{"her", true},
		{"ha", false},
		{"", false},
	}

	for _, tt := range tests {
		got := tr.StartsWith(tt.prefix)
		if got != tt.expected {
			t.Errorf("StartsWith(%q) = %v; want %v", tt.prefix, got, tt.expected)
		}
	}
}

func TestTrieGetWordsWithPrefix(t *testing.T) {
	tr := NewTrie()
	words := []string{"he", "hello", "helium", "hero"}
	for _, w := range words {
		tr.Insert(w)
	}

	prefix := "he"
	expected := []string{"he", "hello", "helium", "hero"}
	got := tr.GetWordsWithPrefix(prefix)

	sort.Strings(expected)
	sort.Strings(got)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("GetWordsWithPrefix(%q) = %v; want %v", prefix, got, expected)
	}

	nonPrefix := "ha"
	got = tr.GetWordsWithPrefix(nonPrefix)
	if len(got) != 0 {
		t.Errorf("GetWordsWithPrefix(%q) = %v; want empty slice", nonPrefix, got)
	}
}

func TestTrieRemove(t *testing.T) {
	tr := NewTrie()
	tr.Insert("he")
	tr.Insert("hello")
	tr.Insert("helium")
	tr.Insert("hero")

	if !tr.Remove("he") {
		t.Errorf("Remove('he') = false; want true")
	}
	if tr.Search("he") {
		t.Errorf("'he' should be removed")
	}

	if !tr.Remove("hello") {
		t.Errorf("Remove('hello') = false; want true")
	}
	if tr.Search("hello") {
		t.Errorf("'hello' should be removed")
	}
	if !tr.Search("helium") {
		t.Errorf("'helium' should still exist")
	}

	if tr.Remove("unknown") {
		t.Errorf("Remove('unknown') = true; want false")
	}
}

func TestTrieSizeAndIsEmpty(t *testing.T) {
	tr := NewTrie()
	if !tr.IsEmpty() {
		t.Errorf("expected trie to be empty")
	}
	if tr.Size() != 0 {
		t.Errorf("expected size 0, got %d", tr.Size())
	}

	tr.Insert("hello")
	if tr.IsEmpty() {
		t.Errorf("expected trie not to be empty")
	}
	if tr.Size() != 1 {
		t.Errorf("expected size 1, got %d", tr.Size())
	}

	tr.Insert("hello")
	if tr.Size() != 1 {
		t.Errorf("expected size 1, got %d", tr.Size())
	}

	tr.Remove("hello")
	if !tr.IsEmpty() {
		t.Errorf("expected trie to be empty after removal")
	}
	if tr.Size() != 0 {
		t.Errorf("expected size 0 after removal, got %d", tr.Size())
	}
}

func TestTrieEmptyString(t *testing.T) {
	tr := NewTrie()
	tr.Insert("")
	if tr.Size() != 0 {
		t.Errorf("Expected size 0, got %v\n", tr.Size())
	}
	result := tr.GetWordsWithPrefix("")
	if result != nil {
		t.Errorf("Expected empty slice, got %v\n", len(result))
	}
	f := tr.Remove("")
	if f {
		t.Errorf("Expected %v, got %v\n", false, f)
	}
}
