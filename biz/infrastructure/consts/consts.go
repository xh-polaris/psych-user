package consts

// 数据库相关
const (
	ID         = "_id"
	StudentId  = "student_id"
	UserID     = "user_id"
	UnitId     = "unit_id"
	ViewId     = "view_id"
	Status     = "status"
	Phone      = "phone"
	Timestamp  = "timestamp"
	LogId      = "log_id"
	Name       = "name"
	Address    = "address"
	Contact    = "contact"
	Level      = "level"
	CreateTime = "create_time"
	UpdateTime = "update_time"
	DeleteTime = "delete_time"
	Password   = "password"
	NotEqual   = "$ne"
)

// createType
const (
	CreateByPhone     = 1
	CreateByStudentId = 2
	CreateByEmail     = 3
)

// authType
const (
	AuthPhoneAndPwd     = 1
	AuthStudentIdAndPwd = 2
	AuthPhoneAndCode    = 3
)

// password
const (
	DefaultPassword = "123456"
	UpdateByOldPwd  = 0
	UpdateByCode    = 1
)

// status
const (
	Active  = 0
	Deleted = 1
)

// gender
const (
	Unknown = 0
	Male    = 1
	Female  = 2
)

// verifyType
const (
	Strong = 0
	Weak   = 1
)
