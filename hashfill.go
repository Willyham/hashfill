package hashfill

import (
	geom "github.com/twpayne/go-geom"
)

// geohashBase32Alphabet is the list of charcters which make up the geohash alphabet.
var geohashBase32Alphabet = []string{
	"0", "1", "2", "3", "4", "5", "6", "7",
	"8", "9", "b", "c", "d", "e", "f", "g",
	"h", "j", "k", "m", "n", "p", "q", "r",
	"s", "t", "u", "v", "w", "x", "y", "z",
}

type FillMode int

const (
	FillIntersects FillMode = 1
	FillContains   FillMode = 2
)

// Filler is anything which can fill a polygon with geohashes.
type Filler interface {
	Fill(*geom.Polygon, FillMode) ([]string, error)
}

type RecursiveFiller struct {
	minPrecision   int
	fixedPrecision bool
}

// Fill fills the polygon with geohashes.
// It works by computing a set of variable length geohashes which are contained
// in the polygon, then extending those hashes out to the specified precision.
func (f RecursiveFiller) Fill(fence *geom.Polygon, mode FillMode) ([]string, error) {
	hashes, err := f.computeVariableHashses(fence, mode, "")
	if err != nil {
		return nil, err
	}

	if !f.fixedPrecision {
		return hashes, nil
	}

	// If we want fixed precision, we have to iterate through each hash and split it down
	// to the precision we want.
	out := make([]string, 0, len(hashes))
	for _, hash := range hashes {
		extended := f.extendHashToMaxPrecision(hash)
		out = append(out, extended...)
	}
	return out, nil
}

// extendHashToMaxPrecision recursively extends out to the max precision.
func (f RecursiveFiller) extendHashToMaxPrecision(hash string) []string {
	if len(hash) == f.minPrecision {
		return []string{hash}
	}
	hashes := make([]string, 0, 32)
	for _, next := range geohashBase32Alphabet {
		out := f.extendHashToMaxPrecision(hash + next)
		hashes = append(hashes, out...)
	}
	return hashes
}

// computeVariableHashses computes the smallest list of hashes which are entirely
// contained by the geofence.
func (f RecursiveFiller) computeVariableHashses(fence *geom.Polygon, mode FillMode, hash string) ([]string, error) {
	cont, err := Contains(fence, hash)
	if err != nil {
		return nil, err
	}
	if cont {
		return []string{hash}, nil
	}

	inter, err := Intersects(fence, hash)
	if err != nil {
		return nil, err
	}
	if !inter {
		return nil, nil
	}

	if len(hash) == f.minPrecision {
		if mode == FillIntersects {
			return []string{hash}, nil
		}
		return nil, nil
	}

	hashes := make([]string, 0)
	for _, next := range geohashBase32Alphabet {
		out, err := f.computeVariableHashses(fence, mode, hash+next)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, out...)
	}
	return hashes, nil
}
