package repository

import (
	"database/sql"
	"errors"
	"filmhub/internal/model"
	"fmt"
	"strings"

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

func (r *ActorPostgres) CreateActor(actor model.Actor) error {
	var existingActorID int
	query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1 AND gender = $2 AND birth_date = $3", actorsTable)
	err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&existingActorID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		return errors.New("actor with the same name, gender, and birth date already exists")
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3)", actorsTable)
	_, err = r.db.Exec(insertQuery, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}

	return nil
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

		err := rows.Scan(&actorID, &actor.Name, &actor.Gender, &actor.BirthDate,
			&movie.ID, &movie.Title, &movie.Description, &releaseDateStr, &movie.Rating)
		if err != nil {
			continue
		}

		movie.ReleaseDate = releaseDateStr

		actorWithMovies, ok := actorMap[actorID]
		if !ok {
			actorWithMovies = model.ActorWithMovies{Name: actor.Name, Gender: actor.Gender, BirthDate: actor.BirthDate, Movies: []model.Movie{}}
			actorMap[actorID] = actorWithMovies
		}

		actorWithMovies.Movies = append(actorWithMovies.Movies, movie)
		actorMap[actorID] = actorWithMovies
	}

	for _, actorWithMovies := range actorMap {
		actorsWithMovies = append(actorsWithMovies, actorWithMovies)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		SELECT a.id, a.name, a.gender, a.birth_date, m.id, m.title, m.description, m.release_date, m.rating
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

	moviesMap := make(map[int]model.Movie)

	for rows.Next() {
		var movie model.Movie
		var movieID int

		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate, &movieID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
		if err != nil {
			continue
		}

		if _, ok := moviesMap[movieID]; !ok {
			moviesMap[movieID] = movie
		}

		actor.Movies = append(actor.Movies, moviesMap[movieID])
	}

	return actor, nil
}

func (r *ActorPostgres) Update(actorID int, data model.ActorWithMovies) error {
	var query string
	var params []interface{}

	query = fmt.Sprintf("UPDATE %s SET ", actorsTable)

	paramIndex := 1

	if data.Name != "" {
		query += fmt.Sprintf("name = $%d, ", paramIndex)
		params = append(params, data.Name)
		paramIndex++
	}

	if data.Gender != "" {
		query += fmt.Sprintf("gender = $%d, ", paramIndex)
		params = append(params, data.Gender)
		paramIndex++
	}

	if data.BirthDate != "" {
		query += fmt.Sprintf("birth_date = $%d, ", paramIndex)
		params = append(params, data.BirthDate)
		paramIndex++
	}

	query = strings.TrimSuffix(query, ", ")

	query += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, actorID)

	_, err := r.db.Exec(query, params...)
	if err != nil {
		return err
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
