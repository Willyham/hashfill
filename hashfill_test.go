package hashfill

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveFill(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")
	filler := RecursiveFiller{minPrecision: 8, fixedPrecision: false}
	hashes, err := filler.Fill(geofence, FillIntersects)
	assert.NoError(t, err)
	fmt.Println(hashes)
}
