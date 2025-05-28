package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/xh-polaris/ActiManage-IDL-gen/kitex_gen/user/userservice"
	"github.com/xh-polaris/ActiManage-user/biz/infrastructure/config"
	"github.com/xh-polaris/ActiManage-user/biz/infrastructure/util/log"
	"github.com/xh-polaris/ActiManage-user/provider"
	"github.com/xh-polaris/gopkg/kitex/middleware"
	logx "github.com/xh-polaris/gopkg/util/log"
	"net"
)

func main() {
	klog.SetLogger(logx.NewKlogLogger())
	s, err := provider.NewProvider()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", config.GetConfig().ListenOn)
	if err != nil {
		panic(err)
	}
	svr := userservice.NewServer(
		s,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.GetConfig().Name}),
		server.WithMiddleware(middleware.LogMiddleware(config.GetConfig().Name)),
	)

	err = svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}
