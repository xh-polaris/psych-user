package controller

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/util/logx"
	"github.com/xh-polaris/psych-user/biz/application/service"
)

type IViewController interface {
	ViewSignUp(ctx context.Context, req *u.ViewSignUpReq) (res *u.ViewSignUpResp, err error)
	ViewGetInfo(ctx context.Context, req *u.ViewGetInfoReq) (res *u.ViewGetInfoResp, err error)
	ViewUpdateInfo(ctx context.Context, req *u.ViewUpdateInfoReq) (res *basic.Response, err error)
	ViewUpdatePassword(ctx context.Context, req *u.ViewUpdatePasswordReq) (res *basic.Response, err error)
	ViewBelongUnit(ctx context.Context, req *u.ViewBelongUnitReq) (res *u.ViewBelongUnitResp, err error)
	ViewSignIn(ctx context.Context, req *u.ViewSignInReq) (res *u.ViewSignInResp, err error)
}

var ViewControllerSet = wire.NewSet(
	wire.Struct(new(ViewController), "*"),
	wire.Bind(new(IViewController), new(*ViewController)),
)

type ViewController struct {
	ViewService *service.ViewService
}

func (u *ViewController) ViewSignUp(ctx context.Context, req *u.ViewSignUpReq) (res *u.ViewSignUpResp, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
	// return u.ViewService.ViewSignUp(ctx, req)
}
func (u *ViewController) ViewGetInfo(ctx context.Context, req *u.ViewGetInfoReq) (res *u.ViewGetInfoResp, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
}
func (u *ViewController) ViewUpdateInfo(ctx context.Context, req *u.ViewUpdateInfoReq) (res *basic.Response, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
}
func (u *ViewController) ViewUpdatePassword(ctx context.Context, req *u.ViewUpdatePasswordReq) (res *basic.Response, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
}
func (u *ViewController) ViewBelongUnit(ctx context.Context, req *u.ViewBelongUnitReq) (res *u.ViewBelongUnitResp, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
}
func (u *ViewController) ViewSignIn(ctx context.Context, req *u.ViewSignInReq) (res *u.ViewSignInResp, err error) {
	logx.Info("ViewSignUp", req)
	panic("implement me")
}
