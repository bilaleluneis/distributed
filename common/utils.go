/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package common

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
)

// GenUUID TODO: do error checking and return error
func GenUUID() UUID {
	b := make([]byte, 16)
	_, _ = rand.Read(b) //FIXME: error check
	i1 := b[0:4]
	i2 := b[4:6]
	i3 := b[6:8]
	i4 := b[8:10]
	i5 := b[10:]
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", i1, i2, i3, i4, i5)
	return uuid
}

func ToBytes[T any](value T) ([]byte, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(value); err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func ToType[T any](data []byte) (T, error) {
	var buffer bytes.Buffer
	buffer.Write(data)
	var result T
	err := gob.NewDecoder(&buffer).Decode(&result)
	return result, err
}

// Contains check if value exist in slice
// and number of times it appears
func Contains[T any](in []T, val T) (int, bool) {
	var exist bool
	var count int
	if v, err := ToBytes[T](val); err == nil {
		for _, item := range in {
			if i, err := ToBytes[T](item); err == nil {
				if bytes.Equal(i, v) {
					exist = true
					count++
				}
			}
		}
	}
	return count, exist
}
