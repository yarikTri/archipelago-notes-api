package models

type Route struct {
	ID          uint32
	Name        string
	Description string
	ImagePath   string
}

func (r *Route) ToTransfer(stations []Station) RouteTransfer {
	return RouteTransfer{
		ID:          r.ID,
		Name:        r.Name,
		Stations:    stations,
		Description: r.Description,
		ImagePath:   r.ImagePath,
	}
}

type RouteTransfer struct {
	ID          uint32
	Name        string
	Stations    []Station
	Description string
	ImagePath   string
}
