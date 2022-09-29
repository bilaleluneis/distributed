/*	@author: Bilal El Uneis
	@since: September 2022
	@email: bilaleluneis@gmail.com	*/

package common

import (
	"log"
)

type logger struct {
	active bool
}

func (l logger) Debug(format string, v ...any) {
	if l.active {
		s := "DEBUG: " + format
		log.Printf(s, v...)
	}
}

func (l logger) Error(format string, v ...any) {
	if l.active {
		s := "ERROR: " + format
		log.Printf(s, v...)
	}
}

// Log by default logger is disabled
var Log = logger{}

func EnableLogger() {
	Log.active = true
}

func DisableLogger() {
	Log.active = false
}
