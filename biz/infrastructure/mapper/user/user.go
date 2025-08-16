package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Phone      string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Password   string             `bson:"password,omitempty" json:"password,omitempty"` // 考虑是否应该存储密码明文
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	Birth      string             `bson:"birth,omitempty" json:"birth,omitempty"` // 建议使用 time.Time 类型
	Gender     int32              `bson:"gender,omitempty" json:"gender,omitempty"`
	Status     int32              `bson:"status,omitempty" json:"status,omitempty"` // 建议使用 int 或 bool 类型
	CreateTime time.Time          `bson:"create_time,omitempty" json:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"update_time,omitempty" json:"updateTime,omitempty"`
	DeleteTime time.Time          `bson:"delete_time,omitempty" json:"deleteTime,omitempty"`
}
