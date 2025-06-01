package unit_user

import (
	"context"
	"errors"
	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
	"github.com/xh-polaris/psych-user/biz/infrastructure/consts"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	prefixUnitCacheKey = "cache:unit_user"
	CollectionName     = "unit_user"
)

type IMongoMapper interface {
	Insert(ctx context.Context, u *UnitUser) error
	FindOneByUU(ctx context.Context, userId string, unitId string) (*UnitUser, error)
	FindOneByUnitAndStu(ctx context.Context, unitId string, studentId string) (*UnitUser, error)
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

func (m MongoMapper) Insert(ctx context.Context, u *UnitUser) error {
	u.Id = primitive.NewObjectID()
	_, err := m.conn.InsertOneNoCache(ctx, u)
	return err
}

func (m MongoMapper) FindOneByUU(ctx context.Context, userId string, unitId string) (*UnitUser, error) {
	var u UnitUser
	err := m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.UserID: userId,
		consts.UnitId: unitId,
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

func (m MongoMapper) FindOneByUnitAndStu(ctx context.Context, unitId string, studentId string) (*UnitUser, error) {
	var u UnitUser
	err := m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.UnitId:    unitId,
		consts.StudentId: studentId,
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
