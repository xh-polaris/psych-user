package common

type Obj struct {
	Key   string `bson:"key,omitempty" json:"key,omitempty"`
	Value string `bson:"value,omitempty" json:"value,omitempty"`
}

type Form struct {
	Options []*Obj `bson:"options,omitempty" json:"options,omitempty"`
}
