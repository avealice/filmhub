package model

type Actor struct {
	ID        int    `json:"-"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type ActorWithMovies struct {
	ID        int     `json:"-"`
	Name      string  `json:"name"`
	Gender    string  `json:"gender"`
	BirthDate string  `json:"birth_date"`
	Movies    []Movie `json:"movies"`
}
