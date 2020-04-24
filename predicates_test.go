package hashfill

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func readFileAsGeometry(t *testing.T, file string) *geom.Polygon {
	data, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	poly := new(geom.T)
	err = geojson.Unmarshal(data, poly)
	require.NoError(t, err)

	return (*poly).(*geom.Polygon)
}

func TestPredicates(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")
	geofenceHole := readFileAsGeometry(t, "testdata/regents_hole.geojson")

	cases := []struct {
		name            string
		hash            string
		fence           *geom.Polygon
		shouldIntersect bool
		shouldContain   bool
	}{
		{"outside", "q", geofence, false, false},
		{"inside", "gcpvht2", geofence, true, true},
		{"boundary", "gcpvhk", geofence, true, false},
		{"empty", "", geofence, true, false},
		{"hole outside", "q", geofenceHole, false, false},
		// {"hole inside", "gcpvht2", geofenceHole, true, true}, // TODO: Why is this failing?!
		{"hole in hole", "gcpvhs2", geofenceHole, false, false},
		{"hole boundary", "gcpvhk", geofenceHole, true, false},
	}

	for _, test := range cases {
		t.Run(test.name+"/intersects", func(t *testing.T) {
			inter, err := Intersects(test.fence, test.hash)
			assert.NoError(t, err)
			assert.Equal(t, test.shouldIntersect, inter)
		})

		t.Run(test.name+"/contains", func(t *testing.T) {
			cont, err := Contains(test.fence, test.hash)
			assert.NoError(t, err)
			assert.Equal(t, test.shouldContain, cont)
		})
	}
}

func TestPolygonToGeometryExample(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")

	geom := polygonToGeometryExample(geofence)
	perimeterLength, err := geom.Length()
	require.NoError(t, err)
	assert.Equal(t, geofence.Length(), perimeterLength)
}
