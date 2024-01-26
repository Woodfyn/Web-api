package psql

import (
	"errors"
	"time"

	"github.com/Woodfyn/Web-api/internal/domain"
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

func (c GameCache) SetCache(game domain.Game) {
	c.cache.Set(game.Id, game, TTLCache)
}

func (c GameCache) GetCache(gameKey int) (domain.Game, error) {
	gameInterface, err := c.cache.Get(gameKey)
	if err != nil {
		return domain.Game{}, err
	}

	game, ok := gameInterface.(domain.Game)
	if !ok {
		return domain.Game{}, errors.New("no this type")
	}

	return game, nil
}
