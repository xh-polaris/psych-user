package user

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
	prefixUserCacheKey = "cache:user"
	CollectionName     = "user"
	UnitUserCollection = "unit_user"
)

type IMongoMapper interface {
	Insert(ctx context.Context, user *User) error
	InsertMany(ctx context.Context, users []*User) error
	Update(ctx context.Context, user *User) error
	FindOne(ctx context.Context, id string) (*User, error)
	FindOneByPhone(ctx context.Context, id string) (*User, error)
	UpdateCount(ctx context.Context, id string, increment int64) error
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

func (m *MongoMapper) Insert(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreateTime = time.Now()
	user.UpdateTime = user.CreateTime
	_, err := m.conn.InsertOneNoCache(ctx, user)
	return err
}

func (m *MongoMapper) InsertWithEcho(ctx context.Context, user *User) (string, error) {
	user.ID = primitive.NewObjectID()
	user.CreateTime = time.Now()
	user.UpdateTime = user.CreateTime
	res, err := m.conn.InsertOneNoCache(ctx, user)
	if err != nil {
		return "", err
	}
	// 获取回显id
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, err
}

func (m *MongoMapper) Update(ctx context.Context, user *User) error {
	user.UpdateTime = time.Now()
	_, err := m.conn.UpdateByIDNoCache(ctx, user.ID, bson.M{"$set": user})
	return err
}

func (m *MongoMapper) FindOne(ctx context.Context, id string) (*User, error) {
	var u User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, consts.ErrInvalidObjectId
	}
	err = m.conn.FindOneNoCache(ctx, &u, bson.M{
		consts.ID: oid,
	})
	if err != nil {
		if errors.Is(err, monc.ErrNotFound) {
			return nil, consts.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (m *MongoMapper) FindOneByPhone(ctx context.Context, phone string) (*User, error) {
	var u User
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

func (m *MongoMapper) UpdateCount(ctx context.Context, id string, increment int64) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return consts.ErrInvalidObjectId
	}
	_, err = m.conn.UpdateByIDNoCache(ctx, oid, bson.M{
		"$inc": bson.M{
			"count": increment,
		},
	})
	return err
}

//func (m *MongoMapper) FindUSULinkBySid(ctx context.Context, sid string) (*UserStudentUnit, error) {
//	oid, err := primitive.ObjectIDFromHex(sid)
//	if err != nil {
//		return nil, consts.ErrInvalidObjectId
//	}
//	var u UserStudentUnit
//	err = m.conn.FindOneNoCache(ctx, &u, bson.M{
//		consts.StudentId: oid,
//	})
//	if err != nil {
//		return nil, consts.ErrNotFound
//	}
//	return &u, nil
//}
