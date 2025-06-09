package consts

// 数据库相关
const (
	ID         = "_id"
	StudentId  = "student_id"
	UserID     = "user_id"
	UnitId     = "unit_id"
	ViewId     = "view_id"
	Status     = "status"
	CreateTime = "create_time"
	Phone      = "phone"
	Timestamp  = "timestamp"
	LogId      = "log_id"
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
