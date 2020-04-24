package hashfill

import (
	"testing"

	"github.com/stretchr/testify/assert"
	geom "github.com/twpayne/go-geom"
)

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if needle == s {
			return true
		}
	}
	return false
}

func TestRecursiveFillIntersects(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")
	filler := NewRecursiveFiller(
		WithMaxPrecision(6),
	)

	expected := []string{"gcpvh7", "gcpvhe", "gcpvhh", "gcpvhj", "gcpvhk", "gcpvhm", "gcpvhs", "gcpvht"}

	hashes, err := filler.Fill(geofence, FillIntersects)
	assert.NoError(t, err)
	assert.Equal(t, expected, hashes)
}

func TestRecursiveFillIntersectsNotFixed(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")
	filler := NewRecursiveFiller(
		WithMaxPrecision(8),
	)

	hashes, err := filler.Fill(geofence, FillIntersects)
	assert.NoError(t, err)
	assert.Len(t, hashes, 948)
	assert.True(t, contains(hashes, "gcpvht0"))  // precision 7
	assert.True(t, contains(hashes, "gcpvhtb0")) // 8
}

func TestRecursiveFillIntersectsFixed(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")
	filler := NewRecursiveFiller(
		WithMaxPrecision(8),
		WithFixedPrecision(),
	)

	hashes, err := filler.Fill(geofence, FillIntersects)
	assert.NoError(t, err)
	assert.Len(t, hashes, 3242)
	assert.False(t, contains(hashes, "gcpvht0")) // No precision 7
	assert.True(t, contains(hashes, "gcpvhtb0")) // 8
}

func TestRecursiveFillContains(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/london.geojson")
	filler := NewRecursiveFiller(
		WithMaxPrecision(5),
	)

	expected := []string{"gcpsz", "gcptn", "gcptp", "gcpu6", "gcpu7", "gcpu8", "gcpu9", "gcpub", "gcpuc", "gcpud", "gcpue", "gcpuf", "gcpug", "gcpuk", "gcpum", "gcpur", "gcpus", "gcput", "gcpuu", "gcpuv", "gcpuw", "gcpux", "gcpuy", "gcpuz", "gcpv0", "gcpv1", "gcpv2", "gcpv3", "gcpv4", "gcpv5", "gcpv6", "gcpv7", "gcpve", "gcpvh", "gcpvj", "gcpvk", "gcpvm", "gcpvn", "gcpvp", "gcpvq", "gcpvr", "gcpvs", "gcpvt", "gcpvw", "gcpvx", "u10h2", "u10h3", "u10h8", "u10h9", "u10hb", "u10hc", "u10hd", "u10he", "u10hf", "u10hg", "u10hs", "u10hu", "u10hv", "u10j0", "u10j1", "u10j2", "u10j3", "u10j4", "u10j5", "u10j6", "u10j7", "u10j8", "u10jh"}

	hashes, err := filler.Fill(geofence, FillContains)
	assert.NoError(t, err)
	assert.Equal(t, expected, hashes)
}

type mockPredicate struct {
	res bool
	err error
}

func (m mockPredicate) Intersects(geofence *geom.Polygon, hash string) (bool, error) {
	return m.res, m.err
}

func (m mockPredicate) Contains(geofence *geom.Polygon, hash string) (bool, error) {
	return m.res, m.err
}

func TestPredicateContainError(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")

	filler := NewRecursiveFiller(
		WithMaxPrecision(8),
		WithPredicates(mockPredicate{false, assert.AnError}, mockPredicate{true, nil}),
	)
	_, err := filler.Fill(geofence, FillIntersects)
	assert.Error(t, err)
}

func TestPredicateIntersectsError(t *testing.T) {
	geofence := readFileAsGeometry(t, "testdata/regents.geojson")

	filler := NewRecursiveFiller(
		WithMaxPrecision(8),
		WithPredicates(mockPredicate{false, nil}, mockPredicate{false, assert.AnError}),
	)
	_, err := filler.Fill(geofence, FillIntersects)
	assert.Error(t, err)
}
