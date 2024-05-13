// Package polygon provides functions for manipulation of polygon
// structures. A polygon is an N straight sided, not
// self-intersecting, shape, where N is greater than 2.
//
// The conventions for this package are x increases to the right, and
// y increases up the page (reverse of typical image formats). This
// convention gives meaning to clockwise and counter-clockwise.
package polygon

import (
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
)

// Zeroish is defined to merge points and avoid rounding error
// problems. The number is chosen to connect anything closer than 0.01
// (which is a convenience default for values representing
// millimeters).
var Zeroish = 1e-4

// Sort two numbers to be in ascending order.
func MinMax(a, b float64) (float64, float64) {
	if a <= b {
		return a, b
	}
	return b, a
}

// Point holds a 2d coordinate value. X increases to the right. Y
// increases up the page. These are the conventions of mathematical
// graph paper and not those of typical image formats.
type Point struct {
	X, Y float64
}

// BB determines the bounding box LL and TR corner points.
func BB(a, b Point) (c0, c1 Point) {
	c0.X, c1.X = MinMax(a.X, b.X)
	c0.Y, c1.Y = MinMax(a.Y, b.Y)
	return
}

// min returns the minimum of a pair of values.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of a pair of values.
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// isLeft determines if point a is left of the line segment (b->c).
func isLeft(a, b, c Point) (left bool) {
	if (a.Y <= b.Y) == (a.Y <= c.Y) {
		return // a is fully above or below (b->c)
	}
	if b.X > a.X && c.X > a.X {
		left = true
		return
	}
	if max(b.X, c.X) <= a.X {
		return
	}
	// a is horizontally within X range of BC.
	// compare a.Y to the y value of line BC at x = a.X.
	if y := b.Y + (c.Y-b.Y)/(c.X-b.X)*(a.X-b.X); math.Abs(a.Y-y) > Zeroish {
		// y below a.Y and c.Y is the lower left edge of bbCD0.
		left = (y < a.Y) == (b.Y < c.Y)
	}
	return
}

// intersect determines if two line segments (a->b) and (c->d)
// intersect (hit) and returns the point that they intersect. It also
// determines if the point a is to the left of the line (c->d).  The
// point c is evaluated for its leftness to (a->b) and this value is
// returned as hold.
func intersect(a, b, c, d Point) (hit bool, left, hold bool, at Point) {
	dABX, dABY := (b.X - a.X), (b.Y - a.Y)
	dCDX, dCDY := (d.X - c.X), (d.Y - c.Y)
	bbAB0, bbAB1 := BB(a, b)
	bbCD0, bbCD1 := BB(c, d)
	left = isLeft(a, c, d)
	hold = isLeft(c, a, b)
	// Do line bounding boxes not come close to overlapping each other?
	if (bbAB0.X > bbCD1.X && math.Abs(bbAB0.X-bbCD1.X) > Zeroish) ||
		(bbAB1.X < bbCD0.X && math.Abs(bbAB1.X-bbCD0.X) > Zeroish) ||
		(bbAB0.Y > bbCD1.Y && math.Abs(bbAB0.Y-bbCD0.Y) > Zeroish) ||
		(bbAB1.Y < bbCD0.Y && math.Abs(bbAB1.Y-bbCD0.Y) > Zeroish) {
		return
	}
	// Overlapping bounding box.
	bb0 := Point{X: max(bbAB0.X, bbCD0.X), Y: max(bbAB0.Y, bbCD0.Y)}
	bb1 := Point{X: min(bbAB1.X, bbCD1.X), Y: min(bbAB1.Y, bbCD1.Y)}
	if r := dABX*dCDY - dABY*dCDX; math.Abs(r) > Zeroish {
		if math.Abs(dABX) < Zeroish {
			at.X = a.X
			mCD := dCDY / dCDX
			cCD := d.Y - mCD*d.X
			at.Y = cCD + mCD*a.X
		} else if math.Abs(dCDX) < Zeroish {
			at.X = d.X
			mAB := dABY / dABX
			cAB := a.Y - mAB*a.X
			at.Y = cAB + mAB*d.X
		} else {
			mAB := dABY / dABX
			mCD := dCDY / dCDX
			cAB := a.Y - mAB*a.X
			cCD := d.Y - mCD*d.X
			at.X = -(cAB - cCD) / (mAB - mCD)
			at.Y = cAB + mAB*at.X
		}
	} else if colinear := (a.Y-d.Y)*dABX - (a.X-d.X)*dABY; math.Abs(colinear) > Zeroish {
		return // parallel but not co-linear.
	} else {
		if a == c {
			// ignore situation where the two lines start from the same place.
			return
		}
		if hit = MatchPoint(a, d); hit {
			at = d
			return
		}
		if hit = MatchPoint(c, b); hit {
			at = c
			return
		}
		if hit = MatchPoint(b, d); hit {
			at = d
			return
		}
		log.Printf("TODO unhandled co-linear lines: %v->%v vs %v->%v", a, b, c, d)
		return
	}
	hit = !(bb0.X > at.X || bb1.X < at.X || bb0.Y > at.Y || bb1.Y < at.Y)
	return
}

