package unit

import (
	"context"
	"errors"
	"time"

	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	prefixUnitCacheKey = "cache:unit"
	CollectionName     = "unit"
)

type IMongoMapper interface {
	Insert(ctx context.Context, unit *Unit) error
	Update(ctx context.Context, unit *Unit) error
	FindOne(ctx context.Context, id string) (*Unit, error)
	FindOneByPhone(ctx context.Context, phone string) (*Unit, error)
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

func (m *MongoMapper) Insert(ctx context.Context, unit *Unit) error {
	if unit.Id == "" {
		objectId := primitive.NewObjectID()
		unit.Id = objectId.Hex()
		unit.CreateTime = time.Now().Unix()
		unit.UpdateTime = unit.CreateTime
	}
	_, err := m.conn.InsertOneNoCache(ctx, unit)
	return err
}

func (m *MongoMapper) Update(ctx context.Context, unit *Unit) error {
	unit.UpdateTime = time.Now().Unix()
	_, err := m.conn.UpdateByIDNoCache(ctx, unit.Id, bson.M{"$set": unit})
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Unit, error) {
	var u Unit
	err := m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.ID: id,
	})
	if err != nil {
		return nil, consts.ErrNotFound
	}
	return &u, nil
}

func (m *MongoMapper) FindOneByPhone(ctx context.Context, phone string) (*Unit, error) {
	var u Unit
	err := m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.Phone: phone,
	})
	switch {
	case err == nil:
		return &u, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}
