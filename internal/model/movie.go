package model

// Movie represents a movie in the system.
type Movie struct {
	ID          int    `json:"-"`            // Unique identifier for the movie
	Title       string `json:"title"`        // Title of the movie
	Description string `json:"description"`  // Description of the movie
	ReleaseDate string `json:"release_date"` // Release date of the movie
	Rating      int    `json:"rating"`       // Rating of the movie
}

// MovieWithActors represents a movie with associated actors in the system.
type MovieWithActors struct {
	ID          int     `json:"-"`            // Unique identifier for the movie
	Title       string  `json:"title"`        // Title of the movie
	Description string  `json:"description"`  // Description of the movie
	ReleaseDate string  `json:"release_date"` // Release date of the movie
	Rating      int     `json:"rating"`       // Rating of the movie
	Actors      []Actor `json:"actors"`       // Actors associated with the movie
}

// MovieActor represents a relationship between a movie and an actor in the system.
type MovieActor struct {
	MovieID int `json:"movie_id"` // ID of the movie
	ActorID int `json:"actor_id"` // ID of the actor
}
