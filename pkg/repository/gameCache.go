package repository

import (
	"errors"
	"time"

	todo "github.com/Woodfyn/Web-api"

	"github.com/Woodfyn/inMemoryCache/cache"
)

const (
	worker   = 2
	TTLCache = 100 * time.Second
)

type GameCache struct {
	cache *cache.Cache
}

func NewGameCache() *GameCache {
	return &GameCache{
		cache: cache.NewCache(worker),
	}
}

func (c GameCache) SetCache(game todo.Game) {
	c.cache.Set(game.Id, game, TTLCache)
}

func (c GameCache) GetCache(gameKey int) (todo.Game, error) {
	gameInterface, err := c.cache.Get(gameKey)
	if err != nil {
		return todo.Game{}, err
	}

	game, ok := gameInterface.(todo.Game)
	if !ok {
		return todo.Game{}, errors.New("no this type")
	}

	return game, nil
}
