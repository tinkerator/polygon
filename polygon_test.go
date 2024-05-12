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
		Point{X: 1, Y: 1},
		Point{X: 2, Y: 1},
		Point{X: 2, Y: 2},
		Point{X: 1, Y: 2},
	}...).Builder([]Point{
		Point{X: 0, Y: 0},
		Point{X: 1.5, Y: 0},
		Point{X: 1.5, Y: 1},
		Point{X: 0, Y: 1},
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
		Point{X: 1, Y: 1},
		Point{X: 2, Y: 1},
		Point{X: 2, Y: 2},
		Point{X: 1, Y: 2},
	}...).Builder([]Point{
		Point{X: 2, Y: 0},
		Point{X: 3, Y: 0},
		Point{X: 3, Y: 2},
		Point{X: 2, Y: 2},
	}...)
	ss.Union()
	if len(ss.P) != 1 {
		t.Fatalf("expecting a single poly, but got %d", len(ss.P))
	}
	us = ss.P[0].PS
	expect = []Point{{1, 1}, {2, 1}, {2, 0}, {3, 0}, {3, 2}, {2, 2}, {1, 2}}
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
		Point{X: 1, Y: 1},
		Point{X: 2, Y: 1},
		Point{X: 2, Y: 2},
		Point{X: 1, Y: 2},
	}...).Builder([]Point{
		Point{X: 2, Y: 0},
		Point{X: 3, Y: 0},
		Point{X: 3, Y: 3},
		Point{X: 2, Y: 3},
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
		Point{X: 1, Y: 0},
		Point{X: 2, Y: 0},
		Point{X: 2, Y: 3},
		Point{X: 1, Y: 3},
	}...).Builder([]Point{
		Point{X: 2, Y: 1},
		Point{X: 3, Y: 1},
		Point{X: 3, Y: 2},
		Point{X: 2, Y: 2},
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
}
