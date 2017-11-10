package hashfill

import (
	"github.com/mmcloughlin/geohash"
	"github.com/paulsmith/gogeos/geos"
	geom "github.com/twpayne/go-geom"
)

type predicate func(geofence *geom.Polygon, hash string) (bool, error)

// Container tests if a hash is contained.
type Container interface {
	Contains(string) (bool, error)
}

// Intersector tests if a hash intersects.
type Intersector interface {
	Intersects(string) (bool, error)
}

func geomToGeosCoord(coord geom.Coord) geos.Coord {
	return geos.Coord{
		X: coord.X(),
		Y: coord.Y(),
	}
}

func geomToGeosCoords(coords []geom.Coord) []geos.Coord {
	out := make([]geos.Coord, len(coords))
	for i := 0; i < len(coords); i++ {
		out[i] = geomToGeosCoord(coords[i])
	}
	return out
}

// hashToGeometry converts a a geohash to a geos polygon by taking its bounding box.
func hashToGeometry(hash string) *geos.Geometry {
	bounds := geohash.BoundingBox(hash)
	return geos.Must(geos.NewPolygon([]geos.Coord{
		geos.NewCoord(bounds.MinLng, bounds.MinLat),
		geos.NewCoord(bounds.MinLng, bounds.MaxLat),
		geos.NewCoord(bounds.MaxLng, bounds.MaxLat),
		geos.NewCoord(bounds.MaxLng, bounds.MinLat),
		geos.NewCoord(bounds.MinLng, bounds.MinLat),
	}))
}

func polygonToGeometry(geofence *geom.Polygon) *geos.Geometry {
	// Convert the outer shell to geos format.
	shell := geofence.LinearRing(0).Coords()
	shellGeos := geomToGeosCoords(shell)

	// Convert each hole to geos format.
	numHoles := geofence.NumLinearRings() - 1
	holes := make([][]geos.Coord, numHoles)
	for i := 0; i < numHoles; i++ {
		holes[i] = geomToGeosCoords(geofence.LinearRing(i).Coords())
	}

	return geos.Must(geos.NewPolygon(shellGeos, holes...))
}

// Intersects tests if the geofence contains the hash by doing a geos intersection.
func Intersects(geofence *geom.Polygon, hash string) (bool, error) {
	hashGeo := hashToGeometry(hash)
	fence := polygonToGeometry(geofence)
	return fence.Intersects(hashGeo)
}

// Contains tests if the geofence contains the hash by doing a geos contains.
func Contains(geofence *geom.Polygon, hash string) (bool, error) {
	hashGeo := hashToGeometry(hash)
	fence := polygonToGeometry(geofence)
	return fence.Contains(hashGeo)
}
