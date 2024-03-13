package repository

import (
	"filmhub/internal/model"
	"fmt"

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

func (r *MoviePostgres) GetAllMovies(sortBy, sortOrder string) ([]model.MovieWithActors, error) {
	var movies []model.MovieWithActors
	query := fmt.Sprintf(`
	SELECT m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, 
		   TO_CHAR(m.release_date, 'YYYY-MM-DD') AS movie_release_date, m.rating AS movie_rating, 
		   a.id AS actor_id, a.name AS actor_name, a.gender AS actor_gender, 
		   TO_CHAR(a.birth_date, 'YYYY-MM-DD') AS actor_birth_date 
	FROM %s m 
	LEFT JOIN %s ma ON m.id = ma.movie_id 
	LEFT JOIN %s a ON ma.actor_id = a.id 
	ORDER BY %s %s
`, moviesTable, movieActorTable, actorsTable, sortBy, sortOrder)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieMap := make(map[int]model.MovieWithActors)

	for rows.Next() {
		var movieID int
		var movie model.Movie
		var actor model.Actor
		var releaseDateStr, actorBirthDateStr string
		err := rows.Scan(&movieID, &movie.Title, &movie.Description, &releaseDateStr, &movie.Rating,
			&actor.ID, &actor.Name, &actor.Gender, &actorBirthDateStr)
		if err != nil {
			return nil, err
		}

		movieWithActors, ok := movieMap[movieID]
		if !ok {
			movie = model.Movie{ID: movieID, Title: movie.Title, Description: movie.Description, ReleaseDate: releaseDateStr, Rating: movie.Rating}
			movieWithActors = model.MovieWithActors{Movie: movie, Actors: []model.Actor{}}
			movieMap[movieID] = movieWithActors
		}

		actor.BirthDate = actorBirthDateStr
		movieWithActors.Actors = append(movieWithActors.Actors, actor)
		movieMap[movieID] = movieWithActors
	}

	for _, movieWithActors := range movieMap {
		movies = append(movies, movieWithActors)
	}

	return movies, nil
}

func (r *MoviePostgres) CreateMovie(movie model.Movie) error {
	query := fmt.Sprintf("INSERT INTO %s (title, description, rating, release_date) VALUES ($1, $2, $3, $4)", moviesTable)
	_, err := r.db.Exec(query, movie.Title, movie.Description, movie.Rating, movie.ReleaseDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoviePostgres) CreateMovieForActor(actorID int, movie model.Movie) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	query := fmt.Sprintf("INSERT INTO %s (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id", moviesTable)

	row := tx.QueryRow(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating)
	var movieID int
	err = row.Scan(&movieID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (movie_id, actor_id) VALUES ($1, $2)", movieActorTable)
	_, err = tx.Exec(query, movieID, actorID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MoviePostgres) GetMovieByID(movieID int) (model.Movie, error) {
	var movie model.Movie

	query := fmt.Sprintf("SELECT id, title, description, TO_CHAR(release_date, 'YYYY-MM-DD') as release_date, rating FROM %s WHERE id = $1", moviesTable)

	err := r.db.QueryRow(query, movieID).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.Rating)
	if err != nil {
		return model.Movie{}, err
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

func (r *MoviePostgres) UpdateMovie(movieID int, data model.Movie) error {
	return nil
}
