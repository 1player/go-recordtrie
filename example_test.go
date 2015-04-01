package recordtrie

import (
	"fmt"
	"strings"
)

func Example() {

	// Input format is key;value1;value2;...
	input := `
		foods;pizza;pie
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

	// Retrieve all fruits
	fmt.Printf("Fruits: %v\n", trie.Find("fruits"))

	// Retrieve all keys starting with "f"
	fmt.Printf("Keys starting with f: %v\n", trie.KeysStartingWith("f"))

	// Output:
	// Fruits: [peach pear apple]
	// Keys starting with f: [fruits fruits fruits foods foods]
}
