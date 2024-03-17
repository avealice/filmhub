package repository

import (
	"database/sql"
	"errors"
	"filmhub/internal/model"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MoviePostgres struct {
	db *sqlx.DB
}

func NewMoviePostgres(db *sqlx.DB) *MoviePostgres {
	return &MoviePostgres{
		db: db,
	}
}

func (r *MoviePostgres) GetAllMovies(sortBy, sortOrder string) ([]model.Movie, error) {
	var movies []model.Movie

	query := fmt.Sprintf("SELECT id, title, description, TO_CHAR(release_date, 'YYYY-MM-DD'), rating FROM %s ORDER BY %s %s", moviesTable, sortBy, sortOrder)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie model.Movie
		var releaseDateStr string
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &releaseDateStr, &movie.Rating)
		if err != nil {
			return nil, err
		}
		movie.ReleaseDate = releaseDateStr
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(movies) == 0 {
		return nil, errors.New("movies not found")
	}

	return movies, nil
}

func (r *MoviePostgres) CreateMovie(movie model.InputMovie) error {
	var existingMovieID int
	query := fmt.Sprintf("SELECT id FROM %s WHERE title = $1 AND description = $2 AND rating = $3 AND release_date = $4", moviesTable)
	err := r.db.QueryRow(query, movie.Title, movie.Description, movie.Rating, movie.ReleaseDate).Scan(&existingMovieID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		return errors.New("movie with the same title, description, rating, and release date already exists")
	}

	movieQuery := fmt.Sprintf("INSERT INTO %s (title, description, rating, release_date) VALUES ($1, $2, $3, $4) RETURNING id", moviesTable)
	var movieID int
	err = r.db.QueryRow(movieQuery, movie.Title, movie.Description, movie.Rating, movie.ReleaseDate).Scan(&movieID)
	if err != nil {
		return err
	}

	for _, actor := range movie.Actors {
		var actorID int
		query := fmt.Sprintf("SELECT id FROM %s WHERE LOWER(name) = LOWER($1) AND LOWER(gender) = LOWER($2) AND birth_date = $3", actorsTable)
		err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&actorID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			query := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id", actorsTable)
			err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&actorID)
			if err != nil {
				return err
			}
		}

		query = fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", movieActorTable)
		_, err = r.db.Exec(query, movieID, actorID)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *MoviePostgres) GetMovieByID(movieID int) (model.MovieWithActors, error) {
	var movie model.MovieWithActors

	query := fmt.Sprintf(`
        SELECT m.id, m.title, m.description, TO_CHAR(m.release_date, 'YYYY-MM-DD'), m.rating, a.id, a.name, a.gender, TO_CHAR(a.birth_date, 'YYYY-MM-DD')
        FROM %s m
        LEFT JOIN %s ma ON m.id = ma.movie_id
        LEFT JOIN %s a ON ma.actor_id = a.id
        WHERE m.id = $1
    `, moviesTable, movieActorTable, actorsTable)

	rows, err := r.db.Query(query, movieID)
	if err != nil {
		return movie, err
	}
	defer rows.Close()

	var movieFound bool
	actorsMap := make(map[int]model.Actor)

	for rows.Next() {
		movieFound = true
		var actor model.Actor
		var actorID int
		var birthDateStr string

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &actorID, &actor.Name, &actor.Gender, &birthDateStr)
		if err != nil {
			continue
		}

		if _, ok := actorsMap[actorID]; !ok {
			actorsMap[actorID] = actor
			actor.ID = actorID
			actor.BirthDate = birthDateStr
			movie.Actors = append(movie.Actors, actor)
		}
	}

	if !movieFound {
		return movie, errors.New("movie not found")
	}

	return movie, nil
}

