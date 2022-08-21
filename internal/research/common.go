/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package research

import (
	"bytes"
	"encoding/gob"
	"errors"
)

/*
types and functions used in understanding the inner
working of type conversion locally and across network
*/

// Common Data Structures used in Tests

type C struct {
	I int
	j string
}

func (c C) getJ() string { return c.j }

type Node[T any] struct {
	Value T // fields must be exported to work with GOB
}

func (n *Node[T]) set(v T) {
	n.Value = v
}

func (n *Node[T]) get() T {
	return n.Value
}

// used with local Generics / Constraint based Tests
var data any

func insertData[T any](val T) {
	data = val
}

func retrieveData[T any]() (T, error) {
	result, ok := data.(T)
	if !ok {
		return result, errors.New("error")
	}
	return result, nil
}

// used with GOB/RPC based conversion Tests
var network bytes.Buffer

func toNetwork[T any](data T) error {
	err := gob.NewEncoder(&network).Encode(data)
	return err
}

func fromNetwork[T any]() (T, error) {
	var result T
	err := gob.NewDecoder(&network).Decode(&result)
	return result, err
}

// Constraints as Interfaces

type TestNodeConstraint[T any] interface {
	set(T)
	get() T
}

type TestInterface interface {
	TestMethod() bool
}

type TestImpl struct {
	Flag bool // fields must be Exported to work with GOB
}

func (tfvi TestImpl) TestMethod() bool { return true }
