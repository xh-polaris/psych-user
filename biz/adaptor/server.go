package adaptor

import "github.com/xh-polaris/psych-user/biz/adaptor/controller"

type Server struct {
	controller.IUserController
	controller.IUnitController
	controller.IViewController
}
