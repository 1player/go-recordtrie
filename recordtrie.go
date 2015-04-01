// This package implements a read-only trie-based key-value store.
// The records are defined when creating the trie, and it is then possible to
// retrieve all the values matching a certain key, or query the list of keys beginning with a specified prefix.
//
// The package expects the keys to be UTF-8 encoded, or more specifically, they should not contain any
// byte with value 0xFF (as specified by the KV_SEPARATOR constant).
// The values however can contain any byte sequence.
package recordtrie

import (
	"fmt"
	"github.com/1player/go-marisa"
	"strings"
)

type RecordTrie struct {
	t marisa.Trie
}

type Record struct {
	Key   string
	Value string
}

// The character we use to separate the key from the value in the trie.
//
// We're using 0xFF because it is non-valid UTF-8 and we're enforcing
// the keys to be UTF-8 encoded
const KV_SEPARATOR = "\xFF"

// Create a new RecordTrie from a list of Records
func New(records []Record) *RecordTrie {
	r := &RecordTrie{
		t: marisa.NewTrie(),
	}
	r.build(records)

	return r
}

// Create a new RecordTrie from file
func NewFromFile(path string) (r *RecordTrie, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	r = &RecordTrie{
		t: marisa.NewTrie(),
	}
	r.t.Mmap(path)

	return
}

func buildTrieKey(key, value string) string {
	return fmt.Sprintf("%s%s%s", key, KV_SEPARATOR, value)
}

func splitTrieKey(trieKey string) (string, string) {
	pieces := strings.SplitN(trieKey, KV_SEPARATOR, 2)
	if len(pieces) < 2 {
		return pieces[0], ""
	}
	return pieces[0], pieces[1]
}

func (r *RecordTrie) build(records []Record) {
	ks := marisa.NewKeyset()

	for _, record := range records {
		trieKey := buildTrieKey(record.Key, record.Value)
		ks.PushBackString(trieKey)
	}

	r.t.Build(ks)
}

// Given a trie key prefix, call iterFunc for each matching (key, value) found
// Stop iterating if iterFunc returns false
func (r *RecordTrie) iter(query string, iterFunc func(k, v string) bool) {
	a := marisa.NewAgent()
	a.SetQueryString(query)

	for r.t.PredictiveSearch(a) {
		trieKey := a.Key().Str()
		k, v := splitTrieKey(trieKey)

		if !iterFunc(k, v) {
			break
		}
	}
}

// Check whether the key exists in the trie
func (r *RecordTrie) Exists(key string) bool {
	exists := false

	r.iter(buildTrieKey(key, ""), func(k, v string) bool {
		exists = true
		return false
	})

	return exists
}

// Retrieve the values list from the trie, given a key
func (r *RecordTrie) Find(key string) []string {
	var values []string

	r.iter(buildTrieKey(key, ""), func(k, v string) bool {
		values = append(values, v)
		return true
	})

	return values
}

// Returns the list of all keys starting with the specified prefix
func (r *RecordTrie) KeysStartingWith(keyPrefix string) []string {
	var keys []string

	r.iter(keyPrefix, func(k, v string) bool {
		keys = append(keys, k)
		return true
	})

	return keys
}

// Returns the list of all records stored in the trie
func (r *RecordTrie) Records() []Record {
	var records []Record

	r.iter("", func(k, v string) bool {
		records = append(records, Record{k, v})
		return true
	})

	return records
}

// Save the trie to file
func (r *RecordTrie) Save(path string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	r.t.Save(path)
	return
}
