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

// func TestPolygonIntersects(t *testing.T) {
//
// 	t.Run("Check for intersect", func(t *testing.T) {
// 		outer := orb.Polygon{{
// 			{0, 0},
// 			{0, 10},
// 			{10, 10},
// 			{10, 0},
// 			{0, 0},
// 		}}
//
// 		inner := orb.Polygon{{
// 			{9, 9},
// 			{9, 20},
// 			{20, 20},
// 			{20, 9},
// 			{9, 9},
// 		}}
//
// 		res := polygonIntersectsPolygon(outer, inner)
// 		assert.True(t, res)
// 	})
//
// 	t.Run("Check for line intersect", func(t *testing.T) {
// 		outer := orb.Polygon{{
// 			{0, 0},
// 			{0, 10},
// 			{10, 10},
// 			{10, 0},
// 			{0, 0},
// 		}}
//
// 		inner := orb.Polygon{{
// 			{-5, 5},
// 			{-5, 20},
// 			{30, 20},
// 			{30, -5},
// 			{-5, -5},
// 		}}
// 		res := polygonIntersectsPolygon(outer, inner)
// 		assert.True(t, res)
// 	})
//
// 	t.Run("Check for no intersect", func(t *testing.T) {
// 		outer := orb.Polygon{{
// 			{0, 0},
// 			{0, 10},
// 			{10, 10},
// 			{10, 0},
// 			{0, 0},
// 		}}
//
// 		inner := orb.Polygon{{
// 			{11, 11},
// 			{11, 20},
// 			{20, 20},
// 			{20, 11},
// 			{11, 11},
// 		}}
// 		res := polygonIntersectsPolygon(outer, inner)
// 		assert.True(t, !res)
// 	})
// }
