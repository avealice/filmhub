package repository

import (
	"filmhub/internal/model"
	"fmt"

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
	query := fmt.Sprintf("INSERT INTO %s (name, gender, birth_date) VALUES ($1, $2, $3)", actorsTable)

	_, err := r.db.Exec(query, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return err
	}

	return nil
}

// func (r *ActorPostgres) GetAllActors() ([]model.ActorWithMovies, error) {
// 	var actors []model.ActorWithMovies

// 	query := fmt.Sprintf("SELECT id, name, gender, TO_CHAR(birth_date, 'YYYY-MM-DD') as birth_date FROM %s", actorsTable)

// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var actor model.Actor
// 		var birthDateStr string
// 		err := rows.Scan(&actor.ID, &actor.Name, &actor.Gender, &birthDateStr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		actor.BirthDate = birthDateStr
// 		actors = append(actors, actor)
// 	}

//		return actors, nil
//	}

func (r *ActorPostgres) GetAllActors() ([]model.ActorWithMovies, error) {
	var actors []model.ActorWithMovies

	query := fmt.Sprintf("SELECT a.id, a.name, a.gender, TO_CHAR(a.birth_date, 'YYYY-MM-DD') as birth_date, m.id as movie_id, m.title as movie_title, m.description as movie_description, TO_CHAR(m.release_date, 'YYYY-MM-DD') as movie_release_date, m.rating as movie_rating FROM %s a LEFT JOIN %s ma ON a.id = ma.actor_id LEFT JOIN %s m ON ma.movie_id = m.id", actorsTable, movieActorTable, moviesTable)

	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	actorMap := make(map[int]model.ActorWithMovies)

	for rows.Next() {
		var actorID int
		var actor model.Actor
		var movie model.Movie
		var birthDateStr, movieReleaseDateStr string
		err := rows.Scan(&actorID, &actor.Name, &actor.Gender, &birthDateStr,
			&movie.ID, &movie.Title, &movie.Description, &movieReleaseDateStr, &movie.Rating)
		if err != nil {
			return nil, err
		}

		actorWithMovies, ok := actorMap[actorID]
		if !ok {
			actor = model.Actor{ID: actorID, Name: actor.Name, Gender: actor.Gender, BirthDate: birthDateStr}
			actorWithMovies = model.ActorWithMovies{Actor: actor, Movies: []model.Movie{}}
			actorMap[actorID] = actorWithMovies
		}

		movie.ReleaseDate = movieReleaseDateStr
		actorWithMovies.Movies = append(actorWithMovies.Movies, movie)
		actorMap[actorID] = actorWithMovies
	}

	for _, actorWithMovies := range actorMap {
		actors = append(actors, actorWithMovies)
	}

	return actors, nil
}

func (r *ActorPostgres) Delete(actorID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", actorsTable)
	_, err := r.db.Exec(query, actorID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ActorPostgres) Get(actorID int) (model.Actor, error) {
	var actor model.Actor

	query := fmt.Sprintf("SELECT id, name, gender, TO_CHAR(birth_date, 'YYYY-MM-DD') as birth_date FROM %s WHERE id = $1", actorsTable)

	err := r.db.QueryRow(query, actorID).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)
	if err != nil {
		return model.Actor{}, err
	}

	return actor, nil
}
