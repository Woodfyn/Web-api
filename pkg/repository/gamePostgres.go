package repository

import (
	"fmt"
	"strings"

	todo "github.com/Woodfyn/Web-api"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type GamePostgres struct {
	db        *sqlx.DB
	mainCache *GameCache
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{
		db:        db,
		mainCache: NewGameCache(),
	}
}

func (r *GamePostgres) Create(game todo.Game) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createGameQuery := fmt.Sprintf(`INSERT INTO %s (title, genre, evaluation) VALUES ($1, $2, $3) RETURNING id`, gameTable)
	err = tx.QueryRow(createGameQuery, game.Title, game.Genre, game.Evaluation).Scan(&game.Id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	go func(game todo.Game) {
		r.mainCache.SetCache(game)
	}(game)

	return game.Id, tx.Commit()
}

func (r *GamePostgres) GetAll() ([]todo.Game, error) {
	var games []todo.Game
	getAllGameQuery := fmt.Sprintf(`SELECT id, title, genre, evaluation FROM %s`, gameTable)
	if err := r.db.Select(&games, getAllGameQuery); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *GamePostgres) GetById(gameId int) (todo.Game, error) {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	cacheGame, err := r.mainCache.GetCache(gameId)
	if err == nil {
		return cacheGame, nil
	}

	var game todo.Game
	getGameQuery := fmt.Sprintf(`SELECT id, title, genre, evaluation FROM %s WHERE id = $1`, gameTable)
	if err := r.db.Get(&game, getGameQuery, gameId); err != nil {
		return game, err
	}

	return game, nil

}

func (r *GamePostgres) UpdateById(gameId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Genre != nil {
		setValues = append(setValues, fmt.Sprintf("genre=$%d", argId))
		args = append(args, *input.Genre)
		argId++
	}

	if input.Evaluation != nil {
		setValues = append(setValues, fmt.Sprintf("evaluation=$%d", argId))
		args = append(args, *input.Evaluation)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, gameTable, setQuery, argId)
	args = append(args, gameId)

	_, err := r.db.Exec(query, args...)

	return err
}

func (r *GamePostgres) DeleteById(gameId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, gameTable)
	_, err := r.db.Exec(query, gameId)

	return err
}
