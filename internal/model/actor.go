package model

// Actor represents an actor in the system.
type Actor struct {
	ID        int    `json:"-"`          // Unique identifier for the actor
	Name      string `json:"name"`       // Name of the actor
	Gender    string `json:"gender"`     // Gender of the actor
	BirthDate string `json:"birth_date"` // Birth date of the actor
}

// ActorWithMovies represents an actor with associated movies in the system.
type ActorWithMovies struct {
	ID        int     `json:"-"`          // Unique identifier for the actor
	Name      string  `json:"name"`       // Name of the actor
	Gender    string  `json:"gender"`     // Gender of the actor
	BirthDate string  `json:"birth_date"` // Birth date of the actor
	Movies    []Movie `json:"movies"`     // Movies associated with the actor
}