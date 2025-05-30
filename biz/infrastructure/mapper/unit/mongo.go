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
	UnitUserCollection = "unit_user"
)

type IMongoMapper interface {
	Insert(ctx context.Context, unit *Unit) error
	FindOneByPhone(ctx context.Context, phone string) (*Unit, error)
	FindOne(ctx context.Context, id string) (*Unit, error)
	LinkUser(ctx context.Context, unitId, userId string) error
	CheckLinkExists(ctx context.Context, unitId, userId string) (bool, error)
}

type MongoMapper struct {
	conn     *monc.Model
	linkConn *monc.Model
}

func NewMongoMapper(config *config.Config) *MongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	linkConn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, UnitUserCollection, config.Cache)
	return &MongoMapper{
		conn:     conn,
		linkConn: linkConn,
	}
}

// Insert 插入新的单位记录
func (m *MongoMapper) Insert(ctx context.Context, unit *Unit) error {
	// 设置创建和更新时间
	now := time.Now()
	unit.CreateTime = now
	unit.UpdateTime = now

	_, err := m.conn.InsertOneNoCache(ctx, unit)
	return err
}

// FindOneByPhone 根据手机号查找单位
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

// FindOne 根据ID查找单位
func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Unit, error) {
	var u Unit
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	err = m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.ID_: oid,
	})

	if err != nil {
		if errors.Is(err, monc.ErrNotFound) {
			return nil, consts.ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

// LinkUser 创建单位和用户的关联
func (m *MongoMapper) LinkUser(ctx context.Context, unitId, userId string) error {
	unitOid, err := primitive.ObjectIDFromHex(unitId)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	userOid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return consts.ErrInvalidObjectId
	}

	link := bson.M{
		"unit_id": unitOid,
		"user_id": userOid,
	}

	_, err = m.linkConn.InsertOneNoCache(ctx, link)
	return err
}

// CheckLinkExists 检查单位和用户的关联是否已存在
func (m *MongoMapper) CheckLinkExists(ctx context.Context, unitId, userId string) (bool, error) {
	unitOid, err := primitive.ObjectIDFromHex(unitId)
	if err != nil {
		return false, consts.ErrInvalidObjectId
	}
	userOid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false, consts.ErrInvalidObjectId
	}

	var result bson.M
	err = m.linkConn.FindOneNoCache(ctx, &result, bson.M{
		"unit_id": unitOid,
		"user_id": userOid,
	})

	if err != nil {
		if errors.Is(err, monc.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
