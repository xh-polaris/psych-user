package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id         primitive.ObjectID `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Phone      string             `protobuf:"bytes,2,opt,name=phone" json:"phone,omitempty"`
	Password   string             `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	Name       string             `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Birth      string             `protobuf:"bytes,5,opt,name=birth" json:"birth,omitempty"`
	Gender     string             `protobuf:"bytes,6,opt,name=gender" json:"gender,omitempty"`
	Status     string             `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
	CreateTime time.Time          `protobuf:"varint,8,opt,name=createTime" json:"createTime,omitempty"`
	UpdateTime time.Time          `protobuf:"varint,9,opt,name=updateTime" json:"updateTime,omitempty"`
	DeleteTime time.Time          `protobuf:"varint,10,opt,name=deleteTime" json:"deleteTime,omitempty"`
}

type UserStudentUnit struct {
	Id        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id,omitempty"`
	StudentId primitive.ObjectID `bson:"student_id" json:"student_id,omitempty"`
	UnitId    primitive.ObjectID `bson:"unit_id" json:"unit_id,omitempty"`
}
