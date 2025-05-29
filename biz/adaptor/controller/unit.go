package controller

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/util/logx"
	"github.com/xh-polaris/psych-user/biz/application/service"
)

type IUnitController interface {
	UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *basic.Response, err error)
	UnitGetInfo(ctx context.Context, req *u.UnitGetInfoReq) (res *u.UnitGetInfoResp, err error)
	UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error)
	UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error)
	UnitStrongVerify(ctx context.Context, req *u.UnitStrongVerifyReq) (res *basic.Response, err error)
	UnitWeakVerify(ctx context.Context, req *u.UnitWeakVerifyReq) (res *basic.Response, err error)
	UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error)
	UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error)
	UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error)
	UnitLinkView(ctx context.Context, req *u.UnitLinkViewReq) (res *basic.Response, err error)
	UnitPageQueryUser(ctx context.Context, req *u.UnitPageQueryUserReq) (res *u.UnitPageQueryUserResp, err error)
	UnitPageQueryView(ctx context.Context, req *u.UnitPageQueryViewReq) (res *u.UnitPageQueryViewResp, err error)
	UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error)
	UnitCreateAndLinkView(ctx context.Context, req *u.UnitCreateAndLinkViewReq) (res *basic.Response, err error)
	UnitGetAppInfo(ctx context.Context, req *u.UnitGetAppInfoReq) (res *u.UnitGetAppInfoResp, err error)
	UnitModelGetInfo(ctx context.Context, req *u.UnitModelGetInfoReq) (res *u.UnitModelGetInfoResp, err error)
	UnitModelUpdateInfo(ctx context.Context, req *u.UnitModelUpdateInfoReq) (res *basic.Response, err error)
}

type UnitController struct {
	UnitService *service.UnitService
}

var UnitControllerSet = wire.NewSet(
	wire.Struct(new(UnitController), "*"),
	wire.Bind(new(IUnitController), new(*UnitController)),
)

func (u UnitController) UnitSignUp(ctx context.Context, req *u.UnitSignUpReq) (res *basic.Response, err error) {
	logx.Info("UnitSignUp", req)
	return u.UnitService.UnitSignUp(ctx, req)
}

func (u UnitController) UnitGetInfo(ctx context.Context, req *u.UnitGetInfoReq) (res *u.UnitGetInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitUpdateInfo(ctx context.Context, req *u.UnitUpdateInfoReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitUpdatePassword(ctx context.Context, req *u.UnitUpdatePasswordReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitStrongVerify(ctx context.Context, req *u.UnitStrongVerifyReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitWeakVerify(ctx context.Context, req *u.UnitWeakVerifyReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitCreateVerify(ctx context.Context, req *u.UnitCreateVerifyReq) (res *u.UnitCreateVerifyResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitUpdateVerifyPassword(ctx context.Context, req *u.UnitUpdateVerifyPasswordReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitLinkUser(ctx context.Context, req *u.UnitLinkUserReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitLinkView(ctx context.Context, req *u.UnitLinkViewReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitPageQueryUser(ctx context.Context, req *u.UnitPageQueryUserReq) (res *u.UnitPageQueryUserResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitPageQueryView(ctx context.Context, req *u.UnitPageQueryViewReq) (res *u.UnitPageQueryViewResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitCreateAndLinkUser(ctx context.Context, req *u.UnitCreateAndLinkUserReq) (res *basic.Response, err error) {
	logx.Info("UnitCreateAndLinkUser", req)
	return u.UnitService.UnitCreateAndLinkUser(ctx, req)
}

func (u UnitController) UnitCreateAndLinkView(ctx context.Context, req *u.UnitCreateAndLinkViewReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitGetAppInfo(ctx context.Context, req *u.UnitGetAppInfoReq) (res *u.UnitGetAppInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitModelGetInfo(ctx context.Context, req *u.UnitModelGetInfoReq) (res *u.UnitModelGetInfoResp, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UnitController) UnitModelUpdateInfo(ctx context.Context, req *u.UnitModelUpdateInfoReq) (res *basic.Response, err error) {
	//TODO implement me
	panic("implement me")
}
