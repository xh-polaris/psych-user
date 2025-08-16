package unit

import (
	"context"
	"errors"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
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
	unit.ID = primitive.NewObjectID()
	// 设置创建和更新时间
	now := time.Now()
	unit.CreateTime = now
	unit.UpdateTime = now

	_, err := m.conn.InsertOneNoCache(ctx, unit)
	return err
}

func (m *MongoMapper) InsertWithEcho(ctx context.Context, unit *Unit) (string, error) {
	unit.ID = primitive.NewObjectID()
	// 设置创建和更新时间
	now := time.Now()
	unit.CreateTime = now
	unit.UpdateTime = now
	res, err := m.conn.InsertOneNoCache(ctx, unit)
	if err != nil {
		return "", err
	}
	// 获取回显id
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

// FindOneByPhone 根据手机号查找单位
func (m *MongoMapper) FindOneByPhone(ctx context.Context, phone string) (*Unit, error) {
	var unit Unit
	err := m.conn.FindOneNoCache(ctx, &unit, bson.M{
		consts.Phone: phone,
	})
	switch {
	case err == nil:
		return &unit, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

// FindOne 根据ID查找单位
func (m *MongoMapper) FindOne(ctx context.Context, id string) (*Unit, error) {
	var unit Unit
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	err = m.conn.FindOneNoCache(ctx, &unit, bson.M{
		consts.ID: oid,
	})

	if err != nil {
		if errors.Is(err, monc.ErrNotFound) {
			return nil, consts.ErrNotFound
		}
		return nil, err
	}

	return &unit, nil
}

func (m *MongoMapper) UpdateBasicInfo(ctx context.Context, unit *Unit) error {
	unit.UpdateTime = time.Now()
	_, err := m.conn.UpdateByIDNoCache(ctx, unit.ID, bson.M{"$set": bson.M{
		consts.Name:       unit.Name,
		consts.Contact:    unit.Contact,
		consts.Address:    unit.Address,
		consts.UpdateTime: unit.UpdateTime,
	}})
	return err
}

func (m *MongoMapper) UpdatePassword(ctx context.Context, unitId, newPassword string) error {
	unitOid, err := primitive.ObjectIDFromHex(unitId)
	_, err = m.conn.UpdateByIDNoCache(ctx, unitOid, bson.M{"$set": bson.M{
		consts.Password:   newPassword,
		consts.UpdateTime: time.Now(),
	}})
	return err
}

func (m *MongoMapper) FindOneByAccount(ctx context.Context, account string) (*Unit, error) {
	var unit Unit
	err := m.conn.FindOneNoCache(ctx, &unit, bson.M{
		consts.Account: account,
	})
	switch {
	case err == nil:
		return &unit, nil
	case errors.Is(err, monc.ErrNotFound):
		return nil, consts.ErrNotFound
	default:
		return nil, err
	}
}

func (m *MongoMapper) UpdateVerify(ctx context.Context, verify *u.UnitVerify) error {
	unitOid, err := primitive.ObjectIDFromHex(verify.UnitId)
	_, err = m.conn.UpdateByIDNoCache(ctx, unitOid, bson.M{"$set": bson.M{
		consts.Account:        verify.Account,
		consts.VerifyPassword: verify.VerifyPassword,
		consts.VerifyType:     verify.VerifyType,
		consts.Form:           verify.Form,
		consts.UpdateTime:     time.Now(),
	}})
	return err
}

func (m *MongoMapper) UpdateVerifyPassword(ctx context.Context, verify *u.UnitVerify) error {
	unitOid, err := primitive.ObjectIDFromHex(verify.UnitId)
	_, err = m.conn.UpdateByIDNoCache(ctx, unitOid, bson.M{"$set": bson.M{
		consts.VerifyPassword: verify.VerifyPassword,
		consts.UpdateTime:     time.Now(),
	}})
	return err
}
