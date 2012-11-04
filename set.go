// Copyright 2012 Kamil Kisiel. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package set provides functions for working with finite sets.
//
// All functions in this package require that the types of their Interface arguments match.
// If the types of the arguments do not match the functions will panic.
// The concrete type of the return value will be the same as that of the inputs.
//
// A type implementing set.Interface must be map, slice, or a struct.
// Other types will result in a panic when trying to call any of the routines expecting an Interface instance.
package set

import (
	"fmt"
	"reflect"
)

// A type, typically a collection, that satisfies set.Interface can be used by the routines in this package.
//
// The concrete type of x is defined by the implementor of the interface. Calling the functions with a value
// of a different type will typically result in a panic.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int

	// Contains returns whether the element x is in the collection.
	Contains(x interface{}) bool

	// Insert inserts a new element in to the collection.
	Insert(x interface{})

	// Remove removes an element from the collection.
	Remove(x interface{})

	// Values returns a channel that produces all of the values in the collection.
	// Implementors must close the channel after sending the last value.
	// Calling Interface.Insert or Interface.Remove or other operations that modify the collection is forbidden
	// until the channel return by Values has been closed, otherwise the result is undefined.
	Values() <-chan interface{}
}

// newSet creates a new set of the same type as s and t, after ensuring they are of the same type.
func newSet(s, t Interface, capacity int) Interface {
	sType := reflect.TypeOf(s)
	tType := reflect.TypeOf(t)
	if sType != tType {
		panic(fmt.Sprintf("Set types %s and %s do not match", sType, tType))
	}

	var r Interface
	switch sType.Kind() {
	case reflect.Map:
		r = reflect.MakeMap(sType).Interface().(Interface)
	case reflect.Slice:
		r = reflect.MakeSlice(sType, 0, capacity).Interface().(Interface)
	case reflect.Struct:
		r = reflect.Zero(sType).Interface().(Interface)
	default:
		panic(fmt.Sprintf("Unsupported set type: %s", sType))
	}
	return r
}

// Union returns a new set containing all the elements of s and t.
func Union(s, t Interface) Interface {
	r := newSet(s, t, s.Len()+t.Len())
	svals := s.Values()
	for v := range svals {
		r.Insert(v)
	}
	tvals := t.Values()
	for v := range tvals {
		r.Insert(v)
	}
	return r
}

// Intersection returns a new set containing the elements that are in both s and t.
func Intersection(s, t Interface) Interface {
	r := newSet(s, t, 0)
	if s.Len() < t.Len() {
		t, s = s, t
	}
	svals := s.Values()
	for v := range svals {
		if t.Contains(v) {
			r.Insert(v)
		}
	}
	return r
}

// Difference returns a new set containing the elements that are in s but not t.
func Difference(s, t Interface) Interface {
	r := newSet(s, t, 0)
	svals := s.Values()
	for v := range svals {
		if !t.Contains(v) {
			r.Insert(v)
		}
	}
	return r
}

// SymmetricDifference returns a new set containing the elements in s that are not in t and the elements in t that are not in s.
func SymmetricDifference(s, t Interface) Interface {
	r := newSet(s, t, 0)
	svals := s.Values()
	for v := range svals {
		if !t.Contains(v) {
			r.Insert(v)
		}
	}
	tvals := t.Values()
	for v := range tvals {
		if !s.Contains(v) {
			r.Insert(v)
		}
	}
	return r
}

// StringSet attaches the methods of Interface to a map[string]bool
type StringSet map[string]bool

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

func (s StringSet) Values() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for v := range s {
			c <- v
		}
		close(c)
	}()
	return c
}

// IntSet attaches the methods of Interface to a map[int]bool
type IntSet map[int]bool

func (s IntSet) Len() int {
	return len(s)
}

func (s IntSet) Contains(x interface{}) bool {
	_, ok := s[x.(int)]
	return ok
}

func (s IntSet) Insert(x interface{}) {
	s[x.(int)] = true
}

func (s IntSet) Remove(x interface{}) {
	delete(s, x.(int))
}

func (s IntSet) Values() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for v := range s {
			c <- v
		}
		close(c)
	}()
	return c
}
