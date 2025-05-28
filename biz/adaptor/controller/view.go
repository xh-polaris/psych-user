package controller

import (
	"context"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	v "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/wirex"
	"github.com/xh-polaris/psych-user/biz/application/service"
)

type IViewController interface {
	ViewSignUp(ctx context.Context, req *v.ViewSignUpReq) (res *basic.Response, err error)
	ViewGetInfo(ctx context.Context, req *v.ViewGetInfoReq) (res *v.ViewGetInfoResp, err error)
	ViewUpdateInfo(ctx context.Context, req *v.ViewUpdateInfoReq) (res *basic.Response, err error)
	ViewUpdatePassword(ctx context.Context, req *v.ViewUpdatePasswordReq) (res *basic.Response, err error)
	ViewBelongUnit(ctx context.Context, req *v.ViewBelongUnitReq) (res *v.ViewBelongUnitResp, err error)
	ViewSignIn(ctx context.Context, req *v.ViewSignInReq) (res *basic.Response, err error)
}

type ViewController struct {
	ViewService *service.ViewService
}

var ViewControllerSet = wirex.NewWireSet[ViewController, IViewController]()