func (r *MoviePostgres) DeleteByID(movieID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", moviesTable)
	_, err := r.db.Exec(query, movieID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MoviePostgres) UpdateMovie(movieID int, data model.InputMovie) error {
	updateQuery := fmt.Sprintf("UPDATE %s SET ", moviesTable)
	var params []interface{}

	paramIndex := 1

	if data.Title != "" {
		updateQuery += fmt.Sprintf("title = $%d, ", paramIndex)
		params = append(params, data.Title)
		paramIndex++
	}

	if data.Description != "" {
		updateQuery += fmt.Sprintf("description = $%d, ", paramIndex)
		params = append(params, data.Description)
		paramIndex++
	}

	if data.ReleaseDate != "" {
		updateQuery += fmt.Sprintf("release_date = $%d, ", paramIndex)
		params = append(params, data.ReleaseDate)
		paramIndex++
	}

	if data.Rating != 0 {
		updateQuery += fmt.Sprintf("rating = $%d, ", paramIndex)
		params = append(params, data.Rating)
		paramIndex++
	}

	updateQuery = strings.TrimSuffix(updateQuery, ", ") + fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, movieID)

	_, err := r.db.Exec(updateQuery, params...)
	if err != nil {
		return err
	}

	for _, actor := range data.Actors {
		var actorID int
		query := fmt.Sprintf("SELECT id FROM %s WHERE LOWER(name) = LOWER($1) AND LOWER(gender) = LOWER($2) AND birth_date = $3", actorsTable)
		err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&actorID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if err == sql.ErrNoRows {
			query := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id", actorsTable)
			err := r.db.QueryRow(query, actor.Name, actor.Gender, actor.BirthDate).Scan(&actorID)
			if err != nil {
				return err
			}
		}

		query = fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", movieActorTable)
		_, err = r.db.Exec(query, movieID, actorID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MoviePostgres) GetMoviesByTitle(titleFragment string) ([]model.MovieWithActors, error) {
	var moviesMap = make(map[int]model.MovieWithActors)
	var movies []model.MovieWithActors

	query := `
            SELECT m.id, m.title, m.description, TO_CHAR(m.release_date, 'YYYY-MM-DD') as release_date, m.rating, a.id, a.name, a.gender, TO_CHAR(a.birth_date, 'YYYY-MM-DD') as birth_date
            FROM movie m
            LEFT JOIN movie_actor ma ON m.id = ma.movie_id
            LEFT JOIN actor a ON ma.actor_id = a.id
            WHERE m.title ILIKE '%' || $1 || '%'
        `
	rows, err := r.db.Query(query, titleFragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := false

	for rows.Next() {
		var movie model.MovieWithActors
		var actor model.Actor
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			continue
		}

		found = true

		if existingMovie, ok := moviesMap[movie.ID]; ok {
			existingMovie.Actors = append(existingMovie.Actors, actor)
			moviesMap[movie.ID] = existingMovie
		} else {
			movie.Actors = append(movie.Actors, actor)
			moviesMap[movie.ID] = movie
		}
	}

	if !found {
		return nil, fmt.Errorf("no movies found with title fragment: %s", titleFragment)
	}

	for _, movie := range moviesMap {
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *MoviePostgres) GetMoviesByActor(actorNameFragment string) ([]model.MovieWithActors, error) {
	var moviesMap = make(map[int]model.MovieWithActors)
	var movies []model.MovieWithActors

	query := `
            SELECT m.id, m.title, m.description, TO_CHAR(m.release_date, 'YYYY-MM-DD') as release_date, m.rating, a.id, a.name, a.gender, TO_CHAR(a.birth_date, 'YYYY-MM-DD') as birth_date
            FROM movie m
            LEFT JOIN movie_actor ma ON m.id = ma.movie_id
            LEFT JOIN actor a ON ma.actor_id = a.id
            WHERE a.name ILIKE '%' || $1 || '%'
        `
	rows, err := r.db.Query(query, actorNameFragment)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var movie model.MovieWithActors
		var actor model.Actor
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating, &actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
		if err != nil {
			continue
		}

		if existingMovie, ok := moviesMap[movie.ID]; ok {
			existingMovie.Actors = append(existingMovie.Actors, actor)
			moviesMap[movie.ID] = existingMovie
		} else {
			movie.Actors = append(movie.Actors, actor)
			moviesMap[movie.ID] = movie
		}
	}

	if len(moviesMap) == 0 {
		return nil, fmt.Errorf("no movies found for actor with name fragment: %s", actorNameFragment)
	}

	for _, movie := range moviesMap {
		movies = append(movies, movie)
	}

	return movies, nil
}
