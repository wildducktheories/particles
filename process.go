package particles

import (
	"fmt"
	"log"
	"sort"
)

// A process characterises a sampling or matching process.
//
// There are two primordial sampling processes. Matching observations from
// two existing processes produces a third process.
//
// The probability of a process is the probability that an observation
// generated by the process will return true if its Confirm() method is called.
//
type process struct {
	name        string   // the cache key for the process
	tokens      []string // the set of tokens representing types of observations generated by the process
	probability float64  // the probability that an observation describes a 'happy' particle
}

var (
	// particleProcess is the process used to sample particles from a pool.
	particleProcess = &process{
		name:        "particle",
		tokens:      []string{"U1u2", "U3d", "dU4", "u1U2", "U3d", "dU4"},
		probability: 1 / 3.0,
	}

	// quarkProcess is the process used to sample quarks from a pool.
	quarkProcess = &process{
		name:        "quark",
		tokens:      []string{"U1u2", "u1U2", "U3d", "dU4"},
		probability: 1 / 2.0,
	}
)

// cache is used to avoid repeating work done by matchingProcess previously
var cache map[string]*process = make(map[string]*process)

// matchingProcess is given two input processes and a matching rule to
// create a new process.
//
// The tokens from the two input process are matched according to the matching
// rule and each token that matches becomes a token in the output process.
//
// The probability of the output process generating a 'happy' particle
// is calculated from the relative frequency of tokens representing 'happy'
// particles in the tokens of the output process.
func matchingProcess(pA *process, pB *process, quarkMatch bool) *process {
	mt := "particle"
	if quarkMatch {
		mt = "quark"
	}
	n := fmt.Sprintf("match(%s, %s, %s)", pA.name, pB.name, mt)

	if p, ok := cache[n]; ok {
		return p
	}

	tokens := []string{}
	c := 0
	uu := 0
	for _, tA := range pA.tokens {
		for _, tB := range pB.tokens {
			if len(tB) != len(tA) { // can't match tokens with different length
				continue
			}
			if len(tA) == 3 && tA != tB { // tokens of length 3 must match exactly
				continue
			}
			if len(tA) == 4 && quarkMatch && tA != tB {
				// tokens of length 4 must match if we are doing a quark match
				continue
			}
			c++
			if len(tA) == 4 { // a token representing two up quarks
				uu++
			}
			tokens = append(tokens, tA)
		}
	}
	p := &process{
		name:        n,
		tokens:      reduce(tokens),
		probability: float64(uu) / float64(c),
	}
	if Verbose {
		log.Printf("name=%s, tokens=%v, probability=%f", p.name, p.tokens, p.probability)
	}
	cache[n] = p
	return p
}

// reduce shrinks an input slice that contains an even number of copies
// of each token by eliminating the duplicate copies from the output slice
func reduce(in []string) []string {
	if len(in)&1 != 0 {
		return in
	}
	out := make([]string, len(in)/2, len(in)/2)
	sort.Sort(sort.StringSlice(in))
	j := 0
	for i, _ := range out {
		if in[j] != in[j+1] {
			return in
		}
		out[i] = in[j]
		j += 2
	}
	return out
}
