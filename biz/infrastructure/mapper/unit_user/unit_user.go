package unit_user

type UnitUser struct {
	UserId    string         `bson:"user_id,omitempty" json:"user_id,omitempty"`
	StudentId string         `bson:"student_id,omitempty" json:"student_id,omitempty"`
	UnitId    string         `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	Options   map[string]any `bson:"options,omitempty" json:"options,omitempty"`
}
