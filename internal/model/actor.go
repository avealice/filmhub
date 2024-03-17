package model

// Actor represents an actor in the system.
type Actor struct {
	ID        int    `json:"-" db:"id"`                  // Unique identifier for the actor
	Name      string `json:"name" db:"name"`             // Name of the actor
	Gender    string `json:"gender" db:"gender"`         // Valid values: "male", "female", "other".
	BirthDate string `json:"birth_date" db:"birth_date"` // Format: "YYYY-M-D".
}

// ActorWithMovies represents an actor with associated movies in the system.
type ActorWithMovies struct {
	ID        int     `json:"id" db:"id"`                 // Unique identifier for the actor
	Name      string  `json:"name" db:"name"`             // Name of the actor
	Gender    string  `json:"gender" db:"gender"`         // Valid values: "male", "female", "other".
	BirthDate string  `json:"birth_date" db:"birth_date"` // Format: "YYYY-M-D".
	Movies    []Movie `json:"movies"`                     // Movies associated with the actor
}

type InputActor struct {
	Name      string  `json:"name" db:"name"`             // Name of the actor
	Gender    string  `json:"gender" db:"gender"`         // Valid values: "male", "female", "other".
	BirthDate string  `json:"birth_date" db:"birth_date"` // Format: "YYYY-M-D".
	Movies    []Movie `json:"movies"`                     // Movies associated with the actor
}
