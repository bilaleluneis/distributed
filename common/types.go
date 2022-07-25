/*	@author: Bilal El Uneis
	@since: June 2022
	@email: bilaleluneis@gmail.com	*/

package common

type UUID = string

type Location struct {
	HostName string
	Port     int
	Uuid     UUID
}

type SearchParams struct {
	Remove  bool
	Address Location
}
