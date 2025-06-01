package controller

import (
	"context"
	"github.com/google/wire"
	"github.com/xh-polaris/psych-idl/kitex_gen/basic"
	u "github.com/xh-polaris/psych-idl/kitex_gen/user"
	"github.com/xh-polaris/psych-pkg/util/logx"
	"github.com/xh-polaris/psych-user/biz/application/service"
)

type IUserController interface {
	UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error)
	UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *basic.Response, err error)
	UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *basic.Response, err error)
	UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error)
	UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error)
	UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error)
}

type UserController struct {
	UserService *service.UserService
}

var UserControllerSet = wire.NewSet(
	wire.Struct(new(UserController), "*"),
	wire.Bind(new(IUserController), new(*UserController)),
)

func (u UserController) UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error) {
	logx.Info("UserGetInfo", req)
	return u.UserService.UserGetInfo(ctx, req)
}

func (u UserController) UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *basic.Response, err error) {
	return u.UserService.UserSignIn(ctx, req)
}

func (u UserController) UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *basic.Response, err error) {
	return u.UserService.UserSignUp(ctx, req)
}

func (u UserController) UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error) {
	return u.UserService.UserUpdateInfo(ctx, req)
}

func (u UserController) UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error) {
	return u.UserService.UserUpdatePassword(ctx, req)
}

func (u UserController) UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error) {
	return u.UserService.UserBelongUnit(ctx, req)
}
