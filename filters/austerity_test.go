package filters

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stripe/unilog/clevels"
)

func TestAusterityFilter(t *testing.T) {
	// Make sure SendSystemAusterityLevel is called before we override
	// the underlying channel below
	AusterityFilter("")

	clevels.SystemAusterityLevel = make(chan clevels.AusterityLevel)

	line := fmt.Sprintf("some random log line! clevel=%s", clevels.SheddablePlus)

	kill := make(chan struct{})

	go func() {
		for {
			select {
			case clevels.SystemAusterityLevel <- clevels.Critical:
			case <-kill:
				return
			}
		}
	}()

	// seed rand deterministically
	rand.Seed(17)

	// count number of lines dropped
	dropped := 0
	var outputtedLine string

	// now sample out the line a bunch!
	for i := 0; i < 10000; i++ {
		outputtedLine = AusterityFilter(line)
		if strings.Contains(outputtedLine, "(shedded)") {
			dropped++
		}
	}

	// this number is deterministic because rand is seeded & deterministic
	// TODO (kiran, 2016-12-06): maybe add an epsilon
	assert.Equal(t, 8983, dropped)
	kill <- struct{}{}
}
