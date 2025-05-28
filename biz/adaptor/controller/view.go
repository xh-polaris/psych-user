package controller

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	v "github.com/xh-polaris/psych-idl/kitex_gen/user"
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

func (v ViewController) ViewSignUp(ctx context.Context, req *v.ViewSignUpReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewController) ViewGetInfo(ctx context.Context, req *v.ViewGetInfoReq) (res *v.ViewGetInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewController) ViewUpdateInfo(ctx context.Context, req *v.ViewUpdateInfoReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewController) ViewUpdatePassword(ctx context.Context, req *v.ViewUpdatePasswordReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewController) ViewBelongUnit(ctx context.Context, req *v.ViewBelongUnitReq) (res *v.ViewBelongUnitResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewController) ViewSignIn(ctx context.Context, req *v.ViewSignInReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

var ViewControllerSet = wire.NewSet(
	wire.Struct(new(ViewController), "*"),
	wire.Bind(new(IViewController), new(*ViewController)),
)
