package filters

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/stripe/unilog/clevels"
)

var startSystemAusterityLevel sync.Once

func AusterityFilter(line string) string {
	// Start austerity level loop sender in goroutine just once
	startSystemAusterityLevel.Do(func() {
		go clevels.SendSystemAusterityLevel()
	})

	criticalityLevel := clevels.Criticality(line)
	austerityLevel := <-clevels.SystemAusterityLevel
	fmt.Printf("austerity level is %s\n", austerityLevel)

	if criticalityLevel >= austerityLevel {
		return line
	}

	if rand.Float64() > samplingRate(austerityLevel, criticalityLevel) {
		return "(shedded)"
	}
	return line
}

// samplingRate calculates the rate at which loglines will be sampled for the
// given criticality level and austerity level. For example, if the austerity level
// is Critical (3), then lines that are Sheddable (0) will be sampled at .001.
func samplingRate(austerityLevel, criticalityLevel clevels.AusterityLevel) float64 {
	if criticalityLevel > austerityLevel {
		return 1
	}

	levelDiff := austerityLevel - criticalityLevel
	samplingRate := math.Pow(10, float64(-levelDiff))

	return samplingRate
}
