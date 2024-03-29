package psql

import (
	"fmt"
	"strings"

	"github.com/Woodfyn/Web-api/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const gameTable = "games"

type GamesPosegres struct {
	db        *sqlx.DB
	mainCache *GameCache
}

func NewGames(db *sqlx.DB) *GamesPosegres {
	return &GamesPosegres{
		db:        db,
		mainCache: NewGameCache(),
	}
}

func (r *GamesPosegres) Create(game domain.Game) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createGameQuery := fmt.Sprintf(`INSERT INTO %s (title, genre, evaluation) VALUES ($1, $2, $3) RETURNING id`, gameTable)
	err = tx.QueryRow(createGameQuery, game.Title, game.Genre, game.Evaluation).Scan(&game.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	go func(game domain.Game) {
		r.mainCache.SetCache(game)
	}(game)

	return tx.Commit()
}

func (r *GamesPosegres) GetAll() ([]domain.Game, error) {
	var games []domain.Game
	getAllGameQuery := fmt.Sprintf(`SELECT id, title, genre, evaluation FROM %s`, gameTable)
	if err := r.db.Select(&games, getAllGameQuery); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *GamesPosegres) GetById(gameId int) (domain.Game, error) {
	cacheGame, err := r.mainCache.GetCache(gameId)
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"handler": "GetById",
			"info":    "fetched from cache",
		}).Info()
		return cacheGame, nil
	}

	var game domain.Game
	getGameQuery := fmt.Sprintf(`SELECT id, title, genre, evaluation FROM %s WHERE id = $1`, gameTable)
	if err := r.db.Get(&game, getGameQuery, gameId); err != nil {
		return game, err
	}

	return game, nil

}

func (r *GamesPosegres) UpdateById(gameId int, input domain.UpdateGameInput) error {
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

func (r *GamesPosegres) DeleteById(gameId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, gameTable)
	_, err := r.db.Exec(query, gameId)

	return err
}
