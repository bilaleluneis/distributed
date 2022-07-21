/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package common

import (
	"crypto/rand"
	"fmt"
)

func GenUUID() UUID {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	i1 := b[0:4]
	i2 := b[4:6]
	i3 := b[6:8]
	i4 := b[8:10]
	i5 := b[10:]
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", i1, i2, i3, i4, i5)
	return uuid
}
