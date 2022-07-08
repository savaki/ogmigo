package chainsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlonzoOrGreater(t *testing.T) {
	expectedResults := []bool{false, false, false, false, true, true}
	gotResults := make([]bool, 0, len(expectedResults))

	for _, era := range Eras {
		alonzoOrGreater := era.AlonzoOrGreater()
		gotResults = append(gotResults, alonzoOrGreater)
	}

	assert.Equal(t, expectedResults, gotResults)
}
