package unit

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UnitVerify struct {
	UnitId     string         `bson:"unit_id,omitempty" json:"userId,omitempty"`
	VerifyType int32          `bson:"verify_type,omitempty"  json:"verifyType,omitempty"`
	Account    string         `bson:"account,omitempty"  json:"account,omitempty"`
	Password   string         `bson:"password,omitempty"  json:"password,omitempty"`
	Form       map[string]any `bson:"form,omitempty" json:"form,omitempty"`
}

type Unit struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Phone          string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Password       string             `bson:"password,omitempty" json:"password,omitempty"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	Address        string             `bson:"address,omitempty" json:"address,omitempty"`
	Contact        string             `bson:"contact,omitempty" json:"contact,omitempty"`
	Level          int32              `bson:"level,omitempty" json:"level,omitempty"`
	Status         int32              `bson:"status,omitempty" json:"status,omitempty"`
	VerifyType     int32              `bson:"verify_type,omitempty"  json:"verifyType,omitempty"`
	Account        string             `bson:"account,omitempty"  json:"account,omitempty"`
	VerifyPassword string             `bson:"verify_password,omitempty"  json:"verify_password,omitempty"`
	Form           map[string]any     `bson:"form,omitempty" json:"form,omitempty"`
	CreateTime     time.Time          `bson:"create_time,omitempty" json:"createTime,omitempty"`
	UpdateTime     time.Time          `bson:"update_time,omitempty" json:"updateTime,omitempty"`
	DeleteTime     time.Time          `bson:"delete_time,omitempty" json:"deleteTime,omitempty"`
}
