package models

type Station struct {
	ID   uint32
	Name string
}

type StationTransfer struct {
	ID   uint32
	Name string
}

type StationOfRouteTransfer struct {
	Station   StationTransfer
	SeqNumber uint32
}