// Shape holds the points in a polygon and some convenience fields,
// such as the properties of its bounding box and whether the
// perimeter is clockwise (by convention a Hole) or counterclockwise
// (by convention a shape).
type Shape struct {
	// MinX etc represent the bounding box for a polygon.
	MinX, MinY, MaxX, MaxY float64
	// Hole indicates the polygon points are ordered (clockwise)
	// to represent a hole instead of an additive shape.
	Hole bool
	// Consecutive points on the perimeter of the polygon. There
	// is an implicit edge joining the last point to the first
	// point.
	PS []Point
}

// Shapes holds a set of polygon shapes each of arrays of (x,y)
// points.
type Shapes struct {
	P []*Shape
}

// Append appends a polygon shape constructed from a series of
// consecutive points. If p is nil, it is allocated. The return value
// is the appended collection of shapes.
func (p *Shapes) Append(pts ...Point) (*Shapes, error) {
	if len(pts) < 3 {
		return p, fmt.Errorf("polygon requires 3 or more points: got=%d", len(pts))
	}
	if p == nil {
		p = &Shapes{}
	}
	var minX, minY, maxX, maxY float64
	var ps []Point
	var zPt int
	for j, v := range pts {
		if minX > v.X || j == 0 {
			minX = v.X
		}
		if maxX < v.X || j == 0 {
			maxX = v.X
		}
		if minY > v.Y || j == 0 {
			minY = v.Y
		}
		if maxY < v.Y || j == 0 {
			maxY = v.Y
		}
		if j != 0 && (v.X < ps[zPt].X || (v.X == ps[zPt].X && v.Y < ps[zPt].Y)) {
			zPt = len(ps)
		}
		ps = append(ps, v)
	}
	tmp := append([]Point{}, ps[zPt:]...)
	ps = append(tmp, ps[:zPt]...)
	d1X, d1Y := ps[0].X-ps[len(ps)-1].X, ps[0].Y-ps[len(ps)-1].Y
	d2X, d2Y := ps[1].X-ps[0].X, ps[1].Y-ps[0].Y
	hole := (d1X*d2Y - d1Y*d2X) < 0
	poly := &Shape{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
		Hole: hole,
		PS:   ps,
	}
	p.P = append(p.P, poly)
	return p, nil
}

// Invert reverses the clockwise <-> counter-clockwise orientation of
// the shape without changing its starting point. The conventions for
// the package are shapes are counter-clockwise and holes are
// clockwise, so the .Hole value for the shape is inverted.
func (p *Shapes) Invert(i int) error {
	if i < 0 || i >= len(p.P) {
		return fmt.Errorf("invalid index %d but %d known shapes", i, len(p.P))
	}
	s := p.P[i]
	s.Hole = !s.Hole
	slices.Reverse(s.PS[1:])
	return nil
}

// Builder turns a set of points into a polygon shape and appends it
// to the provided value, p. If p is nil it is allocated. If the
// operation cannot be performed, the function panics. If you require
// more error control, call p.Append() instead.
func (p *Shapes) Builder(pts ...Point) *Shapes {
	var err error
	p, err = p.Append(pts...)
	if err != nil {
		panic(err)
	}
	return p
}

// Duplicate makes an independent copy of a set of polygon shapes.
func (p *Shapes) Duplicate() *Shapes {
	d := &Shapes{}
	for _, s := range p.P {
		var e []Point
		d.P = append(d.P, &Shape{
			MinX: s.MinX,
			MinY: s.MinY,
			MaxX: s.MaxX,
			MaxY: s.MaxY,
			Hole: s.Hole,
			PS:   append(e, s.PS...),
		})
	}
	return d
}

// MatchPoint recognizes when a is close enough to any of the points b...
func MatchPoint(a Point, b ...Point) bool {
	for _, c := range b {
		if math.Abs(a.X-c.X) < Zeroish && math.Abs(a.Y-c.Y) < Zeroish {
			return true
		}
	}
	return false
}

