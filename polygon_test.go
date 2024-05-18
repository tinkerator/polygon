package polygon

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	vs := []struct{ x, y, a, b float64 }{
		{x: 1, y: 2, a: 1, b: 2},
		{x: 2, y: 1, a: 1, b: 2},
		{x: -1, y: -2, a: -2, b: -1},
		{x: -1, y: 1, a: -1, b: 1},
	}
	for i, v := range vs {
		a, b := MinMax(v.x, v.y)
		if a != v.a || b != v.b {
			t.Errorf("test=%d MinMax(%f,%f) failed: got a=%f, b=%f, wanted a=%f, b=%f", i, v.x, v.y, a, b, v.a, v.b)
		}
	}
}

func TestUnion(t *testing.T) {
	var ss *Shapes
	var err error
	ss, err = ss.Append(Point{0, 0}, Point{2, 0}, Point{2, 2}, Point{0, 2})
	if ss == nil || err != nil {
		t.Fatalf("failed to add first square: %v", err)
	}
	if ss.P[0].Hole {
		t.Fatalf("counter-clockwise shape is a hole: %#v", *ss.P[0])
	}
	ss = ss.Builder(Point{1, 1}, Point{1, 3}, Point{3, 3}, Point{3, 1})
	if len(ss.P) != 2 {
		t.Fatalf("failed to add second shape, only %d shapes recorded", len(ss.P))
	}
	if !ss.P[1].Hole {
		t.Fatalf("clockwise shape is not a hole: %#v", *ss.P[1])
	}
	if err = ss.Invert(3); err == nil {
		t.Fatal("invalid polygon shape inversion performed?")
	}
	if err = ss.Invert(1); err != nil {
		t.Fatalf("polygon shape inversion failed: %v", err)
	}
	if ss.P[1].Hole {
		t.Fatalf("counter-clockwise shape is still a hole: %#v", *ss.P[1])
	}
	ss = ss.Builder(Point{0, 4}, Point{2, 4}, Point{2, 6}, Point{0, 6})
	if len(ss.P) != 3 {
		t.Fatalf("failed to add second shape, only %d shapes recorded", len(ss.P))
	}
	if ss.P[2].Hole {
		t.Fatalf("counter-clockwise shape is a hole: %#v", *ss.P[2])
	}
	ss = ss.Builder(Point{1, 5}, Point{3, 5}, Point{3, 7}, Point{1, 7})
	if len(ss.P) != 4 {
		t.Fatalf("failed to add second shape, only %d shapes recorded", len(ss.P))
	}
	if ss.P[3].Hole {
		t.Fatalf("counter-clockwise shape is a hole: %#v", *ss.P[2])
	}
	ss.Union()
	if len(ss.P) != 2 {
		t.Fatalf("post union shape count != 2, got=%d", len(ss.P))
	}
	us := ss.P[0].PS
	if len(us) != 8 {
		t.Fatalf("expecting 8 post union points: got=%v", us)
	}
	expect := []Point{{0, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 3}, {1, 3}, {1, 2}, {0, 2}}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}
	us = ss.P[1].PS
	expect = []Point{{0, 4}, {2, 4}, {2, 5}, {3, 5}, {3, 7}, {1, 7}, {1, 6}, {0, 6}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[1] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	// Validate coincident heavy overlaps.
	ss = nil
	ss = ss.Builder([]Point{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
		{X: 1, Y: 2},
	}...).Builder([]Point{
		{X: 0, Y: 0},
		{X: 1.5, Y: 0},
		{X: 1.5, Y: 1},
		{X: 0, Y: 1},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{0, 0}, {1.5, 0}, {1.5, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}, {0, 1}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
		{X: 1, Y: 2},
	}...).Builder([]Point{
		{X: 2, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 2},
		{X: 2, Y: 2},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{1, 1}, {2, 1}, {2, 0}, {3, 0}, {3, 2}, {1, 2}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
		{X: 1, Y: 2},
	}...).Builder([]Point{
		{X: 2, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 3},
		{X: 2, Y: 3},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{1, 1}, {2, 1}, {2, 0}, {3, 0}, {3, 3}, {2, 3}, {2, 2}, {1, 2}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 3},
		{X: 1, Y: 3},
	}...).Builder([]Point{
		{X: 2, Y: 1},
		{X: 3, Y: 1},
		{X: 3, Y: 2},
		{X: 2, Y: 2},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{1, 0}, {2, 0}, {2, 1}, {3, 1}, {3, 2}, {2, 2}, {2, 3}, {1, 3}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
		{X: 1, Y: 2},
	}...).Builder([]Point{
		{X: 0, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 3},
		{X: 0, Y: 3},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{0, 0}, {3, 0}, {3, 3}, {0, 3}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 0, Y: 0},
		{X: 2, Y: 0},
		{X: 2, Y: 2},
		{X: 0, Y: 2},
	}...).Builder([]Point{
		{X: 0, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 3},
		{X: 0, Y: 3},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{0, 0}, {3, 0}, {3, 3}, {0, 3}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 0, Y: 0},
		{X: 3, Y: 0},
		{X: 3, Y: 1},
		{X: 0, Y: 1},
	}...).Builder([]Point{
		{X: 2, Y: 0},
		{X: 4, Y: 0},
		{X: 4, Y: 2},
		{X: 2, Y: 2},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{0, 0}, {4, 0}, {4, 2}, {2, 2}, {2, 1}, {0, 1}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	ss = nil
	ss = ss.Builder([]Point{
		{X: 0, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 1},
		{X: 0, Y: 1},
	}...).Builder([]Point{
		{X: 1, Y: 2},
		{X: 5, Y: 2},
		{X: 5, Y: 3},
		{X: 1, Y: 3},
	}...).Builder([]Point{
		{X: 4, Y: 0},
		{X: 5, Y: 0},
		{X: 5, Y: 3},
		{X: 4, Y: 3},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Errorf("expecting a single poly, but got %d", len(ss.P))
		for i := 0; i < len(ss.P); i++ {
			t.Errorf("P[%d] = %#v", i, *ss.P[i])
		}
	}
	us = ss.P[0].PS
	expect = []Point{{0, 0}, {5, 0}, {5, 3}, {1, 3}, {1, 2}, {4, 2}, {4, 1}, {0, 1}}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; got != want {
			t.Errorf("union[0] point[%d]: got=%v, want=%v", i, got, want)
		}
	}
}

func TestIntersect(t *testing.T) {
	vs := []struct {
		a, b, c, d, at  Point
		hit, left, hold bool
	}{
		{
			a:    Point{X: 96.38225424859374, Y: 74.72694631307311},
			b:    Point{X: 96.35022032262698, Y: 74.75022032262696},
			c:    Point{X: 96.35022032262697, Y: 74.75022032262696},
			d:    Point{X: 96.25725424859374, Y: 74.81776412907378},
			hit:  true,
			hold: true,
			at:   Point{X: 96.35022032262697, Y: 74.75022032262696},
		},
		{
			a:    Point{X: 92.0432020322183, Y: 72.27706055336469},
			b:    Point{X: 92.09680528764137, Y: 72.42433428724885},
			c:    Point{X: 91.44, Y: 72.33},
			d:    Point{X: 96.1048, Y: 72.33},
			hit:  true,
			hold: true,
			at:   Point{X: 92.062470, Y: 72.33},
		},
	}
	for i, v := range vs {
		hit, left, hold, at := intersect(v.a, v.b, v.c, v.d)
		if hit != v.hit {
			t.Fatalf("test=%d: hit got=%v, want=%v", i, hit, v.hit)
		}
		if !hit {
			continue
		}
		if hold != v.hold {
			t.Errorf("TODO test=%d: hold got=%v, want=%v", i, hold, v.hold)
		}
		if left != v.left {
			t.Errorf("test=%d: left got=%v, want=%v", i, left, v.left)
		}
		if !MatchPoint(v.at, at) {
			t.Errorf("test=%d got=%v, want=%v", i, at, v.at)
		}
	}
}

func TestTrace(t *testing.T) {
	pts := [][]Point{
		[]Point{{X: 97.35, Y: 68.58}, {X: 97.23, Y: 68.80}, {X: 96.98, Y: 68.80}, {X: 96.85, Y: 68.58}, {X: 96.98, Y: 68.36}, {X: 97.23, Y: 68.36}},
		[]Point{{X: 91.44, Y: 68.33}, {X: 97.10, Y: 68.33}, {X: 97.10, Y: 68.83}, {X: 91.44, Y: 68.83}},
		[]Point{{X: 91.69, Y: 68.58}, {X: 91.57, Y: 68.80}, {X: 91.32, Y: 68.80}, {X: 91.19, Y: 68.58}, {X: 91.32, Y: 68.36}, {X: 91.57, Y: 68.36}},
	}
	var s *Shapes
	for i, ps := range pts {
		var err error
		s, err = s.Append(ps...)
		if err != nil {
			t.Fatalf("shape=%d failed import: %v", i, err)
		}
		s.Union()
		if len(s.P) != 1 {
			t.Fatalf("after shape=%d unioned, got=%d shapes, want=1: %v", i, len(s.P), s.P)
		}
	}
	us := s.P[0].PS
	expect := []Point{
		{X: 91.19, Y: 68.58},
		{X: 91.32, Y: 68.36},
		{X: 91.44, Y: 68.36},
		{X: 91.44, Y: 68.33},
		{X: 97.10, Y: 68.33},
		{X: 97.10, Y: 68.36},
		{X: 97.23, Y: 68.36},
		{X: 97.35, Y: 68.58},
		{X: 97.23, Y: 68.80},
		{X: 97.10, Y: 68.80},
		{X: 97.10, Y: 68.83},
		{X: 91.44, Y: 68.83},
		{X: 91.44, Y: 68.80},
		{X: 91.32, Y: 68.80},
	}
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points: got=%v, want=%v", len(expect), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; !MatchPoint(got, want) {
			t.Errorf("point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	pts = [][]Point{
		[]Point{
			{X: 91.190000, Y: 74.580000},
			{X: 91.237746, Y: 74.433054},
			{X: 91.362746, Y: 74.342236},
			{X: 91.440000, Y: 74.342236},
			{X: 91.440000, Y: 74.330000},
			{X: 96.180000, Y: 74.330000},
			{X: 96.180000, Y: 74.342236},
			{X: 96.257254, Y: 74.342236},
			{X: 96.382254, Y: 74.433054},
			{X: 96.430000, Y: 74.580000},
			{X: 96.382254, Y: 74.726946},
			{X: 96.257254, Y: 74.817764},
			{X: 96.180000, Y: 74.817764},
			{X: 96.180000, Y: 74.830000},
			{X: 91.440000, Y: 74.830000},
			{X: 91.440000, Y: 74.817764},
			{X: 91.362746, Y: 74.817764},
			{X: 91.237746, Y: 74.726946},
		},
		[]Point{
			{X: 95.930000, Y: 74.580000},
			{X: 95.977746, Y: 74.433054},
			{X: 96.009780, Y: 74.409780},
			{X: 96.003223, Y: 74.403223},
			{X: 96.343223, Y: 74.063223},
			{X: 96.349780, Y: 74.069780},
			{X: 96.442746, Y: 74.002236},
			{X: 96.597254, Y: 74.002236},
			{X: 96.722254, Y: 74.093054},
			{X: 96.770000, Y: 74.240000},
			{X: 96.722254, Y: 74.386946},
			{X: 96.690220, Y: 74.410220},
			{X: 96.696777, Y: 74.416777},
			{X: 96.356777, Y: 74.756777},
			{X: 96.350220, Y: 74.750220},
			{X: 96.257254, Y: 74.817764},
			{X: 96.102746, Y: 74.817764},
			{X: 95.977746, Y: 74.726946},
		},
	}
	s = nil
	for i, ps := range pts {
		var err error
		s, err = s.Append(ps...)
		if err != nil {
			t.Fatalf("shape=%d failed import: %v", i, err)
		}
		s.Union()
		if len(s.P) != 1 {
			t.Fatalf("after shape=%d unioned, got=%d shapes, want=1: %v", i, len(s.P), s.P)
		}
	}
	expect = []Point{
		{X: 91.19, Y: 74.58},
		{X: 91.237746, Y: 74.433054},
		{X: 91.362746, Y: 74.342236},
		{X: 91.44, Y: 74.33},
		{X: 96.07644599999999, Y: 74.33},
		{X: 96.343223, Y: 74.063223},
		{X: 96.34978, Y: 74.06978},
		{X: 96.442746, Y: 74.002236},
		{X: 96.597254, Y: 74.002236},
		{X: 96.722254, Y: 74.093054},
		{X: 96.77, Y: 74.24},
		{X: 96.722254, Y: 74.386946},
		{X: 96.696777, Y: 74.416777},
		{X: 96.356777, Y: 74.756777},
		{X: 96.35022006399838, Y: 74.75022006399838},
		{X: 96.257254, Y: 74.817764},
		{X: 96.18, Y: 74.83},
		{X: 91.44, Y: 74.83},
		{X: 91.362746, Y: 74.817764},
		{X: 91.237746, Y: 74.726946},
	}
	us = s.P[0].PS
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points, but see=%d points: got=%#v, want=%#v", len(expect), len(us), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; !MatchPoint(got, want) {
			t.Errorf("point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	pts = [][]Point{
		[]Point{
			{X: 90.7695641085242, Y: 72.50163728296546},
			{X: 90.80570748096952, Y: 72.34913640325517},
			{X: 90.87604572729627, Y: 72.2090814398022},
			{X: 90.9767868944386, Y: 72.08902279193819},
			{X: 91.10249999999999, Y: 71.9954328524455},
			{X: 91.24640781792002, Y: 71.93335707918705},
			{X: 91.40075224048543, Y: 71.9061419931669},
			{X: 91.55721251992517, Y: 71.91525476671676},
			{X: 91.70735384207643, Y: 71.96020412785582},
			{X: 91.84308204939938, Y: 72.03856684489034},
			{X: 91.9570799991053, Y: 72.14611836346158},
			{X: 92.0432020322183, Y: 72.27706055336469},
			{X: 92.09680528764137, Y: 72.42433428724885},
			{X: 92.115, Y: 72.58},
			{X: 92.09680528764137, Y: 72.73566571275114},
			{X: 92.0432020322183, Y: 72.8829394466353},
			{X: 91.9570799991053, Y: 73.01388163653841},
			{X: 91.84308204939938, Y: 73.12143315510966},
			{X: 91.70735384207643, Y: 73.19979587214418},
			{X: 91.55721251992517, Y: 73.24474523328324},
			{X: 91.40075224048543, Y: 73.2538580068331},
			{X: 91.24640781792002, Y: 73.22664292081295},
			{X: 91.10249999999999, Y: 73.1645671475545},
			{X: 90.9767868944386, Y: 73.0709772080618},
			{X: 90.87604572729627, Y: 72.95091856019779},
			{X: 90.80570748096952, Y: 72.81086359674482},
			{X: 90.7695641085242, Y: 72.65836271703454},
		},
		[]Point{
			{X: 91.19, Y: 72.58},
			{X: 91.23774575140627, Y: 72.43305368692688},
			{X: 91.36274575140627, Y: 72.34223587092622},
			{X: 91.44, Y: 72.3422358709262},
			{X: 91.44, Y: 72.33},
			{X: 96.1048, Y: 72.33},
			{X: 96.1048, Y: 72.34223587092622},
			{X: 96.18205424859373, Y: 72.3422358709262},
			{X: 96.30705424859373, Y: 72.43305368692688},
			{X: 96.3548, Y: 72.58},
			{X: 96.30705424859373, Y: 72.72694631307311},
			{X: 96.18205424859373, Y: 72.81776412907378},
			{X: 96.1048, Y: 72.8177641290738},
			{X: 96.1048, Y: 72.83},
			{X: 91.44, Y: 72.83},
			{X: 91.44, Y: 72.81776412907378},
			{X: 91.36274575140627, Y: 72.8177641290738},
			{X: 91.23774575140627, Y: 72.72694631307311},
		},
	}
	s = nil
	for i, ps := range pts {
		var err error
		s, err = s.Append(ps...)
		if err != nil {
			t.Fatalf("shape=%d failed import: %v", i, err)
		}
		s.Union()
		if len(s.P) != 1 {
			t.Fatalf("after shape=%d unioned, got=%d shapes, want=1: %v", i, len(s.P), s.P)
		}
	}
	expect = []Point{
		{X: 90.7695641085242, Y: 72.50163728296546},
		{X: 90.80570748096952, Y: 72.34913640325517},
		{X: 90.87604572729627, Y: 72.2090814398022},
		{X: 90.9767868944386, Y: 72.08902279193819},
		{X: 91.10249999999999, Y: 71.9954328524455},
		{X: 91.24640781792002, Y: 71.93335707918705},
		{X: 91.40075224048543, Y: 71.9061419931669},
		{X: 91.55721251992517, Y: 71.91525476671676},
		{X: 91.70735384207643, Y: 71.96020412785582},
		{X: 91.84308204939938, Y: 72.03856684489034},
		{X: 91.9570799991053, Y: 72.14611836346158},
		{X: 92.0432020322183, Y: 72.27706055336469},
		{X: 92.06247041501207, Y: 72.32999999999998},
		{X: 96.1048, Y: 72.33},
		{X: 96.18205424859373, Y: 72.3422358709262},
		{X: 96.30705424859373, Y: 72.43305368692688},
		{X: 96.3548, Y: 72.58},
		{X: 96.30705424859373, Y: 72.72694631307311},
		{X: 96.18205424859373, Y: 72.81776412907378},
		{X: 96.1048, Y: 72.83},
		{X: 92.06247041501207, Y: 72.82999999999998},
		{X: 92.0432020322183, Y: 72.8829394466353},
		{X: 91.9570799991053, Y: 73.01388163653841},
		{X: 91.84308204939938, Y: 73.12143315510966},
		{X: 91.70735384207643, Y: 73.19979587214418},
		{X: 91.55721251992517, Y: 73.24474523328324},
		{X: 91.40075224048543, Y: 73.2538580068331},
		{X: 91.24640781792002, Y: 73.22664292081295},
		{X: 91.10249999999999, Y: 73.1645671475545},
		{X: 90.9767868944386, Y: 73.0709772080618},
		{X: 90.87604572729627, Y: 72.95091856019779},
		{X: 90.80570748096952, Y: 72.81086359674482},
		{X: 90.7695641085242, Y: 72.65836271703454},
	}
	us = s.P[0].PS
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points, but see=%d points: got=%#v, want=%#v", len(expect), len(us), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; !MatchPoint(got, want) {
			t.Errorf("point[%d]: got=%v, want=%v", i, got, want)
		}
	}

	pts = [][]Point{
		[]Point{
			{X: 90.765, Y: 67.905},
			{X: 92.115, Y: 67.905},
			{X: 92.115, Y: 68.33},
			{X: 97.1, Y: 68.33},
			{X: 97.1, Y: 68.34223587092622},
			{X: 97.17725424859373, Y: 68.3422358709262},
			{X: 97.30225424859373, Y: 68.43305368692688},
			{X: 97.35, Y: 68.58},
			{X: 97.30225424859373, Y: 68.72694631307311},
			{X: 97.17725424859373, Y: 68.81776412907378},
			{X: 97.1, Y: 68.8177641290738},
			{X: 97.1, Y: 68.83},
			{X: 92.115, Y: 68.83},
			{X: 92.115, Y: 69.255},
			{X: 90.765, Y: 69.255},
		},
		[]Point{
			{X: 96.85, Y: 68.58},
			{X: 96.89774575140626, Y: 68.43305368692688},
			{X: 97.02274575140626, Y: 68.34223587092622},
			{X: 97.17725424859373, Y: 68.3422358709262},
			{X: 97.27022032262697, Y: 68.40977967737305},
			{X: 97.27677669529663, Y: 68.40322330470336},
			{X: 97.91177669529664, Y: 69.03822330470337},
			{X: 97.90522032262695, Y: 69.04477967737304},
			{X: 97.93725424859373, Y: 69.06805368692689},
			{X: 97.985, Y: 69.215},
			{X: 97.93725424859373, Y: 69.36194631307312},
			{X: 97.81225424859373, Y: 69.45276412907378},
			{X: 97.65774575140627, Y: 69.45276412907378},
			{X: 97.56477967737304, Y: 69.38522032262695},
			{X: 97.55822330470336, Y: 69.39177669529664},
			{X: 96.92322330470336, Y: 68.75677669529664},
			{X: 96.92977967737303, Y: 68.75022032262696},
			{X: 96.89774575140626, Y: 68.72694631307311},
		},
	}
	s = nil
	for i, ps := range pts {
		var err error
		s, err = s.Append(ps...)
		if err != nil {
			t.Fatalf("shape=%d failed import: %v", i, err)
		}
		s.Union()
		if len(s.P) != 1 {
			t.Fatalf("after shape=%d unioned, got=%d shapes, want=1: %v", i, len(s.P), s.P)
		}
	}
	expect = []Point{
		{X: 90.765, Y: 67.905},
		{X: 92.115, Y: 67.905},
		{X: 92.115, Y: 68.33},
		{X: 97.1, Y: 68.33},
		{X: 97.17725424859373, Y: 68.3422358709262},
		{X: 97.27022032262697, Y: 68.40977967737305},
		{X: 97.27677669529663, Y: 68.40322330470336},
		{X: 97.91177669529664, Y: 69.03822330470337},
		{X: 97.93725424859373, Y: 69.06805368692689},
		{X: 97.985, Y: 69.215},
		{X: 97.93725424859373, Y: 69.36194631307312},
		{X: 97.81225424859373, Y: 69.45276412907378},
		{X: 97.65774575140627, Y: 69.45276412907378},
		{X: 97.56477967737304, Y: 69.38522032262695},
		{X: 97.55822330470336, Y: 69.39177669529664},
		{X: 96.99644660940672, Y: 68.83},
		{X: 92.115, Y: 68.83},
		{X: 92.115, Y: 69.255},
		{X: 90.765, Y: 69.255},
	}
	us = s.P[0].PS
	if len(us) != len(expect) {
		t.Fatalf("expecting %d post union points, but see=%d points: got=%#v, want=%#v", len(expect), len(us), us, expect)
	}
	for i, got := range us {
		if want := expect[i]; !MatchPoint(got, want) {
			t.Errorf("point[%d]: got=%v, want=%v", i, got, want)
		}
	}
}
