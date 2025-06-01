package unit

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Unit struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Phone      string             `bson:"phone" json:"phone,omitempty"`
	Password   string             `bson:"password,omitempty" json:"password,omitempty"`
	Name       string             `bson:"name" json:"name,omitempty"`
	Address    string             `bson:"address,omitempty" json:"address,omitempty"`
	Contact    string             `bson:"contact,omitempty" json:"contact,omitempty"`
	Level      int32              `bson:"level" json:"level,omitempty"`
	Status     int32              `bson:"status" json:"status,omitempty"`
	CreateTime time.Time          `bson:"create_time,omitempty" json:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"update_time,omitempty" json:"updateTime,omitempty"`
	DeleteTime time.Time          `bson:"delete_time,omitempty" json:"deleteTime,omitempty"`
}
