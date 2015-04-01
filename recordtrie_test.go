package recordtrie

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestKeyEncoding(t *testing.T) {
	tests := []struct {
		trieKey string
		key     string
		value   string
	}{
		{"test", "test", ""},
		{"test\xFF", "test", ""},
		{"foo\xFFbar", "foo", "bar"},
		{"foo\xFFbar\xFF", "foo", "bar\xFF"},
	}

	for _, test := range tests {
		k, v := splitTrieKey(test.trieKey)
		if k != test.key || v != test.value {
			t.Errorf("RecordTrie.splitTrieKey(%q): expected (%q, %q), got (%q, %q)\n",
				test.trieKey,
				test.key, test.value,
				k, v,
			)
		}
	}
}

func TestTrieExists(t *testing.T) {
	r := New([]Record{
		{"foo", "bar"},
		{"foobar", "baz"},
	})

	if exists := r.Exists("foo"); !exists {
		t.Error("RecordTrie.Exists(\"foo\"): expected true, got false")
	}

	if exists := r.Exists("key"); exists {
		t.Error("RecordTrie.Exists(\"key\"): expected false, got true")
	}
}

func TestTrieFind(t *testing.T) {
	records := []Record{
		{"foo", "bar"},
		{"abc", "def"},
		{"foo", "baz"},
	}

	r := New(records)

	tests := []struct {
		key    string
		values []string
	}{
		{"foo", []string{"bar", "baz"}},
		{"abc", []string{"def"}},
		{"def", []string{}},
	}

	for _, test := range tests {
		v := r.Find(test.key)
		if !compareStringSlices(v, test.values) {
			t.Errorf("RecordTrie.Find(%q): got %v expected %v\n", test.key,
				v, test.values)
		}
	}
}

func TestTrieKeysStartingWith(t *testing.T) {
	records := []Record{
		{"foo", "bar"},
		{"abc", "def"},
		{"a", "apple"},
		{"ac", "acorn"},
		{"foo", "baz"},
	}

	r := New(records)

	tests := []struct {
		query string
		keys  []string
	}{
		{"foo", []string{"foo", "foo"}},
		{"a", []string{"a", "ac", "abc"}},
		{"ab", []string{"abc"}},
		{"def", []string{}},
	}

	for _, test := range tests {
		keys := r.KeysStartingWith(test.query)
		if !compareStringSlices(keys, test.keys) {
			t.Errorf("RecordTrie.KeysStartingWith(%q): got %v expected %v\n", test.query,
				keys, test.keys)
		}
	}
}

func TestTrieLoadSave(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "testTrie")
	if err != nil {
		t.Fatal(err)
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)

	r := New([]Record{
		{"abc", "def"},
	})
	err = r.Save(tmpFilePath)
	if err != nil {
		t.Fatal("RecordTrie.Save", err)
	}

	r, err = NewFromFile(tmpFilePath)
	if err != nil {
		t.Fatal("RecordTrie.NewFromFile", err)
	}
	v := r.Find("abc")
	if len(v) != 1 || v[0] != "def" {
		t.Errorf("RecordTrie.Save(): unexpected data")
	}
}

func ExampleParse() {

	// Input format is key;value1;value2;...
	input := `
		fruits;apple;pear;peach
		animals;dog;cat;horse;fish
	`

	var records []Record

	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, ";")

		key := parts[0]
		values := parts[1:]

		// Build the flat list of records
		for _, value := range values {
			records = append(records, Record{key, value})
		}
	}

	// Build the trie
	trie := New(records)

	fmt.Printf("%v\n", trie.Find("fruits"))
	// [apple pear peach]

	fmt.Printf("%v\n", trie.Find("foo"))
	// []
}
