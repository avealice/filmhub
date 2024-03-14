package model

type Movie struct {
	ID          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating      int    `json:"rating"`
}

type MovieWithActors struct {
	ID          int     `json:"-"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ReleaseDate string  `json:"release_date"`
	Rating      int     `json:"rating"`
	Actors      []Actor `json:"actors"`
}

type MovieActor struct {
	MovieID int `json:"movie_id"`
	ActorID int `json:"actor_id"`
}
