// particles
package main

import (
	"flag"
	"github.com/wildducktheories/particles"
	"log"
	"os"
)

func main() {

	processA := "particle"
	processB := "quark"
	matchType := "particle"
	size := 10000
	maxMatches := 0
	verbose := false

	flag.StringVar(&processA, "process-A", "quark", "Sampling process A: particle or quark")
	flag.StringVar(&processB, "process-B", "particle", "Sampling process B: particle or quark")
	flag.StringVar(&matchType, "match-type", "particle", "Type of match: particle or quark")
	flag.IntVar(&size, "pool-size", 10000, "Size of the particle pool.")
	flag.IntVar(&maxMatches, "max-matches", 1000, "Maximum number of matches")
	flag.BoolVar(&verbose, "verbose", false, "Be verbose about statistics")
	flag.Parse()

	particles.Verbose = verbose

	pool := particles.NewPool(size)

	observerProcess := func(n string) particles.ObserverProcess {
		switch n {
		case "particle":
			return particles.ProcessFromObserver(pool.SampleParticle)
		case "quark":
			return particles.ProcessFromObserver(pool.SampleQuark)
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}
		panic("unreachable")
	}

	fA := observerProcess(processA)
	fB := observerProcess(processB)
	fC := func() particles.MatcherProcess {
		switch matchType {
		case "particle":
			return particles.ProcessFromMatcher(particles.ParticleMatcher)
		case "quark":
			return particles.ProcessFromMatcher(particles.QuarkMatcher)
		case "a-side":
			return particles.ProcessFromMatcher(particles.ASideMatcher)
		case "b-side":
			return particles.ProcessFromMatcher(particles.BSideMatcher)
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}
		panic("unreachable")
	}()

	done := make(chan struct{})
	ac := make(chan particles.Observation)
	bc := make(chan particles.Observation)
	cd := make(chan particles.Observation)

	go fA(done, ac)
	go fB(done, bc)
	go fC(ac, bc, cd)

	verboseLimit := 10
	growth := 1.2
	total := 0
	totalTrue := 0

	for o := range cd {
		if total >= maxMatches {
			if done != nil {
				close(done)
				done = nil
			}
			continue
		}
		total++
		if o.Confirm() {
			totalTrue++
		}
		if total == verboseLimit {
			if verbose {
				log.Printf("n=%d, total=%d, happy=%d, ratio=%f", size, total, totalTrue, float64(totalTrue)/float64(total))
			}
			verboseLimit = int(float64(verboseLimit) * growth)
		}
	}

	log.Printf("n=%d, total=%d, happy=%d, ratio=%f", size, total, totalTrue, float64(totalTrue)/float64(total))
}
