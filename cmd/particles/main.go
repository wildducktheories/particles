// particles
package main

import (
	"flag"
	"github.com/wildducktheories/particles"
	"log"
	"os"
)

func main() {

	size := 10000
	matchType := "particle"
	verbose := false
	maxMatches := 0
	processA := "particle"
	processB := "quark"

	flag.IntVar(&size, "pool-size", 10000, "Size of the particle pool.")
	flag.StringVar(&matchType, "match-type", "particle", "Type of match: particle or quark")
	flag.BoolVar(&verbose, "verbose", false, "Be verbose about statistics")
	flag.IntVar(&maxMatches, "max-matches", 1000, "Maximum number of matches")
	flag.StringVar(&processA, "process-A", "quark", "Sampling process A: particle or quark")
	flag.StringVar(&processB, "process-B", "particle", "Sampling process B: particle or quark")
	flag.Parse()

	particles.Verbose = verbose

	pool := particles.NewPool(size)

	var fA, fB particles.ObserverProcess
	var fC particles.MatcherProcess

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

	fA = observerProcess(processA)
	fB = observerProcess(processB)

	switch matchType {
	case "particle":
		fC = particles.ProcessFromMatcher(particles.ParticleMatcher)
	case "quark":
		fC = particles.ProcessFromMatcher(particles.QuarkMatcher)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

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
		total++
		if o.Read() {
			totalTrue++
		}
		if total == verboseLimit {
			if verbose {
				log.Printf("n=%d, total=%d, true=%d, ratio=%f", size, total, totalTrue, float64(totalTrue)/float64(total))
			}
			verboseLimit = int(float64(verboseLimit) * growth)
		}
		if total >= maxMatches && done != nil {
			close(done)
			done = nil
		}
	}

	log.Printf("n=%d, total=%d, true=%d, ratio=%f", size, total, totalTrue, float64(totalTrue)/float64(total))
}
