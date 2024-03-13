package model

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
}

type ActorWithMovies struct {
	Actor  Actor   `json:"actor"`
	Movies []Movie `json:"movies"`
}
