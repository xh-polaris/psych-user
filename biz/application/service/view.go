package service

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	v "github.com/xh-polaris/psych-idl/kitex_gen/user"
)

type IViewService interface {
	ViewSignUp(ctx context.Context, req *v.ViewSignUpReq) (res *basic.Response, err error)
	ViewGetInfo(ctx context.Context, req *v.ViewGetInfoReq) (res *v.ViewGetInfoResp, err error)
	ViewUpdateInfo(ctx context.Context, req *v.ViewUpdateInfoReq) (res *basic.Response, err error)
	ViewUpdatePassword(ctx context.Context, req *v.ViewUpdatePasswordReq) (res *basic.Response, err error)
	ViewBelongUnit(ctx context.Context, req *v.ViewBelongUnitReq) (res *v.ViewBelongUnitResp, err error)
	ViewSignIn(ctx context.Context, req *v.ViewSignInReq) (res *basic.Response, err error)
}
type ViewService struct{}

var ViewServiceSet = wire.NewSet(
	wire.Struct(new(ViewService), "*"),
	wire.Bind(new(IViewService), new(*ViewService)),
)

func (v ViewService) ViewSignUp(ctx context.Context, req *v.ViewSignUpReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewService) ViewGetInfo(ctx context.Context, req *v.ViewGetInfoReq) (res *v.ViewGetInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewService) ViewUpdateInfo(ctx context.Context, req *v.ViewUpdateInfoReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewService) ViewUpdatePassword(ctx context.Context, req *v.ViewUpdatePasswordReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewService) ViewBelongUnit(ctx context.Context, req *v.ViewBelongUnitReq) (res *v.ViewBelongUnitResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (v ViewService) ViewSignIn(ctx context.Context, req *v.ViewSignInReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}
