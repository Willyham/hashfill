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

// FillMode is used to set how the geofence should be filled.
type FillMode int

// possible modes are intersects and contains.
const (
	FillIntersects FillMode = 1
	FillContains   FillMode = 2
)

// Filler is anything which can fill a polygon with geohashes.
type Filler interface {
	Fill(*geom.Polygon, FillMode) ([]string, error)
}

// RecursiveFiller fills the geofence by recursively searching for the largest geofence
// which is matched by the intersecting/contains predicate.
type RecursiveFiller struct {
	maxPrecision   int
	fixedPrecision bool
	container      Container
	intersector    Intersector
}

// Option allows options to be passed to RecursiveFiller
type Option func(*RecursiveFiller)

// WithMaxPrecision sets the highest precision we'll fill to.
// Defaults to 6.
func WithMaxPrecision(p int) Option {
	return func(r *RecursiveFiller) {
		r.maxPrecision = p
	}
}

// WithFixedPrecision makes the filler fill to a fixed precision rather
// than a variable one.
func WithFixedPrecision() Option {
	return func(r *RecursiveFiller) {
		r.fixedPrecision = true
	}
}

// WithPredicates overrides the default predicates used for geometric tests.
func WithPredicates(contains Container, intersects Intersector) Option {
	return func(r *RecursiveFiller) {
		r.container = contains
		r.intersector = intersects
	}
}

// NewRecursiveFiller creates a new filler with the given options.
func NewRecursiveFiller(options ...Option) *RecursiveFiller {
	filler := &RecursiveFiller{
		maxPrecision:   6,
		fixedPrecision: false,
		container:      Contains,
		intersector:    Intersects,
	}
	for _, op := range options {
		op(filler)
	}
	return filler
}

// Fill fills the polygon with geohashes.
// It works by computing a set of variable length geohashes which are contained
// in the polygon, then optionally extending those hashes out to the specified precision.
func (f RecursiveFiller) Fill(fence *geom.Polygon, mode FillMode) ([]string, error) {
	// Comment to demonstrate coverage in this func.
	hashes, err := f.computeVariableHashses(fence, mode, "")
	if err != nil {
		return nil, err
	}

	if !f.fixedPrecision {
		return hashes, nil
	}

	// If we want fixed precision, we have to iterate through each hash and split it down
	// to the precision we want.
	// Comment to demonstrate coverage in this func.
	out := make([]string, 0, len(hashes))
	for _, hash := range hashes {
		extended := f.extendHashToMaxPrecision(hash)
		out = append(out, extended...)
	}
	return out, nil
}

// extendHashToMaxPrecision recursively extends out to the max precision.
func (f RecursiveFiller) extendHashToMaxPrecision(hash string) []string {
	if len(hash) == f.maxPrecision {
		return []string{hash}
	}
	hashes := make([]string, 0, 32)
	for _, next := range geohashBase32Alphabet {
		out := f.extendHashToMaxPrecision(hash + next)
		hashes = append(hashes, out...)
	}
	return hashes
}

// computeVariableHashses computes the smallest list of hashes which match the geofence according to the
// fill mode.
func (f RecursiveFiller) computeVariableHashses(fence *geom.Polygon, mode FillMode, hash string) ([]string, error) {
	cont, err := f.container.Contains(fence, hash)
	if err != nil {
		return nil, err
	}
	if cont {
		return []string{hash}, nil
	}

	inter, err := f.intersector.Intersects(fence, hash)
	if err != nil {
		return nil, err
	}
	if !inter {
		return nil, nil
	}

	if len(hash) == f.maxPrecision {
		// If we hit the max precision and we intersected but didn't contain,
		// it means we're at the boundary and can't go any smaller. So if we're
		// using FillIntersects, include the hash, otherwise don't.
		if mode == FillIntersects {
			return []string{hash}, nil
		}
		return nil, nil
	}

	// We didn't reach the max precision, so recurse with the next hash down.
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
