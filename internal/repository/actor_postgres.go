package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/avealice/filmhub/internal/model"

	"github.com/jmoiron/sqlx"
)

type ActorPostgres struct {
	db *sqlx.DB
}

func NewActorPostgres(db *sqlx.DB) *ActorPostgres {
	return &ActorPostgres{
		db: db,
	}
}

// func (r *ActorPostgres) CreateActor(actor model.InputActor) (int, error) {
// 	var existingActorID int
// 	query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1 AND gender = $2 AND birth_date = $3", actorsTable)
// 	err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&existingActorID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return 0, err
// 	}

// 	if err == nil {
// 		return 0, errors.New("actor with the same name, gender, and birth date already exists")
// 	}

// 	insertQuery := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id", actorsTable)
// 	var insertedID int
// 	err = r.db.QueryRow(insertQuery, actor.Name, actor.Gender, actor.BirthDate).Scan(&insertedID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return insertedID, nil
// }

func (r *ActorPostgres) CreateActor(actor model.InputActor) (int, error) {
	var existingActorID int
	query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1 AND gender = $2 AND birth_date = $3", actorsTable)
	err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&existingActorID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == nil {
		return 0, errors.New("actor with the same name, gender, and birth date already exists")
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id", actorsTable)
	var insertedID int
	err = r.db.QueryRow(insertQuery, actor.Name, actor.Gender, actor.BirthDate).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	for _, movie := range actor.Movies {
		var existingMovieID int
		query := fmt.Sprintf("SELECT id FROM %s WHERE title = $1 AND description = $2 AND rating = $3 AND release_date = $4", moviesTable)
		err := r.db.QueryRow(query, movie.Title, movie.Description, movie.Rating, movie.ReleaseDate).Scan(&existingMovieID)
		if err != nil && err != sql.ErrNoRows {
			return 0, err
		}

		if err == nil {
			continue
		}

		movieQuery := fmt.Sprintf("INSERT INTO %s (title, description, rating, release_date) VALUES ($1, $2, $3, $4) RETURNING id", moviesTable)
		var movieID int
		err = r.db.QueryRow(movieQuery, movie.Title, movie.Description, movie.Rating, movie.ReleaseDate).Scan(&movieID)
		if err != nil {
			return 0, err
		}

		query = fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", movieActorTable)
		_, err = r.db.Exec(query, movieID, insertedID)
		if err != nil {
			return 0, err
		}
	}

	return insertedID, nil
}

func (r *ActorPostgres) GetAllActors() ([]model.ActorWithMovies, error) {
	var actorsWithMovies []model.ActorWithMovies

	query := fmt.Sprintf(`
		SELECT a.id, a.name, a.gender, a.birth_date,
			   m.id AS movie_id, m.title AS movie_title, m.description AS movie_description,
			   TO_CHAR(m.release_date, 'YYYY-MM-DD') AS movie_release_date, m.rating AS movie_rating
		FROM %s a
		LEFT JOIN %s ma ON a.id = ma.actor_id
		LEFT JOIN %s m ON ma.movie_id = m.id
	`, actorsTable, movieActorTable, moviesTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actorMap := make(map[int]model.ActorWithMovies)

	for rows.Next() {
		var actorID int
		var actor model.Actor
		var movie model.Movie
		var releaseDateStr string

		rows.Scan(&actorID, &actor.Name, &actor.Gender, &actor.BirthDate,
			&movie.ID, &movie.Title, &movie.Description, &releaseDateStr, &movie.Rating)

		actorWithMovies, ok := actorMap[actorID]
		if !ok {
			actorWithMovies = model.ActorWithMovies{
				ID:        actorID,
				Name:      actor.Name,
				Gender:    actor.Gender,
				BirthDate: actor.BirthDate,
				Movies:    []model.Movie{},
			}
		}

		if movie.ID != 0 {
			movie.ReleaseDate = releaseDateStr
			actorWithMovies.Movies = append(actorWithMovies.Movies, movie)
		}

		actorMap[actorID] = actorWithMovies
	}

	for _, actorWithMovies := range actorMap {
		actorsWithMovies = append(actorsWithMovies, actorWithMovies)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(actorsWithMovies) == 0 {
		return nil, errors.New("movies not found")
	}

	return actorsWithMovies, nil
}

func (r *ActorPostgres) Delete(actorID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", actorsTable)
	_, err := r.db.Exec(query, actorID)
	if err != nil {
		return err
	}
	return nil
}
func (r *ActorPostgres) Get(actorID int) (model.ActorWithMovies, error) {
	var actor model.ActorWithMovies

	query := fmt.Sprintf(`
        SELECT a.id, a.name, a.gender, a.birth_date, m.id, m.title, m.description, TO_CHAR(m.release_date, 'YYYY-MM-DD'), m.rating
        FROM %s a
        LEFT JOIN %s ma ON a.id = ma.actor_id
        LEFT JOIN %s m ON ma.movie_id = m.id
        WHERE a.id = $1
    `, actorsTable, movieActorTable, moviesTable)

	rows, err := r.db.Query(query, actorID)
	if err != nil {
		return actor, err
	}
	defer rows.Close()

	var actorFound bool

	for rows.Next() {
		actorFound = true
		var movie model.Movie
		var movieID int
		var releaseDateStr string

		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate, &movieID, &movie.Title, &movie.Description, &releaseDateStr, &movie.Rating)
		if err != nil {
			continue
		}

		movie.ID = movieID
		movie.ReleaseDate = releaseDateStr
		actor.Movies = append(actor.Movies, movie)
	}

	if !actorFound {
		return actor, errors.New("actor not found")
	}

	return actor, nil
}

func (r *ActorPostgres) Update(actorID int, data model.InputActor) error {
	var query strings.Builder
	var params []interface{}

	query.WriteString("UPDATE ")
	query.WriteString(actorsTable)
	query.WriteString(" SET ")

	paramIndex := 1

	if data.Name != "" {
		query.WriteString(fmt.Sprintf("name = $%d, ", paramIndex))
		params = append(params, data.Name)
		paramIndex++
	}

	if data.Gender != "" {
		query.WriteString(fmt.Sprintf("gender = $%d, ", paramIndex))
		params = append(params, data.Gender)
		paramIndex++
	}

	if data.BirthDate != "" {
		query.WriteString(fmt.Sprintf("birth_date = $%d, ", paramIndex))
		params = append(params, data.BirthDate)
		paramIndex++
	}

	// Remove trailing comma and space
	if query.Len() > len("UPDATE "+actorsTable+" SET ") {
		query.WriteString(fmt.Sprintf(" WHERE id = $%d", paramIndex))
		params = append(params, actorID)

		_, err := r.db.Exec(query.String(), params...)
		if err != nil {
			return err
		}
	}

	for _, movie := range data.Movies {
		var existingMovieID int
		query := fmt.Sprintf("SELECT id FROM %s WHERE title = $1 AND description = $2 AND release_date = $3 AND rating = $4", moviesTable)
		err := r.db.QueryRow(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&existingMovieID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == nil {
			continue
		}

		movieQuery := fmt.Sprintf("INSERT INTO %s (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id", moviesTable)
		err = r.db.QueryRow(movieQuery, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&movie.ID)
		if err != nil {
			return err
		}

		movieActorQuery := fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", movieActorTable)
		_, err = r.db.Exec(movieActorQuery, movie.ID, actorID)
		if err != nil {
			return err
		}
	}

	return nil
}
