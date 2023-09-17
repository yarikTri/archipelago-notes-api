package models

type Route struct {
	ID   uint32
	Name string
}

type RouteTransfer struct {
	ID       uint32
	Name     string
	Stations []StationOfRouteTransfer
}
