package hashfill

import (
	"github.com/engelsjk/polygol"
	"github.com/mmcloughlin/geohash"
	geom "github.com/twpayne/go-geom"
)

// Container tests if a hash is contained.
type Container interface {
	Contains(*geom.Polygon, string) (bool, error)
}

// Intersector tests if a hash intersects.
type Intersector interface {
	Intersects(*geom.Polygon, string) (bool, error)
}

type predicate func(geofence *geom.Polygon, hash string) (bool, error)

type containsFunc predicate

func (f containsFunc) Contains(p *geom.Polygon, hash string) (bool, error) {
	return f(p, hash)
}

type intersectsFunc predicate

func (f intersectsFunc) Intersects(p *geom.Polygon, hash string) (bool, error) {
	return f(p, hash)
}

// Intersects tests if the geofence contains the hash.
var Intersects = intersectsFunc(func(geofence *geom.Polygon, hash string) (bool, error) {
	hashGeom := hashToGeom(hash)
	out, err := polygol.Intersection(hashGeom, g2p(geofence))
	if err != nil {
		return false, err
	}
	if len(out) == 0 {
		return false, nil
	} else {
		return true, nil
	}
})

// Contains tests if the geofence contains the hash.
var Contains = containsFunc(func(geofence *geom.Polygon, hash string) (bool, error) {
	hashGeom := hashToGeom(hash)
	fenceGeom := g2p(geofence)
	out, err := polygol.Intersection(fenceGeom, hashGeom)
	if err != nil {
		return false, err
	}
	if len(out) == 0 {
		return false, nil
	} else {
		o, err := polygol.XOR(hashGeom, out)
		if err != nil {
			return false, err
		}
		if len(o) == 0 {
			return true, err
		}
		return false, nil
	}
})

func hashToGeom(hash string) polygol.Geom {
	bounds := geohash.BoundingBox(hash)
	return polygol.Geom{{{
		{bounds.MinLng, bounds.MinLat},
		{bounds.MinLng, bounds.MaxLat},
		{bounds.MaxLng, bounds.MaxLat},
		{bounds.MaxLng, bounds.MinLat},
		{bounds.MinLng, bounds.MinLat},
	}}}
}

func g2p(g geom.T) [][][][]float64 {

	var coords [][][]geom.Coord

	switch v := g.(type) {
	case *geom.Polygon:
		coords = [][][]geom.Coord{v.Coords()}
	case *geom.MultiPolygon:
		coords = v.Coords()
	}

	p := make([][][][]float64, len(coords))

	for i := range coords {
		p[i] = make([][][]float64, len(coords[i]))
		for j := range coords[i] {
			p[i][j] = make([][]float64, len(coords[i][j]))
			for k := range coords[i][j] {
				coord := coords[i][j][k]
				pt := make([]float64, 2)
				pt[0], pt[1] = coord.X(), coord.Y()
				p[i][j][k] = pt
			}
		}
	}

	return p
}
