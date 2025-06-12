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
// problems. The number is chosen to connect anything closer than
// 0.001 (which is a convenience default for values representing
// millimeters).
var Zeroish = 1e-6

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

// Line holds a 2d line between 2 Points.
type Line struct {
	From, To Point
}

// AddX adds a to x*b.
func (a Point) AddX(b Point, x float64) Point {
	return Point{
		X: a.X + b.X*x,
		Y: a.Y + b.Y*x,
	}
}

// BB determines the bounding box LL and TR corner points.
func BB(a, b Point) (ll, tr Point) {
	ll.X, tr.X = MinMax(a.X, b.X)
	ll.Y, tr.Y = MinMax(a.Y, b.Y)
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

// Shape holds the points in a polygon and some convenience fields,
// such as the properties of its bounding box and whether the
// perimeter is clockwise (by convention a Hole) or counterclockwise
// (by convention a shape).
type Shape struct {
	// Index is a string assigned when the shape is defined. It
	// can be really useful when trying to debug why polygons have
	// been merged.
	Index string
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

// Return the bounding box lower left and top right corner points for
// the shape.
func (s *Shape) BB() (ll, tr Point) {
	return Point{s.MinX, s.MinY}, Point{s.MaxX, s.MaxY}
}

// Shapes holds a set of polygon shapes each of arrays of (x,y)
// points.
type Shapes struct {
	// index is used to assign unique Index values to Shape
	// members of P.
	index int
	// P holds the polygon Shape data.
	P []*Shape
}

// Return the bounding box lower left and top right corner points for
// the shapes.
func (p *Shapes) BB() (ll, tr Point) {
	for i, s := range p.P {
		if i == 0 {
			ll, tr = s.BB()
		} else {
			ll2, tr2 := s.BB()
			ll, _ = BB(ll, ll2)
			_, tr = BB(tr, tr2)
		}
	}
	return
}

// Rationalize builds a properly constructed shape.
func Rationalize(pts []Point) (*Shape, error) {
	if len(pts) < 3 {
		return nil, fmt.Errorf("polygon requires 3 or more points: got=%d", len(pts))
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
	return &Shape{
		MinX: minX,
		MinY: minY,
		MaxX: maxX,
		MaxY: maxY,
		Hole: hole,
		PS:   ps,
	}, nil
}

// Transform returns a rotated Shapes structure, p is rotated by theta
// radians (+ve = counterclockwise) around a fixed Point, pt, and
// scales the rotated shape by a factor of scale. The scaled and
// rotated shape is then translated from pt to to.
func (p *Shapes) Transform(at, to Point, theta, scale float64) *Shapes {
	if p == nil {
		return nil
	}
	var sh *Shapes
	s := math.Sin(theta) * scale
	c := math.Cos(theta) * scale
	for _, v := range p.P {
		var pts []Point
		for _, pt := range v.PS {
			dX, dY := pt.X-at.X, pt.Y-at.Y
			pts = append(pts, Point{
				X: to.X + c*dX - s*dY,
				Y: to.Y + s*dX + c*dY,
			})
		}
		sh = sh.Builder(pts...)
	}
	return sh
}

// Append appends a polygon shape constructed from a series of
// consecutive points. If p is nil, it is allocated. The return value
// is the appended collection of shapes. The newly added polygon is
// the last one, and it's zeroth point is guaranteed to be leftmost
// and lowest.
func (p *Shapes) Append(pts ...Point) (*Shapes, error) {
	poly, err := Rationalize(pts)
	if err != nil {
		return p, err
	}
	if p == nil {
		poly.Index = "0"
		return &Shapes{
			P: []*Shape{poly},
		}, nil
	}
	p.index++
	poly.Index = fmt.Sprint(p.index)
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

// Dot computes the dot product of two vectors.
func (a Point) Dot(b Point) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Unit returns a unit vector in the direction of a towards b, or an
// error if the points are too close.
func (a Point) Unit(b Point) (u Point, err error) {
	v := b.AddX(a, -1)
	l2 := v.Dot(v)
	if l2 < Zeroish {
		err = fmt.Errorf("a=%v and b=%v too close", a, b)
		return
	}
	isqrt := 1.0 / math.Sqrt(l2)
	u = Point{X: v.X * isqrt, Y: v.Y * isqrt}
	return
}

// moreClockwise confirms that c is more clockwise than d from b.
func moreClockwise(b, c, d Point) bool {
	bc := c.AddX(b, -1)
	bd := d.AddX(b, -1)
	crossBCBD := bc.X*bd.Y - bc.Y*bd.X
	return crossBCBD >= 0
}

// isLeft determines if point a is left of the line segment (b->c). By
// "to the left of" we mean looking along the line (b->c) towards c,
// do we see a on the left of this line?
func (a Point) isLeft(b, c Point) bool {
	return moreClockwise(b, c, a)
}

// Narrows computes the polygon corners where two (non-crossing) lines
// (a->b) (c->d) fall within some threshold distance, delta.
func Narrows(a, b, c, d Point, delta float64) (hit bool, w, x, y, z Point) {
	hit = false
	u1, err := a.Unit(b)
	if err != nil {
		return
	}
	u2, err := c.Unit(d)
	if err != nil {
		return
	}
	phi := u1.Dot(u2)
	if phi > 0 {
		return // more parallel than anti-parallel
	}
	delta2 := delta * delta
	if phi*phi > 1-Zeroish {
		// anti co-linear: calculate separation.
		v := c.AddX(a, -1)
		shift := v.Dot(u1)
		v = v.AddX(u1, -shift)
		v2 := v.Dot(v)
		if v2 > delta2 {
			return
		}
		// overlap extending on line.
		excess := math.Sqrt(delta2 - v2)
		// in u1 direction, compute a sortable offset
		oa := a.Dot(u1)
		ob := b.Dot(u1)
		oc := c.Dot(u1)
		od := d.Dot(u1)
		if oa-excess > oc || ob+excess < od {
			return
		}
		w = a
		z = d
		if od < oa {
			if od+excess < oa {
				z.X = a.X + v.X - excess*u1.X
				z.Y = a.Y + v.Y - excess*u1.Y
			}
		} else {
			if oa+excess < od {
				w.X = d.X - v.X - excess*u1.X
				w.Y = d.Y - v.Y - excess*u1.Y
			}
		}
		x = b
		y = c
		if oc < ob {
			if oc+excess < ob {
				x.X = c.X - v.X + excess*u1.X
				x.Y = c.Y - v.Y + excess*u1.Y
			}
		} else {
			if ob+excess < oc {
				x.X = b.X + v.X + excess*u1.X
				x.Y = b.Y + v.Y + excess*u1.Y
			}
		}
		hit = true
		return
	}
	// non co-linear, converging on point, P.
	ds := c.AddX(a, -1)
	du := Point{
		X: (u1.X - phi*u2.X) / (1 - phi*phi),
		Y: (u1.Y - phi*u2.Y) / (1 - phi*phi),
	}
	alpha := ds.Dot(du)
	p := a.AddX(u1, alpha)
	r := delta / (2 * math.Cos(0.5*math.Acos(phi)))
	// short of B?
	bp := p.AddX(b, -1)
	if bp.Dot(u1) > r {
		return
	}
	pc := c.AddX(p, -1)
	if pc.Dot(u2) > r {
		return
	}
	w, x, y, z = a, b, c, d
	if alpha > r {
		w = p.AddX(u1, -r)
	}
	pd := d.AddX(p, -1)
	if beta := pd.Dot(u2); beta > r {
		z = p.AddX(u2, r)
	}
	hit = true
	return
}

// intersect determines if two line segments (a->b) and (c->d)
// intersect (hit) and returns the point that they intersect. It also
// determines if the point a is to the 'left' of the line (c->d). See
// isLeft() for calculation. The point c is evaluated for its leftness
// to (a->b) and this value is returned as hold.
func intersect(a, b, c, d Point) (hit bool, left, hold bool, at Point) {
	dABX, dABY := (b.X - a.X), (b.Y - a.Y)
	dCDX, dCDY := (d.X - c.X), (d.Y - c.Y)
	bbAB0, bbAB1 := BB(a, b)
	bbCD0, bbCD1 := BB(c, d)
	left = a.isLeft(c, d)
	hold = c.isLeft(a, b)
	// Do line bounding boxes not come close to overlapping each other?
	if (bbAB0.X > bbCD1.X && math.Abs(bbAB0.X-bbCD1.X) > Zeroish) ||
		(bbAB1.X < bbCD0.X && math.Abs(bbAB1.X-bbCD0.X) > Zeroish) ||
		(bbAB0.Y > bbCD1.Y && math.Abs(bbAB0.Y-bbCD1.Y) > Zeroish) ||
		(bbAB1.Y < bbCD0.Y && math.Abs(bbAB1.Y-bbCD0.Y) > Zeroish) {
		return
	}
	// Overlapping bounding box (extended slightly by the rounding error protection).
	bb0 := Point{X: max(bbAB0.X, bbCD0.X), Y: max(bbAB0.Y, bbCD0.Y)}
	bb1 := Point{X: min(bbAB1.X, bbCD1.X), Y: min(bbAB1.Y, bbCD1.Y)}
	if bb0.X == bb1.X {
		bb0.X -= Zeroish / 2
		bb1.X += Zeroish / 2
	}
	if bb0.Y == bb1.Y {
		bb0.Y -= Zeroish / 2
		bb1.Y += Zeroish / 2
	}
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
		if MatchPoint(a, at) {
			at = a
		} else if MatchPoint(b, at) {
			at = b
		}
		// Confirm at falls within the bounding box of both lines.
		hit = !((bb0.X-Zeroish) > at.X || (bb1.X+Zeroish) < at.X || (bb0.Y-Zeroish) > at.Y || (bb1.Y+Zeroish) < at.Y)
		return
	}
	if colinear := (a.Y-d.Y)*dABX - (a.X-d.X)*dABY; math.Abs(colinear) > Zeroish {
		return // parallel but not co-linear.
	}
	if a == c {
		// ignore situation where the two lines start from the same place.
		return
	}
	if hit = MatchPoint(a, d); hit {
		at = a
		return
	}
	if hit = MatchPoint(b, d); hit {
		at = b
		return
	}
	if hit = MatchPoint(b, c); hit {
		at = b
		return
	}
	return
}

// dissolve eliminates collinear points from a polygon.
func (s *Shape) dissolve() (poly *Shape, dissolved bool) {
	if s == nil {
		return
	}
	pts := s.PS
	for i := 0; i < len(pts); {
		a := pts[i]
		bI := (i + 1) % len(pts)
		b := pts[bI] // evaluate whether to delete this
		bad := false
		if MatchPoint(a, b) {
			bad = true
		} else if u, err := a.Unit(b); err != nil {
			bad = true
		} else if c := pts[(i+2)%len(pts)]; MatchPoint(b, c) {
			bad = true
		} else if v, err := a.Unit(c); err != nil {
			bad = true
		} else if math.Abs(u.Dot(v)-1) < Zeroish {
			bad = true
		}
		if !bad {
			i++
			continue
		}
		pts = append(pts[:bI], pts[bI+1:]...)
		dissolved = true
	}
	var err error
	poly, err = Rationalize(pts)
	if err != nil {
		log.Fatalf("failed to rationalize: %v", err)
	}
	return
}

// Inside confirms that a is fully inside some polygon.
func (a Point) Inside(p *Shape) bool {
	if a.X < p.MinX || a.X > p.MaxX || a.Y < p.MinY || a.Y > p.MaxY {
		return false
	}
	// Point is inside the bounding box for p.  Consider how many
	// times the line (a->to) intersects with a line from p.
	// odd = inside, even = outside.
	to := a
	to.X = p.MaxX + 1
	inside := false
	prev := p.PS[len(p.PS)-1]
	for i, next := range p.PS {
		hit, _, _, _ := intersect(a, to, prev, next)
		if hit && math.Abs(a.Y-prev.Y) < Zeroish {
			// The prev point lies on the line a->to, so
			// we only consider it to be worth double
			// counting if the line doubles back on
			// itself.
			pprev := p.PS[(i+len(p.PS)-2)%len(p.PS)]
			if cf := (prev.Y - pprev.Y) * (next.Y - prev.Y); cf > 0 {
				prev = next
				continue
			}
		}
		if hit {
			inside = !inside
		}
		prev = next
	}
	return inside
}

// crossings evaluates p1 and p2 for common points of intersection. It
// returns n1 and n2 as the same shapes but with all of the hit points
// inserted into both shapes.
func crossings(p1, p2 *Shape) (hits map[Point]bool, n1, n2 *Shape) {
	n1, _ = Rationalize(p1.PS)
	n2, _ = Rationalize(p2.PS)
	hits = make(map[Point]bool)
	for i := 0; i < len(n1.PS); i++ {
		a := n1.PS[i]
		b := n1.PS[(i+1)%len(n1.PS)]
		if MatchPoint(a, b) {
			// trim out points that are too close together
			if i == 0 {
				n1.PS = append(n1.PS[:1], n1.PS[2:]...)
			} else {
				n1.PS = append(n1.PS[:i], n1.PS[i+1:]...)
			}
			i--
			continue
		}
		for j := 0; j < len(n2.PS); j++ {
			c := n2.PS[j]
			d := n2.PS[(j+1)%len(n2.PS)]
			// Close but not equal is a source of
			// problems, so given a close match treat a as
			// the anchor point and move c and/or d to it.
			if MatchPoint(a, c) && a != c {
				n2.PS[j] = a
				c = a
			}
			if MatchPoint(a, d) && a != d {
				n2.PS[(j+1)%len(n2.PS)] = a
				d = a
			}
			if MatchPoint(c, d) {
				// trim out points that are too close together
				// preserve the 0th point.
				if j == 0 {
					n2.PS = append(n2.PS[:1], n2.PS[2:]...)
				} else {
					n2.PS = append(n2.PS[:j], n2.PS[j+1:]...)
				}
				j--
				continue
			}
			hit, _, _, e := intersect(a, b, c, d)
			if hit {
				// Prefer canonical points vs derived ones.
				// Above we've confirmed that a != b.
				if MatchPoint(e, a) && e != a {
					e = a
				} else if MatchPoint(e, b) && e != b {
					e = b
				}
				// For this polygon we nudge the
				// points themselves. This is needed to
				// make use of the hits map later.
				if MatchPoint(e, c) && e != c {
					c = e
					n2.PS[j] = e
				} else if MatchPoint(e, d) && e != d {
					d = e
					n2.PS[(j+1)%len(n2.PS)] = e
				}
				hits[e] = true
				if !MatchPoint(e, c, d) {
					tmp := append([]Point{e}, n2.PS[j+1:]...)
					n2.PS = append(n2.PS[:j+1], tmp...)
					// possible the next intersection will be "before" this hit.
					j--
				}
				if !MatchPoint(e, a, b) {
					tmp := append([]Point{e}, n1.PS[i+1:]...)
					n1.PS = append(n1.PS[:i+1], tmp...)
					b = e
				}
			}
		}
	}
	return
}

// combine computes the union of two Polygon shapes, indexed in p as n
// and m. This is either a no-op, or will generate one polygon and
// zero or more holes. The return value, banked, indicates how many
// additional shapes from index m have been resolved. This value can
// be negative.
func (p *Shapes) combine(n, m int) (banked int) {
	banked = m + 1
	p1, p2 := p.P[n], p.P[m]
	if p1.MinX > p2.MaxX || p1.MaxX < p2.MinX || p1.MinY > p2.MaxY || p1.MaxY < p2.MinY {
		// Bounding boxes do not overlap.
		return
	}
	// Explore polygons p1, p2 for overlaps. Consider pairs of each
	// polygon at a time. Record each overlapping point with a
	// lookup table entry.
	hits, p1, p2 := crossings(p1, p2)
	if len(hits) == 0 {
		if p1.Hole != p2.Hole {
			banked = m + 1
			return
		}
		// No intersections, but one polygon might consume other.
		if p1.PS[0].Inside(p2) {
			p2.Index = fmt.Sprint("(", p2.Index, "!", p1.Index, ")")
			p.P = append(p.P[:n], p.P[n+1:]...)
			banked = n + 1
			return
		}
		if p2.PS[0].Inside(p1) {
			p1.Index = fmt.Sprint("(", p1.Index, "!", p2.Index, ")")
			p.P = append(p.P[:m], p.P[m+1:]...)
			banked = m
			return
		}
		return
	}

	src1, src2 := p1.PS, p2.PS
	var pts, extra1, extra2 []Point
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
			lockedOn = true
			ptKeep := src1[(offset1+i+1)%len(src1)]
			ptSwap := src2[(offset2+j+1)%len(src2)]
			if moreClockwise(pt1, ptSwap, ptKeep) {
				i++
				src1, src2 = src2, src1
				i, j = j, i
				offset1, offset2 = offset2, offset1
				extra1, extra2 = extra2, extra1
			}
		}
		i++
		pts = append(pts, pt1)
	}
	union, err := Rationalize(pts)
	if err != nil {
		log.Fatalf("union polygon failed to rationailze: %v", err)
	}
	was := len(union.PS)
	union, dissolved := union.dissolve()
	if dissolved && was < len(union.PS) {
		log.Printf("dissolved negative points was=%d, is=%d", was, len(union.PS))
	}
	union.Index = fmt.Sprint("(", p1.Index, "+", p2.Index, ")")
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
	// The merged polygon may overlap with a previously
	// non-overlapping polygon, so backtrack to the one
	// immediately after this merged polygon.
	banked = n + 1
	if poly != nil {
		keep = append(poly.P, keep...)
	}
	keep = append(append([]*Shape{union}, keep...), rest...)
	p.P = append(p.P[:n], keep...)
	return
}

// Reorder sorts all of the polygons by their bounding boxes left to
// right, down to up. This guarantees that the left most point of the
// 0th polygon is an outer point.
func (p *Shapes) Reorder() {
	cf := func(a, b int) bool {
		if cmp := p.P[a].MinX - p.P[b].MinX; cmp < 0 {
			return true
		} else if cmp > 0 {
			return false
		}
		if cmp := p.P[a].MinY - p.P[b].MinY; cmp < 0 {
			return true
		} else if cmp > 0 {
			return false
		}
		if cmp := p.P[a].MaxX - p.P[b].MaxX; cmp > 0 {
			return true
		} else if cmp < 0 {
			return false
		}
		return p.P[a].MaxY > p.P[b].MaxY
	}
	sort.Slice(p.P, cf)
}

// Add enhances p by importing s into it. No effort is made to
// unionize overlapping outlines. Call Union on the returned shapes
// for that. This function alters p, but not s.
func (p *Shapes) Add(s *Shapes) *Shapes {
	for _, o := range s.P {
		p = p.Builder(o.PS...)
	}
	p.Reorder()
	return p
}

// Union tries to combine all of the overlapping shape outlines into
// union outlines.
func (p *Shapes) Union() {
	p.Reorder()
	for i := 1; i < len(p.P); i++ {
		for j := i; j < len(p.P); {
			j = p.combine(i-1, j)
			if j < len(p.P) && p.P[i-1].MaxX < p.P[j].MinX {
				break // next polygon too far right to overlap
			}
		}
	}
}

// Inflate inflates an indexed shape by distance, d. Holes are
// deflated by this amount. If we inflate a circle by d, its radius
// will increase by that much.
func (p *Shapes) Inflate(n int, d float64) error {
	if n < 0 || n >= len(p.P) {
		return fmt.Errorf("invalid polygon=%d in shapes (%d known)", n, len(p.P))
	}
	if d == 0 {
		return nil // nothing needed
	}
	s := p.P[n]
	first := s.PS[0]
	last := s.PS[len(s.PS)-1]
	d *= 0.5 // Since we add an offset twice per point.
	var pts []Point
	for i, this := range s.PS {
		pre := this
		next := first
		if i < len(s.PS)-1 {
			next = s.PS[i+1]
		}

		dX, dY := this.X-last.X, this.Y-last.Y
		r := math.Sqrt(dX*dX + dY*dY)
		dX, dY = d*dX/r, d*dY/r
		this.X += dY
		this.Y -= dX

		dX, dY = next.X-pre.X, next.Y-pre.Y
		r = math.Sqrt(dX*dX + dY*dY)
		dX, dY = d*dX/r, d*dY/r
		this.X += dY
		this.Y -= dX

		pts = append(pts, this)
		last = pre
	}
	poly, err := Rationalize(pts)
	if err != nil {
		return err
	}
	poly.Index = s.Index
	p.P[n] = poly
	return nil
}

// Slice returns an array of horizontal (dy=0) lines to render the
// filled polygon. This can be used to rasterize a shape in some
// output format. The radial width of a rendered line is d. The lines
// are drawn from d/2 inside the shape to allow for this imprecision.
// If s is known to contain holes, and the indices of the holes are
// provided, then the corresponding polygon holes are used to further
// shorten the returned lines.
func (s *Shapes) Slice(i int, d float64, holeI ...int) (lines []Line, err error) {
	if s == nil || i < 0 || i >= len(s.P) {
		err = fmt.Errorf("invalid index %d for shapes", i)
		return
	}
	// Walk from least Y+d/2, to largest Y-d/2.
	p := s.P[i]
	if p.Hole {
		err = fmt.Errorf("no overlap with (shape %d) a hole", i)
		return
	}
	half := d / 2
	bottom, top := p.MinY, p.MaxY
	if top < bottom {
		bottom = (top + bottom) / 2
	}
	// X range guaranteed to extend outside of polygon.
	left, right := p.MinX-half, p.MaxX+half
	for level := bottom + half; level < top; level += half {
		a := Point{X: left, Y: level}
		b := Point{X: right, Y: level}
		var ats []float64
		for j := 0; j < len(p.PS); j++ {
			from := p.PS[j]
			to := p.PS[(j+1)%len(p.PS)]
			hit, _, _, e := intersect(a, b, from, to)
			if !hit {
				continue
			}
			ats = append(ats, e.X)
		}
		if len(ats) == 0 {
			continue
		}
		if len(ats)&1 == 1 {
			err = fmt.Errorf("shape %d has odd crossings at %f", i, level)
			return
		}
		sort.Slice(ats, func(i, j int) bool { return ats[i] < ats[j] })
		for j := 0; j < len(ats); j += 2 {
			line := Line{
				From: Point{X: ats[j] + half, Y: level},
				To:   Point{X: ats[j+1] - half, Y: level},
			}
			if line.From.X > line.To.X {
				continue // too short to render
			}
			// cut line if it overlaps a hole. Because the
			// holes do not intersect the the perimeter of
			// any non-hold polygon, the lines are either
			// broken by a hole into two, or do not
			// overlap at all.
			var hits []float64
			for _, hi := range holeI {
				hole := s.P[hi]
				if hole.MaxY < level || hole.MinY > level || hole.MinX > line.To.X || hole.MaxX < line.From.X {
					continue
				}
				for k := 0; k < len(hole.PS); k++ {
					a := hole.PS[k]
					b := hole.PS[(k+1)%len(hole.PS)]
					hit, _, _, e := intersect(line.From, line.To, a, b)
					if hit {
						hits = append(hits, e.X)
					}
				}
			}
			if len(hits) == 0 {
				lines = append(lines, line)
				continue
			}
			sort.Slice(hits, func(i, j int) bool { return hits[i] < hits[j] })
			hits = append(append([]float64{line.From.X - half}, hits...), line.To.X+half)
			for hi := 0; hi < len(hits); hi += 2 {
				from := hits[hi] + half
				to := hits[hi+1] - half
				if from+half > to-half {
					continue
				}
				lines = append(lines, Line{
					From: Point{X: from, Y: level},
					To:   Point{X: to, Y: level},
				})
			}
		}
	}
	return
}

// OptimizeLines rearranges the result of (*Shapes) Slice() into lines
// that can be plotted in a shorter time. It works by reordering
// consecutive lines when that minimizes the flight time of the
// plotter head between lines.
func OptimizeLines(lines []Line) {
	var last Point
	for i, line := range lines {
		dF := line.From.AddX(last, -1)
		dT := line.To.AddX(last, -1)
		cf := dT.Dot(dT) - dF.Dot(dF)
		if cf < 0 {
			lines[i] = Line{
				From: line.To,
				To:   line.From,
			}
			last = line.From
		} else {
			last = line.To
		}
	}
}
