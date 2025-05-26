package provider

import (
	"github.com/google/wire"
	"github.com/xh-polaris/psych-user/biz/infrastructure/config"
)

var provider *Provider

func Init() {
	var err error
	provider, err = NewProvider()
	if err != nil {
		panic(err)
	}
}

// Provider 提供controller依赖的对象
type Provider struct {
	Config *config.Config
}

func Get() *Provider {
	return provider
}

var RpcSet = wire.NewSet()

var ApplicationSet = wire.NewSet()

var InfrastructureSet = wire.NewSet(
	config.NewConfig,
)

var AllProvider = wire.NewSet(
	ApplicationSet,
	InfrastructureSet,
)
