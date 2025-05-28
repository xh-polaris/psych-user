package provider

import (
	"github.com/google/wire"
	"github.com/xh-polaris/psych-user/biz/adaptor/controller"
	"github.com/xh-polaris/psych-user/biz/application/service"
	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/unit"
	"github.com/xh-polaris/psych-user/biz/infrastructure/mapper/user"
)

var ApplicationSet = wire.NewSet(
	service.UserServiceSet,
	service.ViewServiceSet,
	service.UnitServiceSet,
)

var MapperSet = wire.NewSet(
	unit.NewMongoMapper,
	user.NewMongoMapper,
	// view.NewMongoMapper,
)

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
	MapperSet,
)

var ControllerSet = wire.NewSet(
	controller.UnitControllerSet,
	controller.UserControllerSet,
	controller.ViewControllerSet,
)

var ServerProvider = wire.NewSet(
	ControllerSet,
	ApplicationSet,
	InfrastructureSet,
)
