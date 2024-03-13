package model

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating      int    `json:"rating"`
}

type MovieActor struct {
	MovieID int `json:"movie_id"`
	ActorID int `json:"actor_id"`
}

type MovieWithActors struct {
	Movie  Movie   `json:"movie"`
	Actors []Actor `json:"actors"`
}
