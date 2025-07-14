package consts

// 数据库相关
const (
	ID             = "_id"
	StudentId      = "student_id"
	UserID         = "user_id"
	UnitId         = "unit_id"
	ViewId         = "view_id"
	Status         = "status"
	Phone          = "phone"
	Timestamp      = "timestamp"
	LogId          = "log_id"
	Name           = "name"
	Address        = "address"
	Contact        = "contact"
	Level          = "level"
	CreateTime     = "create_time"
	UpdateTime     = "update_time"
	DeleteTime     = "delete_time"
	Password       = "password"
	NotEqual       = "$ne"
	Account        = "account"
	VerifyPassword = "verify_password"
	VerifyType     = "verify_type"
	Form           = "form"
	Options        = "options"
)

// createType
const (
	CreateByPhone     = 0
	CreateByStudentId = 1
	CreateByEmail     = 2
)

// authType
const (
	AuthPhoneAndPwd       = 0
	AuthStudentIdAndPwd   = 1
	AuthPhoneAndCode      = 2
	AuthWeakAccountAndPwd = 3
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

// verify type
const (
	Strong = 0
	Weak   = 1
)
