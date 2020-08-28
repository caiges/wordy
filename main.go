package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var debug = flag.Bool("debug", false, "Print debug messages")
var grouping = flag.Int("grouping", 3, "Number of words per group")
var top = flag.Int("top", 100, "Number of groupings to return")

// TokenGroup contains a list of string tokens used to collects sequences.
type TokenGroup struct {
	tokens []string
}

// Accumulator manages state and data when accumulating groups of tokens.
type Accumulator struct {
	groupings   map[string]int
	groupSize   int
	groupIndex  int
	tokenGroups []TokenGroup
}

// NewAccumulator creates a new Accumulator.
func NewAccumulator(g int) Accumulator {
	a := Accumulator{groupSize: g, groupIndex: 0, groupings: map[string]int{}, tokenGroups: []TokenGroup{}}

	for i := 0; i < g; i++ {
		a.tokenGroups = append(a.tokenGroups, TokenGroup{})
	}

	return a
}

// Add adds a string token to the collection to be reported.
func (a *Accumulator) Add(s string) {
	// Append to other token groups.
	for i, t := range a.tokenGroups {
		// Add the current token to all other non-empty groups.
		if len(t.tokens) > 0 {
			a.tokenGroups[i].tokens = append(t.tokens, s)
		}

		// If the token group contains the target groupSize, build the combined
		// string and add it to the collection.
		if len(a.tokenGroups[i].tokens) == a.groupSize {
			grouping := strings.Join(a.tokenGroups[i].tokens, " ")

			// Add new word tokens if necessary.
			if _, ok := a.groupings[grouping]; !ok {
				a.groupings[grouping] = 0
			}

			// Add the current grouping to the existing count.
			a.groupings[grouping]++

			// Clear current token group since we've accumulated our target groupSize.
			a.tokenGroups[i] = TokenGroup{}
		}
	}

	// Add token to the current token group.
	a.tokenGroups[a.groupIndex].tokens = append(a.tokenGroups[a.groupIndex].tokens, s)

	if *debug {
		fmt.Printf("adding word '%s' on token group: %d >> %+v\n", s, a.groupIndex, a.tokenGroups)
	}

	a.groupIndex++
	// Reset group index if we should continue the next combination of groupings.
	if a.groupIndex == a.groupSize {
		a.groupIndex = 0
	}
}

// Groupings is a collection of word groupings. This is used mainly for sorting.
type Groupings []Grouping

// Grouping represents a single word grouping and its count in the source data.
type Grouping struct {
	Grouping string
	Count    int
}

// Len returns the number of total groupings.
func (g Groupings) Len() int {
	return len(g)
}

// Less returns whether the grouping at index i is less than at index j.
func (g Groupings) Less(i, j int) bool {
	return g[i].Count < g[j].Count
}

// Swap swaps the grouping at index i with the grouping at index j.
func (g Groupings) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

// GroupingsReport returns a nice sorted list of groupings.
func GroupingsReport(a Accumulator) Groupings {
	groupings := make(Groupings, len(a.groupings))
	i := 0

	for k, v := range a.groupings {
		groupings[i] = Grouping{k, v}
		i++
	}

	sort.Sort(sort.Reverse(groupings))

	return groupings
}

// ScanWords is a split function for a Scanner that returns each
// space-separated word of text, with surrounding spaces deleted. It will
// never return an empty string. The definition of space is set by
// unicode.IsSpace.
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	// Count punctuation width.
	punctWidth := 0

	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])

		if !unicode.IsSpace(r) {
			break
		}
	}

	// Scan until space, marking end of word.
	tData := []byte{}

	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		if unicode.IsSpace(r) {
			return i + width, tData, nil
		}

		// We care about punctuation width even if we don't want it in our data.
		if unicode.IsPunct(r) {
			punctWidth += width
		}

		// Only allow letters in our word.
		if unicode.IsLetter(r) {
			tData = append(tData, data[i:i+width]...)
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(tData) > start {
		return len(tData) + punctWidth, tData, nil
	}

	// Request more data.
	return start, nil, nil
}

// help displays usage information.
func help() {
	fmt.Println("Usage: wordy [-debug] [FILE]")
	fmt.Println("wordy accepts from stdin as well.")
}

func main() {
	flag.Parse()

	// Parse arguments and prepare source.
	args := flag.Args()

	var r io.Reader

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		if len(args) > 0 {
			help()
			os.Exit(1)
		}

		r = bufio.NewReader(os.Stdin)
	}

	if len(args) > 0 {
		i, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Could not open file: %s\n", args[0])
			os.Exit(1)
		}

		r = i
	}

	if r == nil {
		fmt.Println("No input provided")
		os.Exit(1)
	}

	// Initialize accumulator and start scanning source.
	a := NewAccumulator(*grouping)

	scanner := bufio.NewScanner(r)
	scanner.Split(ScanWords)

	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if word != "" {
			a.Add(word)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	groupings := GroupingsReport(a)
	fetch := *top

	if len(groupings) < *top {
		fetch = len(groupings)
	}

	for _, g := range groupings[0:fetch] {
		fmt.Printf("%d - %s\n", g.Count, g.Grouping)
	}
}
