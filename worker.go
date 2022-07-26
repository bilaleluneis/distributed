/*	@author: Bilal El Uneis
	@since: July 2022
	@email: bilaleluneis@gmail.com	*/

package distributed

import (
	"distributed/internal"
	"log"
)

func InitWorker(port int) {
	err := internal.InitServer(port)
	if err != nil {
		log.Panic("worker initalization failed")
	}
}
