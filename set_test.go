package set

import (
	"testing"
)

func checkSet(s Interface, vals []interface{}, t *testing.T) {
	// Check that Len() returns the right thing
	if l := s.Len(); l != len(vals) {
		t.Errorf("expected len == 3, got %d", l)
	}
	for _, v := range vals {
		if !s.Contains(v) {
			t.Errorf("set does not contains %v", v)
		}
	}
}

func TestSet(t *testing.T) {
	a := make(StringSet)
	a.Insert("a")
	a.Insert("b")
	a.Insert("c")
	b := make(StringSet)
	b.Insert("b")
	b.Insert("c")
	b.Insert("d")

	u := Union(a, b)
	uvals := []interface{}{"a", "b", "c", "d"}
	checkSet(u, uvals, t)

	i := Intersection(a, b)
	ivals := []interface{}{"b", "c"}
	checkSet(i, ivals, t)

	d := Difference(a, b)
	dvals := []interface{}{"a"}
	checkSet(d, dvals, t)

	sd := SymmetricDifference(a, b)
	sdvals := []interface{}{"a", "d"}
	checkSet(sd, sdvals, t)
}
