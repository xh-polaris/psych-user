package unit_user

import (
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitUser struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    string             `bson:"user_id" json:"user_id,omitempty"`
	StudentId string             `bson:"student_id" json:"student_id,omitempty"`
	UnitId    string             `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	Options   *common.Form       `bson:"options" json:"options,omitempty"`
}
