package view

import (
	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

const (
	prefixUserCacheKey = "cache:view"
	CollectionName     = "view"
)

type IMongoMapper interface {
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}
