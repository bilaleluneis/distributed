/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"distributed/internal"
	"log"
)

// InitWorker TODO: need to refactor to allow open port 0.0.0.0
func InitWorker(port int) {
	err := internal.InitServer(port)
	if err != nil {
		log.Panic("worker initalization failed")
	}
}
