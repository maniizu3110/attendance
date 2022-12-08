package mycache

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/collection"
)

type cache struct {
	cache *collection.Cache
}

type Cache interface {
	Init()
	CreatePreloadCache()
	CreateJoinsCache()
	IsPreloadable(modelName string, key string) bool
	IsJoinable(modelName string, key string) bool
}

func NewCache() Cache {
	c, err := collection.NewCache(time.Hour * 24 * 365 * 100)
	if err != nil {
		logrus.Warn("キャッシュの作成に失敗しました", err)
	}
	return &cache{
		cache: c,
	}
}

func (c *cache) Init() {
	c.CreatePreloadCache()
	c.CreateJoinsCache()
}
