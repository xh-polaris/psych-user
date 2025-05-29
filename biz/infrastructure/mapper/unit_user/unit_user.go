package unit_user

import "github.com/xh-polaris/psych-idl/kitex_gen/user"

type UnitUser struct {
	UserId    string
	StudentId string
	UnitId    string
	Options   *user.Option
}
