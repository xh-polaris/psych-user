package unit

type Unit struct {
	Id         string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Phone      string `protobuf:"bytes,2,opt,name=phone" json:"phone,omitempty"`
	Password   string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	Name       string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Address    string `protobuf:"bytes,5,opt,name=address" json:"address,omitempty"`
	Contact    string `protobuf:"bytes,6,opt,name=contact" json:"contact,omitempty"`
	Level      int32  `protobuf:"varint,7,opt,name=level" json:"level,omitempty"`
	Status     string `protobuf:"bytes,8,opt,name=status" json:"status,omitempty"`
	CreateTime int64  `protobuf:"varint,9,opt,name=createTime" json:"createTime,omitempty"`
	UpdateTime int64  `protobuf:"varint,10,opt,name=updateTime" json:"updateTime,omitempty"`
	DeleteTime int64  `protobuf:"varint,11,opt,name=deleteTime" json:"deleteTime,omitempty"`
}
