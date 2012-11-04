package set

import (
	"fmt"
	"reflect"
)

type Interface interface {
	Len() int
	Contains(x interface{}) bool
	Insert(x interface{})
	Remove(x interface{})
}

type StringSet map[string]bool
type stringset map[string]bool

func newSet(s Interface) Interface {
	t := reflect.TypeOf(s)
	switch t.Kind() {
	case reflect.Map:
		r := reflect.MakeMap(t).Interface().(Interface)
		return r
	default:
		panic(fmt.Sprintf("Unsupported set type: %s", t))
	}
	return nil
}

func Union(s, t stringset) stringset {
	r := make(stringset, len(s)+len(t))
	for k := range s {
		r[k] = true
	}
	for k := range t {
		r[k] = true
	}
	return r
}

func Intersection(s, t stringset) stringset {
	r := make(stringset)
	if len(s) < len(t) {
		t, s = s, t
	}
	for k := range s {
		if _, ok := t[k]; ok {
			r[k] = true
		}
	}
	return r
}

func Difference(s, t stringset) stringset {
	r := make(stringset)
	for k := range s {
		if _, ok := t[k]; !ok {
			r[k] = true
		}
	}
	return r
}

func SymmetricDifference(s, t stringset) stringset {
	r := make(stringset)
	for k := range s {
		if _, ok := t[k]; !ok {
			r[k] = true
		}
	}
	for k := range t {
		if _, ok := s[k]; !ok {
			r[k] = true
		}
	}
	return r
}

func (s StringSet) Len() int {
	return len(s)
}

func (s StringSet) Contains(x interface{}) bool {
	_, ok := s[x.(string)]
	return ok
}

func (s StringSet) Insert(x interface{}) {
	s[x.(string)] = true
}

func (s StringSet) Remove(x interface{}) {
	delete(s, x.(string))
}