// combine computes the union of two Polygon shapes, indexed in p as n
// and m. This is either a no-op, or will generate one polygon and
// zero or more holes. The return value, banked, indicates how many
// additional shapes from index m have been resolved. This value can
// be negative.
func (p *Shapes) combine(n, m int) (banked int) {
	p1, p2 := p.P[n], p.P[m]
	if p1.MinX > p2.MaxX || p1.MaxX < p1.MinX || p1.MinY > p2.MaxY || p1.MaxY < p2.MinY {
		// Bounding boxes do not overlap.
		return
	}
	// Explore polygons p1, p2 for overlaps. Consider pairs of each
	// polygon at a time. Record each overlapping point with a
	// lookup table entry.
	hits := make(map[Point]bool)
	holds := true
	outside := true
	for i := 0; i < len(p1.PS); i++ {
		a := p1.PS[i]
		b := p1.PS[(i+1)%len(p1.PS)]
		for j := 0; j < len(p2.PS); j++ {
			c := p2.PS[j]
			d := p2.PS[(j+1)%len(p2.PS)]
			hit, left, hold, e := intersect(a, b, c, d)
			outside = outside != left
			holds = holds != hold
			if hit {
				hits[e] = true
				if !MatchPoint(e, c, d) {
					tmp := append([]Point{e}, p2.PS[j+1:]...)
					p2.PS = append(p2.PS[:j+1], tmp...)
					// possible the next intersection will be "before" this hit.
					j--
				}
				if !MatchPoint(e, a, b) {
					tmp := append([]Point{e}, p1.PS[i+1:]...)
					p1.PS = append(p1.PS[:i+1], tmp...)
					b = e
				}
			}
		}
	}
	if len(hits) == 0 {
		if holds && p1.Hole == p2.Hole {
			// since no hits and p1 holds p2, p2 is fully inside p1 - delete it
			p.P = append(p.P[:m], p.P[m+1:]...)
			banked = -1
		}
		if !outside && p1.Hole == p2.Hole {
			log.Printf("TODO n=%d should be swallowed by m=%d %v %v", n, m, p1, p2)
		}
		return
	}

	// Need to start from a point that is guaranteed to be on the
	// perimeter. Append() rotates the built shape to guarantee
	// that the 0th point is on the outer hull of the shape
	// (leftmost lower corner).
	union := &Shape{
		MinX: min(p1.MinX, p2.MinX),
		MinY: min(p1.MinY, p2.MinY),
		MaxX: max(p1.MaxX, p2.MaxX),
		MaxY: max(p1.MaxY, p2.MaxY),
	}
	src1, src2 := p1.PS, p2.PS
	var extra1, extra2 []Point
	var offset1, offset2 int

	// Initially, we step around p2 until we find the intersection
	// point of interest, and then we increment j instead to find
	// subsequent intersection points in p2.
	lockedOn := false
	for i, j := 0, 0; i < len(src1); {
		pt1 := src1[(offset1+i)%len(src1)]
		if hits[pt1] {
			// crossing point need to find it.
			cmp := src2[(offset2+j)%len(src2)]
			if cmp != pt1 {
				if lockedOn {
					extra2 = append(extra2, cmp)
					j++
				} else {
					offset2++
				}
				continue
			}
			if !lockedOn {
				lockedOn = true
			}
			i++
			src1, src2 = src2, src1
			i, j = j, i
			offset1, offset2 = offset2, offset1
			extra1, extra2 = extra2, extra1
		}
		i++
		union.PS = append(union.PS, pt1)
	}
	rest := p.P[m+1:]
	keep := append([]*Shape{}, p.P[n+1:m]...)
	var poly *Shapes
	for since, i := -1, 0; i < len(extra1); i++ {
		if hits[extra1[i]] {
			if since < 0 {
				since = i
				continue
			} else {
				if i+1-since > 2 {
					poly = poly.Builder(extra1[since : i+1]...)
				}
				since = -1
				continue
			}
		}
	}
	for since, i := -1, 0; i < len(extra2); i++ {
		if hits[extra2[i]] {
			if since < 0 {
				since = i
				continue
			} else {
				if i+1-since > 2 {
					poly = poly.Builder(extra2[since : i+1]...)
				}
				since = -1
				continue
			}
		}
	}
	banked = -1
	if poly != nil {
		keep = append(poly.P, keep...)
	}
	keep = append(append([]*Shape{union}, keep...), rest...)
	p.P = append(p.P[:n], keep...)
	return
}

// Union tries to combine all of the shape outlines into union outlines.
func (p *Shapes) Union() {
	// sort all of the polygons by their bounding boxes left to
	// right, down to up. This guarantees that the left most point
	// of the 0th polygon is an outer point.
	cf := func(a, b int) bool {
		if cmp := p.P[a].MinX - p.P[b].MinX; cmp < 0 {
			return true
		} else if cmp > 0 {
			return false
		}
		return p.P[a].MinY < p.P[b].MinY
	}
	sort.Slice(p.P, cf)
	for i := 1; i < len(p.P); i++ {
		for j := i; j < len(p.P); j++ {
			j += p.combine(i-1, j)
			if j+1 < len(p.P) && p.P[i-1].MaxX < p.P[j+1].MinX {
				break // next polygon too far right to overlap
			}
		}
	}
}
