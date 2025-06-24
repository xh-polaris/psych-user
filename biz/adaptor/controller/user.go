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
	UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *u.UserSignUpResp, err error)
	UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error)
	UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error)
	UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error)
	UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error)
	UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *u.UserSignInResp, err error)
}

type UserController struct {
	UserService *service.UserService
}

var UserControllerSet = wire.NewSet(
	wire.Struct(new(UserController), "*"),
	wire.Bind(new(IUserController), new(*UserController)),
)

func (u *UserController) UserSignUp(ctx context.Context, req *u.UserSignUpReq) (res *u.UserSignUpResp, err error) {
	logx.Info("UserSignUp", req)
	return u.UserService.UserSignUp(ctx, req)
}

func (u *UserController) UserGetInfo(ctx context.Context, req *u.UserGetInfoReq) (res *u.UserGetInfoResp, err error) {
	logx.Info("UserGetInfo", req)
	return u.UserService.UserGetInfo(ctx, req)
}

func (u *UserController) UserUpdateInfo(ctx context.Context, req *u.UserUpdateInfoReq) (res *basic.Response, err error) {
	logx.Info("UserUpdateInfo", req)
	return u.UserService.UserUpdateInfo(ctx, req)
}
func (u *UserController) UserUpdatePassword(ctx context.Context, req *u.UserUpdatePasswordReq) (res *basic.Response, err error) {
	logx.Info("UserUpdatePassword", req)
	return u.UserService.UserUpdatePassword(ctx, req)
}
func (u *UserController) UserBelongUnit(ctx context.Context, req *u.UserBelongUnitReq) (res *u.UserBelongUnitResp, err error) {
	logx.Info("UserBelongUnit", req)
	return u.UserService.UserBelongUnit(ctx, req)
}
func (u *UserController) UserSignIn(ctx context.Context, req *u.UserSignInReq) (res *u.UserSignInResp, err error) {
	logx.Info("UserSignIn", req)
	return u.UserService.UserSignIn(ctx, req)
}
