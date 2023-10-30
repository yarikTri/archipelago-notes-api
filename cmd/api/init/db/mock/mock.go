package mock

import "github.com/yarikTri/web-transport-cards/internal/models"

type MockDB struct {
	Routes        map[string]models.Route
	Stations      map[string]models.Station
	RoutesStation map[string][]string
	Tickets       map[string]models.Ticket
}

var MockDBImpl = MockDB{
	Routes: map[string]models.Route{
		"1": {
			ID:          1,
			Name:        "Автобус № 12",
			Description: "Автобус, позволяющий из окна насладиться всеми главными городскими достопримечательностями",
			ImagePath:   "/static/image/12.jpeg",
		},
		"2": {
			ID:          2,
			Name:        "Троллейбус № 6",
			Description: "От Мариевки до Ворошилова с божьей помощью",
			ImagePath:   "/static/image/6.jpeg",
		},
		"3": {
			ID:          3,
			Name:        "Автобус № 3",
			Description: "Комфортный всегда пустой автобус, едущий из Копаево до Переборов",
			ImagePath:   "/static/image/3.jpeg",
		},
	},
	Stations: map[string]models.Station{
		"1": {
			ID:   1,
			Name: "Магазин \"Восток\"",
		},
		"2": {
			ID:   2,
			Name: "Улица Гагарина",
		},
		"3": {
			ID:   3,
			Name: "Улица Труда",
		},
		"4": {
			ID:   4,
			Name: "Железнодорожный вокзал",
		},
		"5": {
			ID:   5,
			Name: "Сенная площадь",
		},
		"6": {
			ID:   6,
			Name: "Универмаг \"Юбилейный\"",
		},
		"7": {
			ID:   7,
			Name: "Троллейбусный парк",
		},
		"8": {
			ID:   8,
			Name: "Улица Ворошилова",
		},
		"9": {
			ID:   9,
			Name: "Школа № 27",
		},
		"10": {
			ID:   10,
			Name: "Улица Расторгуева",
		},
		"11": {
			ID:   11,
			Name: "Улица Черепанова",
		},
		"12": {
			ID:   12,
			Name: "Завод \"Призма\"",
		},
		"13": {
			ID:   13,
			Name: "Школа № 20",
		},
		"14": {
			ID:   14,
			Name: "Проспект Батова",
		},
	},
	RoutesStation: map[string][]string{
		"1": {"1", "2", "3", "4", "5", "6"},
		"2": {"12", "13", "14"},
		"3": {"7", "8", "9", "10", "11"},
	},
	Tickets: map[string]models.Ticket{},
}
