package model

// Movie represents a movie in the system.
type Movie struct {
	ID          int    `json:"-" db:"id"`                      // Unique identifier for the movie
	Title       string `json:"title" db:"title"`               // Title of the movie
	Description string `json:"description" db:"description"`   // Description of the movie
	ReleaseDate string `json:"release_date" db:"release_date"` // Format: "YYYY-M-D".
	Rating      int    `json:"rating" db:"rating"`             // Rating of the movie
}

// MovieWithActors represents a movie with associated actors in the system.
type MovieWithActors struct {
	ID          int     `json:"id" db:"id"`                     // Unique identifier for the movie
	Title       string  `json:"title" db:"title"`               // Title of the movie
	Description string  `json:"description" db:"description"`   // Description of the movie
	ReleaseDate string  `json:"release_date" db:"release_date"` // Format: "YYYY-M-D".
	Rating      int     `json:"rating" db:"rating"`             // Rating of the movie
	Actors      []Actor `json:"actors"`                         // Actors associated with the movie
}

type InputMovie struct {
	Title       string  `json:"title" db:"title"`               // Title of the movie
	Description string  `json:"description" db:"description"`   // Description of the movie
	ReleaseDate string  `json:"release_date" db:"release_date"` // Format: "YYYY-M-D".
	Rating      int     `json:"rating" db:"rating"`             // Rating of the movie
	Actors      []Actor `json:"actors"`                         // Actors associated with the movie
}

// MovieActor represents a relationship between a movie and an actor in the system.
type MovieActor struct {
	MovieID int `json:"movie_id" db:"movie_id"` // ID of the movie
	ActorID int `json:"actor_id" db:"actor_id"` // ID of the actor
}
