package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Prefix is a Markov chain prefix of one or more words
type Prefix []string

// Chain maps prefixes to a list of suffixes
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// String returns the Prefix as a string
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends a new word
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// NewChain returns a chain of prefixLen prefixes
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided reader (impl. interface)
// parses it into the prefix and suffixes that are stored in Chain
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	prefix := make(Prefix, c.prefixLen)
	for {
		// no while loops, only for constructs can loop
		var s string
		// Fscan takes till a space or so
		if _, err := fmt.Fscan(br, &s); err != nil {
			break // break out on error from scanf
		}
		key := prefix.String()
		c.chain[key] = append(c.chain[key], s)
		prefix.Shift(s) // "" "I" "I am" "am a" "a good"
	}
}

// Generate returns a string of at most n words generated from Chain
func (c *Chain) Generate(n int) string {
	prefix := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[prefix.String()]
		if len(choices) == 0 {
			break
		}

		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		prefix.Shift(next) // Go to next prefix in Chain
	}

	return strings.Join(words, " ")
}

func main() {
	// Register cmd flags
	numWords := flag.Int("words", 100, "maximum number of words to print")
	prefixLen := flag.Int("prefix", 2, "prefix length in words")

	flag.Parse()
	fmt.Println("We parsed successfully")
	rand.Seed(time.Now().UnixNano())

	c := NewChain(*prefixLen)
	c.Build(os.Stdin)
	text := c.Generate(*numWords)
	fmt.Println(text)

}
